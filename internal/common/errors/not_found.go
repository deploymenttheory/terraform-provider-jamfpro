package errors

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IsNotFoundError returns true when err represents a resource-not-found condition.
// It matches HTTP 404 status codes and common SDK not-found message strings.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "404") ||
		strings.Contains(msg, "resource with name does not exist") ||
		strings.Contains(msg, "resource with id does not exist") ||
		strings.Contains(msg, "does not exist")
}

// HandleResourceNotFoundError is a helper function to handle 404 and 410 errors and remove the resource from Terraform state
func HandleResourceNotFoundError(err error, d *schema.ResourceData, cleanup bool) diag.Diagnostics {
	var diags diag.Diagnostics
	ErrorTypeIsNotFound := IsNotFoundError(err)

	if cleanup && ErrorTypeIsNotFound {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource not found and will be redeployed",
		})

	} else {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}
