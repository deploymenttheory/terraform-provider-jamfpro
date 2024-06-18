// waitfor.go
package waitfor

import (
	"context"
	"fmt"
	"log"
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
// for a resource to become available. This function is useful in scenarios where resource
// creation is asynchronous and may not be immediately available after a create API call.
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
// After successfully locating the resource, the function initiates a customizable stabilization
// period, allowing time for the resource to reach a steady state. The duration of this period
// is specified by the stabilizationTime parameter. Both the number of retries and the duration
// of the stabilization period are included in the logging to provide insight into the wait process.
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
//   - stabilizationTime: The duration of the stabilization period to wait after the resource is found before concluding the wait process.
//
// Returns:
//   - interface{}: The successfully fetched resource if available, needing type assertion to the expected resource type by the caller.
//   - diag.Diagnostics: Diagnostic information including any errors encountered during the wait operation, or warnings related to the resource's availability state.
func ResourceIsAvailable(ctx context.Context, d *schema.ResourceData, resourceType string, resourceID interface{}, checkResourceExists APICallFunc, defaultPropagationTime time.Duration) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var resource interface{}
	var retryCount int
	var propagationTime time.Duration

	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoffFactor := 2.0
	jitterFactor := 0.5

	currentBackoff := initialBackoff

	log.Printf("Starting to wait for %s resource with ID '%v'", resourceType, resourceID)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retryCount++
		log.Printf("Attempting to fetch %s resource with ID '%v' (Retry #%d)", resourceType, resourceID, retryCount)
		var apiErr error
		resource, apiErr = checkResourceExists(resourceID)
		if apiErr != nil {
			log.Printf("Error fetching %s resource with ID '%v': %v (Retry #%d)", resourceType, resourceID, apiErr, retryCount)

			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				log.Printf("Resource with ID '%v' not found, retrying with backoff of %v (Retry #%d)", resourceID, currentBackoff, retryCount)
				time.Sleep(currentBackoff + time.Duration(rand.Float64()*jitterFactor*float64(currentBackoff)))
				currentBackoff = time.Duration(float64(currentBackoff) * backoffFactor)
				if currentBackoff > maxBackoff {
					currentBackoff = maxBackoff
				}
				log.Printf("Adjusted backoff for resource with ID '%v' to %v (Retry #%d)", resourceID, currentBackoff, retryCount)
				return retry.RetryableError(apiErr)
			}

			// Return a non-retryable error to exit the retry loop.
			return retry.NonRetryableError(apiErr)
		}

		propagationTime = defaultPropagationTime
		log.Printf("%s resource with ID '%v' found after %d retries. Cookie Jar is disabled, initiating the Jamf Cloud default propagation wait period of %v.", resourceType, resourceID, retryCount, propagationTime)

		time.Sleep(propagationTime)

		log.Printf("Concluding wait process for %s resource with ID '%v' after a propagation period of %v.", resourceType, resourceID, propagationTime)
		return nil
	})

	if err != nil {
		// Generate a more descriptive error message based on the context of the retry operation.
		errorMessage := fmt.Sprintf("Error waiting for %s resource with ID '%v' to become available after %d retries. Last error: %v", resourceType, resourceID, retryCount, err)
		log.Print(errorMessage)
		errorDiags := diag.FromErr(fmt.Errorf(errorMessage))
		diags = append(diags, errorDiags...)
		return nil, diags
	}

	log.Printf("Successfully waited for %s resource with ID '%v' after %d retries", resourceType, resourceID, retryCount)
	return resource, diags
}
