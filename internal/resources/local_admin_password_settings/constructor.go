package local_admin_password_settings

// Package local_admin_password_settings provides the schema and CRUD operations for managing Jamf Pro Local Admin Password Settings in Terraform.
// This package includes functions to create, read, update, and delete the Local Admin Password Settings configuration.
// It also includes a function to construct the ResourceLocalAdminPasswordSettings object from the schema data.

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructLocalAdminPasswordSettings constructs a ResourceLocalAdminPasswordSettings object from the provided schema data
func constructLocalAdminPasswordSettings(d *schema.ResourceData) (*jamfpro.ResourceLocalAdminPasswordSettings, error) {
	resource := &jamfpro.ResourceLocalAdminPasswordSettings{
		AutoDeployEnabled:        d.Get("auto_deploy_enabled").(bool),
		PasswordRotationTime:     d.Get("password_rotation_time_seconds").(int),
		AutoRotateEnabled:        d.Get("auto_rotate_enabled").(bool),
		AutoRotateExpirationTime: d.Get("auto_rotate_expiration_time_seconds").(int),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Local Admin Password Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Local Admin Password Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
