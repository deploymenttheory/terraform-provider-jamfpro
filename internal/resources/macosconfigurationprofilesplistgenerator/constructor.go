// macosconfigurationprofilesplistgenerator_constructor.go
package macosconfigurationprofilesplistgenerator

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
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
		resource.SelfService = constructMacOSConfigurationProfileSubsetSelfService(selfServiceData)
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro macOS Configuration Profile '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro macOS Configuration Profile XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// constructMacOSConfigurationProfileSubsetScope constructs a MacOSConfigurationProfileSubsetScope object from the provided schema data.
func constructMacOSConfigurationProfileSubsetScope(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetScope {
	scope := jamfpro.MacOSConfigurationProfileSubsetScope{
		AllComputers: data["all_computers"].(bool),
		AllJSSUsers:  data["all_jss_users"].(bool),
	}

	if computerIDs, ok := data["computer_ids"]; ok {
		scope.Computers = constructComputers(computerIDs.([]interface{}))
	}
	if computerGroupIDs, ok := data["computer_group_ids"]; ok {
		scope.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDs.([]interface{}))
	}
	if buildingIDs, ok := data["building_ids"]; ok {
		scope.Buildings = constructScopeEntitiesFromIds(buildingIDs.([]interface{}))
	}
	if departmentIDs, ok := data["department_ids"]; ok {
		scope.Departments = constructScopeEntitiesFromIds(departmentIDs.([]interface{}))
	}
	if jssUserIDs, ok := data["jss_user_ids"]; ok {
		scope.JSSUsers = constructScopeEntitiesFromIds(jssUserIDs.([]interface{}))
	}
	if jssUserGroupIDs, ok := data["jss_user_group_ids"]; ok {
		scope.JSSUserGroups = constructScopeEntitiesFromIds(jssUserGroupIDs.([]interface{}))
	}

	if limitations, ok := data["limitations"]; ok && len(limitations.([]interface{})) > 0 {
		limitationData := limitations.([]interface{})[0].(map[string]interface{})
		scope.Limitations = constructLimitations(limitationData)
	}

	if exclusions, ok := data["exclusions"]; ok && len(exclusions.([]interface{})) > 0 {
		exclusionData := exclusions.([]interface{})[0].(map[string]interface{})
		scope.Exclusions = constructExclusions(exclusionData)
	}

	return scope
}

