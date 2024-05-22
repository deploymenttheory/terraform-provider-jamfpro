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

	// Check if the level is "System" and set it to "Device Level", otherwise use the value from resource
	levelValue := resource.General.Level
	if levelValue == "System" {
		levelValue = "Device Level"
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

	// Sanitize and set the payloads using the plist processor function removing the mdm server unique identifiers
	keysToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization"}
	processedProfile, err := configurationprofiles.ProcessConfigurationProfile(resource.General.Payloads, keysToRemove)
	if err != nil {
		log.Printf("Error processing configuration profile: %v\n", err)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	log.Printf("Processed profile payload: %s\n", processedProfile)

	// Set the processed payloads field
	if err := d.Set("payloads", processedProfile); err != nil {
		log.Printf("Error setting payload: %v\n", err)
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
	} else if err := d.Set("self_service", []interface{}{selfServiceData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
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
	scopeData["limitations"] = limitationsData

	// Gather exclusions
	exclusionsData, err := setExclusions(resource.Scope.Exclusions)
	if err != nil {
		return nil, err
	}
	scopeData["exclusions"] = exclusionsData

	return scopeData, nil
}

// setLimitations collects and formats limitations data for the Terraform state.
func setLimitations(limitations jamfpro.MacOSConfigurationProfileSubsetLimitations) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(limitations.NetworkSegments) > 0 {
		result["network_segment_ids"] = flattenAndSortNetworkSegmentIds(limitations.NetworkSegments)
	}

	if len(limitations.IBeacons) > 0 {
		result["ibeacon_ids"] = flattenAndSortScopeEntityIds(limitations.IBeacons)
	}

	if len(limitations.Users) > 0 {
		result["directory_service_or_local_usernames"] = flattenAndSortScopeEntityNames(limitations.Users)
	}

	if len(limitations.UserGroups) > 0 {
		result["directory_service_usergroup_ids"] = flattenAndSortScopeEntityIds(limitations.UserGroups)
	}

	return []map[string]interface{}{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MacOSConfigurationProfileSubsetExclusions) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(exclusions.Computers) > 0 {
		result["computer_ids"] = flattenAndSortComputerIds(exclusions.Computers)
	}

	if len(exclusions.ComputerGroups) > 0 {
		result["computer_group_ids"] = flattenAndSortScopeEntityIds(exclusions.ComputerGroups)
	}

	if len(exclusions.Buildings) > 0 {
		result["building_ids"] = flattenAndSortScopeEntityIds(exclusions.Buildings)
	}

	if len(exclusions.JSSUsers) > 0 {
		result["jss_user_ids"] = flattenAndSortScopeEntityIds(exclusions.JSSUsers)
	}

	if len(exclusions.JSSUserGroups) > 0 {
		result["jss_user_group_ids"] = flattenAndSortScopeEntityIds(exclusions.JSSUserGroups)
	}

	if len(exclusions.Departments) > 0 {
		result["department_ids"] = flattenAndSortScopeEntityIds(exclusions.Departments)
	}

	if len(exclusions.NetworkSegments) > 0 {
		result["network_segment_ids"] = flattenAndSortNetworkSegmentIds(exclusions.NetworkSegments)
	}

	if len(exclusions.Users) > 0 {
		result["directory_service_or_local_usernames"] = flattenAndSortScopeEntityNames(exclusions.Users)
	}

	if len(exclusions.UserGroups) > 0 {
		result["directory_service_or_local_usergroup_ids"] = flattenAndSortScopeEntityIds(exclusions.UserGroups)
	}

	if len(exclusions.IBeacons) > 0 {
		result["ibeacon_ids"] = flattenAndSortScopeEntityIds(exclusions.IBeacons)
	}

	return []map[string]interface{}{result}, nil
}

// setSelfService converts the self-service structure into a format suitable for setting in the Terraform state.
func setSelfService(selfService jamfpro.MacOSConfigurationProfileSubsetSelfService) (map[string]interface{}, error) {
	selfServiceData := map[string]interface{}{
		"install_button_text":             selfService.InstallButtonText,
		"self_service_description":        selfService.SelfServiceDescription,
		"force_users_to_view_description": selfService.ForceUsersToViewDescription,
		"feature_on_main_page":            selfService.FeatureOnMainPage,
		"notification_subject":            selfService.NotificationSubject,
		"notification_message":            selfService.NotificationMessage,
	}

	// Fix the notification field
	if len(selfService.Notification) > 0 {
		correctNotifValue, err := FixDuplicateNotificationKey(&jamfpro.ResourceMacOSConfigurationProfile{
			SelfService: selfService,
		})
		if err != nil {
			return nil, err
		}
		selfServiceData["notification"] = correctNotifValue
	} else {
		selfServiceData["notification"] = false // Default to false if no valid boolean value is found
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

	return selfServiceData, nil
}

// helper functions

// flattenAndSortScopeEntityIds converts a slice of general scope entities (like user groups, buildings) to a format suitable for Terraform state.
func flattenAndSortScopeEntityIds(entities []jamfpro.MacOSConfigurationProfileSubsetScopeEntity) []int {
	var ids []int
	for _, entity := range entities {
		ids = append(ids, entity.ID)
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortScopeEntityNames converts a slice of RestrictedSoftwareSubsetScopeEntity into a sorted slice of strings.
func flattenAndSortScopeEntityNames(entities []jamfpro.MacOSConfigurationProfileSubsetScopeEntity) []string {
	var names []string
	for _, entity := range entities {
		names = append(names, entity.Name)
	}
	sort.Strings(names)
	return names
}

// flattenAndSortComputerIds converts a slice of MacOSConfigurationProfileSubsetComputer into a sorted slice of integers.
func flattenAndSortComputerIds(computers []jamfpro.MacOSConfigurationProfileSubsetComputer) []int {
	var ids []int
	for _, computer := range computers {
		ids = append(ids, computer.ID)
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortNetworkSegmentIds converts a slice of MacOSConfigurationProfileSubsetNetworkSegment into a sorted slice of integers.
func flattenAndSortNetworkSegmentIds(segments []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment) []int {
	var ids []int
	for _, segment := range segments {
		ids = append(ids, segment.ID)
	}
	sort.Ints(ids)
	return ids
}
