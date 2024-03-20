// waitfor.go
package waitfor

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// APICallFunc is a generic function type for API calls that can handle different types of IDs.
type APICallFunc func(interface{}) (interface{}, error)

// ResourceIsAvailable employs a retry mechanism with exponential backoff and jitter to wait
// for a resource to become available. This function is particularly useful in scenarios
// where a resource creation is asynchronous and may not be immediately available after a
// create API call.
//
// The function uses an APICallFunc to repeatedly check for the existence of the resource,
// retrying in the face of "resource not found" errors, which are common immediately after
// resource creation. Other types of errors are not retried and lead to an immediate return.
//
// Exponential backoff helps in efficiently spacing out retry attempts to reduce load on the
// server and minimize the chance of failures due to rate limiting or server overload. Jitter
// is added to the backoff duration to prevent retry storms in scenarios with many concurrent
// operations.
//
// The retry process respects the provided context's deadline, ensuring that the function does
// not exceed the overall timeout specified for the resource creation operation in Terraform.
// This approach ensures robustness in transient network issues or temporary server-side
// unavailability.
//
// Parameters:
//   - ctx: The context governing the retry operation, carrying timeout and cancellation signals.
//   - d: The Terraform resource data schema instance, providing access to the resource's operational timeout settings.
//   - resourceID: The unique identifier of the resource being waited on.
//   - checkResourceExists: A function conforming to the APICallFunc type that attempts to fetch the resource by its ID, returning the resource or an error.
//
// Returns:
//   - interface{}: The successfully fetched resource if available, needing type assertion to the expected resource type by the caller.
//   - diag.Diagnostics: Diagnostic information including any errors encountered during the wait operation, or warnings related to the resource's availability state.
func ResourceIsAvailable(ctx context.Context, d *schema.ResourceData, resourceID interface{}, checkResourceExists APICallFunc) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var lastError error
	var resource interface{}

	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoffFactor := 2.0
	jitterFactor := 0.5

	currentBackoff := initialBackoff

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = checkResourceExists(resourceID)
		if apiErr != nil {
			lastError = apiErr

			// Check specifically for "resource not found" errors to retry
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Apply exponential backoff with jitter
				time.Sleep(currentBackoff + time.Duration(rand.Float64()*jitterFactor*float64(currentBackoff)))
				currentBackoff = time.Duration(float64(currentBackoff) * backoffFactor) // Corrected line
				if currentBackoff > maxBackoff {
					currentBackoff = maxBackoff
				}
				return retry.RetryableError(apiErr)
			}

			// For other types of errors, do not retry and return the error
			return retry.NonRetryableError(apiErr)
		}

		// If no error, the resource exists, stop retrying
		lastError = nil
		return nil
	})

	// If an error occurred during retries (other than the resource not found),
	// add it to diagnostics
	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("error waiting for resource with ID '%v' to become available: %v", resourceID, lastError))...)
		return nil, diags // Return nil as the resource and the diagnostics
	}

	// Wait for an additional 5 seconds before concluding the wait process
	// This can be helpful in scenarios where the resource might need a few extra moments to stabilize or propagate
	time.Sleep(5 * time.Second)

	// Return the successfully fetched resource and any diagnostics
	return resource, diags
}
