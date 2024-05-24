// common/state/state.go
// This package contains shared / common resource functions for stating

package state

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HandleResourceNotFoundError is a helper function to handle 404 and 410 errors and remove the resource from Terraform state
func HandleResourceNotFoundError(err error, d *schema.ResourceData) diag.Diagnostics {
	if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "410") {
		d.SetId("") // Remove the resource from Terraform state
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   "The resource was not found on the remote server. It has been removed from the Terraform state.",
			},
		}
	} else {
		// For other errors, return a diagnostic error
		return diag.FromErr(err)
	}
}
