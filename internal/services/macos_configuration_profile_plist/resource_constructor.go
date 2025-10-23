package macos_configuration_profile_plist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	helpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"
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
func constructJamfProMacOSConfigurationProfilePlist(d *schema.ResourceData, mode string, meta any) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
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

	if _, ok := d.GetOk("scope"); ok {
		resource.Scope = constructMacOSConfigurationProfileSubsetScope(d)
	}

	if v, ok := d.GetOk("self_service"); ok {
		selfServiceData := v.([]any)[0].(map[string]any)
		if selfServiceData["notification"] != nil {
			log.Println("[WARN] Self Service notification bool key is temporarily disabled, please review the docs.")

		}
		resource.SelfService = constructMacOSConfigurationProfileSubsetSelfService(selfServiceData)
	}

	if mode != "update" {
		resource.General.Payloads = html.EscapeString(d.Get("payloads").(string))

	} else if mode == "update" {
		var existingPlist map[string]any
		var newPlist map[string]any

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

		// Ensure nested UUIDs and PayloadIdentifiers are also matched properly
		uuidMap := make(map[string]string)
		identifierMap := make(map[string]string)
		helpers.ExtractUUIDs(existingPlist, uuidMap, true)
		helpers.ExtractPayloadIdentifiers(existingPlist, identifierMap, true)
		helpers.UpdateUUIDs(newPlist, uuidMap, identifierMap, true)

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

// constructMacOSConfigurationProfileSubsetScope reads the scope map and calls helpers
func constructMacOSConfigurationProfileSubsetScope(d *schema.ResourceData) jamfpro.MacOSConfigurationProfileSubsetScope {
	scope := jamfpro.MacOSConfigurationProfileSubsetScope{}

	scopeData := d.Get("scope").([]any)[0].(map[string]any)

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
func constructMacOSConfigurationProfileSubsetSelfService(data map[string]any) jamfpro.MacOSConfigurationProfileSubsetSelfService {
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
func constructSelfServiceCategories(categoryList []any) []jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory {
	selfServiceCategories := make([]jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory, 0, len(categoryList))
	for i, category := range categoryList {
		catData, mapOk := category.(map[string]any)
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
