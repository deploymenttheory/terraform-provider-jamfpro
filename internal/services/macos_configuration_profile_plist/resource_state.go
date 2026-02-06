package macos_configuration_profile_plist

import (
	"log"
	"reflect"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/collections"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMacOSConfigurationProfile
// information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMacOSConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]any{
		"name":                resp.General.Name,
		"description":         resp.General.Description,
		"uuid":                resp.General.UUID,
		"distribution_method": resp.General.DistributionMethod,
		"user_removable":      resp.General.UserRemovable,
		//"redeploy_on_update":  resp.General.RedeployOnUpdate, Intentionally not stating this value and using the hcl value defined with the create func. this
		// key is broken and always returns an 'any' value. But the api call still accepts 'Newly Assigned' as well.

	}

	// Check if the level is "Computer" and set it to "System", otherwise use the value from resource
	// This is done to match the Jamf Pro API behavior
	levelValue := resp.General.Level
	if levelValue == "Computer" {
		levelValue = "System"
	}
	resourceData["level"] = levelValue

	d.Set("site_id", resp.General.Site.ID)

	profile := plist.NormalizePayloadState(resp.General.Payloads)
	if err := d.Set("payloads", profile); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.Set("category_id", resp.General.Category.ID)

	if scopeData, err := setScope(resp); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []any{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	defaultSelfService := jamfpro.MacOSConfigurationProfileSubsetSelfService{}
	removeSelfService := reflect.DeepEqual(resp.SelfService, defaultSelfService) || resp.General.DistributionMethod == "Install Automatically"
	if !removeSelfService {
		if selfServiceData, err := setSelfService(resp.SelfService); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		} else if selfServiceData != nil {
			if err := d.Set("self_service", []any{selfServiceData}); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else {
		log.Println("Self-service block is empty, default, or set to 'Install Automatically', removing from state")
		if err := d.Set("self_service", []any{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMacOSConfigurationProfile) (map[string]any, error) {
	scopeData := map[string]any{
		"all_computers": resp.Scope.AllComputers,
		"all_jss_users": resp.Scope.AllJSSUsers,
	}

	scopeData["computer_ids"] = collections.FlattenSortIDs(
		resp.Scope.Computers,
		func(computer jamfpro.MacOSConfigurationProfileSubsetComputer) int { return computer.ID },
	)
	scopeData["computer_group_ids"] = collections.FlattenSortIDs(
		resp.Scope.ComputerGroups,
		func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["jss_user_ids"] = collections.FlattenSortIDs(
		resp.Scope.JSSUsers,
		func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["jss_user_group_ids"] = collections.FlattenSortIDs(
		resp.Scope.JSSUserGroups,
		func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["building_ids"] = collections.FlattenSortIDs(
		resp.Scope.Buildings,
		func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["department_ids"] = collections.FlattenSortIDs(
		resp.Scope.Departments,
		func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)

	limitationsData, err := setLimitations(resp.Scope.Limitations)
	if err != nil {
		return nil, err
	}
	if limitationsData != nil {
		scopeData["limitations"] = limitationsData
	}

	exclusionsData, err := setExclusions(resp.Scope.Exclusions)
	if err != nil {
		return nil, err
	}
	if exclusionsData != nil {
		scopeData["exclusions"] = exclusionsData
	}

	return scopeData, nil
}

// setLimitations collects and formats limitations data for the Terraform state.
func setLimitations(limitations jamfpro.MacOSConfigurationProfileSubsetLimitations) ([]map[string]any, error) {
	result := map[string]any{}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := collections.FlattenSortIDs(
			limitations.NetworkSegments,
			func(segment jamfpro.MacOSConfigurationProfileSubsetNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(limitations.IBeacons) > 0 {
		ibeaconIDs := collections.FlattenSortIDs(
			limitations.IBeacons,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(limitations.Users) > 0 {
		userNames := collections.FlattenSortStrings(
			limitations.Users,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) string { return entity.Name },
		)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupNames := collections.FlattenSortStrings(
			limitations.UserGroups,
			func(userGroup jamfpro.MacOSConfigurationProfileSubsetScopeUserGroup) string { return userGroup.Name },
		)
		if len(userGroupNames) > 0 {
			result["directory_service_usergroup_names"] = userGroupNames
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]any{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MacOSConfigurationProfileSubsetExclusions) ([]map[string]any, error) {
	result := map[string]any{}

	if len(exclusions.Computers) > 0 {
		computerIDs := collections.FlattenSortIDs(
			exclusions.Computers,
			func(computer jamfpro.MacOSConfigurationProfileSubsetComputer) int { return computer.ID },
		)
		if len(computerIDs) > 0 {
			result["computer_ids"] = computerIDs
		}
	}

	if len(exclusions.ComputerGroups) > 0 {
		computerGroupIDs := collections.FlattenSortIDs(
			exclusions.ComputerGroups,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(computerGroupIDs) > 0 {
			result["computer_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := collections.FlattenSortIDs(
			exclusions.Buildings,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := collections.FlattenSortIDs(
			exclusions.JSSUsers,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := collections.FlattenSortIDs(
			exclusions.JSSUserGroups,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := collections.FlattenSortIDs(
			exclusions.Departments,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := collections.FlattenSortIDs(
			exclusions.NetworkSegments,
			func(segment jamfpro.MacOSConfigurationProfileSubsetNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userNames := collections.FlattenSortStrings(
			exclusions.Users,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) string { return entity.Name },
		)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupNames := collections.FlattenSortStrings(
			exclusions.UserGroups,
			func(userGroup jamfpro.MacOSConfigurationProfileSubsetScopeUserGroup) string { return userGroup.Name },
		)
		if len(userGroupNames) > 0 {
			result["directory_service_usergroup_names"] = userGroupNames
		}
	}

	if len(exclusions.IBeacons) > 0 {
		ibeaconIDs := collections.FlattenSortIDs(
			exclusions.IBeacons,
			func(entity jamfpro.MacOSConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]any{result}, nil
}

// setSelfService converts the self-service structure into a format suitable for setting in the Terraform state.
func setSelfService(selfService jamfpro.MacOSConfigurationProfileSubsetSelfService) (map[string]any, error) {
	// Define default values
	defaults := map[string]any{
		"self_service_display_name":       "",
		"install_button_text":             "Install",
		"self_service_description":        "",
		"force_users_to_view_description": false,
		"feature_on_main_page":            false,
		"notification":                    nil,
		"notification_subject":            "",
		"notification_message":            "",
		"self_service_icon_id":            0,
		"self_service_icon":               nil,
		"self_service_category":           nil,
	}

	selfServiceBlock := map[string]any{
		"self_service_display_name":       selfService.SelfServiceDisplayName,
		"install_button_text":             selfService.InstallButtonText,
		"self_service_description":        selfService.SelfServiceDescription,
		"force_users_to_view_description": selfService.ForceUsersToViewDescription,
		"feature_on_main_page":            selfService.FeatureOnMainPage,
		"notification_subject":            selfService.NotificationSubject,
		"notification_message":            selfService.NotificationMessage,
	}

	// Handle self service icon
	if selfService.SelfServiceIcon.ID != 0 {
		selfServiceBlock["self_service_icon_id"] = selfService.SelfServiceIcon.ID
		selfServiceBlock["self_service_icon"] = []any{
			map[string]any{
				"id":       selfService.SelfServiceIcon.ID,
				"uri":      selfService.SelfServiceIcon.URI,
				"data":     selfService.SelfServiceIcon.Data,
				"filename": selfService.SelfServiceIcon.Filename,
			},
		}
	}

	// Temporarily set to nil in all runs due to issues with the API.
	// Will be reimplemented once those are fixed.
	selfServiceBlock["notification"] = nil

	// Handle self service categories
	if len(selfService.SelfServiceCategories) > 0 {
		categories := make([]any, len(selfService.SelfServiceCategories))
		for i, category := range selfService.SelfServiceCategories {
			categories[i] = map[string]any{
				"id":         category.ID,
				"name":       category.Name,
				"display_in": category.DisplayIn,
				"feature_in": category.FeatureIn,
			}
		}
		selfServiceBlock["self_service_category"] = categories
	}

	// Check if all values are default
	allDefault := true
	for key, value := range selfServiceBlock {
		if defaultVal, ok := defaults[key]; ok && !reflect.DeepEqual(value, defaultVal) {
			allDefault = false
			break
		}
	}

	if allDefault {
		log.Println("All self service values are default, skipping state")
		return nil, nil
	}

	log.Println("Initializing self service in state")
	log.Printf("Final state self service: %+v\n", selfServiceBlock)

	return selfServiceBlock, nil
}
