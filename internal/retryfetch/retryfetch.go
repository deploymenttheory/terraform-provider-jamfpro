// retryfetch.go
package retryfetch

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

// APICallFuncInt is specifically for API calls that require an integer ID.
type APICallFuncInt func(int) (interface{}, error)

// APICallFuncString is specifically for API calls that require a string ID.
type APICallFuncString func(string) (interface{}, error)

// RetryAPIReadCallInt is a wrapper function that facilitates retrying API Read calls which require an integer ID.
// It adapts an API call expecting an integer ID to the generic retry mechanism provided by RetryAPIReadCallCommon.
//
// Parameters:
//   - ctx: A context.Context instance that carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
//   - d: The Terraform resource data schema instance, providing access to the operations timeout settings and the resource's state.
//   - resourceID: The unique integer identifier of the resource to be fetched.
//   - apiCall: The specific API call function that accepts an integer ID and returns the resource along with any error encountered during the fetch operation.
//
// Returns:
//   - interface{}: The resource fetched by the API call if successful. This will need to be type-asserted to the specific resource type expected by the caller.
//   - diag.Diagnostics: A collection of diagnostic information including any errors encountered during the operation or warnings related to the resource's state.
//
// Note: If the resource cannot be found or if an error occurs, appropriate diagnostics are returned to Terraform, potentially marking the resource for deletion from the state if not found.
func ByIntID(ctx context.Context, d *schema.ResourceData, resourceID int, apiCall APICallFuncInt) (interface{}, diag.Diagnostics) {
	genericAPICall := func(id interface{}) (interface{}, error) {
		intID, ok := id.(int)
		if !ok {
			return nil, fmt.Errorf("expected int ID, got %T", id)
		}
		return apiCall(intID)
	}
	return RetryAPIReadCallCommon(ctx, d, resourceID, genericAPICall)
}

// ByStringID is a wrapper function that facilitates retrying API Read calls which require a string ID.
// It adapts an API call expecting a string ID to the generic retry mechanism provided by RetryAPIReadCallCommon.
//
// Parameters:
//   - ctx: A context.Context instance that carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
//   - d: The Terraform resource data schema instance, providing access to the operations timeout settings and the resource's state.
//   - resourceID: The unique string identifier of the resource to be fetched.
//   - apiCall: The specific API call function that accepts a string ID and returns the resource along with any error encountered during the fetch operation.
//
// Returns:
//   - interface{}: The resource fetched by the API call if successful. This will need to be type-asserted to the specific resource type expected by the caller.
//   - diag.Diagnostics: A collection of diagnostic information including any errors encountered during the operation or warnings related to the resource's state.
//
// Note: If the resource cannot be found or if an error occurs, appropriate diagnostics are returned to Terraform, potentially marking the resource for deletion from the state if not found.
func ByStringID(ctx context.Context, d *schema.ResourceData, resourceID string, apiCall APICallFuncString) (interface{}, diag.Diagnostics) {
	genericAPICall := func(id interface{}) (interface{}, error) {
		strID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("expected string ID, got %T", id)
		}
		return apiCall(strID)
	}
	return RetryAPIReadCallCommon(ctx, d, resourceID, genericAPICall)
}

// RetryAPIReadCallCommon executes an API read call with retry logic to fetch a resource by its ID.
// This function incorporates exponential backoff with jitter to efficiently manage retry attempts
// in the face of transient errors or temporary server unavailability, ensuring robust error handling
// and resource state synchronization.
//
// Exponential backoff increases the delay between retry attempts exponentially, which helps to
// alleviate load on the server and reduce the likelihood of cascading failures. Jitter adds a
// random variation to the backoff periods, further helping to spread out retry attempts and
// prevent synchronization issues ("thundering herd" problem) that could occur when many clients
// retry simultaneously.
//
// The function respects the context's deadline, making repeated attempts to retrieve the resource
// until success or until the context expires. This is particularly useful in environments where
// network instability or temporary server issues might prevent a successful resource fetch on the
// first attempt.
//
// Parameters:
//   - ctx: The context governing the API call, which includes timeout/cancellation signals.
//   - d: The Terraform resource data schema, used here primarily to fetch the operation's timeout setting.
//   - resourceID: The unique identifier of the resource to be fetched.
//   - apiCall: A function that encapsulates the specific API call needed to fetch the resource by its ID.
//     This function must conform to the APICallFunc type, accepting an integer ID and returning
//     an interface{} (the resource) and an error.
//
// Returns:
//   - interface{}: The resource fetched by the API call, if successful. This will need to be type-asserted
//     to the specific resource type expected by the caller.
//   - diag.Diagnostics: A collection of diagnostic information that includes errors encountered during
//     the operation or warnings about the resource's state (e.g., if the resource is not found).
//
// Note: If the resource cannot be found (indicated by 404 or 410 HTTP status codes in the error message),
// the function marks the resource for deletion from the Terraform state by clearing its ID. This signifies
// to Terraform that the resource no longer exists and should be removed from the state file.
func RetryAPIReadCallCommon(ctx context.Context, d *schema.ResourceData, resourceID interface{}, apiCall APICallFunc) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var lastError error
	var resource interface{}

	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoffFactor := 2.0
	jitterFactor := 0.5

	currentBackoff := initialBackoff

	retryErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = apiCall(resourceID)
		if apiErr != nil {
			lastError = apiErr

			time.Sleep(currentBackoff + time.Duration(rand.Float64()*jitterFactor*float64(currentBackoff)))

			currentBackoff *= time.Duration(backoffFactor)
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}

			return retry.RetryableError(apiErr)
		}
		lastError = nil
		return nil
	})

	if retryErr != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("retry logic failed: %v", retryErr))...)
	}

	if lastError != nil {
		if strings.Contains(lastError.Error(), "404") || strings.Contains(lastError.Error(), "410") {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Resource with ID '%v' was not found on the server after all retries and is marked for deletion from Terraform state.", resourceID),
			})
		} else {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to read resource with ID '%v' after all retries: %v", resourceID, lastError))...)
		}
	}

	return resource, diags
}
