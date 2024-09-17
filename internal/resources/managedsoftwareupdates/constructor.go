package managedsoftwareupdates

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceManagedSoftwareUpdatePlan object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceManagedSoftwareUpdatePlan, error) {
	resource := &jamfpro.ResourceManagedSoftwareUpdatePlan{
		Config: jamfpro.ResourcManagedSoftwareUpdatePlanConfig{
			UpdateAction:              d.Get("config.0.update_action").(string),
			VersionType:               d.Get("config.0.version_type").(string),
			SpecificVersion:           d.Get("config.0.specific_version").(string),
			MaxDeferrals:              d.Get("config.0.max_deferrals").(int),
			ForceInstallLocalDateTime: d.Get("config.0.force_install_local_date_time").(string),
		},
	}

	// Handle group
	if v, ok := d.GetOk("group"); ok {
		groups := v.([]interface{})
		if len(groups) > 0 {
			group := groups[0].(map[string]interface{})
			resource.Group = jamfpro.ResourcManagedSoftwareUpdatePlanObject{
				GroupId:    group["group_id"].(string),
				ObjectType: group["object_type"].(string),
			}
		}
	}

	// Optional: Add build_version if it's set
	if v, ok := d.GetOk("config.0.build_version"); ok {
		resource.Config.BuildVersion = v.(string)
	}

	// Validate that specific_version is set when required
	if resource.Config.VersionType == "SPECIFIC_VERSION" || resource.Config.VersionType == "CUSTOM_VERSION" {
		if resource.Config.SpecificVersion == "" {
			return nil, fmt.Errorf("specific_version is required when version_type is SPECIFIC_VERSION or CUSTOM_VERSION")
		}
	}

	// Validate that max_deferrals is set when required
	if resource.Config.UpdateAction == "DOWNLOAD_INSTALL_ALLOW_DEFERRAL" && resource.Config.MaxDeferrals == 0 {
		return nil, fmt.Errorf("max_deferrals is required when update_action is DOWNLOAD_INSTALL_ALLOW_DEFERRAL")
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Managed Software Update Plan: %+v\n", resource)

	return resource, nil
}
