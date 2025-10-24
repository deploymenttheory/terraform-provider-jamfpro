package managed_software_update

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
		groupList := v.([]any)
		if len(groupList) > 0 {
			group := groupList[0].(map[string]any)
			resource.Group = jamfpro.ResourcManagedSoftwareUpdatePlanObject{
				GroupId:    group["group_id"].(string),
				ObjectType: group["object_type"].(string),
			}
		}
	}

	// Handle device
	if v, ok := d.GetOk("device"); ok {
		deviceList := v.([]any)
		if len(deviceList) > 0 {
			device := deviceList[0].(map[string]any)
			resource.Devices = []jamfpro.ResourcManagedSoftwareUpdatePlanObject{
				{
					DeviceId:   device["device_id"].(string),
					ObjectType: device["object_type"].(string),
				},
			}
		}
	}

	// Now handle the fields that were previously inside the config block
	resource.Config.UpdateAction = d.Get("update_action").(string)
	resource.Config.VersionType = d.Get("version_type").(string)

	if v, ok := d.GetOk("specific_version"); ok {
		resource.Config.SpecificVersion = v.(string)
	} else {
		resource.Config.SpecificVersion = "NO_SPECIFIC_VERSION"
	}

	if v, ok := d.GetOk("build_version"); ok {
		resource.Config.BuildVersion = v.(string)
	} else {
		resource.Config.BuildVersion = ""
	}

	if v, ok := d.GetOk("max_deferrals"); ok {
		resource.Config.MaxDeferrals = v.(int)
	} else {
		resource.Config.MaxDeferrals = 0
	}

	if v, ok := d.GetOk("force_install_local_date_time"); ok {
		resource.Config.ForceInstallLocalDateTime = v.(string)
	} else {
		resource.Config.ForceInstallLocalDateTime = ""
	}

	// Debug logging for config specifically
	configJSON, err := json.MarshalIndent(resource.Config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro managed software update config to JSON: %v", err)
	}
	log.Printf("[DEBUG] Constructed Jamf Pro managed software update config JSON:\n%s\n", string(configJSON))

	return resource, nil
}
