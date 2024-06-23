// allowedfileextensions_state.go
package allowedfileextensions

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Allowed File Extension information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceAllowedFileExtension) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("extension", resp.Extension); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
