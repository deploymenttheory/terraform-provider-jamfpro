package mobiledeviceconfigurationprofiles

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProMobileDeviceConfigurationProfile constructs a ResourceMobileDeviceConfigurationProfile object from the provided schema data.
func constructJamfProMobileDeviceConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error) {
	profile := &jamfpro.ResourceMobileDeviceConfigurationProfile{
		General: jamfpro.MobileDeviceConfigurationProfileSubsetGeneral{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Level:            d.Get("level").(string),
			UUID:             d.Get("uuid").(string),
			DeploymentMethod: d.Get("deployment_method").(string),
			RedeployOnUpdate: d.Get("redeploy_on_update").(string),
			// Use html.EscapeString to escape the payloads content
			Payloads: html.EscapeString(d.Get("payloads").(string)),
		},
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		profile.General.Site = constructobject.ConstructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		profile.General.Site = constructobject.ConstructSharedResourceSite([]interface{}{})
	}

	// Handle Category
	if v, ok := d.GetOk("category"); ok {
		profile.General.Category = constructobject.ConstructSharedResourceCategory(v.([]interface{}))
	} else {
		// Set default values if 'category' data is not provided
		profile.General.Category = constructobject.ConstructSharedResourceCategory([]interface{}{})
	}

	// Handle Scope
	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		profile.Scope = constructMobileDeviceConfigurationProfileSubsetScope(scopeData)
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

// getNestedMap retrieves a nested map from the provided data.
func getNestedMap(data map[string]interface{}, key string) map[string]interface{} {
	if v, ok := data[key]; ok {
		if nestedMap, ok := v.(map[string]interface{}); ok {
			return nestedMap
		}
	}
	return map[string]interface{}{}
}

// getSlice retrieves a slice from the provided data.
func getSlice(data map[string]interface{}, key string) []interface{} {
	if v, ok := data[key]; ok {
		if slice, ok := v.([]interface{}); ok {
			return slice
		}
	}
	return []interface{}{}
}

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
