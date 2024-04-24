// mobiledeviceconfigurationprofiles_state.go
package mobiledeviceconfigurationprofiles

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest ResourceMobileDeviceConfigurationProfile
// information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceMobileDeviceConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a map to hold the resource data
	resourceData := map[string]interface{}{
		"name":              resource.General.Name,
		"description":       resource.General.Description,
		"uuid":              resource.General.UUID,
		"deployment_method": resource.General.DeploymentMethod,

		// Skipping the 'distribution_method' attribute as appears to be deprecated but still in documenation
		//"redeploy_on_update":                resource.General.RedeployOnUpdate,
		//"redeploy_days_before_cert_expires": resource.General.RedeployDaysBeforeCertExpires,

		// Skipping stating payloads and let terraform handle it directly
		// "payloads": html.UnescapeString(resource.General.Payloads),
	}

	// Check if the level is "System" and set it to "Device Level", otherwise use the value from resource
	levelValue := resource.General.Level
	if levelValue == "System" {
		levelValue = "Device Level"
	}
	resourceData["level"] = levelValue

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	site := []interface{}{}
	if resource.General.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.General.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set the 'category' attribute in the state only if it's not empty (i.e., not default values)
	category := []interface{}{}
	if resource.General.Category.ID != -1 {
		category = append(category, map[string]interface{}{
			"id": resource.General.Category.ID,
		})
	}
	if len(category) > 0 {
		if err := d.Set("category", category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Preparing and setting scope data
	if scopeData, err := prepareScopeData(resource); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}

// prepareScopeData prepares the scope data for the Terraform state.
func prepareScopeData(resource *jamfpro.ResourceMobileDeviceConfigurationProfile) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_mobile_devices": resource.Scope.AllMobileDevices,
		"all_jss_users":      resource.Scope.AllJSSUsers,
	}

	// Gather mobile devices, groups, etc.
	mobileDevices, err := setMobileDevices(resource.Scope.MobileDevices)
	if err != nil {
		return nil, err
	}
	scopeData["mobile_devices"] = mobileDevices

	mobileDeviceGroups, err := setScopeEntities(resource.Scope.MobileDeviceGroups)
	if err != nil {
		return nil, err
	}
	scopeData["mobile_device_groups"] = mobileDeviceGroups

	jssUsers, err := setScopeEntities(resource.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}
	scopeData["jss_users"] = jssUsers

	jssUserGroups, err := setScopeEntities(resource.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}
	scopeData["jss_user_groups"] = jssUserGroups

	buildings, err := setScopeEntities(resource.Scope.Buildings)
	if err != nil {
		return nil, err
	}
	scopeData["buildings"] = buildings

	departments, err := setScopeEntities(resource.Scope.Departments)
	if err != nil {
		return nil, err
	}
	scopeData["departments"] = departments

	// Gather limitations
	limitationsData, err := setLimitations(resource.Scope.Limitations)
	if err != nil {
		return nil, err
	}
	scopeData["limitations"] = limitationsData

	// Gather exclusions
	exclusionsData, err := setExclusions(resource.Scope.Exclusions)
	if err != nil {
		return nil, err
	}
	scopeData["exclusions"] = exclusionsData

	return scopeData, nil
}

// setLimitations collects and formats limitations data for the Terraform state.
func setLimitations(limitations jamfpro.MobileDeviceConfigurationProfileSubsetLimitation) ([]interface{}, error) {
	result := make(map[string]interface{})

	// Iterate through each type of exclusion and set them
	limitationTypes := map[string][]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
		"users":       limitations.Users,
		"user_groups": limitations.UserGroups,
		"ibeacons":    limitations.Ibeacons,
	}

	for key, entities := range limitationTypes {
		if len(entities) > 0 {
			entityData, err := setScopeEntities(entities)
			if err != nil {
				return nil, fmt.Errorf("error setting %s: %v", key, err)
			}
			result[key] = entityData
		}
	}

	// Handle Network Segments specifically if needed
	if len(limitations.NetworkSegments) > 0 {
		networkSegments, err := setNetworkSegments(limitations.NetworkSegments)
		if err != nil {
			return nil, fmt.Errorf("error setting network segments: %v", err)
		}
		result["network_segments"] = networkSegments
	}

	// Ensure to wrap the map into a list to match the expected TypeList structure in Terraform
	return []interface{}{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MobileDeviceConfigurationProfileSubsetExclusion) ([]interface{}, error) {
	result := make(map[string]interface{})

	// Iterate through each type of exclusion and set them
	exclusionTypes := map[string][]jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity{
		"mobile_device_groups": exclusions.MobileDeviceGroups,
		"users":                exclusions.Users,
		"user_groups":          exclusions.UserGroups,
		"buildings":            exclusions.Buildings,
		"departments":          exclusions.Departments,
		"jss_users":            exclusions.JSSUsers,
		"jss_user_groups":      exclusions.JSSUserGroups,
		"ibeacons":             exclusions.IBeacons,
	}

	// This loop will ensure each exclusion type is seted and added to the result map correctly
	for key, entities := range exclusionTypes {
		if len(entities) > 0 {
			entitiesData, err := setScopeEntities(entities)
			if err != nil {
				return nil, fmt.Errorf("error setting %s for exclusions: %v", key, err)
			}
			result[key] = entitiesData
		}
	}

	// Handle Mobile Devices specifically if needed
	if len(exclusions.MobileDevices) > 0 {
		mobileDevices, err := setMobileDevices(exclusions.MobileDevices)
		if err != nil {
			return nil, fmt.Errorf("error setting mobile devices for exclusions: %v", err)
		}
		result["mobile_devices"] = mobileDevices
	}

	// Handle Network Segments specifically if needed
	if len(exclusions.NetworkSegments) > 0 {
		networkSegments, err := setNetworkSegments(exclusions.NetworkSegments)
		if err != nil {
			return nil, fmt.Errorf("error setting network segments for exclusions: %v", err)
		}
		result["network_segments"] = networkSegments
	}

	// Wrap the map in a slice to match the TypeList expectation
	return []interface{}{result}, nil
}

// setScopeEntities converts a slice of general scope entities (like user groups, buildings) to a format suitable for Terraform state.
func setScopeEntities(entities []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) ([]interface{}, error) {
	var entityList []interface{}
	for _, entity := range entities {
		entityMap := map[string]interface{}{
			"id":   entity.ID,
			"name": entity.Name,
		}
		entityList = append(entityList, entityMap)
	}
	return entityList, nil
}

// setMobileDevices converts a slice of MobileDevice entities to a format suitable for Terraform state.
func setMobileDevices(devices []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice) ([]interface{}, error) {
	var deviceList []interface{}
	for _, device := range devices {
		deviceMap := map[string]interface{}{
			"id":               device.ID,
			"name":             device.Name,
			"udid":             device.UDID,
			"wifi_mac_address": device.WifiMacAddress,
		}
		deviceList = append(deviceList, deviceMap)
	}
	return deviceList, nil
}

// Helper function specific to network segments if they have an additional field such as 'uid'.
func setNetworkSegments(segments []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment) ([]interface{}, error) {
	var segmentList []interface{}
	for _, segment := range segments {
		segmentMap := map[string]interface{}{
			"id":   segment.ID,
			"name": segment.Name,
		}
		segmentList = append(segmentList, segmentMap)
	}
	return segmentList, nil
}
