package managedsoftwareupdates

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceManagedSoftwareUpdatePlan object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceManagedSoftwareUpdatePlan, error) {
	resource := &jamfpro.ResourceManagedSoftwareUpdatePlan{
		Config: jamfpro.ResourcManagedSoftwareUpdatePlanConfig{},
	}

	// Handle group
	if v, ok := d.GetOk("group"); ok {
		groupList := v.([]interface{})
		if len(groupList) > 0 {
			group := groupList[0].(map[string]interface{})
			resource.Group = jamfpro.ResourcManagedSoftwareUpdatePlanObject{
				GroupId:    group["group_id"].(string),
				ObjectType: group["object_type"].(string),
			}
		}
	}

	// Handle device
	if v, ok := d.GetOk("device"); ok {
		deviceList := v.([]interface{})
		if len(deviceList) > 0 {
			device := deviceList[0].(map[string]interface{})
			resource.Devices = []jamfpro.ResourcManagedSoftwareUpdatePlanObject{
				{
					DeviceId:   device["device_id"].(string),
					ObjectType: device["object_type"].(string),
				},
			}
		}
	}

	// Handle config
	if v, ok := d.GetOk("config"); ok {
		configList := v.([]interface{})
		if len(configList) > 0 {
			config := configList[0].(map[string]interface{})

			// Set common required fields
			resource.Config = jamfpro.ResourcManagedSoftwareUpdatePlanConfig{
				UpdateAction: config["update_action"].(string),
				VersionType:  config["version_type"].(string),
			}

			// Conditionally set specificVersion
			if config["version_type"].(string) == "SPECIFIC_VERSION" || config["version_type"].(string) == "CUSTOM_VERSION" {
				resource.Config.SpecificVersion = config["specific_version"].(string)
			} else {
				resource.Config.SpecificVersion = "NO_SPECIFIC_VERSION"
			}

			// Conditionally set buildVersion if version_type is CUSTOM_VERSION
			if config["version_type"].(string) == "CUSTOM_VERSION" {
				resource.Config.BuildVersion = config["build_version"].(string)
			}

			// Conditionally set maxDeferrals if update_action is DOWNLOAD_INSTALL_ALLOW_DEFERRAL
			if config["update_action"].(string) == "DOWNLOAD_INSTALL_ALLOW_DEFERRAL" {
				resource.Config.MaxDeferrals = config["max_deferrals"].(int)
			}

			// Conditionally set forceInstallLocalDateTime if provided
			if v, ok := config["force_install_local_date_time"]; ok && v.(string) != "" {
				resource.Config.ForceInstallLocalDateTime = v.(string)
			}
		}
	}

	// Debug logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro managed software update to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro managed software update JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
