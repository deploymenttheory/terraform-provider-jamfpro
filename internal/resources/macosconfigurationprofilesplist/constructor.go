package macosconfigurationprofilesplist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/constructors"
	helpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// constructJamfProMacOSConfigurationProfilePlist constructs a ResourceMacOSConfigurationProfile object from schema data.
// It supports two modes:
//   - create: Builds profile from schema data only
//   - update: Fetches existing profile from Jamf Pro, extracts PayloadUUID/PayloadIdentifier values from existing plist,
//     injects them into the new plist to maintain UUID continuity
//
// The function:
// 1. For update mode:
//   - Retrieves existing profile from Jamf Pro API
//   - Decodes existing plist to extract UUIDs
//
// 2. Constructs base profile from schema data (name, description, etc)
// 3. Builds scope and self-service sections if configured
// 4. For update mode:
//   - Maps existing UUIDs by PayloadDisplayName
//   - Updates PayloadUUID/PayloadIdentifier in new plist to match existing
//   - Re-encodes updated plist
//
// Parameters:
// - d: Schema ResourceData containing configuration
// - mode: "create" or "update" to control UUID handling
// - meta: Provider meta containing client for API calls
//
// Returns:
// - Constructed ResourceMacOSConfigurationProfile
// - Error if construction or API calls fail
//
// Jamf Pro modifies the top-level PayloadUUID and PayloadIdentifier upon profile creation.
// Nested payload identifiers and UUIDs remain unchanged from the original request.
// Therefore, when performing profile updates, only top-level PayloadUUID and PayloadIdentifier
// need to be synced from Jamf Pro's existing profile state.
func constructJamfProMacOSConfigurationProfilePlist(d *schema.ResourceData, mode string, meta interface{}) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
	var existingProfile *jamfpro.ResourceMacOSConfigurationProfile
	var buf bytes.Buffer

	resource := &jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string),
			RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
			// We'll handle payloads att differently based on mode
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

	if mode != "update" {
		resource.General.Payloads = html.EscapeString(d.Get("payloads").(string))

	} else if mode == "update" {
		var existingPlist map[string]interface{}
		var newPlist map[string]interface{}

		client := meta.(*jamfpro.Client)
		resourceID := d.Id()
		var err error
		existingProfile, err = client.GetMacOSConfigurationProfileByID(resourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing configuration profile by ID for update operation: %v", err)
		}

		// Decode existing payload from Jamf Pro which has the jamf pro post processed uuid's etc
		existingPayload := existingProfile.General.Payloads
		if err := plist.NewDecoder(strings.NewReader(existingPayload)).Decode(&existingPlist); err != nil {
			return nil, fmt.Errorf("failed to decode existing plist payload stored in jamf pro for update operation: %v", err)
		}

		// Decode payloads field from Terraform state ready for injection
		newPayload := d.Get("payloads").(string)
		if err := plist.NewDecoder(strings.NewReader(newPayload)).Decode(&newPlist); err != nil {
			return nil, fmt.Errorf("failed to decode new plist payload from terraform state for update operation: %v", err)
		}

		// Jamf Pro modifies only the top-level PayloadUUID and PayloadIdentifier upon profile creation.
		// All nested payload UUIDs/identifiers remain unchanged.
		// Copy top-level PayloadUUID and PayloadIdentifier from existing (Jamf Pro) to new (Terraform)
		newPlist["PayloadUUID"] = existingPlist["PayloadUUID"]
		newPlist["PayloadIdentifier"] = existingPlist["PayloadIdentifier"]

		// Ensure nested UUIDs are also matched properly
		uuidMap := make(map[string]string)
		helpers.ExtractUUIDs(existingPlist, uuidMap, true)
		helpers.UpdateUUIDs(newPlist, uuidMap, true)

		var mismatches []string
		helpers.ValidatePayloadUUIDsMatch(existingPlist, newPlist, "Payload", &mismatches)

		if len(mismatches) > 0 {
			return nil, fmt.Errorf("configuration profile UUID mismatch found:\n%s", strings.Join(mismatches, "\n"))
		}

		// Encode the plist with injections

		encoder := plist.NewEncoder(&buf)
		encoder.Indent("    ")
		if err := encoder.Encode(newPlist); err != nil {
			return nil, fmt.Errorf("failed to encode plist payload with injected PayloadUUID and PayloadIdentifier: %v", err)
		}

		// Since we're embedding a Plist (which is XML) inside another XML document (the request),
		// we need to properly correctly normalize the XML for the xml.MarshalIndent and also for jamf pro.
		if buf.Len() > 0 {
			unquotedContent := preMarshallingXMLPayloadUnescaping(buf.String())
			resource.General.Payloads = preMarshallingXMLPayloadEscaping(unquotedContent)
		}
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro macOS Configuration Profile '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro macOS Configuration Profile XML:\n%s\n", resourceXML)

	return resource, nil
}

// preMarshallingXMLPayloadUnescaping unescapes content ready for jamf pro based on plist reqs
func preMarshallingXMLPayloadUnescaping(input string) string {
	input = strings.ReplaceAll(input, "&#34;", "\"")
	return input
}

// preMarshallingXMLPayloadEscaping ensures that the XML marshaller (used in xml.MarshalIndent)
// doesn't choke on special XML characters (&) inside the payload
func preMarshallingXMLPayloadEscaping(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	return input
}

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
		if intID, ok := constructors.ParseResourceID(idRaw, "computer", i); ok {
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
		if intID, ok := constructors.ParseResourceID(idRaw, "network segment", i); ok {
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
		if intID, ok := constructors.ParseResourceID(idRaw, "scope entity", i); ok {
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
			if cat.ID, idOk = constructors.ParseResourceID(idVal, "self service category", i); !idOk {
				continue
			}
		} else {
			log.Printf("[WARN] constructSelfServiceCategories: Missing 'id' key in category map at index %d. Skipping.", i)
			continue
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
