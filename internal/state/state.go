package state

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HandleResourceNotFound(ctx context.Context, d *schema.ResourceData, resourceID string, err error, diags *diag.Diagnostics) diag.Diagnostics {
	if err.Error() == "resource not found, marked for deletion" {
		// Resource not found, remove from Terraform state
		d.SetId("")
		// Append a warning diagnostic and return
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource not found",
			Detail:   fmt.Sprintf("Resource with ID '%s' was not found on the server and is marked for deletion from Terraform state.", resourceID),
		})
		return *diags
	}
	// For other errors, return an error diagnostic
	*diags = append(*diags, diag.FromErr(fmt.Errorf("failed to read resource with ID '%s' after retries: %v", resourceID, err))...)
	return *diags
}
