// macosconfigurationprofilesplistgenerator_constructor.go
package macosconfigurationprofilesplistgenerator

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProMacOSConfigurationProfilesPlistGenerator constructs a ResourceMacOSConfigurationProfile object from the provided schema data.
func constructJamfProMacOSConfigurationProfilesPlistGenerator(d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
	var resource *jamfpro.ResourceMacOSConfigurationProfile

	plistXML, err := plist.ConvertHCLToPlist(d)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plist from payloads: %v", err)
	}

	resource = &jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string),
			RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
			Payloads:           html.EscapeString(plistXML),
		},
	}

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		resource.Scope = constructMacOSConfigurationProfileSubsetScope(scopeData)
	}

	if v, ok := d.GetOk("self_service"); ok {
		selfServiceData := v.([]interface{})[0].(map[string]interface{})
		if selfServiceData["notification"] != nil {
			log.Println("[WARN] Self Service notification bool key is temporarily disabled, please review the docs.")

		}
		resource.SelfService = constructMacOSConfigurationProfileSubsetSelfService(selfServiceData)
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro macOS Configuration Profile '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro macOS Configuration Profile XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// --- Scope Construction Functions (Adapted for TypeSet) ---

// constructMacOSConfigurationProfileSubsetScope reads the scope map and calls helpers.
func constructMacOSConfigurationProfileSubsetScope(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetScope {
	scope := jamfpro.MacOSConfigurationProfileSubsetScope{
		AllComputers: data["all_computers"].(bool),
		AllJSSUsers:  data["all_jss_users"].(bool),
	}

	// Use constructors.GetListFromSet for *Set fields
	if computerIDsList := constructors.GetListFromSet(data, "computer_ids"); len(computerIDsList) > 0 {
		scope.Computers = constructComputers(computerIDsList)
	}
	if computerGroupIDsList := constructors.GetListFromSet(data, "computer_group_ids"); len(computerGroupIDsList) > 0 {
		scope.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDsList)
	}
	if buildingIDsList := constructors.GetListFromSet(data, "building_ids"); len(buildingIDsList) > 0 {
		scope.Buildings = constructScopeEntitiesFromIds(buildingIDsList)
	}
	if departmentIDsList := constructors.GetListFromSet(data, "department_ids"); len(departmentIDsList) > 0 {
		scope.Departments = constructScopeEntitiesFromIds(departmentIDsList)
	}
	if jssUserIDsList := constructors.GetListFromSet(data, "jss_user_ids"); len(jssUserIDsList) > 0 {
		scope.JSSUsers = constructScopeEntitiesFromIds(jssUserIDsList)
	}
	if jssUserGroupIDsList := constructors.GetListFromSet(data, "jss_user_group_ids"); len(jssUserGroupIDsList) > 0 {
		scope.JSSUserGroups = constructScopeEntitiesFromIds(jssUserGroupIDsList)
	}

	// Handle Limitations Block (Outer TypeList, Inner Map with TypeSet fields)
	if limitationsVal, ok := data["limitations"]; ok && limitationsVal != nil {
		limitationList, listOk := limitationsVal.([]interface{})
		if listOk && len(limitationList) > 0 && limitationList[0] != nil {
			limitationData, mapOk := limitationList[0].(map[string]interface{})
			if mapOk {
				scope.Limitations = constructLimitations(limitationData)
			} else {
				log.Printf("[WARN] constructMacOSConfigurationProfileSubsetScope: Could not cast limitations element to map.")
			}
		} else if !listOk {
			log.Printf("[WARN] constructMacOSConfigurationProfileSubsetScope: Unexpected type for limitations: %T", limitationsVal)
		}
	}

	// Handle Exclusions Block (Outer TypeList, Inner Map with TypeSet fields)
	if exclusionsVal, ok := data["exclusions"]; ok && exclusionsVal != nil {
		exclusionList, listOk := exclusionsVal.([]interface{})
		if listOk && len(exclusionList) > 0 && exclusionList[0] != nil {
			exclusionData, mapOk := exclusionList[0].(map[string]interface{})
			if mapOk {
				scope.Exclusions = constructExclusions(exclusionData)
			} else {
				log.Printf("[WARN] constructMacOSConfigurationProfileSubsetScope: Could not cast exclusions element to map.")
			}
		} else if !listOk {
			log.Printf("[WARN] constructMacOSConfigurationProfileSubsetScope: Unexpected type for exclusions: %T", exclusionsVal)
		}
	}

	return scope
}

// constructLimitations reads the limitations map and calls helpers.
func constructLimitations(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetLimitations {
	limitations := jamfpro.MacOSConfigurationProfileSubsetLimitations{}

	// Use constructors.GetListFromSet for *Set fields inside the limitations map
	if userNamesList := constructors.GetListFromSet(data, "directory_service_or_local_usernames"); len(userNamesList) > 0 {
		limitations.Users = constructScopeEntitiesFromIdsFromNames(userNamesList)
	}
	if userGroupIDsList := constructors.GetListFromSet(data, "directory_service_usergroup_ids"); len(userGroupIDsList) > 0 {
		limitations.UserGroups = constructScopeEntitiesFromIds(userGroupIDsList)
	}
	if networkSegmentIDsList := constructors.GetListFromSet(data, "network_segment_ids"); len(networkSegmentIDsList) > 0 {
		limitations.NetworkSegments = constructNetworkSegments(networkSegmentIDsList)
	}
	if ibeaconIDsList := constructors.GetListFromSet(data, "ibeacon_ids"); len(ibeaconIDsList) > 0 {
		limitations.IBeacons = constructScopeEntitiesFromIds(ibeaconIDsList)
	}

	return limitations
}

// constructExclusions reads the exclusions map and calls helpers.
func constructExclusions(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetExclusions {
	exclusions := jamfpro.MacOSConfigurationProfileSubsetExclusions{}

	// Use constructors.GetListFromSet for *Set fields inside the exclusions map
	if computerIDsList := constructors.GetListFromSet(data, "computer_ids"); len(computerIDsList) > 0 {
		exclusions.Computers = constructComputers(computerIDsList)
	}
	if computerGroupIDsList := constructors.GetListFromSet(data, "computer_group_ids"); len(computerGroupIDsList) > 0 {
		exclusions.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDsList)
	}
	if userIDsList := constructors.GetListFromSet(data, "jss_user_ids"); len(userIDsList) > 0 {
		exclusions.JSSUsers = constructScopeEntitiesFromIds(userIDsList)
	}
	if userGroupIDsList := constructors.GetListFromSet(data, "jss_user_group_ids"); len(userGroupIDsList) > 0 {
		exclusions.JSSUserGroups = constructScopeEntitiesFromIds(userGroupIDsList)
	}
	if buildingIDsList := constructors.GetListFromSet(data, "building_ids"); len(buildingIDsList) > 0 {
		exclusions.Buildings = constructScopeEntitiesFromIds(buildingIDsList)
	}
	if departmentIDsList := constructors.GetListFromSet(data, "department_ids"); len(departmentIDsList) > 0 {
		exclusions.Departments = constructScopeEntitiesFromIds(departmentIDsList)
	}
	if networkSegmentIDsList := constructors.GetListFromSet(data, "network_segment_ids"); len(networkSegmentIDsList) > 0 {
		exclusions.NetworkSegments = constructNetworkSegments(networkSegmentIDsList)
	}
	if userNamesList := constructors.GetListFromSet(data, "directory_service_or_local_usernames"); len(userNamesList) > 0 {
		exclusions.Users = constructScopeEntitiesFromIdsFromNames(userNamesList)
	}
	if userGroupIDsList := constructors.GetListFromSet(data, "directory_service_usergroup_ids"); len(userGroupIDsList) > 0 {
		exclusions.UserGroups = constructScopeEntitiesFromIds(userGroupIDsList)
	}
	if ibeaconIDsList := constructors.GetListFromSet(data, "ibeacon_ids"); len(ibeaconIDsList) > 0 {
		exclusions.IBeacons = constructScopeEntitiesFromIds(ibeaconIDsList)
	}

	return exclusions
}

// --- Helpers

// constructComputers uses robust conversion.
func constructComputers(ids []interface{}) []jamfpro.MacOSConfigurationProfileSubsetComputer {
	computers := make([]jamfpro.MacOSConfigurationProfileSubsetComputer, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "computer", i); ok {
			computers = append(computers, jamfpro.MacOSConfigurationProfileSubsetComputer{
				MacOSConfigurationProfileSubsetScopeEntity: jamfpro.MacOSConfigurationProfileSubsetScopeEntity{ID: intID},
			})
		}
	}
	log.Printf("[DEBUG] constructComputers: Input count %d, Output count %d", len(ids), len(computers))
	return computers
}

// constructNetworkSegments uses robust conversion.
func constructNetworkSegments(ids []interface{}) []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment {
	networkSegments := make([]jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "network segment", i); ok {
			networkSegments = append(networkSegments, jamfpro.MacOSConfigurationProfileSubsetNetworkSegment{
				MacOSConfigurationProfileSubsetScopeEntity: jamfpro.MacOSConfigurationProfileSubsetScopeEntity{ID: intID},
			})
		}
	}
	log.Printf("[DEBUG] constructNetworkSegments: Input count %d, Output count %d", len(ids), len(networkSegments))
	return networkSegments
}

