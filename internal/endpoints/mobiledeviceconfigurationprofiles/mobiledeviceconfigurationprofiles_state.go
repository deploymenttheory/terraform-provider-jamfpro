// mobiledeviceconfigurationprofiles_state.go
package mobiledeviceconfigurationprofiles

import (
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
		"name":                              resource.General.Name,
		"description":                       resource.General.Description,
		"level":                             resource.General.Level,
		"uuid":                              resource.General.UUID,
		"deployment_method":                 resource.General.DeploymentMethod,
		"redeploy_on_update":                resource.General.RedeployOnUpdate,
		"redeploy_days_before_cert_expires": resource.General.RedeployDaysBeforeCertExpires,
		// Skipping stating payloads and let terraform handle it directly
		// "payloads": html.UnescapeString(resource.General.Payloads),
	}

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

	// Create a map to hold the scope data
	scope := []map[string]interface{}{
		{
			"all_mobile_devices":   resource.Scope.AllMobileDevices,
			"all_jss_users":        resource.Scope.AllJSSUsers,
			"mobile_devices":       make([]interface{}, len(resource.Scope.MobileDevices)),
			"mobile_device_groups": make([]interface{}, len(resource.Scope.MobileDeviceGroups)),
			"jss_users":            make([]interface{}, len(resource.Scope.JSSUsers)),
			"jss_user_groups":      make([]interface{}, len(resource.Scope.JSSUserGroups)),
			"buildings":            make([]interface{}, len(resource.Scope.Buildings)),
			"departments":          make([]interface{}, len(resource.Scope.Departments)),
		},
	}

	// Fill the slices in the scope data
	for i, device := range resource.Scope.MobileDevices {
		scope[0]["mobile_devices"].([]interface{})[i] = device
	}

	for i, group := range resource.Scope.MobileDeviceGroups {
		scope[0]["mobile_device_groups"].([]interface{})[i] = group
	}

	for i, user := range resource.Scope.JSSUsers {
		scope[0]["jss_users"].([]interface{})[i] = user
	}

	for i, userGroup := range resource.Scope.JSSUserGroups {
		scope[0]["jss_user_groups"].([]interface{})[i] = userGroup
	}

	for i, building := range resource.Scope.Buildings {
		scope[0]["buildings"].([]interface{})[i] = building
	}

	for i, department := range resource.Scope.Departments {
		scope[0]["departments"].([]interface{})[i] = department
	}

	// Create a map to hold the limitations data
	// Initialize a slice for limitations
	limitations := make([]map[string]interface{}, 0)

	// Add network segments to limitations
	networkSegments := make([]map[string]interface{}, len(resource.Scope.Limitations.NetworkSegments))
	for i, segment := range resource.Scope.Limitations.NetworkSegments {
		networkSegments[i] = map[string]interface{}{
			"id":   segment.ID,
			"name": segment.Name,
		}
	}
	if len(networkSegments) > 0 {
		limitations = append(limitations, map[string]interface{}{"network_segments": networkSegments})
	}

	// Add users to limitations
	users := make([]map[string]interface{}, len(resource.Scope.Limitations.Users))
	for i, user := range resource.Scope.Limitations.Users {
		users[i] = map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		}
	}
	if len(users) > 0 {
		limitations = append(limitations, map[string]interface{}{"users": users})
	}

	// Add user groups to limitations
	userGroups := make([]map[string]interface{}, len(resource.Scope.Limitations.UserGroups))
	for i, group := range resource.Scope.Limitations.UserGroups {
		userGroups[i] = map[string]interface{}{
			"id":   group.ID,
			"name": group.Name,
		}
	}
	if len(userGroups) > 0 {
		limitations = append(limitations, map[string]interface{}{"user_groups": userGroups})
	}

	// Add iBeacons to limitations
	ibeacons := make([]map[string]interface{}, len(resource.Scope.Limitations.Ibeacons))
	for i, ibeacon := range resource.Scope.Limitations.Ibeacons {
		ibeacons[i] = map[string]interface{}{
			"id":   ibeacon.ID,
			"name": ibeacon.Name,
		}
	}
	if len(ibeacons) > 0 {
		limitations = append(limitations, map[string]interface{}{"ibeacons": ibeacons})
	}

	// Assign limitations to the scope
	if len(limitations) > 0 {
		scope[0]["limitations"] = limitations
	}

	// Initialize a slice for exclusions
	exclusions := make([]map[string]interface{}, 0)

	// Add exclusions to the slice
	exclusionsData := map[string]interface{}{
		"mobile_devices":       make([]interface{}, len(resource.Scope.Exclusions.MobileDevices)),
		"mobile_device_groups": make([]interface{}, len(resource.Scope.Exclusions.MobileDeviceGroups)),
		"users":                make([]interface{}, len(resource.Scope.Exclusions.Users)),
		"user_groups":          make([]interface{}, len(resource.Scope.Exclusions.UserGroups)),
		"buildings":            make([]interface{}, len(resource.Scope.Exclusions.Buildings)),
		"departments":          make([]interface{}, len(resource.Scope.Exclusions.Departments)),
		"network_segments":     make([]interface{}, len(resource.Scope.Exclusions.NetworkSegments)),
		"jss_users":            make([]interface{}, len(resource.Scope.Exclusions.JSSUsers)),
		"jss_user_groups":      make([]interface{}, len(resource.Scope.Exclusions.JSSUserGroups)),
		"ibeacons":             make([]interface{}, len(resource.Scope.Exclusions.IBeacons)),
	}

	// Populate the exclusions data
	for i, device := range resource.Scope.Exclusions.MobileDevices {
		exclusionsData["mobile_devices"].([]interface{})[i] = map[string]interface{}{
			"id":               device.ID,
			"name":             device.Name,
			"udid":             device.UDID,
			"wifi_mac_address": device.WifiMacAddress,
		}
	}

	for i, group := range resource.Scope.Exclusions.MobileDeviceGroups {
		exclusionsData["mobile_device_groups"].([]interface{})[i] = map[string]interface{}{
			"id":   group.ID,
			"name": group.Name,
		}
	}

	for i, user := range resource.Scope.Exclusions.Users {
		exclusionsData["users"].([]interface{})[i] = map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		}
	}

	for i, group := range resource.Scope.Exclusions.UserGroups {
		exclusionsData["user_groups"].([]interface{})[i] = map[string]interface{}{
			"id":   group.ID,
			"name": group.Name,
		}
	}

	for i, building := range resource.Scope.Exclusions.Buildings {
		exclusionsData["buildings"].([]interface{})[i] = map[string]interface{}{
			"id":   building.ID,
			"name": building.Name,
		}
	}

	for i, department := range resource.Scope.Exclusions.Departments {
		exclusionsData["departments"].([]interface{})[i] = map[string]interface{}{
			"id":   department.ID,
			"name": department.Name,
		}
	}

	for i, segment := range resource.Scope.Exclusions.NetworkSegments {
		exclusionsData["network_segments"].([]interface{})[i] = map[string]interface{}{
			"id":   segment.ID,
			"name": segment.Name,
			"uid":  segment.UID,
		}
	}

	for i, jssUser := range resource.Scope.Exclusions.JSSUsers {
		exclusionsData["jss_users"].([]interface{})[i] = map[string]interface{}{
			"id":   jssUser.ID,
			"name": jssUser.Name,
		}
	}

	for i, jssUserGroup := range resource.Scope.Exclusions.JSSUserGroups {
		exclusionsData["jss_user_groups"].([]interface{})[i] = map[string]interface{}{
			"id":   jssUserGroup.ID,
			"name": jssUserGroup.Name,
		}
	}

	for i, ibeacon := range resource.Scope.Exclusions.IBeacons {
		exclusionsData["ibeacons"].([]interface{})[i] = map[string]interface{}{
			"id":   ibeacon.ID,
			"name": ibeacon.Name,
		}
	}

	// Add exclusionsData to exclusions
	if len(exclusionsData) > 0 {
		exclusions = append(exclusions, exclusionsData)
	}

	// Assign exclusions to the scope
	if len(exclusions) > 0 {
		scope[0]["exclusions"] = exclusions
	}

	// Add scope data to the resource data
	resourceData["scope"] = scope

	// Set resource data into the Terraform schema
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
