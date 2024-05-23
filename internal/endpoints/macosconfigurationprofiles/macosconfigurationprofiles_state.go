package macosconfigurationprofiles

import (
	"log"
	"sort"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest ResourceMacOSConfigurationProfile
// information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceMacOSConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a map to hold the resource data
	resourceData := map[string]interface{}{
		"name":                resource.General.Name,
		"description":         resource.General.Description,
		"uuid":                resource.General.UUID,
		"distribution_method": resource.General.DistributionMethod,
		"user_removable":      resource.General.UserRemovable,
		"redeploy_on_update":  resource.General.RedeployOnUpdate,
	}

	// Check if the level is "Computer" and set it to "System", otherwise use the value from resource
	// This is done to match the Jamf Pro API behavior
	levelValue := resource.General.Level
	if levelValue == "Computer" {
		levelValue = "System"
	}
	resourceData["level"] = levelValue

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	site := []interface{}{}
	if resource.General.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.General.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Sanitize and set the payloads using the plist processor function
	processedProfile, err := configurationprofiles.ProcessConfigurationProfileForState(resource.General.Payloads)
	if err != nil {
		log.Printf("Error processing configuration profile: %v\n", err)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	log.Printf("Processed profile payload: %s\n", processedProfile)

	// Set the processed payloads field
	if err := d.Set("payloads", processedProfile); err != nil {
		log.Printf("Error setting payloads: %v\n", err)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'category' attribute in the state only if it's not empty (i.e., not default values)
	category := []interface{}{}
	if resource.General.Category.ID != -1 {
		category = append(category, map[string]interface{}{
			"id": resource.General.Category.ID,
		})
	}
	if len(category) > 0 {
		if err := d.Set("category", category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Preparing and setting scope data
	if scopeData, err := setScope(resource); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Preparing and setting self-service data
	if selfServiceData, err := setSelfService(resource.SelfService); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if selfServiceData != nil {
		if err := d.Set("self_service", []interface{}{selfServiceData}); err != nil {
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
func setScope(resource *jamfpro.ResourceMacOSConfigurationProfile) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_computers": resource.Scope.AllComputers,
		"all_jss_users": resource.Scope.AllJSSUsers,
	}

	// Gather computers, groups, etc.
	scopeData["computer_ids"] = flattenAndSortComputerIds(resource.Scope.Computers)
	scopeData["computer_group_ids"] = flattenAndSortScopeEntityIds(resource.Scope.ComputerGroups)
	scopeData["jss_user_ids"] = flattenAndSortScopeEntityIds(resource.Scope.JSSUsers)
	scopeData["jss_user_group_ids"] = flattenAndSortScopeEntityIds(resource.Scope.JSSUserGroups)
	scopeData["building_ids"] = flattenAndSortScopeEntityIds(resource.Scope.Buildings)
	scopeData["department_ids"] = flattenAndSortScopeEntityIds(resource.Scope.Departments)

	// Gather limitations
	limitationsData, err := setLimitations(resource.Scope.Limitations)
	if err != nil {
		return nil, err
	}
	if limitationsData != nil {
		scopeData["limitations"] = limitationsData
	}

	// Gather exclusions
	exclusionsData, err := setExclusions(resource.Scope.Exclusions)
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

	// Fix the notification field
	if len(selfService.Notification) > 0 {
		correctNotifValue, err := FixDuplicateNotificationKey(&jamfpro.ResourceMacOSConfigurationProfile{
			SelfService: selfService,
		})
		if err != nil {
			return nil, err
		}
		// Only set notification if it's true
		if correctNotifValue {
			selfServiceData["notification"] = correctNotifValue
		}
	}

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
