package mac_application

import (
	"sort"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMacApplications
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMacApplications) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("name", resp.General.Name)
	d.Set("version", resp.General.Version)
	d.Set("bundle_id", resp.General.BundleID)
	d.Set("url", resp.General.URL)
	d.Set("is_free", resp.General.IsFree)
	d.Set("deployment_type", resp.General.DeploymentType)
	d.Set("site_id", resp.General.Site.ID)
	d.Set("category_id", resp.General.Category.ID)

	if resp.SelfService.SelfServiceDescription != "" || resp.SelfService.NotificationMessage != "" {
		selfService := []map[string]interface{}{
			{
				"self_service_description":        resp.SelfService.SelfServiceDescription,
				"install_button_text":             resp.SelfService.InstallButtonText,
				"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
				"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
				"notification":                    resp.SelfService.Notification,
				"notification_subject":            resp.SelfService.NotificationSubject,
				"notification_message":            resp.SelfService.NotificationMessage,
			},
		}

		// Handle self service icon
		if resp.SelfService.SelfServiceIcon.ID != 0 {
			selfService[0]["self_service_icon"] = []map[string]interface{}{
				{
					"id":   resp.SelfService.SelfServiceIcon.ID,
					"data": resp.SelfService.SelfServiceIcon.Data,
					"uri":  resp.SelfService.SelfServiceIcon.URI,
				},
			}
		}

		// Handle self service categories
		if len(resp.SelfService.SelfServiceCategories) > 0 {
			var categories []map[string]interface{}
			for _, cat := range resp.SelfService.SelfServiceCategories {
				categories = append(categories, map[string]interface{}{
					"id":         cat.ID,
					"name":       cat.Name,
					"display_in": cat.DisplayIn,
					"feature_in": cat.FeatureIn,
				})
			}
			selfService[0]["self_service_category"] = categories
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

	if scopeData, err := setScope(resp); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMacApplications) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_computers": resp.Scope.AllComputers,
		"all_jss_users": resp.Scope.AllJSSUsers,
	}

	scopeData["computer_ids"] = flattenAndSortComputerIDs(resp.Scope.Computers)
	scopeData["computer_group_ids"] = flattenAndSortComputerGroupIDs(resp.Scope.ComputerGroups)
	scopeData["jss_user_ids"] = flattenAndSortUserIDs(resp.Scope.JSSUsers)
	scopeData["jss_user_group_ids"] = flattenAndSortUserGroupIDs(resp.Scope.JSSUserGroups)
	scopeData["building_ids"] = flattenAndSortBuildingIDs(resp.Scope.Buildings)
	scopeData["department_ids"] = flattenAndSortDepartmentIDs(resp.Scope.Departments)

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
func setLimitations(limitations jamfpro.MacAppScopeLimitations) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(limitations.Users) > 0 {
		userIDs := flattenAndSortUserIDs(limitations.Users)
		if len(userIDs) > 0 {
			result["users"] = userIDs
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupIDs := flattenAndSortUserGroupIDs(limitations.UserGroups)
		if len(userGroupIDs) > 0 {
			result["user_groups"] = userGroupIDs
		}
	}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIDs(limitations.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segments"] = networkSegmentIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MacAppScopeExclusions) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(exclusions.Computers) > 0 {
		computerIDs := flattenAndSortComputerIDs(exclusions.Computers)
		if len(computerIDs) > 0 {
			result["computer_ids"] = computerIDs
		}
	}

	if len(exclusions.ComputerGroups) > 0 {
		computerGroupIDs := flattenAndSortComputerGroupIDs(exclusions.ComputerGroups)
		if len(computerGroupIDs) > 0 {
			result["computer_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userIDs := flattenAndSortUserIDs(exclusions.Users)
		if len(userIDs) > 0 {
			result["users"] = userIDs
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupIDs := flattenAndSortUserGroupIDs(exclusions.UserGroups)
		if len(userGroupIDs) > 0 {
			result["user_groups"] = userGroupIDs
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := flattenAndSortBuildingIDs(exclusions.Buildings)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := flattenAndSortDepartmentIDs(exclusions.Departments)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIDs(exclusions.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segments"] = networkSegmentIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := flattenAndSortUserIDs(exclusions.JSSUsers)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := flattenAndSortUserGroupIDs(exclusions.JSSUserGroups)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// Helper functions

func flattenAndSortComputerIDs(computers []jamfpro.MacAppSubsetScopeComputer) []int {
	var ids []int
	for _, computer := range computers {
		if computer.ID != 0 {
			ids = append(ids, computer.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortComputerGroupIDs(groups []jamfpro.MacAppSubsetScopeComputerGroup) []int {
	var ids []int
	for _, group := range groups {
		if group.ID != 0 {
			ids = append(ids, group.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortUserIDs(users []jamfpro.MacAppSubsetScopeUser) []int {
	var ids []int
	for _, user := range users {
		if user.ID != 0 {
			ids = append(ids, user.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortUserGroupIDs(groups []jamfpro.MacAppSubsetScopeUserGroup) []int {
	var ids []int
	for _, group := range groups {
		if group.ID != 0 {
			ids = append(ids, group.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortBuildingIDs(buildings []jamfpro.MacAppSubsetScopeBuilding) []int {
	var ids []int
	for _, building := range buildings {
		if building.ID != 0 {
			ids = append(ids, building.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortDepartmentIDs(departments []jamfpro.MacAppSubsetScopeDepartment) []int {
	var ids []int
	for _, department := range departments {
		if department.ID != 0 {
			ids = append(ids, department.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func flattenAndSortNetworkSegmentIDs(segments []jamfpro.MacAppSubsetScopeNetworkSegment) []int {
	var ids []int
	for _, segment := range segments {
		if segment.ID != 0 {
			ids = append(ids, segment.ID)
		}
	}
	sort.Ints(ids)
	return ids
}
