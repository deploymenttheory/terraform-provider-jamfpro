package managedsoftwareupdates

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceManagedSoftwareUpdatePlan object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceManagedSoftwareUpdatePlan, error) {
	resource := &jamfpro.ResourceManagedSoftwareUpdatePlan{
		Config: jamfpro.ResourcManagedSoftwareUpdatePlanConfig{},
	}

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

	if v, ok := d.GetOk("config"); ok {
		configList := v.([]interface{})
		if len(configList) > 0 {
			config := configList[0].(map[string]interface{})
			resource.Config = jamfpro.ResourcManagedSoftwareUpdatePlanConfig{
				UpdateAction: config["update_action"].(string),
				VersionType:  config["version_type"].(string),
			}

			if v, ok := config["specific_version"]; ok {
				resource.Config.SpecificVersion = v.(string)
			}
			if v, ok := config["build_version"]; ok {
				resource.Config.BuildVersion = v.(string)
			}
			if v, ok := config["max_deferrals"]; ok {
				resource.Config.MaxDeferrals = v.(int)
			}
			if v, ok := config["force_install_local_date_time"]; ok {
				resource.Config.ForceInstallLocalDateTime = v.(string)
			}
		}
	}

	// Validation
	if resource.Config.VersionType == "SPECIFIC_VERSION" || resource.Config.VersionType == "CUSTOM_VERSION" {
		if resource.Config.SpecificVersion == "" {
			return nil, fmt.Errorf("specific_version is required when version_type is SPECIFIC_VERSION or CUSTOM_VERSION")
		}
	}

	if resource.Config.UpdateAction == "DOWNLOAD_INSTALL_ALLOW_DEFERRAL" && resource.Config.MaxDeferrals == 0 {
		return nil, fmt.Errorf("max_deferrals is required when update_action is DOWNLOAD_INSTALL_ALLOW_DEFERRAL")
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro managed software update '%s' to JSON: %v", resource.Group.GroupId, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro managed software update JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
