// mobiledeviceconfigurationprofilesplist_object.go
package mobiledeviceconfigurationprofilesplist

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

// constructJamfProMobileDeviceConfigurationProfilePlist constructs a ResourceMobileDeviceConfigurationProfile object from schema data.
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
func constructJamfProMobileDeviceConfigurationProfilePlist(d *schema.ResourceData, mode string, meta interface{}) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
	var existingProfile *jamfpro.ResourceMobileDeviceConfigurationProfile

	resource := &jamfpro.ResourceMobileDeviceConfigurationProfile{
		General: jamfpro.MobileDeviceConfigurationProfileSubsetGeneral{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Level:            d.Get("level").(string),
			UUID:             d.Get("uuid").(string),
			DeploymentMethod: d.Get("deployment_method").(string),
			RedeployOnUpdate: d.Get("redeploy_on_update").(string),
			// We'll handle payloads att differently based on mode
		},
	}

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	// Handle Scope
	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		resource.Scope = constructMobileDeviceConfigurationProfileSubsetScope(scopeData)
	}

	if mode != "update" {

		resource.General.Payloads = html.EscapeString(d.Get("payloads").(string))

	} else if mode == "update" {
		var existingPlist map[string]interface{}
		var newPlist map[string]interface{}

		client := meta.(*jamfpro.Client)
		resourceID := d.Id()
		var err error
		existingProfile, err = client.GetMobileDeviceConfigurationProfileByID(resourceID)
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
		var buf bytes.Buffer
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

	log.Printf("[DEBUG] Constructed Jamf Pro macOS Configuration Profile XML:\n%s\n", string(resourceXML))

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

// constructMobileDeviceConfigurationProfileSubsetScope constructs a MobileDeviceConfigurationProfileSubsetScope object from the provided schema data.
func constructMobileDeviceConfigurationProfileSubsetScope(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetScope {
	scope := jamfpro.MobileDeviceConfigurationProfileSubsetScope{
		AllMobileDevices: data["all_mobile_devices"].(bool),
		AllJSSUsers:      data["all_jss_users"].(bool),
	}

	if mobileDeviceIDs, ok := data["mobile_device_ids"]; ok {
		scope.MobileDevices = constructMobileDevices(mobileDeviceIDs.([]interface{}))
	}
	if mobileDeviceGroupIDs, ok := data["mobile_device_group_ids"]; ok {
		scope.MobileDeviceGroups = constructScopeEntitiesFromIds(mobileDeviceGroupIDs.([]interface{}))
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

	// Handle Limitations
	if limitations, ok := data["limitations"]; ok && len(limitations.([]interface{})) > 0 {
		limitationData := limitations.([]interface{})[0].(map[string]interface{})
		scope.Limitations = constructLimitations(limitationData)
	}

	// Handle Exclusions
	if exclusions, ok := data["exclusions"]; ok && len(exclusions.([]interface{})) > 0 {
		exclusionData := exclusions.([]interface{})[0].(map[string]interface{})
		scope.Exclusions = constructExclusions(exclusionData)
	}

	return scope
}

// constructLimitations constructs a MobileDeviceConfigurationProfileSubsetLimitation object from the provided schema data.
func constructLimitations(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetLimitation {
	limitations := jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{}

	if userNames, ok := data["directory_service_or_local_usernames"]; ok {
		limitations.Users = constructScopeEntitiesFromIdsFromNames(userNames.([]interface{}))
	}
	if userGroupIDs, ok := data["user_group_ids"]; ok {
		limitations.UserGroups = constructScopeEntitiesFromIds(userGroupIDs.([]interface{}))
	}
	if networkSegmentIDs, ok := data["network_segment_ids"]; ok {
		limitations.NetworkSegments = constructNetworkSegments(networkSegmentIDs.([]interface{}))
	}
	if ibeaconIDs, ok := data["ibeacon_ids"]; ok {
		limitations.Ibeacons = constructScopeEntitiesFromIds(ibeaconIDs.([]interface{}))
	}

	return limitations
}

// constructExclusions constructs a MobileDeviceConfigurationProfileSubsetExclusion object from the provided schema data.
func constructExclusions(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetExclusion {
	exclusions := jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{}

	if mobileDeviceIDs, ok := data["mobile_device_ids"]; ok {
		exclusions.MobileDevices = constructMobileDevices(mobileDeviceIDs.([]interface{}))
	}
	if mobileDeviceGroupIDs, ok := data["mobile_device_group_ids"]; ok {
		exclusions.MobileDeviceGroups = constructScopeEntitiesFromIds(mobileDeviceGroupIDs.([]interface{}))
	}
	if userIDs, ok := data["user_ids"]; ok {
		exclusions.Users = constructScopeEntitiesFromIds(userIDs.([]interface{}))
	}
	if userGroupIDs, ok := data["user_group_ids"]; ok {
		exclusions.UserGroups = constructScopeEntitiesFromIds(userGroupIDs.([]interface{}))
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
	if jssUserIDs, ok := data["jss_user_ids"]; ok {
		exclusions.JSSUsers = constructScopeEntitiesFromIds(jssUserIDs.([]interface{}))
	}
	if jssUserGroupIDs, ok := data["jss_user_group_ids"]; ok {
		exclusions.JSSUserGroups = constructScopeEntitiesFromIds(jssUserGroupIDs.([]interface{}))
	}

	return exclusions
}

// constructMobileDevices constructs a slice of MobileDeviceConfigurationProfileSubsetMobileDevice from the provided schema data.
func constructMobileDevices(ids []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice {
	mobileDevices := make([]jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice, len(ids))
	for i, id := range ids {
		mobileDevices[i] = jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice{
			ID: id.(int),
		}
	}
	return mobileDevices
}

// constructNetworkSegments constructs a slice of MobileDeviceConfigurationProfileSubsetNetworkSegment from the provided schema data.
func constructNetworkSegments(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment {
	networkSegments := make([]jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, len(data))
	for i, id := range data {
		networkSegments[i] = jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment{
			MobileDeviceConfigurationProfileSubsetScopeEntity: jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
				ID: id.(int),
			},
		}
	}
	return networkSegments
}

// Helper functions for nested structures

// constructScopeEntitiesFromIds constructs a slice of MobileDeviceConfigurationProfileSubsetScopeEntity from a list of IDs.
func constructScopeEntitiesFromIds(ids []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, len(ids))
	for i, id := range ids {
		scopeEntities[i] = jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
			ID: id.(int),
		}
	}
	return scopeEntities
}

// constructScopeEntitiesFromIdsFromNames constructs a slice of MobileDeviceConfigurationProfileSubsetScopeEntity from a list of names.
func constructScopeEntitiesFromIdsFromNames(names []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity {
	scopeEntities := make([]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, len(names))
	for i, name := range names {
		scopeEntities[i] = jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
			Name: name.(string),
		}
	}
	return scopeEntities
}
