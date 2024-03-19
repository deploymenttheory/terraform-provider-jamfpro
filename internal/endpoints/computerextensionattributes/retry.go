// computerextensionattributes_resource.go
package computerextensionattributes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// APICallFunc is a function type that represents an API call to fetch a resource by its ID.
type APICallFunc func(int) (interface{}, error)

// retryAPIReadCall executes an API read call with retry logic and handles errors and resource state.
func retryAPIReadCall(ctx context.Context, d *schema.ResourceData, resourceID int, apiCall APICallFunc) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var lastError error
	var resource interface{}

	// Execute the retry logic
	retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = apiCall(resourceID)
		if apiErr != nil {
			lastError = apiErr
			// Treat all errors as retryable until the context deadline is reached
			return retry.RetryableError(apiErr)
		}
		lastError = nil // Reset last error on success
		return nil      // Success, exit retry loop
	})

	// Check the last error after retries are completed
	if lastError != nil {
		if strings.Contains(lastError.Error(), "404") || strings.Contains(lastError.Error(), "410") {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Resource with ID '%d' was not found on the server after all retries and is marked for deletion from Terraform state.", resourceID),
			})
		} else {
			// For other errors after all retries, return an error diagnostic
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to read resource with ID '%d' after all retries: %v", resourceID, lastError))...)
		}
	}

	return resource, diags
}
