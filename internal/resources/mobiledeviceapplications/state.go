package mobiledeviceapplications

import (
	"sort"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMobileDeviceApplication
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDeviceApplication) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("name", resp.General.Name)
	d.Set("display_name", resp.General.DisplayName)
	d.Set("description", normalizeWhitespace(resp.General.Description))
	d.Set("bundle_id", resp.General.BundleID)
	d.Set("version", resp.General.Version)
	d.Set("internal_app", resp.General.InternalApp)
	d.Set("itunes_store_url", resp.General.ITunesStoreURL)
	d.Set("make_available_after_install", resp.General.MakeAvailableAfterInstall)
	d.Set("itunes_country_region", resp.General.ITunesCountryRegion)
	d.Set("itunes_sync_time", resp.General.ITunesSyncTime)
	d.Set("deployment_type", resp.General.DeploymentType)
	d.Set("deploy_automatically", resp.General.DeployAutomatically)
	d.Set("deploy_as_managed_app", resp.General.DeployAsManagedApp)
	d.Set("remove_app_when_mdm_profile_is_removed", resp.General.RemoveAppWhenMDMProfileIsRemoved)
	d.Set("prevent_backup_of_app_data", resp.General.PreventBackupOfAppData)
	d.Set("keep_description_and_icon_up_to_date", resp.General.KeepDescriptionAndIconUpToDate)
	d.Set("keep_app_updated_on_devices", resp.General.KeepAppUpdatedOnDevices)
	d.Set("free", resp.General.Free)
	d.Set("take_over_management", resp.General.TakeOverManagement)
	d.Set("host_externally", resp.General.HostExternally)
	d.Set("external_url", resp.General.ExternalURL)
	d.Set("mobile_device_provisioning_profile", resp.General.ProvisioningProfile)
	d.Set("site_id", resp.General.Site.ID)
	d.Set("category_id", resp.General.Category.ID)

	if resp.General.IPA.Name != "" || resp.General.IPA.URI != "" || resp.General.IPA.Data != "" {
		ipa := []map[string]interface{}{
			{
				"name": resp.General.IPA.Name,
				"uri":  resp.General.IPA.URI,
				"data": resp.General.IPA.Data,
			},
		}
		d.Set("ipa", ipa)
	}

	if resp.General.Icon.ID != 0 || resp.General.Icon.Name != "" || resp.General.Icon.URI != "" {
		icon := []map[string]interface{}{
			{
				"id":   resp.General.Icon.ID,
				"name": resp.General.Icon.Name,
				"uri":  resp.General.Icon.URI,
			},
		}
		d.Set("icon", icon)
	}

	if resp.SelfService.SelfServiceDescription != "" || resp.SelfService.NotificationMessage != "" {
		selfService := []map[string]interface{}{
			{
				"self_service_description": normalizeWhitespace(resp.SelfService.SelfServiceDescription),
				"feature_on_main_page":     resp.SelfService.FeatureOnMainPage,
				"notification":             resp.SelfService.Notification,
				"notification_subject":     resp.SelfService.NotificationSubject,
				"notification_message":     resp.SelfService.NotificationMessage,
			},
		}

		if resp.SelfService.SelfServiceIcon.ID != 0 || resp.SelfService.SelfServiceIcon.Name != "" || resp.SelfService.SelfServiceIcon.URI != "" {
			selfService[0]["self_service_icon"] = []map[string]interface{}{
				{
					"id":       resp.SelfService.SelfServiceIcon.ID,
					"filename": resp.SelfService.SelfServiceIcon.Name,
					"uri":      resp.SelfService.SelfServiceIcon.URI,
				},
			}
		}

		d.Set("self_service", selfService)
	}

	if resp.VPP.VPPAdminAccountID != 0 {
		vpp := []map[string]interface{}{
			{
				"assign_vpp_device_based_licenses": resp.VPP.AssignVPPDeviceBasedLicenses,
				"vpp_admin_account_id":             resp.VPP.VPPAdminAccountID,
			},
		}
		d.Set("vpp", vpp)
	}

	if resp.AppConfiguration.Preferences != "" {
		appConfig := []map[string]interface{}{
			{
				"preferences": normalizeWhitespace(resp.AppConfiguration.Preferences),
			},
		}
		d.Set("app_configuration", appConfig)
	}

	if scopeData, err := setScope(resp); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMobileDeviceApplication) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_mobile_devices": resp.Scope.AllMobileDevices,
		"all_jss_users":      resp.Scope.AllJSSUsers,
	}

	scopeData["mobile_device_ids"] = flattenAndSortMobileDeviceIDs(resp.Scope.MobileDevices)
	scopeData["mobile_device_group_ids"] = flattenAndSortScopeEntityIds(resp.Scope.MobileDeviceGroups)
	scopeData["jss_user_ids"] = flattenAndSortScopeEntityIds(resp.Scope.JSSUsers)
	scopeData["jss_user_group_ids"] = flattenAndSortScopeEntityIds(resp.Scope.JSSUserGroups)
	scopeData["building_ids"] = flattenAndSortScopeEntityIds(resp.Scope.Buildings)
	scopeData["department_ids"] = flattenAndSortScopeEntityIds(resp.Scope.Departments)

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
func setLimitations(limitations jamfpro.MobileDeviceApplicationSubsetLimitation) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIds(limitations.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(limitations.Users) > 0 {
		userNames := flattenAndSortScopeEntityNames(limitations.Users)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupIDs := flattenAndSortScopeEntityIds(limitations.UserGroups)
		if len(userGroupIDs) > 0 {
			result["directory_service_usergroup_ids"] = userGroupIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MobileDeviceApplicationSubsetExclusion) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(exclusions.MobileDevices) > 0 {
		computerIDs := flattenAndSortMobileDeviceIDs(exclusions.MobileDevices)
		if len(computerIDs) > 0 {
			result["mobile_device_ids"] = computerIDs
		}
	}

	if len(exclusions.MobileDeviceGroups) > 0 {
		computerGroupIDs := flattenAndSortScopeEntityIds(exclusions.MobileDeviceGroups)
		if len(computerGroupIDs) > 0 {
			result["mobile_device_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := flattenAndSortScopeEntityIds(exclusions.Buildings)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := flattenAndSortScopeEntityIds(exclusions.JSSUsers)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := flattenAndSortScopeEntityIds(exclusions.JSSUserGroups)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := flattenAndSortScopeEntityIds(exclusions.Departments)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIds(exclusions.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userNames := flattenAndSortScopeEntityNames(exclusions.Users)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupIDs := flattenAndSortScopeEntityIds(exclusions.UserGroups)
		if len(userGroupIDs) > 0 {
			result["directory_service_usergroup_ids"] = userGroupIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// helper functions

// flattenAndSortScopeEntityIds converts a slice of general scope entities (like user groups, buildings) to a format suitable for Terraform state.
func flattenAndSortScopeEntityIds(entities []jamfpro.MobileDeviceApplicationSubsetScopeEntity) []int {
	var ids []int
	for _, entity := range entities {
		if entity.ID != 0 {
			ids = append(ids, entity.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortScopeEntityNames converts a slice of RestrictedSoftwareSubsetScopeEntity into a sorted slice of strings.
func flattenAndSortScopeEntityNames(entities []jamfpro.MobileDeviceApplicationSubsetScopeEntity) []string {
	var names []string
	for _, entity := range entities {
		if entity.Name != "" {
			names = append(names, entity.Name)
		}
	}
	sort.Strings(names)
	return names
}

// flattenAndSortMobileDeviceIDs converts a slice of MobileDeviceApplicationSubsetMobileDevice into a sorted slice of integers.
func flattenAndSortMobileDeviceIDs(devices []jamfpro.MobileDeviceApplicationSubsetMobileDevice) []int {
	var ids []int
	for _, device := range devices {
		if device.ID != 0 {
			ids = append(ids, device.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortNetworkSegmentIds converts a slice of MobileDeviceApplicationSubsetNetworkSegment into a sorted slice of integers.
func flattenAndSortNetworkSegmentIds(segments []jamfpro.MobileDeviceApplicationSubsetNetworkSegment) []int {
	var ids []int
	for _, segment := range segments {
		if segment.ID != 0 {
			ids = append(ids, segment.ID)
		}
	}
	sort.Ints(ids)
	return ids
}
