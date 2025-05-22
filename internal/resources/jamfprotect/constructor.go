package jamfprotect

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProtectResources struct {
	Registration *jamfpro.ResourceJamfProtectRegistration
	Settings     *jamfpro.ResourceJamfProtectSettings
}

// construct creates a new instance of Jamf Protect registration and settings based on the provided schema.
func construct(d *schema.ResourceData) (*ProtectResources, error) {
	registration := &jamfpro.ResourceJamfProtectRegistration{
		ProtectURL: d.Get("protect_url").(string),
		ClientID:   d.Get("client_id").(string),
		Password:   d.Get("password").(string),
	}

	registrationJSON, err := json.MarshalIndent(registration, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Protect Registration to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Protect Registration JSON:\n%s\n", string(registrationJSON))

	settings := &jamfpro.ResourceJamfProtectSettings{
		AutoInstall: d.Get("auto_install").(bool),
	}

	settingsJSON, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Protect Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Protect Settings JSON:\n%s\n", string(settingsJSON))

	return &ProtectResources{
		Registration: registration,
		Settings:     settings,
	}, nil
}
