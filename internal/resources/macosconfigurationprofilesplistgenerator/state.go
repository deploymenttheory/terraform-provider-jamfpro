// macosconfigurationprofilesplistgenerator_state.go
package macosconfigurationprofilesplistgenerator

import (
	"reflect"
	"sort"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest ResourceMacOSConfigurationProfile
// information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceMacOSConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"name":                resp.General.Name,
		"description":         resp.General.Description,
		"uuid":                resp.General.UUID,
		"distribution_method": resp.General.DistributionMethod,
		"user_removable":      resp.General.UserRemovable,
		"redeploy_on_update":  resp.General.RedeployOnUpdate,
	}

	// Check if the level is "Computer" and set it to "System", otherwise use the value from resource
	// This is done to match the Jamf Pro API behavior
	levelValue := resp.General.Level
	if levelValue == "Computer" {
		levelValue = "System"
	}
	resourceData["level"] = levelValue

	d.Set("site_id", resp.General.Site.ID)

	// Convert the plist payloads back to HCL format and set them in the state
	payloadsList, err := plist.ConvertPlistToHCL(resp.General.Payloads)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		if err := d.Set("payloads", payloadsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	d.Set("category_id", resp.General.Category.ID)

	// Preparing and setting scope data
	if scopeData, err := setScope(resp); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Check if the self_service block is provided and set it in the state accordingly
	defaultSelfService := jamfpro.MacOSConfigurationProfileSubsetSelfService{}
	if !compareSelfService(resp.SelfService, defaultSelfService) {
		if selfServiceData, err := setSelfService(resp.SelfService); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		} else if selfServiceData != nil {
			if err := d.Set("self_service", []interface{}{selfServiceData}); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

	} else {
		// TODO why?
		// If self_service block is not provided, set it to an empty array
		if err := d.Set("self_service", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Update the resource data
	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resp *jamfpro.ResourceMacOSConfigurationProfile) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_computers": resp.Scope.AllComputers,
		"all_jss_users": resp.Scope.AllJSSUsers,
	}

	// Gather computers, groups, etc.
	scopeData["computer_ids"] = flattenAndSortComputerIds(resp.Scope.Computers)
	scopeData["computer_group_ids"] = flattenAndSortScopeEntityIds(resp.Scope.ComputerGroups)
	scopeData["jss_user_ids"] = flattenAndSortScopeEntityIds(resp.Scope.JSSUsers)
	scopeData["jss_user_group_ids"] = flattenAndSortScopeEntityIds(resp.Scope.JSSUserGroups)
	scopeData["building_ids"] = flattenAndSortScopeEntityIds(resp.Scope.Buildings)
	scopeData["department_ids"] = flattenAndSortScopeEntityIds(resp.Scope.Departments)

	// Gather limitations
	limitationsData, err := setLimitations(resp.Scope.Limitations)
	if err != nil {
		return nil, err
	}
	if limitationsData != nil {
		scopeData["limitations"] = limitationsData
	}

	// Gather exclusions
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
func setLimitations(limitations jamfpro.MacOSConfigurationProfileSubsetLimitations) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIds(limitations.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(limitations.IBeacons) > 0 {
		ibeaconIDs := flattenAndSortScopeEntityIds(limitations.IBeacons)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
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
func setExclusions(exclusions jamfpro.MacOSConfigurationProfileSubsetExclusions) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(exclusions.Computers) > 0 {
		computerIDs := flattenAndSortComputerIds(exclusions.Computers)
		if len(computerIDs) > 0 {
			result["computer_ids"] = computerIDs
		}
	}

	if len(exclusions.ComputerGroups) > 0 {
		computerGroupIDs := flattenAndSortScopeEntityIds(exclusions.ComputerGroups)
		if len(computerGroupIDs) > 0 {
			result["computer_group_ids"] = computerGroupIDs
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

	if len(exclusions.IBeacons) > 0 {
		ibeaconIDs := flattenAndSortScopeEntityIds(exclusions.IBeacons)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// setSelfService converts the self-service structure into a format suitable for setting in the Terraform state.
func setSelfService(selfService jamfpro.MacOSConfigurationProfileSubsetSelfService) (map[string]interface{}, error) {
	selfServiceData := make(map[string]interface{})

	// Set real values only, avoiding defaults
	if selfService.InstallButtonText != "" && selfService.InstallButtonText != "Install" {
		selfServiceData["install_button_text"] = selfService.InstallButtonText
	}
	if selfService.SelfServiceDescription != "" && selfService.SelfServiceDescription != "no description set" {
		selfServiceData["self_service_description"] = selfService.SelfServiceDescription
	}
	if selfService.ForceUsersToViewDescription {
		selfServiceData["force_users_to_view_description"] = selfService.ForceUsersToViewDescription
	}
	if selfService.FeatureOnMainPage {
		selfServiceData["feature_on_main_page"] = selfService.FeatureOnMainPage
	}
	if selfService.NotificationSubject != "" && selfService.NotificationSubject != "no message subject set" {
		selfServiceData["notification_subject"] = selfService.NotificationSubject
	}
	if selfService.NotificationMessage != "" {
		selfServiceData["notification_message"] = selfService.NotificationMessage
	}

	// Temporarily set to nil in all runs due to issues with the API.
	// Will be reimplemented once those are fixed.
	selfServiceData["notification"] = nil

	if selfService.SelfServiceIcon.ID != 0 {
		selfServiceData["self_service_icon"] = []interface{}{
			map[string]interface{}{
				"id":       selfService.SelfServiceIcon.ID,
				"uri":      selfService.SelfServiceIcon.URI,
				"data":     selfService.SelfServiceIcon.Data,
				"filename": selfService.SelfServiceIcon.Filename,
			},
		}
	}

	if len(selfService.SelfServiceCategories) > 0 {
		categories := []interface{}{}
		for _, category := range selfService.SelfServiceCategories {
			categories = append(categories, map[string]interface{}{
				"id":         category.ID,
				"name":       category.Name,
				"display_in": category.DisplayIn,
				"feature_in": category.FeatureIn,
			})
		}
		selfServiceData["self_service_categories"] = categories
	}

	// Return nil map if there are no real values to avoid setting the self_service block
	if len(selfServiceData) == 0 {
		return nil, nil
	}

	return selfServiceData, nil
}

// TODO what is going on here?
// compareSelfService compares two MacOSConfigurationProfileSubsetSelfService structs
func compareSelfService(a, b jamfpro.MacOSConfigurationProfileSubsetSelfService) bool {
	return a.InstallButtonText == b.InstallButtonText &&
		a.SelfServiceDescription == b.SelfServiceDescription &&
		a.ForceUsersToViewDescription == b.ForceUsersToViewDescription &&
		reflect.DeepEqual(a.SelfServiceIcon, b.SelfServiceIcon) &&
		a.FeatureOnMainPage == b.FeatureOnMainPage &&
		reflect.DeepEqual(a.SelfServiceCategories, b.SelfServiceCategories) &&
		a.Notification == b.Notification &&
		a.NotificationSubject == b.NotificationSubject &&
		a.NotificationMessage == b.NotificationMessage
}

// helper functions

// flattenAndSortScopeEntityIds converts a slice of general scope entities (like user groups, buildings) to a format suitable for Terraform state.
func flattenAndSortScopeEntityIds(entities []jamfpro.MacOSConfigurationProfileSubsetScopeEntity) []int {
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
func flattenAndSortScopeEntityNames(entities []jamfpro.MacOSConfigurationProfileSubsetScopeEntity) []string {
	var names []string
	for _, entity := range entities {
		if entity.Name != "" {
			names = append(names, entity.Name)
		}
	}
	sort.Strings(names)
	return names
}

// flattenAndSortComputerIds converts a slice of MacOSConfigurationProfileSubsetComputer into a sorted slice of integers.
func flattenAndSortComputerIds(computers []jamfpro.MacOSConfigurationProfileSubsetComputer) []int {
	var ids []int
	for _, computer := range computers {
		if computer.ID != 0 {
			ids = append(ids, computer.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortNetworkSegmentIds converts a slice of MacOSConfigurationProfileSubsetNetworkSegment into a sorted slice of integers.
func flattenAndSortNetworkSegmentIds(segments []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment) []int {
	var ids []int
	for _, segment := range segments {
		if segment.ID != 0 {
			ids = append(ids, segment.ID)
		}
	}
	sort.Ints(ids)
	return ids
}
