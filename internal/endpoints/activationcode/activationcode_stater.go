// activationcode_state.go
package activationcode

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Activation Code information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceActivationCode) diag.Diagnostics {
	var diags diag.Diagnostics

	activationCodeData := map[string]interface{}{
		"organization_name": resource.OrganizationName,
		"code":              resource.Code,
	}

	for key, val := range activationCodeData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}
