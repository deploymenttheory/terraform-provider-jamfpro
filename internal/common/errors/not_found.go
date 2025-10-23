package errors

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HandleResourceNotFoundError is a helper function to handle 404 and 410 errors and remove the resource from Terraform state
func HandleResourceNotFoundError(err error, d *schema.ResourceData, cleanup bool) diag.Diagnostics {
	var diags diag.Diagnostics
	ErrorTypeIsNotFound := strings.Contains(err.Error(), "404")

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
