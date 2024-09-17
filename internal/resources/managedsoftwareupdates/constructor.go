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
			resource.Config.UpdateAction = config["update_action"].(string)
			resource.Config.VersionType = config["version_type"].(string)

			// Set SpecificVersion, default to "NO_SPECIFIC_VERSION"
			if v, ok := config["specific_version"]; ok && v.(string) != "" {
				resource.Config.SpecificVersion = v.(string)
			} else {
				resource.Config.SpecificVersion = "NO_SPECIFIC_VERSION"
			}

			// Set BuildVersion, default to empty string if not present
			if v, ok := config["build_version"]; ok && v.(string) != "" {
				resource.Config.BuildVersion = v.(string)
			} else {
				resource.Config.BuildVersion = ""
			}

			// Set MaxDeferrals, default to 0 if not applicable
			if v, ok := config["max_deferrals"]; ok {
				resource.Config.MaxDeferrals = v.(int)
			} else {
				resource.Config.MaxDeferrals = 0
			}

			// Set ForceInstallLocalDateTime, default to empty string if not present
			if v, ok := config["force_install_local_date_time"]; ok && v.(string) != "" {
				resource.Config.ForceInstallLocalDateTime = v.(string)
			} else {
				resource.Config.ForceInstallLocalDateTime = ""
			}
		}
	}

	// Debug logging for config specifically
	configJSON, err := json.MarshalIndent(resource.Config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro managed software update config to JSON: %v", err)
	}
	log.Printf("[DEBUG] Constructed Jamf Pro managed software update config JSON:\n%s\n", string(configJSON))

	return resource, nil
}
