package appinstallerglobalsettings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructAppInstallerGlobalSettings constructs a ResponseJamfAppCatalogGlobalSettings object from the provided schema data.
func constructAppInstallerGlobalSettings(d *schema.ResourceData) (*jamfpro.ResponseJamfAppCatalogGlobalSettings, error) {
	settings := &jamfpro.JamfAppCatalogDeploymentSubsetNotificationSettings{}

	if v, ok := d.GetOk("notification_message"); ok {
		settings.NotificationMessage = v.(string)
	}
	if v, ok := d.GetOk("notification_interval"); ok {
		settings.NotificationInterval = v.(int)
	}
	if v, ok := d.GetOk("deadline_message"); ok {
		settings.DeadlineMessage = v.(string)
	}
	if v, ok := d.GetOk("deadline"); ok {
		settings.Deadline = v.(int)
	}
	if v, ok := d.GetOk("quit_delay"); ok {
		settings.QuitDelay = v.(int)
	}
	if v, ok := d.GetOk("complete_message"); ok {
		settings.CompleteMessage = v.(string)
	}
	if v, ok := d.GetOk("relaunch"); ok {
		val := v.(bool)
		settings.Relaunch = &val
	}
	if v, ok := d.GetOk("suppress"); ok {
		val := v.(bool)
		settings.Suppress = &val
	}

	resource := &jamfpro.ResponseJamfAppCatalogGlobalSettings{
		EndUserExperienceSettings: *settings,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Jamf App Catalog Global Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf App Catalog Global Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
