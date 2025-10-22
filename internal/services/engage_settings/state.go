package engage_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Engage settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceEngageSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	engageSettingsConfig := map[string]any{
		"is_enabled": resp.IsEnabled,
	}

	for key, val := range engageSettingsConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
