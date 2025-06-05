package account_driven_user_enrollment_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds the request for UpdateADUESessionTokenSettings
func construct(d *schema.ResourceData) (*jamfpro.ResourceADUETokenSettings, error) {
	resource := &jamfpro.ResourceADUETokenSettings{
		Enabled: d.Get("enabled").(bool),
	}

	// Add optional fields if they exist in the schema
	if v, ok := d.GetOk("expiration_interval_days"); ok {
		resource.ExpirationIntervalDays = v.(int)
	}

	if v, ok := d.GetOk("expiration_interval_seconds"); ok {
		resource.ExpirationIntervalSeconds = v.(int)
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Account Driven User Enrollment Settings to JSON: %v", err)
	}
	log.Printf("[DEBUG] Constructed Account Driven User Enrollment Settings resource:\n%s", string(resourceJSON))

	return resource, nil
}
