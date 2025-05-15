package engagesettings

// engagesettings_resource.go
// Package engagesettings provides the schema and CRUD operations for managing Jamf Pro Engage Settings in Terraform.
// This package includes functions to create, read, update, and delete the Engage Settings configuration.
// It also includes a function to construct the ResourceEngageSettings object from the schema data.

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructEngageSettings constructs a ResourceEngageSettings object from the provided schema data
func constructEngageSettings(d *schema.ResourceData) (*jamfpro.ResourceEngageSettings, error) {
	resource := &jamfpro.ResourceEngageSettings{
		IsEnabled: d.Get("is_enabled").(bool),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Engage Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Engage Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
