// activationcode_state.go
package activationcode

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Activation Code information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resource *jamfpro.ResourceActivationCode) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	err = d.Set("organization_name", resource.OrganizationName)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("code", resource.Code)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
