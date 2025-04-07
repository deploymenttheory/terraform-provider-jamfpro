package appinstallerglobalsettings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Client Check-In information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseJamfAppCatalogGlobalSettings) diag.Diagnostics {
	var diags diag.Diagnostics
	settings := resp.EndUserExperienceSettings

	stateData := map[string]interface{}{
		"notification_message":  settings.NotificationMessage,
		"notification_interval": settings.NotificationInterval,
		"deadline_message":      settings.DeadlineMessage,
		"deadline":              settings.Deadline,
		"quit_delay":            settings.QuitDelay,
		"complete_message":      settings.CompleteMessage,
		"relaunch":              settings.Relaunch,
		"suppress":              settings.Suppress,
	}

	for key, val := range stateData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
