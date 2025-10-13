package crud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadWithRetry executes a read operation with retry logic within the context timeout
// It repeatedly calls the provided read function until success or context timeout
func ReadWithRetry(
	ctx context.Context,
	readFunc func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse),
	readReq resource.ReadRequest,
	stateContainer StateContainer,
	opts ReadWithRetryOptions,
) error {
	resourceID := extractResourceID(ctx, stateContainer.GetState())
	resourceType := opts.ResourceTypeName
	if resourceType == "" {
		resourceType = "resource"
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting read with retry for %s operation", opts.Operation), map[string]any{
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})

	// Ensure we have reasonable defaults
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = 30
	}
	if opts.InitialRetryInterval <= 0 {
		opts.InitialRetryInterval = 2 * time.Second
	}
	if opts.MaxRetryInterval <= 0 {
		opts.MaxRetryInterval = 30 * time.Second
	}
	if opts.BackoffMultiplier <= 0 {
		opts.BackoffMultiplier = 1.5
	}
	if opts.Operation == "" {
		opts.Operation = "Operation"
	}

	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		return fmt.Errorf("context must have a deadline for retry operations")
	}

	timeRemaining := time.Until(deadline) - time.Second
	if timeRemaining <= 0 {
		return fmt.Errorf("insufficient time remaining in context for retry operation")
	}

	tflog.Debug(ctx, fmt.Sprintf("Will attempt up to %d retries with exponential backoff", opts.MaxRetries), map[string]any{
		"resource_id":   resourceID,
		"resource_type": resourceType,
		"initial_delay": opts.InitialRetryInterval,
		"max_delay":     opts.MaxRetryInterval,
		"multiplier":    opts.BackoffMultiplier,
	})

	var lastErr error
	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry attempt %d: %w", attempt, ctx.Err())
		default:
		}

		delay := calculateBackoffDelay(attempt-1, opts)
		if attempt > 0 && time.Until(deadline) < delay {
			tflog.Debug(ctx, "Insufficient time remaining for another retry attempt", map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"next_delay":    delay,
			})
			break
		}

		tflog.Debug(ctx, fmt.Sprintf("Read retry attempt %d/%d", attempt+1, opts.MaxRetries+1), map[string]any{
			"resource_id":   resourceID,
			"resource_type": resourceType,
		})

		readResp := &resource.ReadResponse{State: stateContainer.GetState()}
		ctxWithOp := context.WithValue(ctx, "retry_operation", opts.Operation)

		readFunc(ctxWithOp, readReq, readResp)

		// No error
		if !readResp.Diagnostics.HasError() {
			tflog.Debug(ctx, fmt.Sprintf("Read successful on attempt %d", attempt+1), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			stateContainer.SetState(readResp.State)
			return nil
		}

		lastErr = fmt.Errorf(
			"error reading resource state after %s method on attempt %d: %s",
			opts.Operation,
			attempt+1,
			readResp.Diagnostics.Errors(),
		)

		errorInfo := extractErrorFromDiagnostics(readResp.Diagnostics)

		if !isRetryableReadError(errorInfo.StatusCode) {
			tflog.Error(ctx, fmt.Sprintf("Read failed on attempt %d (non-retryable error)", attempt+1), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"status_code":   errorInfo.StatusCode,
				"error_code":    errorInfo.ErrorCode,
				"diagnostics":   readResp.Diagnostics.Errors(),
			})
			return fmt.Errorf("read operation failed with non-retryable error: %w", lastErr)
		}

		if attempt < opts.MaxRetries {
			tflog.Warn(ctx, fmt.Sprintf("Read failed on attempt %d (retryable error), waiting %s before retry", attempt+1, delay), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"status_code":   errorInfo.StatusCode,
				"error_code":    errorInfo.ErrorCode,
				"diagnostics":   readResp.Diagnostics.Errors(),
				"next_delay":    delay,
			})

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
			}
		} else {
			tflog.Error(ctx, fmt.Sprintf("Read failed on final attempt %d", attempt+1), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"status_code":   errorInfo.StatusCode,
				"error_code":    errorInfo.ErrorCode,
				"diagnostics":   readResp.Diagnostics.Errors(),
			})
		}

	}

	if lastErr != nil {
		return fmt.Errorf("failed to read resource state for %s after %d attempts: %w", resourceType, opts.MaxRetries+1, lastErr)
	}

	return fmt.Errorf("failed to read resource state for %s after %d attempts", resourceType, opts.MaxRetries+1)
}
