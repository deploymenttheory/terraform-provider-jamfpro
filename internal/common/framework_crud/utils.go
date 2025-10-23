package framework_crud

import (
	"context"
	"math"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// extractResourceID attempts to extract the ID from the state for logging purposes
func extractResourceID(ctx context.Context, state tfsdk.State) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "unknown"
		}
	}()

	if state.Raw.IsNull() || !state.Raw.IsKnown() {
		return "unknown"
	}

	var idValue types.String
	diags := state.GetAttribute(ctx, path.Root("id"), &idValue)
	if diags.HasError() || idValue.IsNull() || idValue.IsUnknown() {
		return "unknown"
	}
	return idValue.ValueString()
}

// calculateBackoffDelay calculates the delay for exponential backoff
func calculateBackoffDelay(attempt int, opts ReadWithRetryOptions) time.Duration {
	if attempt == 0 {
		return opts.InitialRetryInterval
	}

	delay := float64(opts.InitialRetryInterval) * math.Pow(opts.BackoffMultiplier, float64(attempt))
	if delay > float64(opts.MaxRetryInterval) {
		delay = float64(opts.MaxRetryInterval)
	}

	return time.Duration(delay)
}
