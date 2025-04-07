// macosconfigurationprofilesplist_object.go
package macosconfigurationprofilesplist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
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
		exclusionData := exclusions.([]interface{})[0]
		if exclusionMap, ok := exclusionData.(map[string]interface{}); ok {
			scope.Exclusions = constructExclusions(exclusionMap)
		}
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

	if computerIDs, ok := data["computer_ids"]; ok && len(computerIDs.([]interface{})) > 0 {
		exclusions.Computers = constructComputers(computerIDs.([]interface{}))
	}
	if computerGroupIDs, ok := data["computer_group_ids"]; ok && len(computerGroupIDs.([]interface{})) > 0 {
		exclusions.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDs.([]interface{}))
	}
	if userIDs, ok := data["jss_user_ids"]; ok && len(userIDs.([]interface{})) > 0 {
		exclusions.JSSUsers = constructScopeEntitiesFromIds(userIDs.([]interface{}))
	}
	if userGroupIDs, ok := data["jss_user_group_ids"]; ok && len(userGroupIDs.([]interface{})) > 0 {
		exclusions.JSSUserGroups = constructScopeEntitiesFromIds(userGroupIDs.([]interface{}))
	}
	if buildingIDs, ok := data["building_ids"]; ok && len(buildingIDs.([]interface{})) > 0 {
		exclusions.Buildings = constructScopeEntitiesFromIds(buildingIDs.([]interface{}))
	}
	if departmentIDs, ok := data["department_ids"]; ok && len(departmentIDs.([]interface{})) > 0 {
		exclusions.Departments = constructScopeEntitiesFromIds(departmentIDs.([]interface{}))
	}
	if networkSegmentIDs, ok := data["network_segment_ids"]; ok && len(networkSegmentIDs.([]interface{})) > 0 {
		exclusions.NetworkSegments = constructNetworkSegments(networkSegmentIDs.([]interface{}))
	}
	if ibeaconIDs, ok := data["ibeacon_ids"]; ok && len(ibeaconIDs.([]interface{})) > 0 {
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

		// Removed because there are several issues with this payload in the API
		// Will be reimplemented once those have been fixed.
		// Notification:                data["notification"].(string),

		NotificationSubject: data["notification_subject"].(string),
		NotificationMessage: data["notification_message"].(string),
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
