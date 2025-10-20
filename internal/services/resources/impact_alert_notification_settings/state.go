package impact_alert_notification_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Impact Alert Notification Settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceImpactAlertNotificationSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	impactAlertSettingsConfig := map[string]interface{}{
		"scopeable_objects_alert_enabled":              resp.ScopeableObjectsAlertEnabled,
		"scopeable_objects_confirmation_code_enabled":  resp.ScopeableObjectsConfirmationCodeEnabled,
		"deployable_objects_alert_enabled":             resp.DeployableObjectsAlertEnabled,
		"deployable_objects_confirmation_code_enabled": resp.DeployableObjectsConfirmationCodeEnabled,
	}

	for key, val := range impactAlertSettingsConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