// constructLimitations constructs a MacOSConfigurationProfileSubsetLimitation object from the provided schema data.
func constructLimitations(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetLimitations {
	limitations := jamfpro.MacOSConfigurationProfileSubsetLimitations{}

	if userNames, ok := data["directory_service_or_local_usernames"]; ok {
		limitations.Users = constructScopeEntitiesFromIdsFromNames(userNames.([]interface{}))
	}
	if userGroupIDs, ok := data["directory_service_usergroup_ids"]; ok {
		limitations.UserGroups = constructScopeEntitiesFromIds(userGroupIDs.([]interface{}))
	}
	if networkSegmentIDs, ok := data["network_segment_ids"]; ok {
		limitations.NetworkSegments = constructNetworkSegments(networkSegmentIDs.([]interface{}))
	}
	if ibeaconIDs, ok := data["ibeacon_ids"]; ok {
		limitations.IBeacons = constructScopeEntitiesFromIds(ibeaconIDs.([]interface{}))
	}

	return limitations
}

// constructExclusions constructs a MacOSConfigurationProfileSubsetExclusion object from the provided schema data.
func constructExclusions(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetExclusions {
	exclusions := jamfpro.MacOSConfigurationProfileSubsetExclusions{}

	if computerIDs, ok := data["computer_ids"]; ok {
		exclusions.Computers = constructComputers(computerIDs.([]interface{}))
	}
	if computerGroupIDs, ok := data["computer_group_ids"]; ok {
		exclusions.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDs.([]interface{}))
	}
	if userIDs, ok := data["jss_user_ids"]; ok {
		exclusions.JSSUsers = constructScopeEntitiesFromIds(userIDs.([]interface{}))
	}
	if userGroupIDs, ok := data["jss_user_group_ids"]; ok {
		exclusions.JSSUserGroups = constructScopeEntitiesFromIds(userGroupIDs.([]interface{}))
	}
	if buildingIDs, ok := data["building_ids"]; ok {
		exclusions.Buildings = constructScopeEntitiesFromIds(buildingIDs.([]interface{}))
	}
	if departmentIDs, ok := data["department_ids"]; ok {
		exclusions.Departments = constructScopeEntitiesFromIds(departmentIDs.([]interface{}))
	}
	if networkSegmentIDs, ok := data["network_segment_ids"]; ok {
		exclusions.NetworkSegments = constructNetworkSegments(networkSegmentIDs.([]interface{}))
	}
	if ibeaconIDs, ok := data["ibeacon_ids"]; ok {
		exclusions.IBeacons = constructScopeEntitiesFromIds(ibeaconIDs.([]interface{}))
	}

	return exclusions
}

// constructComputers constructs a slice of MacOSConfigurationProfileSubsetComputer from the provided schema data.
func constructComputers(ids []interface{}) []jamfpro.MacOSConfigurationProfileSubsetComputer {
	computers := make([]jamfpro.MacOSConfigurationProfileSubsetComputer, len(ids))
	for i, id := range ids {
		computers[i] = jamfpro.MacOSConfigurationProfileSubsetComputer{
			MacOSConfigurationProfileSubsetScopeEntity: jamfpro.MacOSConfigurationProfileSubsetScopeEntity{
				ID: id.(int),
			},
		}
	}
	return computers
}

// constructNetworkSegments constructs a slice of MacOSConfigurationProfileSubsetNetworkSegment from the provided schema data.
func constructNetworkSegments(data []interface{}) []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment {
	networkSegments := make([]jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, len(data))
	for i, id := range data {
		networkSegments[i] = jamfpro.MacOSConfigurationProfileSubsetNetworkSegment{
			MacOSConfigurationProfileSubsetScopeEntity: jamfpro.MacOSConfigurationProfileSubsetScopeEntity{
				ID: id.(int),
			},
		}
	}
	return networkSegments
}

// constructMacOSConfigurationProfileSubsetSelfService constructs a MacOSConfigurationProfileSubsetSelfService object from the provided schema data.
func constructMacOSConfigurationProfileSubsetSelfService(data map[string]interface{}) jamfpro.MacOSConfigurationProfileSubsetSelfService {
	selfService := jamfpro.MacOSConfigurationProfileSubsetSelfService{
		SelfServiceDisplayName:      data["self_service_display_name"].(string),
		InstallButtonText:           data["install_button_text"].(string),
		SelfServiceDescription:      data["self_service_description"].(string),
		ForceUsersToViewDescription: data["force_users_to_view_description"].(bool),
		FeatureOnMainPage:           data["feature_on_main_page"].(bool),
		Notification:                data["notification"].(string),
		NotificationSubject:         data["notification_subject"].(string),
		NotificationMessage:         data["notification_message"].(string),
	}

	if iconID, ok := data["self_service_icon_id"].(int); ok && iconID != 0 {
		selfService.SelfServiceIcon = jamfpro.SharedResourceSelfServiceIcon{
			ID: iconID,
		}
	}

	if categories, ok := data["self_service_category"]; ok {
		selfService.SelfServiceCategories = constructSelfServiceCategories(categories.(*schema.Set))
	}

	return selfService
}

// constructSelfServiceCategories constructs a slice of MacOSConfigurationProfileSubsetSelfServiceCategory from the provided schema data.
func constructSelfServiceCategories(categories *schema.Set) []jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory {
	categoryList := categories.List()
	selfServiceCategories := make([]jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory, len(categoryList))
	for i, category := range categoryList {
		catData := category.(map[string]interface{})
		selfServiceCategories[i] = jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory{
			ID:        catData["id"].(int),
			DisplayIn: catData["display_in"].(bool),
			FeatureIn: catData["feature_in"].(bool),
		}
	}
	return selfServiceCategories
}

// Helper functions for nested structures

// constructScopeEntitiesFromIds constructs a slice of MacOSConfigurationProfileSubsetScopeEntity from a list of IDs.
func constructScopeEntitiesFromIds(ids []interface{}) []jamfpro.MacOSConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MacOSConfigurationProfileSubsetScopeEntity, len(ids))
	for i, id := range ids {
		scopeEntities[i] = jamfpro.MacOSConfigurationProfileSubsetScopeEntity{
			ID: id.(int),
		}
	}
	return scopeEntities
}

// constructScopeEntitiesFromIdsFromNames constructs a slice of MacOSConfigurationProfileSubsetScopeEntity from a list of names.
func constructScopeEntitiesFromIdsFromNames(names []interface{}) []jamfpro.MacOSConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MacOSConfigurationProfileSubsetScopeEntity, len(names))
	for i, name := range names {
		scopeEntities[i] = jamfpro.MacOSConfigurationProfileSubsetScopeEntity{
			Name: name.(string),
		}
	}
	return scopeEntities
}
