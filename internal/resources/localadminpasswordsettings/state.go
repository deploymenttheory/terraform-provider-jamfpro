package localadminpasswordsettings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest local account password settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceLocalAdminPasswordSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	localAdminPasswordSettingsConfig := map[string]interface{}{
		"auto_deploy_enabled":                 resp.AutoDeployEnabled,
		"password_rotation_time_seconds":      resp.PasswordRotationTime,
		"auto_rotate_enabled":                 resp.AutoRotateEnabled,
		"auto_rotate_expiration_time_seconds": resp.AutoRotateExpirationTime,
	}

	for key, val := range localAdminPasswordSettingsConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
