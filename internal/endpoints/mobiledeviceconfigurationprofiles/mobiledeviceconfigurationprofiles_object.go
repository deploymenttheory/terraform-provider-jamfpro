// mobiledeviceconfigurationprofiles_resource.go
package mobiledeviceconfigurationprofiles

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProMobileDeviceConfigurationProfile constructs a ResourceMobileDeviceConfigurationProfile object from the provided schema data.
func constructJamfProMobileDeviceConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
	profile := &jamfpro.ResourceMobileDeviceConfigurationProfile{
		General: jamfpro.MobileDeviceConfigurationProfileSubsetGeneral{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Level:            d.Get("level").(string),
			Site:             constructSharedResourceSite(d.Get("site").([]interface{})),
			Category:         constructSharedResourceCategory(d.Get("category").([]interface{})),
			UUID:             d.Get("uuid").(string),
			DeploymentMethod: d.Get("deployment_method").(string),
			RedeployOnUpdate: d.Get("redeploy_on_update").(string),
			// Use html.EscapeString to escape the payloads content
			Payloads: html.EscapeString(d.Get("payloads").(string)),
		},
		Scope: constructMobileDeviceConfigurationProfileSubsetScope(d.Get("scope").([]interface{})),
	}

	// Serialize and pretty-print the Mobile Device Configuration Profile object as XML for logging
	resourceXML, err := xml.MarshalIndent(profile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Configuration Profile '%s' to XML: %v", profile.General.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Configuration Profile XML:\n%s\n", string(resourceXML))

	return profile, nil
}

// constructMobileDeviceConfigurationProfileSubsetScope constructs a MobileDeviceConfigurationProfileSubsetScope object from the provided schema data, using the consolidated helper function for entity construction.
func constructMobileDeviceConfigurationProfileSubsetScope(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetScope {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetScope{}
	}
	scopeData := data[0].(map[string]interface{})

	// Use constructScopeEntities for all entities except those needing special handling like NetworkSegments
	return jamfpro.MobileDeviceConfigurationProfileSubsetScope{
		AllMobileDevices:   scopeData["all_mobile_devices"].(bool),
		AllJSSUsers:        scopeData["all_jss_users"].(bool),
		MobileDevices:      constructMobileDevices(scopeData["mobile_devices"].(*schema.Set).List()),
		MobileDeviceGroups: constructScopeEntities(scopeData["mobile_device_groups"].(*schema.Set).List()),
		Buildings:          constructScopeEntities(scopeData["buildings"].(*schema.Set).List()),
		Departments:        constructScopeEntities(scopeData["departments"].(*schema.Set).List()),
		JSSUsers:           constructScopeEntities(scopeData["jss_users"].(*schema.Set).List()),
		JSSUserGroups:      constructScopeEntities(scopeData["jss_user_groups"].(*schema.Set).List()),

		Limitations: constructLimitations(scopeData["limitations"].([]interface{})),
		Exclusions:  constructExclusions(scopeData["exclusions"].([]interface{})),
	}
}

// constructLimitations constructs a MobileDeviceConfigurationProfileSubsetLimitation object from the provided schema data.
func constructLimitations(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetLimitation {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{}
	}
	limitationData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{
		Users:           constructScopeEntities(limitationData["users"].(*schema.Set).List()),
		UserGroups:      constructScopeEntities(limitationData["user_groups"].(*schema.Set).List()),
		NetworkSegments: constructNetworkSegments(limitationData["network_segments"].(*schema.Set).List()),
		Ibeacons:        constructScopeEntities(limitationData["ibeacons"].(*schema.Set).List()),
	}
}

// constructExclusions constructs a MobileDeviceConfigurationProfileSubsetExclusion object from the provided schema data.
func constructExclusions(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetExclusion {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{}
	}
	exclusionData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{
		MobileDevices:      constructMobileDevices(exclusionData["mobile_devices"].(*schema.Set).List()),
		MobileDeviceGroups: constructScopeEntities(exclusionData["mobile_device_groups"].(*schema.Set).List()),
		Users:              constructScopeEntities(exclusionData["users"].(*schema.Set).List()),
		UserGroups:         constructScopeEntities(exclusionData["user_groups"].(*schema.Set).List()),
		Buildings:          constructScopeEntities(exclusionData["buildings"].(*schema.Set).List()),
		Departments:        constructScopeEntities(exclusionData["departments"].(*schema.Set).List()),
		NetworkSegments:    constructNetworkSegments(exclusionData["network_segments"].(*schema.Set).List()),
		IBeacons:           constructScopeEntities(exclusionData["ibeacons"].(*schema.Set).List()),
		JSSUsers:           constructScopeEntities(exclusionData["jss_users"].(*schema.Set).List()),
		JSSUserGroups:      constructScopeEntities(exclusionData["jss_user_groups"].(*schema.Set).List()),
	}
}

// Helper functions for nested structures

// constructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func constructSharedResourceSite(data []interface{}) jamfpro.SharedResourceSite {
	// Check if 'site' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		site := data[0].(map[string]interface{})

		// Return the 'site' object with data from the schema
		return jamfpro.SharedResourceSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	// Return default 'site' values if no data is provided or it is empty
	return jamfpro.SharedResourceSite{
		ID:   -1,     // Default ID
		Name: "None", // Default name
	}
}

// constructSharedResourceCategory constructs a SharedResourceCategory object from the provided schema data,
// setting default values if none are presented.
func constructSharedResourceCategory(data []interface{}) jamfpro.SharedResourceCategory {
	// Check if 'category' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		category := data[0].(map[string]interface{})

		// Return the 'category' object with data from the schema
		return jamfpro.SharedResourceCategory{
			ID:   category["id"].(int),
			Name: category["name"].(string),
		}
	}

	// Return default 'category' values if no data is provided or it is empty
	return jamfpro.SharedResourceCategory{
		ID:   -1,                     // Default ID
		Name: "No category assigned", // Default name
	}
}

// constructMobileDevices constructs a slice of MobileDeviceConfigurationProfileSubsetMobileDevice from the provided schema data.
func constructMobileDevices(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice {
	mobileDevices := make([]jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice, len(data))
	for i, item := range data {
		deviceData := item.(map[string]interface{})
		mobileDevices[i] = jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice{
			ID:             deviceData["id"].(int),
			Name:           deviceData["name"].(string),
			UDID:           deviceData["udid"].(string),
			WifiMacAddress: deviceData["wifi_mac_address"].(string),
		}
	}
	return mobileDevices
}

// constructScopeEntities constructs a slice of MobileDeviceConfigurationProfileSubsetEntity from the provided schema data.
func constructScopeEntities(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity {
	entities := make([]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity, len(data))
	for i, item := range data {
		entityData := item.(map[string]interface{})
		entities[i] = jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
			ID:   entityData["id"].(int),
			Name: entityData["name"].(string),
		}
	}
	return entities
}

// constructNetworkSegments constructs a slice of MobileDeviceConfigurationProfileSubsetNetworkSegment from the provided schema data.
func constructNetworkSegments(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment {
	networkSegments := make([]jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, len(data))
	for i, item := range data {
		segmentData := item.(map[string]interface{})
		networkSegments[i] = jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment{
			MobileDeviceConfigurationProfileSubsetScopeEntity: jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
				ID:   segmentData["id"].(int),
				Name: segmentData["name"].(string),
			},
		}
	}
	return networkSegments
}
