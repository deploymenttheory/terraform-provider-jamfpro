package mobile_device_configuration_profile_plist

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
func constructJamfProMobileDeviceConfigurationProfilePlist(d *schema.ResourceData, mode string, meta any) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
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

	if _, ok := d.GetOk("scope"); ok {
		// Pass the ResourceData object directly instead of extracting the scope data
		resource.Scope = constructMobileDeviceConfigurationProfileSubsetScope(d)
	} else {
		log.Printf("[DEBUG] constructJamfProMobileDeviceConfigurationProfilePlist: No scope block found or it's empty.")
	}

	// Handle Payloads based on mode
	if mode != "update" {
		resource.General.Payloads = html.EscapeString(d.Get("payloads").(string))
	} else if mode == "update" {
		var existingPlist map[string]any
		var newPlist map[string]any

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
		identifierMap := make(map[string]string)
		helpers.ExtractUUIDs(existingPlist, uuidMap, true)
		helpers.ExtractPayloadIdentifiers(existingPlist, identifierMap, true)
		helpers.UpdateUUIDs(newPlist, uuidMap, identifierMap, true)

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
func constructMobileDeviceConfigurationProfileSubsetScope(d *schema.ResourceData) jamfpro.MobileDeviceConfigurationProfileSubsetScope {
	scope := jamfpro.MobileDeviceConfigurationProfileSubsetScope{}
	// Get the scope data from the resource data
	scopeData := d.Get("scope").([]any)[0].(map[string]any)

	scope.AllMobileDevices = scopeData["all_mobile_devices"].(bool)
	scope.AllJSSUsers = scopeData["all_jss_users"].(bool)

	// Use MapSetToStructs for mobile devices
	var mobileDevices []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice, int](
		"scope.0.mobile_device_ids", "ID", d, &mobileDevices); err == nil {
		scope.MobileDevices = mobileDevices
	}

	// Use MapSetToStructs for mobile device groups
	var mobileDeviceGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.mobile_device_group_ids", "ID", d, &mobileDeviceGroups); err == nil {
		scope.MobileDeviceGroups = mobileDeviceGroups
	}

	// Buildings
	var buildings []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.building_ids", "ID", d, &buildings); err == nil {
		scope.Buildings = buildings
	}

	// Departments
	var departments []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.department_ids", "ID", d, &departments); err == nil {
		scope.Departments = departments
	}

	// JSS Users
	var jssUsers []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		scope.JSSUsers = jssUsers
	}

	// JSS User Groups
	var jssUserGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		scope.JSSUserGroups = jssUserGroups
	}

	// Handle Limitations and Exclusions which may need special handling
	if _, ok := d.GetOk("scope.0.limitations"); ok {
		scope.Limitations = constructLimitations(d)
	}

	if _, ok := d.GetOk("scope.0.exclusions"); ok {
		scope.Exclusions = constructExclusions(d)
	}

	return scope
}

// Refactor constructLimitations to use MapSetToStructs
func constructLimitations(d *schema.ResourceData) jamfpro.MobileDeviceConfigurationProfileSubsetLimitation {
	limitations := jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{}

	// User names (strings)
	var users []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, string](
		"scope.0.limitations.0.directory_service_or_local_usernames", "Name", d, &users); err == nil {
		limitations.Users = users
	}

	// User groups
	var userGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, &userGroups); err == nil {
		limitations.UserGroups = userGroups
	}

	// Network segments
	var networkSegments []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, int](
		"scope.0.limitations.0.network_segment_ids", "ID", d, &networkSegments); err == nil {
		limitations.NetworkSegments = networkSegments
	}

	// iBeacons
	var ibeacons []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.limitations.0.ibeacon_ids", "ID", d, &ibeacons); err == nil {
		limitations.Ibeacons = ibeacons
	}

	return limitations
}

func constructExclusions(d *schema.ResourceData) jamfpro.MobileDeviceConfigurationProfileSubsetExclusion {
	exclusions := jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{}

	// Mobile devices
	var mobileDevices []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice, int](
		"scope.0.exclusions.0.mobile_device_ids", "ID", d, &mobileDevices); err == nil {
		exclusions.MobileDevices = mobileDevices
	}

	// Mobile device groups
	var mobileDeviceGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.mobile_device_group_ids", "ID", d, &mobileDeviceGroups); err == nil {
		exclusions.MobileDeviceGroups = mobileDeviceGroups
	}

	// User names (strings)
	var users []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, string](
		"scope.0.exclusions.0.directory_service_or_local_usernames", "Name", d, &users); err == nil {
		exclusions.Users = users
	}

	// User groups
	var userGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.directory_service_usergroup_ids", "ID", d, &userGroups); err == nil {
		exclusions.UserGroups = userGroups
	}

	// Buildings
	var buildings []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.building_ids", "ID", d, &buildings); err == nil {
		exclusions.Buildings = buildings
	}

	// Departments
	var departments []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.department_ids", "ID", d, &departments); err == nil {
		exclusions.Departments = departments
	}

	// Network segments
	var networkSegments []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, int](
		"scope.0.exclusions.0.network_segment_ids", "ID", d, &networkSegments); err == nil {
		exclusions.NetworkSegments = networkSegments
	}

	// JSS Users
	var jssUsers []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		exclusions.JSSUsers = jssUsers
	}

	// JSS User Groups
	var jssUserGroups []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		exclusions.JSSUserGroups = jssUserGroups
	}

	// IBeacons
	var iBeacons []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, int](
		"scope.0.exclusions.0.ibeacon_ids", "ID", d, &iBeacons); err == nil {
		exclusions.IBeacons = iBeacons
	}

	return exclusions
}
