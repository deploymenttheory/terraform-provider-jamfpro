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

// constructMobileDeviceConfigurationProfileSubsetScope constructs a MobileDeviceConfigurationProfileSubsetScope object from the provided schema data.
func constructMobileDeviceConfigurationProfileSubsetScope(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetScope {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetScope{}
	}
	scopeData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetScope{
		AllMobileDevices:   scopeData["all_mobile_devices"].(bool),
		AllJSSUsers:        scopeData["all_jss_users"].(bool),
		MobileDevices:      constructMobileDevices(scopeData["mobile_devices"].(*schema.Set).List()),
		MobileDeviceGroups: constructMobileDeviceGroups(scopeData["mobile_device_groups"].(*schema.Set).List()),
		Buildings:          constructBuildings(scopeData["buildings"].(*schema.Set).List()),
		Departments:        constructDepartments(scopeData["departments"].(*schema.Set).List()),
		JSSUsers:           constructJSSUsers(scopeData["jss_users"].(*schema.Set).List()),
		JSSUserGroups:      constructJSSUserGroups(scopeData["jss_user_groups"].(*schema.Set).List()),
		Limitations:        constructLimitations(scopeData["limitations"].([]interface{})),
		Exclusions:         constructExclusions(scopeData["exclusions"].([]interface{})),
	}
}

// constructMobileDeviceConfigurationProfileSubsetSelfService constructs a MobileDeviceConfigurationProfileSubsetSelfService object from the provided schema data.
func constructMobileDeviceConfigurationProfileSubsetSelfService(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetSelfService {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetSelfService{}
	}
	selfServiceData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetSelfService{
		SelfServiceDescription: selfServiceData["self_service_description"].(string),
		// Further extraction for fields like SecurityName, SelfServiceIcon, FeatureOnMainPage, etc.
	}
}

// Helper functions for nested structures

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