// constructScopeEntitiesFromIds uses robust conversion.
func constructScopeEntitiesFromIds(ids []interface{}) []jamfpro.MacOSConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MacOSConfigurationProfileSubsetScopeEntity, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "scope entity", i); ok {
			scopeEntities = append(scopeEntities, jamfpro.MacOSConfigurationProfileSubsetScopeEntity{ID: intID})
		}
	}
	log.Printf("[DEBUG] constructScopeEntitiesFromIds: Input count %d, Output count %d", len(ids), len(scopeEntities))
	return scopeEntities
}

// constructScopeEntitiesFromIdsFromNames uses robust conversion.
func constructScopeEntitiesFromIdsFromNames(names []interface{}) []jamfpro.MacOSConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MacOSConfigurationProfileSubsetScopeEntity, 0, len(names))
	for i, nameRaw := range names {
		if strName, ok := nameRaw.(string); ok && strName != "" {
			scopeEntities = append(scopeEntities, jamfpro.MacOSConfigurationProfileSubsetScopeEntity{Name: strName})
		} else if !ok {
			log.Printf("[WARN] constructScopeEntitiesFromIdsFromNames: Unexpected type %T for scope entity name: %v at index %d. Skipping.", nameRaw, nameRaw, i)
		} // Don't warn for empty string, just skip it
	}
	log.Printf("[DEBUG] constructScopeEntitiesFromIdsFromNames: Input count %d, Output count %d", len(names), len(scopeEntities))
	return scopeEntities
}

