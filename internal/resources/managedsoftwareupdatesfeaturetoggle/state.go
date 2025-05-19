package managedsoftwareupdatesfeaturetoggle

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Managed Software Updates Feature Toggle information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceManagedSoftwareUpdateFeatureToggle) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("toggle", resp.Toggle); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
