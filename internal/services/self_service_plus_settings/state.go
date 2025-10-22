package self_service_plus_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Self Service Plus settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSelfServicePlusSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	selfServicePlusSettingsConfig := map[string]any{
		"enabled": resp.Enabled,
	}

	for key, val := range selfServicePlusSettingsConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