// --- Self Service Construction Functions (Adapted for TypeSet) ---

// constructMacOSConfigurationProfileSubsetSelfService reads the self_service map.
func constructMacOSConfigurationProfileSubsetSelfService(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetSelfService {
	selfService := jamfpro.MacOSConfigurationProfileSubsetSelfService{}

	// Use type assertion with ok check for safety
	if val, ok := data["self_service_display_name"].(string); ok {
		selfService.SelfServiceDisplayName = val
	}
	if val, ok := data["install_button_text"].(string); ok {
		selfService.InstallButtonText = val
	}
	if val, ok := data["self_service_description"].(string); ok {
		selfService.SelfServiceDescription = val
	}
	if val, ok := data["force_users_to_view_description"].(bool); ok {
		selfService.ForceUsersToViewDescription = val
	}
	if val, ok := data["feature_on_main_page"].(bool); ok {
		selfService.FeatureOnMainPage = val
	}
	// Removed because there are several issues with this payload in the API
	// Will be reimplemented once those have been fixed.
	// if val, ok := data["notification"].(string); ok {
	// 	selfService.NotificationSubject = val
	// }
	if val, ok := data["notification_subject"].(string); ok {
		selfService.NotificationSubject = val
	}
	if val, ok := data["notification_message"].(string); ok {
		selfService.NotificationMessage = val
	}

	if iconID, ok := data["self_service_icon_id"].(int); ok && iconID != 0 {
		selfService.SelfServiceIcon = jamfpro.SharedResourceSelfServiceIcon{ID: iconID}
	}

	if categoriesList := constructors.GetListFromSet(data, "self_service_category"); len(categoriesList) > 0 {
		selfService.SelfServiceCategories = constructSelfServiceCategories(categoriesList)
	}

	return selfService
}

// constructSelfServiceCategories processes the list from constructors.GetListFromSet.
func constructSelfServiceCategories(categoryList []interface{}) []jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory {
	selfServiceCategories := make([]jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory, 0, len(categoryList))
	for i, category := range categoryList {
		catData, mapOk := category.(map[string]interface{})
		if !mapOk {
			log.Printf("[WARN] constructSelfServiceCategories: Could not cast category data to map at index %d. Skipping.", i)
			continue
		}

		cat := jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory{}
		idOk := false
		if idVal, ok := catData["id"]; ok {
			if cat.ID, idOk = constructors.ConvertToInt(idVal, "self service category", i); !idOk {
				continue // Skip if ID conversion fails
			}
		} else {
			log.Printf("[WARN] constructSelfServiceCategories: Missing 'id' key in category map at index %d. Skipping.", i)
			continue // Skip if ID is missing
		}

		if displayInVal, ok := catData["display_in"].(bool); ok {
			cat.DisplayIn = displayInVal
		}
		if featureInVal, ok := catData["feature_in"].(bool); ok {
			cat.FeatureIn = featureInVal
		}
		// Name is read-only from API, not set here

		selfServiceCategories = append(selfServiceCategories, cat)
	}
	log.Printf("[DEBUG] constructSelfServiceCategories: Input count %d, Output count %d", len(categoryList), len(selfServiceCategories))
	return selfServiceCategories
}
