// mobiledeviceconfigurationprofilesplist_state.go
package mobile_device_configuration_profile_plist

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMobileDeviceConfigurationProfile
// information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDeviceConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]any{
		"name":              resp.General.Name,
		"description":       resp.General.Description,
		"uuid":              resp.General.UUID,
		"deployment_method": resp.General.DeploymentMethod,
		// Skipping the 'distribution_method' attribute as it appears to be deprecated but still in documentation
		"redeploy_on_update":                resp.General.RedeployOnUpdate,
		"redeploy_days_before_cert_expires": resp.General.RedeployDaysBeforeCertExpires,
	}

	// Check if the level is "System" and set it to "Device Level", otherwise use the value from resource
	// This is done to match the Jamf Pro API behavior
	levelValue := resp.General.Level
	if levelValue == "System" {
		levelValue = "Device Level"
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

	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMobileDeviceConfigurationProfile) (map[string]any, error) {
	scopeData := map[string]any{
		"all_mobile_devices": resp.Scope.AllMobileDevices,
		"all_jss_users":      resp.Scope.AllJSSUsers,
	}

	scopeData["mobile_device_ids"] = utils.FlattenSortIDs(
		resp.Scope.MobileDevices,
		func(device jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice) int { return device.ID },
	)
	scopeData["mobile_device_group_ids"] = utils.FlattenSortIDs(
		resp.Scope.MobileDeviceGroups,
		func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["jss_user_ids"] = utils.FlattenSortIDs(
		resp.Scope.JSSUsers,
		func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["jss_user_group_ids"] = utils.FlattenSortIDs(
		resp.Scope.JSSUserGroups,
		func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["building_ids"] = utils.FlattenSortIDs(
		resp.Scope.Buildings,
		func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
	)
	scopeData["department_ids"] = utils.FlattenSortIDs(
		resp.Scope.Departments,
		func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
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
func setLimitations(limitations jamfpro.MobileDeviceConfigurationProfileSubsetLimitation) ([]map[string]any, error) {
	result := map[string]any{}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := utils.FlattenSortIDs(
			limitations.NetworkSegments,
			func(segment jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(limitations.Ibeacons) > 0 {
		ibeaconIDs := utils.FlattenSortIDs(
			limitations.Ibeacons,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(limitations.Users) > 0 {
		userNames := utils.FlattenSortStrings(
			limitations.Users,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) string { return entity.Name },
		)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupNames := utils.FlattenSortStrings(
			limitations.UserGroups,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) string { return entity.Name },
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
func setExclusions(exclusions jamfpro.MobileDeviceConfigurationProfileSubsetExclusion) ([]map[string]any, error) {
	result := map[string]any{}

	if len(exclusions.MobileDevices) > 0 {
		computerIDs := utils.FlattenSortIDs(
			exclusions.MobileDevices,
			func(device jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice) int { return device.ID },
		)
		if len(computerIDs) > 0 {
			result["mobile_device_ids"] = computerIDs
		}
	}

	if len(exclusions.MobileDeviceGroups) > 0 {
		computerGroupIDs := utils.FlattenSortIDs(
			exclusions.MobileDeviceGroups,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(computerGroupIDs) > 0 {
			result["mobile_device_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := utils.FlattenSortIDs(
			exclusions.Buildings,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := utils.FlattenSortIDs(
			exclusions.JSSUsers,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := utils.FlattenSortIDs(
			exclusions.JSSUserGroups,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := utils.FlattenSortIDs(
			exclusions.Departments,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
		)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := utils.FlattenSortIDs(
			exclusions.NetworkSegments,
			func(segment jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userNames := utils.FlattenSortStrings(
			exclusions.Users,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) string { return entity.Name },
		)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupNames := utils.FlattenSortStrings(
			exclusions.UserGroups,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) string { return entity.Name },
		)
		if len(userGroupNames) > 0 {
			result["directory_service_usergroup_names"] = userGroupNames
		}
	}

	if len(exclusions.IBeacons) > 0 {
		ibeaconIDs := utils.FlattenSortIDs(
			exclusions.IBeacons,
			func(entity jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) int { return entity.ID },
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
