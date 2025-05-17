package managedsoftwareupdatesfeaturetoggle

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Managed Software Updates Feature Toggle information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceManagedSoftwareUpdateFeatureToggle) diag.Diagnostics {
	var diags diag.Diagnostics

	managedSoftwareUpdatesFeatureToggleConfig := map[string]interface{}{
		"toggle": resp.Toggle,
	}

	for key, val := range managedSoftwareUpdatesFeatureToggleConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
