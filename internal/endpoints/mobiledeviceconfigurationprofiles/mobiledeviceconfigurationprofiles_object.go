// mobiledeviceconfigurationprofiles_resource.go
package mobiledeviceconfigurationprofiles

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProMobileDeviceConfigurationProfile constructs a ResourceMobileDeviceConfigurationProfile object from the provided schema data.
func constructJamfProMobileDeviceConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
	profile := &jamfpro.ResourceMobileDeviceConfigurationProfile{
		General: jamfpro.MobileDeviceConfigurationProfileSubsetGeneral{
			ID:               d.Get("id").(int),
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Site:             constructSharedResourceSite(d.Get("site").([]interface{})),
			Category:         constructSharedResourceCategory(d.Get("category").([]interface{})),
			UUID:             d.Get("uuid").(string),
			DeploymentMethod: d.Get("deployment_method").(string),
			RedeployOnUpdate: d.Get("redeploy_on_update").(string),
			Payloads:         d.Get("payloads").(string),
		},
		Scope:       constructMobileDeviceConfigurationProfileSubsetScope(d.Get("scope").([]interface{})),
		SelfService: constructMobileDeviceConfigurationProfileSubsetSelfService(d.Get("self_service").([]interface{})),
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

// constructMobileDeviceConfigurationProfileSubsetSelfService constructs a MobileDeviceConfigurationProfileSubsetSelfService object from the provided schema data.
func constructMobileDeviceConfigurationProfileSubsetSelfService(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetSelfService {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetSelfService{}
	}
	selfServiceData := data[0].(map[string]interface{})

	// Construct SecurityName
	var securityName jamfpro.MobileDeviceConfigurationProfileSubsetSelfServiceSecurityName
	if securityNameData, ok := selfServiceData["security_name"].(map[string]interface{}); ok {
		securityName.RemovalDisallowed = securityNameData["removal_disallowed"].(string)
	}

	// Construct SelfServiceIcon
	var selfServiceIcon jamfpro.SharedResourceSelfServiceIcon
	if iconData, ok := selfServiceData["self_service_icon"].(map[string]interface{}); ok {
		selfServiceIcon.ID = iconData["id"].(int)
		selfServiceIcon.URI = iconData["uri"].(string)
		selfServiceIcon.Data = iconData["data"].(string)
		selfServiceIcon.Filename = iconData["filename"].(string)
	}

	// Construct SelfServiceCategories
	var selfServiceCategories jamfpro.SharedResourceSelfServiceCategories
	if categoriesData, ok := selfServiceData["self_service_categories"].([]interface{}); ok {
		for _, category := range categoriesData {
			catData := category.(map[string]interface{})
			categoryStruct := jamfpro.SharedResourceSelfServiceCategory{
				ID:       catData["id"].(int),
				Name:     catData["name"].(string),
				Priority: catData["priority"].(int),
			}
			selfServiceCategories.Category = append(selfServiceCategories.Category, categoryStruct)
		}
	}

	return jamfpro.MobileDeviceConfigurationProfileSubsetSelfService{
		SelfServiceDescription: selfServiceData["self_service_description"].(string),
		SecurityName:           securityName,
		SelfServiceIcon:        selfServiceIcon,
		FeatureOnMainPage:      selfServiceData["feature_on_main_page"].(bool),
		SelfServiceCategories:  []jamfpro.SharedResourceSelfServiceCategories{selfServiceCategories}, // Wrap in a slice since the struct field expects a slice
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

// constructSharedResourceSite constructs a SharedResourceSite object from the provided schema data.
func constructSharedResourceSite(data []interface{}) jamfpro.SharedResourceSite {
	if len(data) == 0 {
		return jamfpro.SharedResourceSite{}
	}
	site := data[0].(map[string]interface{})
	return jamfpro.SharedResourceSite{
		ID:   site["id"].(int),
		Name: site["name"].(string),
	}
}

// constructSharedResourceCategory constructs a SharedResourceCategory object from the provided schema data.
func constructSharedResourceCategory(data []interface{}) jamfpro.SharedResourceCategory {
	if len(data) == 0 {
		return jamfpro.SharedResourceCategory{}
	}
	category := data[0].(map[string]interface{})
	return jamfpro.SharedResourceCategory{
		ID:   category["id"].(int),
		Name: category["name"].(string),
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
			UID: segmentData["uid"].(string),
		}
	}
	return networkSegments
}
