package mac_application

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/collections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMacApplications
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMacApplications) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resp.General.Name); err != nil {
		log.Printf("[ERROR] Failed to set name: %v", err)
	}

	if err := d.Set("version", resp.General.Version); err != nil {
		log.Printf("[ERROR] Failed to set version: %v", err)
	}

	if err := d.Set("bundle_id", resp.General.BundleID); err != nil {
		log.Printf("[ERROR] Failed to set bundle_id: %v", err)
	}

	if err := d.Set("url", resp.General.URL); err != nil {
		log.Printf("[ERROR] Failed to set url: %v", err)
	}

	if err := d.Set("is_free", resp.General.IsFree); err != nil {
		log.Printf("[ERROR] Failed to set is_free: %v", err)
	}

	if err := d.Set("deployment_type", resp.General.DeploymentType); err != nil {
		log.Printf("[ERROR] Failed to set deployment_type: %v", err)
	}

	if err := d.Set("site_id", resp.General.Site.ID); err != nil {
		log.Printf("[ERROR] Failed to set site_id: %v", err)
	}

	if err := d.Set("category_id", resp.General.Category.ID); err != nil {
		log.Printf("[ERROR] Failed to set category_id: %v", err)
	}

	if resp.SelfService.SelfServiceDescription != "" || resp.SelfService.NotificationMessage != "" {
		selfService := []map[string]any{
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
			selfService[0]["self_service_icon"] = []map[string]any{
				{
					"id":   resp.SelfService.SelfServiceIcon.ID,
					"data": resp.SelfService.SelfServiceIcon.Data,
					"uri":  resp.SelfService.SelfServiceIcon.URI,
				},
			}
		}

		// Handle self service categories
		if len(resp.SelfService.SelfServiceCategories) > 0 {
			var categories []map[string]any
			for _, cat := range resp.SelfService.SelfServiceCategories {
				categories = append(categories, map[string]any{
					"id":         cat.ID,
					"name":       cat.Name,
					"display_in": cat.DisplayIn,
					"feature_in": cat.FeatureIn,
				})
			}
			selfService[0]["self_service_category"] = categories
		}

		if err := d.Set("self_service", selfService); err != nil {
			log.Printf("[ERROR] Failed to set self_service: %v", err)
		}
	}

	if resp.VPP.VPPAdminAccountID != 0 {
		vpp := []map[string]any{
			{
				"assign_vpp_device_based_licenses": resp.VPP.AssignVPPDeviceBasedLicenses,
				"vpp_admin_account_id":             resp.VPP.VPPAdminAccountID,
			},
		}

		if err := d.Set("vpp", vpp); err != nil {
			log.Printf("[ERROR] Failed to set vpp: %v", err)
		}

	}

	if scopeData, err := setScope(resp); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []any{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMacApplications) (map[string]any, error) {
	scopeData := map[string]any{
		"all_computers": resp.Scope.AllComputers,
		"all_jss_users": resp.Scope.AllJSSUsers,
	}

	scopeData["computer_ids"] = collections.FlattenSortIDs(
		resp.Scope.Computers,
		func(computer jamfpro.MacAppSubsetScopeComputer) int { return computer.ID },
	)
	scopeData["computer_group_ids"] = collections.FlattenSortIDs(
		resp.Scope.ComputerGroups,
		func(group jamfpro.MacAppSubsetScopeComputerGroup) int { return group.ID },
	)
	scopeData["jss_user_ids"] = collections.FlattenSortIDs(
		resp.Scope.JSSUsers,
		func(user jamfpro.MacAppSubsetScopeUser) int { return user.ID },
	)
	scopeData["jss_user_group_ids"] = collections.FlattenSortIDs(
		resp.Scope.JSSUserGroups,
		func(group jamfpro.MacAppSubsetScopeJSSUserGroup) int { return group.ID },
	)
	scopeData["building_ids"] = collections.FlattenSortIDs(
		resp.Scope.Buildings,
		func(building jamfpro.MacAppSubsetScopeBuilding) int { return building.ID },
	)
	scopeData["department_ids"] = collections.FlattenSortIDs(
		resp.Scope.Departments,
		func(department jamfpro.MacAppSubsetScopeDepartment) int { return department.ID },
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
func setLimitations(limitations jamfpro.MacAppScopeLimitations) ([]map[string]any, error) {
	result := map[string]any{}

	if len(limitations.Users) > 0 {
		userIDs := collections.FlattenSortIDs(
			limitations.Users,
			func(user jamfpro.MacAppSubsetScopeUser) int { return user.ID },
		)
		if len(userIDs) > 0 {
			result["users"] = userIDs
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupNames := collections.FlattenSortStrings(
			limitations.UserGroups,
			func(userGroup jamfpro.MacAppSubsetScopeUserGroup) string { return userGroup.Name },
		)
		if len(userGroupNames) > 0 {
			result["directory_service_usergroup_names"] = userGroupNames
		}
	}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := collections.FlattenSortIDs(
			limitations.NetworkSegments,
			func(segment jamfpro.MacAppSubsetScopeNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segments"] = networkSegmentIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]any{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MacAppScopeExclusions) ([]map[string]any, error) {
	result := map[string]any{}

	if len(exclusions.Computers) > 0 {
		computerIDs := collections.FlattenSortIDs(
			exclusions.Computers,
			func(computer jamfpro.MacAppSubsetScopeComputer) int { return computer.ID },
		)
		if len(computerIDs) > 0 {
			result["computer_ids"] = computerIDs
		}
	}

	if len(exclusions.ComputerGroups) > 0 {
		computerGroupIDs := collections.FlattenSortIDs(
			exclusions.ComputerGroups,
			func(group jamfpro.MacAppSubsetScopeComputerGroup) int { return group.ID },
		)
		if len(computerGroupIDs) > 0 {
			result["computer_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userIDs := collections.FlattenSortIDs(
			exclusions.Users,
			func(user jamfpro.MacAppSubsetScopeUser) int { return user.ID },
		)
		if len(userIDs) > 0 {
			result["users"] = userIDs
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupNames := collections.FlattenSortStrings(
			exclusions.UserGroups,
			func(userGroup jamfpro.MacAppSubsetScopeUserGroup) string { return userGroup.Name },
		)
		if len(userGroupNames) > 0 {
			result["directory_service_usergroup_names"] = userGroupNames
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := collections.FlattenSortIDs(
			exclusions.Buildings,
			func(building jamfpro.MacAppSubsetScopeBuilding) int { return building.ID },
		)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := collections.FlattenSortIDs(
			exclusions.Departments,
			func(department jamfpro.MacAppSubsetScopeDepartment) int { return department.ID },
		)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := collections.FlattenSortIDs(
			exclusions.NetworkSegments,
			func(segment jamfpro.MacAppSubsetScopeNetworkSegment) int { return segment.ID },
		)
		if len(networkSegmentIDs) > 0 {
			result["network_segments"] = networkSegmentIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := collections.FlattenSortIDs(
			exclusions.JSSUsers,
			func(user jamfpro.MacAppSubsetScopeUser) int { return user.ID },
		)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := collections.FlattenSortIDs(
			exclusions.JSSUserGroups,
			func(group jamfpro.MacAppSubsetScopeJSSUserGroup) int { return group.ID },
		)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]any{result}, nil
}
