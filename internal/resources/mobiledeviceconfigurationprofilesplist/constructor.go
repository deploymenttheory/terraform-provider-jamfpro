// mobiledeviceconfigurationprofilesplist_object.go
package mobiledeviceconfigurationprofilesplist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"log" // Import strconv for robust conversion
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/constructors"
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
// The function now reads scope data assuming TypeSet in the schema.
//
// Parameters:
// - d: Schema ResourceData containing configuration
// - mode: "create" or "update" to control UUID handling
// - meta: Provider meta containing client for API calls
//
// Returns:
// - Constructed ResourceMobileDeviceConfigurationProfile
// - Error if construction or API calls fail
func constructJamfProMobileDeviceConfigurationProfilePlist(d *schema.ResourceData, mode string, meta interface{}) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
	var existingProfile *jamfpro.ResourceMobileDeviceConfigurationProfile
	var buf bytes.Buffer

	resource := &jamfpro.ResourceMobileDeviceConfigurationProfile{
		General: jamfpro.MobileDeviceConfigurationProfileSubsetGeneral{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Level:            d.Get("level").(string),
			UUID:             d.Get("uuid").(string),
			DeploymentMethod: d.Get("deployment_method").(string),
			RedeployOnUpdate: d.Get("redeploy_on_update").(string),
			// Payloads handled below based on mode
		},
	}

	if v, ok := d.GetOk("redeploy_days_before_cert_expires"); ok {
		resource.General.RedeployDaysBeforeCertExpires = v.(int)
	}

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	if v, ok := d.GetOk("scope"); ok {
		scopeList := v.([]interface{})
		if len(scopeList) > 0 && scopeList[0] != nil {
			scopeData, mapOk := scopeList[0].(map[string]interface{})
			if mapOk {
				resource.Scope = constructMobileDeviceConfigurationProfileSubsetScope(scopeData)
			} else {
				log.Printf("[WARN] constructJamfProMobileDeviceConfigurationProfilePlist: Could not cast scope element to map.")
			}
		}
	} else {
		log.Printf("[DEBUG] constructJamfProMobileDeviceConfigurationProfilePlist: No scope block found or it's empty.")
	}

	// Handle Payloads based on mode
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
			return nil, fmt.Errorf("failed to get existing mobile device configuration profile by ID %s for update: %v", resourceID, err)
		}

		existingPayload := existingProfile.General.Payloads
		existingPayload = html.UnescapeString(existingPayload)
		if err := plist.NewDecoder(strings.NewReader(existingPayload)).Decode(&existingPlist); err != nil {
			return nil, fmt.Errorf("failed to decode existing plist payload from Jamf Pro for update (ID: %s): %v\nPayload attempted:\n%s", resourceID, err, existingPayload)
		}

		newPayload := d.Get("payloads").(string)
		if err := plist.NewDecoder(strings.NewReader(newPayload)).Decode(&newPlist); err != nil {
			return nil, fmt.Errorf("failed to decode new plist payload from Terraform state for update: %v", err)
		}

		// Jamf Pro modifies only the top-level PayloadUUID and PayloadIdentifier upon profile creation.
		// All nested payload UUIDs/identifiers remain unchanged.
		// Copy top-level PayloadUUID and PayloadIdentifier from existing (Jamf Pro) to new (Terraform)
		newPlist["PayloadUUID"] = existingPlist["PayloadUUID"]
		newPlist["PayloadIdentifier"] = existingPlist["PayloadIdentifier"]

		uuidMap := make(map[string]string)
		helpers.ExtractUUIDs(existingPlist, uuidMap, true)
		helpers.UpdateUUIDs(newPlist, uuidMap, true)

		var mismatches []string
		helpers.ValidatePayloadUUIDsMatch(existingPlist, newPlist, "Payload", &mismatches)
		if len(mismatches) > 0 {
			log.Printf("[WARN] Mobile device configuration profile (ID: %s) UUID mismatches found after update attempt:\n%s", resourceID, strings.Join(mismatches, "\n"))
		}

		// Encode the plist with injections

		encoder := plist.NewEncoder(&buf)
		encoder.Indent("    ")
		if err := encoder.Encode(newPlist); err != nil {
			return nil, fmt.Errorf("failed to encode updated plist payload: %v", err)
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
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Configuration Profile '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Configuration Profile XML:\n%s\n", string(resourceXML))

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

// constructMobileDeviceConfigurationProfileSubsetScope constructs the scope using TypeSet from schema.
func constructMobileDeviceConfigurationProfileSubsetScope(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetScope {
	scope := jamfpro.MobileDeviceConfigurationProfileSubsetScope{
		AllMobileDevices: data["all_mobile_devices"].(bool),
		AllJSSUsers:      data["all_jss_users"].(bool),
	}

	// Use constructors.GetListFromSet for *Set fields
	if mobileDeviceIDsList := constructors.GetListFromSet(data, "mobile_device_ids"); len(mobileDeviceIDsList) > 0 {
		scope.MobileDevices = constructMobileDevices(mobileDeviceIDsList)
	}
	if mobileDeviceGroupIDsList := constructors.GetListFromSet(data, "mobile_device_group_ids"); len(mobileDeviceGroupIDsList) > 0 {
		scope.MobileDeviceGroups = constructMobileDeviceScopeEntitiesFromIds(mobileDeviceGroupIDsList)
	}
	if buildingIDsList := constructors.GetListFromSet(data, "building_ids"); len(buildingIDsList) > 0 {
		scope.Buildings = constructMobileDeviceScopeEntitiesFromIds(buildingIDsList)
	}
	if departmentIDsList := constructors.GetListFromSet(data, "department_ids"); len(departmentIDsList) > 0 {
		scope.Departments = constructMobileDeviceScopeEntitiesFromIds(departmentIDsList)
	}
	if jssUserIDsList := constructors.GetListFromSet(data, "jss_user_ids"); len(jssUserIDsList) > 0 {
		scope.JSSUsers = constructMobileDeviceScopeEntitiesFromIds(jssUserIDsList)
	}
	if jssUserGroupIDsList := constructors.GetListFromSet(data, "jss_user_group_ids"); len(jssUserGroupIDsList) > 0 {
		scope.JSSUserGroups = constructMobileDeviceScopeEntitiesFromIds(jssUserGroupIDsList)
	}

	// Handle Limitations Block (Outer TypeSet, Inner Map with TypeSet fields)
	if limitationsSet := constructors.GetListFromSet(data, "limitations"); len(limitationsSet) > 0 {
		// TypeSet with MaxItems: 1 means limitationsSet will have 0 or 1 element
		if limitationsSet[0] != nil {
			limitationData, mapOk := limitationsSet[0].(map[string]interface{})
			if mapOk {
				scope.Limitations = constructLimitations(limitationData)
			} else {
				log.Printf("[WARN] constructMobileDeviceConfigurationProfileSubsetScope: Could not cast limitations element to map.")
			}
		}
	}

	// Handle Exclusions Block (Outer TypeSet, Inner Map with TypeSet fields)
	if exclusionsSet := constructors.GetListFromSet(data, "exclusions"); len(exclusionsSet) > 0 {
		// TypeSet with MaxItems: 1 means exclusionsSet will have 0 or 1 element
		if exclusionsSet[0] != nil {
			exclusionData, mapOk := exclusionsSet[0].(map[string]interface{})
			if mapOk {
				scope.Exclusions = constructExclusions(exclusionData)
			} else {
				log.Printf("[WARN] constructMobileDeviceConfigurationProfileSubsetScope: Could not cast exclusions element to map.")
			}
		}
	}

	return scope
}

// constructLimitations constructs limitations using TypeSet from schema.
func constructLimitations(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetLimitation {
	limitations := jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{}

	// Use constructors.GetListFromSet for *Set fields inside the limitations map
	if userNamesList := constructors.GetListFromSet(data, "directory_service_or_local_usernames"); len(userNamesList) > 0 {
		limitations.Users = constructMobileDeviceScopeEntitiesFromNames(userNamesList)
	}
	if userGroupIDsList := constructors.GetListFromSet(data, "directory_service_usergroup_ids"); len(userGroupIDsList) > 0 {
		limitations.UserGroups = constructMobileDeviceScopeEntitiesFromIds(userGroupIDsList)
	}
	if networkSegmentIDsList := constructors.GetListFromSet(data, "network_segment_ids"); len(networkSegmentIDsList) > 0 {
		limitations.NetworkSegments = constructMobileDeviceNetworkSegments(networkSegmentIDsList)
	}
	if ibeaconIDsList := constructors.GetListFromSet(data, "ibeacon_ids"); len(ibeaconIDsList) > 0 {
		limitations.Ibeacons = constructMobileDeviceScopeEntitiesFromIds(ibeaconIDsList) // SDK uses 'Ibeacons' (lowercase b)
	}

	return limitations
}

// constructExclusions constructs exclusions using TypeSet from schema.
func constructExclusions(data map[string]interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetExclusion {
	exclusions := jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{}

	// Use constructors.GetListFromSet for *Set fields inside the exclusions map
	if mobileDeviceIDsList := constructors.GetListFromSet(data, "mobile_device_ids"); len(mobileDeviceIDsList) > 0 {
		exclusions.MobileDevices = constructMobileDevices(mobileDeviceIDsList)
	}
	if mobileDeviceGroupIDsList := constructors.GetListFromSet(data, "mobile_device_group_ids"); len(mobileDeviceGroupIDsList) > 0 {
		exclusions.MobileDeviceGroups = constructMobileDeviceScopeEntitiesFromIds(mobileDeviceGroupIDsList)
	}
	if userNamesList := constructors.GetListFromSet(data, "directory_service_or_local_usernames"); len(userNamesList) > 0 {
		exclusions.Users = constructMobileDeviceScopeEntitiesFromNames(userNamesList)
	}
	if userGroupIDsList := constructors.GetListFromSet(data, "directory_service_usergroup_ids"); len(userGroupIDsList) > 0 {
		exclusions.UserGroups = constructMobileDeviceScopeEntitiesFromIds(userGroupIDsList)
	}
	if buildingIDsList := constructors.GetListFromSet(data, "building_ids"); len(buildingIDsList) > 0 {
		exclusions.Buildings = constructMobileDeviceScopeEntitiesFromIds(buildingIDsList)
	}
	if departmentIDsList := constructors.GetListFromSet(data, "department_ids"); len(departmentIDsList) > 0 {
		exclusions.Departments = constructMobileDeviceScopeEntitiesFromIds(departmentIDsList)
	}
	if networkSegmentIDsList := constructors.GetListFromSet(data, "network_segment_ids"); len(networkSegmentIDsList) > 0 {
		exclusions.NetworkSegments = constructMobileDeviceNetworkSegments(networkSegmentIDsList)
	}
	if jssUserIDsList := constructors.GetListFromSet(data, "jss_user_ids"); len(jssUserIDsList) > 0 {
		exclusions.JSSUsers = constructMobileDeviceScopeEntitiesFromIds(jssUserIDsList)
	}
	if jssUserGroupIDsList := constructors.GetListFromSet(data, "jss_user_group_ids"); len(jssUserGroupIDsList) > 0 {
		exclusions.JSSUserGroups = constructMobileDeviceScopeEntitiesFromIds(jssUserGroupIDsList)
	}
	if ibeaconIDsList := constructors.GetListFromSet(data, "ibeacon_ids"); len(ibeaconIDsList) > 0 {
		exclusions.IBeacons = constructMobileDeviceScopeEntitiesFromIds(ibeaconIDsList) // SDK uses 'IBeacons' (uppercase B)
	}

	return exclusions
}

// --- Mobile Device Specific Helper Functions (Accepting []interface{}) ---
// These functions remain the same as they accept the output of GetListFromSet

// constructMobileDevices uses robust conversion from a list.
func constructMobileDevices(ids []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice {
	if ids == nil {
		return nil
	}
	mobileDevices := make([]jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "mobile device", i); ok {
			mobileDevices = append(mobileDevices, jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice{ID: intID})
		}
	}
	log.Printf("[DEBUG] constructMobileDevices: Input count %d, Output count %d", len(ids), len(mobileDevices))
	return mobileDevices
}

// constructMobileDeviceNetworkSegments uses robust conversion from a list.
func constructMobileDeviceNetworkSegments(ids []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment {
	if ids == nil {
		return nil
	}
	networkSegments := make([]jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "network segment", i); ok {
			networkSegments = append(networkSegments, jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment{
				MobileDeviceConfigurationProfileSubsetScopeEntity: jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{ID: intID},
			})
		}
	}
	log.Printf("[DEBUG] constructMobileDeviceNetworkSegments: Input count %d, Output count %d", len(ids), len(networkSegments))
	return networkSegments
}

// constructMobileDeviceScopeEntitiesFromIds uses robust conversion from a list for generic ID-based entities.
func constructMobileDeviceScopeEntitiesFromIds(ids []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity {
	if ids == nil {
		return nil
	}
	scopeEntities := make([]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, 0, len(ids))
	for i, idRaw := range ids {
		if intID, ok := constructors.ConvertToInt(idRaw, "scope entity (ID based)", i); ok {
			scopeEntities = append(scopeEntities, jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{ID: intID})
		}
	}
	log.Printf("[DEBUG] constructMobileDeviceScopeEntitiesFromIds: Input count %d, Output count %d", len(ids), len(scopeEntities))
	return scopeEntities
}

// constructMobileDeviceScopeEntitiesFromNames uses robust conversion from a list for generic name-based entities.
func constructMobileDeviceScopeEntitiesFromNames(names []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity {
	if names == nil {
		return nil
	}
	scopeEntities := make([]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, 0, len(names))
	for i, nameRaw := range names {
		if strName, ok := nameRaw.(string); ok && strName != "" {
			scopeEntities = append(scopeEntities, jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{Name: strName})
		} else if !ok {
			log.Printf("[WARN] constructMobileDeviceScopeEntitiesFromNames: Unexpected type %T for scope entity name: %v at index %d. Skipping.", nameRaw, nameRaw, i)
		} // Skip empty strings silently
	}
	log.Printf("[DEBUG] constructMobileDeviceScopeEntitiesFromNames: Input count %d, Output count %d", len(names), len(scopeEntities))
	return scopeEntities
}
