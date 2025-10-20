// macosconfigurationprofilesplistgenerator_constructor.go
package macos_configuration_profile_plist_generator

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
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

	if _, ok := d.GetOk("scope"); ok {
		resource.Scope = constructMacOSConfigurationProfileSubsetScope(d)
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

// constructMacOSConfigurationProfileSubsetScope reads the scope map and calls helpers
func constructMacOSConfigurationProfileSubsetScope(d *schema.ResourceData) jamfpro.MacOSConfigurationProfileSubsetScope {
	scope := jamfpro.MacOSConfigurationProfileSubsetScope{}

	// Get the scope data from the resource data
	scopeData := d.Get("scope").([]interface{})[0].(map[string]interface{})

	// Set basic boolean fields
	scope.AllComputers = scopeData["all_computers"].(bool)
	scope.AllJSSUsers = scopeData["all_jss_users"].(bool)

	// Use MapSetToStructs for computer IDs
	var computers []jamfpro.MacOSConfigurationProfileSubsetComputer
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.computer_ids", "ID", d, &computers); err != nil {
		log.Printf("[WARN] Error mapping computer IDs: %v", err)
	}
	scope.Computers = computers

	// Computer groups
	var computerGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.computer_group_ids", "ID", d, &computerGroups); err != nil {
		log.Printf("[WARN] Error mapping computer group IDs: %v", err)
	}
	scope.ComputerGroups = computerGroups

	// Buildings
	var buildings []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.building_ids", "ID", d, &buildings); err != nil {
		log.Printf("[WARN] Error mapping building IDs: %v", err)
	}
	scope.Buildings = buildings

	// Departments
	var departments []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.department_ids", "ID", d, &departments); err != nil {
		log.Printf("[WARN] Error mapping department IDs: %v", err)
	}
	scope.Departments = departments

	// JSS Users
	var jssUsers []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.jss_user_ids", "ID", d, &jssUsers); err != nil {
		log.Printf("[WARN] Error mapping JSS user IDs: %v", err)
	}
	scope.JSSUsers = jssUsers

	// JSS User Groups
	var jssUserGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.jss_user_group_ids", "ID", d, &jssUserGroups); err != nil {
		log.Printf("[WARN] Error mapping JSS user group IDs: %v", err)
	}
	scope.JSSUserGroups = jssUserGroups

	// Handle Limitations
	if _, ok := d.GetOk("scope.0.limitations"); ok {
		scope.Limitations = constructLimitations(d)
	}

	// Handle Exclusions
	if _, ok := d.GetOk("scope.0.exclusions"); ok {
		scope.Exclusions = constructExclusions(d)
	}

	return scope
}

// constructLimitations builds the limitations object using MapSetToStructs
func constructLimitations(d *schema.ResourceData) jamfpro.MacOSConfigurationProfileSubsetLimitations {
	limitations := jamfpro.MacOSConfigurationProfileSubsetLimitations{}

	// Directory service users (strings)
	var users []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, string](
		"scope.0.limitations.0.directory_service_or_local_usernames", "Name", d, &users); err != nil {
		log.Printf("[WARN] Error mapping user names: %v", err)
	}
	limitations.Users = users

	// Directory service user groups
	var userGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, &userGroups); err != nil {
		log.Printf("[WARN] Error mapping user group IDs: %v", err)
	}
	limitations.UserGroups = userGroups

	// Network segments
	var networkSegments []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int](
		"scope.0.limitations.0.network_segment_ids", "ID", d, &networkSegments); err != nil {
		log.Printf("[WARN] Error mapping network segment IDs: %v", err)
	}
	limitations.NetworkSegments = networkSegments

	// iBeacons
	var iBeacons []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.limitations.0.ibeacon_ids", "ID", d, &iBeacons); err != nil {
		log.Printf("[WARN] Error mapping iBeacon IDs: %v", err)
	}
	limitations.IBeacons = iBeacons

	return limitations
}

// constructExclusions builds the exclusions object using MapSetToStructs
func constructExclusions(d *schema.ResourceData) jamfpro.MacOSConfigurationProfileSubsetExclusions {
	exclusions := jamfpro.MacOSConfigurationProfileSubsetExclusions{}

	// Computers
	var computers []jamfpro.MacOSConfigurationProfileSubsetComputer
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetComputer, int](
		"scope.0.exclusions.0.computer_ids", "ID", d, &computers); err != nil {
		log.Printf("[WARN] Error mapping excluded computer IDs: %v", err)
	}
	exclusions.Computers = computers

	// Computer groups
	var computerGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.computer_group_ids", "ID", d, &computerGroups); err != nil {
		log.Printf("[WARN] Error mapping excluded computer group IDs: %v", err)
	}
	exclusions.ComputerGroups = computerGroups

	// JSS Users
	var jssUsers []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_ids", "ID", d, &jssUsers); err != nil {
		log.Printf("[WARN] Error mapping excluded JSS user IDs: %v", err)
	}
	exclusions.JSSUsers = jssUsers

	// JSS User Groups
	var jssUserGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_group_ids", "ID", d, &jssUserGroups); err != nil {
		log.Printf("[WARN] Error mapping excluded JSS user group IDs: %v", err)
	}
	exclusions.JSSUserGroups = jssUserGroups

	// Buildings
	var buildings []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.building_ids", "ID", d, &buildings); err != nil {
		log.Printf("[WARN] Error mapping excluded building IDs: %v", err)
	}
	exclusions.Buildings = buildings

	// Departments
	var departments []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.department_ids", "ID", d, &departments); err != nil {
		log.Printf("[WARN] Error mapping excluded department IDs: %v", err)
	}
	exclusions.Departments = departments

	// Network segments
	var networkSegments []jamfpro.MacOSConfigurationProfileSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int](
		"scope.0.exclusions.0.network_segment_ids", "ID", d, &networkSegments); err != nil {
		log.Printf("[WARN] Error mapping excluded network segment IDs: %v", err)
	}
	exclusions.NetworkSegments = networkSegments

	// User names (strings)
	var users []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, string](
		"scope.0.exclusions.0.directory_service_or_local_usernames", "Name", d, &users); err != nil {
		log.Printf("[WARN] Error mapping excluded user names: %v", err)
	}
	exclusions.Users = users

	// User groups
	var userGroups []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.directory_service_usergroup_ids", "ID", d, &userGroups); err != nil {
		log.Printf("[WARN] Error mapping excluded user group IDs: %v", err)
	}
	exclusions.UserGroups = userGroups

	// iBeacons
	var iBeacons []jamfpro.MacOSConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.ibeacon_ids", "ID", d, &iBeacons); err != nil {
		log.Printf("[WARN] Error mapping excluded iBeacon IDs: %v", err)
	}
	exclusions.IBeacons = iBeacons

	return exclusions
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
			if cat.ID, idOk = constructors.ParseResourceID(idVal, "self service category", i); !idOk {
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
