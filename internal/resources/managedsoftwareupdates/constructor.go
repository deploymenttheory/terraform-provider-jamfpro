package managedsoftwareupdates

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceCreateManagedSoftwareUpdatePlan object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceCreateManagedSoftwareUpdatePlan, error) {
	resource := &jamfpro.ResourceCreateManagedSoftwareUpdatePlan{
		Config: jamfpro.ManagedSoftwareUpdatePlanConfig{
			UpdateAction:              d.Get("config.0.update_action").(string),
			VersionType:               d.Get("config.0.version_type").(string),
			SpecificVersion:           d.Get("config.0.specific_version").(string),
			MaxDeferrals:              d.Get("config.0.max_deferrals").(int),
			ForceInstallLocalDateTime: d.Get("config.0.force_install_local_date_time").(string),
		},
	}

	// Handle devices
	if v, ok := d.GetOk("devices"); ok {
		devices := v.([]interface{})
		if len(devices) > 0 {
			device := devices[0].(map[string]interface{})
			resource.Devices = []jamfpro.ManagedSoftwareUpdatePlanObject{
				{
					DeviceId:   device["device_id"].(string),
					ObjectType: device["object_type"].(string),
				},
			}
		}
	}

	// Handle group
	if v, ok := d.GetOk("group"); ok {
		groups := v.([]interface{})
		if len(groups) > 0 {
			group := groups[0].(map[string]interface{})
			resource.Group = jamfpro.ManagedSoftwareUpdatePlanObject{
				GroupId:    group["group_id"].(string),
				ObjectType: group["object_type"].(string),
			}
		}
	}

	// Optional: Add build_version if it's set
	if v, ok := d.GetOk("config.0.build_version"); ok {
		resource.Config.BuildVersion = v.(string)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Managed Software Update Plan: %+v\n", resource)

	return resource, nil
}
