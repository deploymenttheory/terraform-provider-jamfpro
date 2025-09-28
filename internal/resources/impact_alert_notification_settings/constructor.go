package impact_alert_notification_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrMarshalImpactAlertSettings = fmt.Errorf("failed to marshal Jamf Pro Impact Alert Notification Settings to JSON")
)

// constructImpactAlertNotificationSettings constructs a ResourceImpactAlertNotificationSettings object from the provided schema data
func constructImpactAlertNotificationSettings(d *schema.ResourceData) (*jamfpro.ResourceImpactAlertNotificationSettings, error) {
	resource := &jamfpro.ResourceImpactAlertNotificationSettings{
		ScopeableObjectsAlertEnabled:             d.Get("scopeable_objects_alert_enabled").(bool),
		ScopeableObjectsConfirmationCodeEnabled:  d.Get("scopeable_objects_confirmation_code_enabled").(bool),
		DeployableObjectsAlertEnabled:            d.Get("deployable_objects_alert_enabled").(bool),
		DeployableObjectsConfirmationCodeEnabled: d.Get("deployable_objects_confirmation_code_enabled").(bool),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMarshalImpactAlertSettings, err.Error())
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Impact Alert Notification Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