func constructMobileDeviceGroups(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDeviceGroup {
	groups := make([]jamfpro.MobileDeviceConfigurationProfileSubsetMobileDeviceGroup, len(data))
	for i, item := range data {
		groupData := item.(map[string]interface{})
		groups[i] = jamfpro.MobileDeviceConfigurationProfileSubsetMobileDeviceGroup{
			ID:   groupData["id"].(int),
			Name: groupData["name"].(string),
		}
	}
	return groups
}

func constructBuildings(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetBuilding {
	buildings := make([]jamfpro.MobileDeviceConfigurationProfileSubsetBuilding, len(data))
	for i, item := range data {
		buildingData := item.(map[string]interface{})
		buildings[i] = jamfpro.MobileDeviceConfigurationProfileSubsetBuilding{
			ID:   buildingData["id"].(int),
			Name: buildingData["name"].(string),
		}
	}
	return buildings
}

func constructDepartments(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetDepartment {
	departments := make([]jamfpro.MobileDeviceConfigurationProfileSubsetDepartment, len(data))
	for i, item := range data {
		departmentData := item.(map[string]interface{})
		departments[i] = jamfpro.MobileDeviceConfigurationProfileSubsetDepartment{
			ID:   departmentData["id"].(int),
			Name: departmentData["name"].(string),
		}
	}
	return departments
}

func constructJSSUsers(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetJSSUser {
	jssUsers := make([]jamfpro.MobileDeviceConfigurationProfileSubsetJSSUser, len(data))
	for i, item := range data {
		userData := item.(map[string]interface{})
		jssUsers[i] = jamfpro.MobileDeviceConfigurationProfileSubsetJSSUser{
			User: jamfpro.MobileDeviceConfigurationProfileSubsetUser{
				ID:   userData["id"].(int),
				Name: userData["name"].(string),
			},
		}
	}
	return jssUsers
}

func constructJSSUserGroups(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetJSSUserGroup {
	jssUserGroups := make([]jamfpro.MobileDeviceConfigurationProfileSubsetJSSUserGroup, len(data))
	for i, item := range data {
		groupData := item.(map[string]interface{})
		jssUserGroups[i] = jamfpro.MobileDeviceConfigurationProfileSubsetJSSUserGroup{
			UserGroup: jamfpro.MobileDeviceConfigurationProfileSubsetUserGroup{
				ID:   groupData["id"].(int),
				Name: groupData["name"].(string),
			},
		}
	}
	return jssUserGroups
}

func constructLimitations(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetLimitation {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{}
	}
	limitationData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetLimitation{
		Users:           constructUsers(limitationData["users"].(*schema.Set).List()),
		UserGroups:      constructUserGroups(limitationData["user_groups"].(*schema.Set).List()),
		NetworkSegments: constructNetworkSegments(limitationData["network_segments"].(*schema.Set).List()),
		IBeacons:        constructIBeacons(limitationData["ibeacons"].(*schema.Set).List()),
	}
}

func constructUsers(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetUser {
	users := make([]jamfpro.MobileDeviceConfigurationProfileSubsetUser, len(data))
	for i, item := range data {
		userData := item.(map[string]interface{})
		users[i] = jamfpro.MobileDeviceConfigurationProfileSubsetUser{
			ID:   userData["id"].(int),
			Name: userData["name"].(string),
		}
	}
	return users
}

func constructUserGroups(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetUserGroup {
	userGroups := make([]jamfpro.MobileDeviceConfigurationProfileSubsetUserGroup, len(data))
	for i, item := range data {
		groupData := item.(map[string]interface{})
		userGroups[i] = jamfpro.MobileDeviceConfigurationProfileSubsetUserGroup{
			ID:   groupData["id"].(int),
			Name: groupData["name"].(string),
		}
	}
	return userGroups
}

func constructNetworkSegments(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment {
	networkSegments := make([]jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment, len(data))
	for i, item := range data {
		segmentData := item.(map[string]interface{})
		networkSegments[i] = jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment{
			ID:   segmentData["id"].(int),
			UID:  segmentData["uid"].(string),
			Name: segmentData["name"].(string),
		}
	}
	return networkSegments
}

func constructIBeacons(data []interface{}) []jamfpro.MobileDeviceConfigurationProfileSubsetIbeacon {
	ibeacons := make([]jamfpro.MobileDeviceConfigurationProfileSubsetIbeacon, len(data))
	for i, item := range data {
		beaconData := item.(map[string]interface{})
		ibeacons[i] = jamfpro.MobileDeviceConfigurationProfileSubsetIbeacon{
			ID:   beaconData["id"].(int),
			Name: beaconData["name"].(string),
		}
	}
	return ibeacons
}

func constructExclusions(data []interface{}) jamfpro.MobileDeviceConfigurationProfileSubsetExclusion {
	if len(data) == 0 {
		return jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{}
	}
	exclusionData := data[0].(map[string]interface{})

	return jamfpro.MobileDeviceConfigurationProfileSubsetExclusion{
		MobileDevices:      constructMobileDevices(exclusionData["mobile_devices"].(*schema.Set).List()),
		MobileDeviceGroups: constructMobileDeviceGroups(exclusionData["mobile_device_groups"].(*schema.Set).List()),
		Users:              constructUsers(exclusionData["users"].(*schema.Set).List()),
		UserGroups:         constructUserGroups(exclusionData["user_groups"].(*schema.Set).List()),
		Buildings:          constructBuildings(exclusionData["buildings"].(*schema.Set).List()),
		Departments:        constructDepartments(exclusionData["departments"].(*schema.Set).List()),
		NetworkSegments:    constructNetworkSegments(exclusionData["network_segments"].(*schema.Set).List()),
		JSSUsers:           constructJSSUsers(exclusionData["jss_users"].(*schema.Set).List()),
		JSSUserGroups:      constructJSSUserGroups(exclusionData["jss_user_groups"].(*schema.Set).List()),
		IBeacons:           constructIBeacons(exclusionData["ibeacons"].(*schema.Set).List()),
	}
}
