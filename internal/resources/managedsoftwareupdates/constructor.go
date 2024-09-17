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
			resource.Config = jamfpro.ResourcManagedSoftwareUpdatePlanConfig{
				UpdateAction:              config["update_action"].(string),
				VersionType:               config["version_type"].(string),
				SpecificVersion:           config["specific_version"].(string),
				BuildVersion:              config["build_version"].(string),
				MaxDeferrals:              config["max_deferrals"].(int),
				ForceInstallLocalDateTime: config["force_install_local_date_time"].(string),
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
