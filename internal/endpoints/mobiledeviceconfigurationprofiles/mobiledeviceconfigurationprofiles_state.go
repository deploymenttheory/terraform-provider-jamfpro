// mobiledeviceconfigurationprofiles_state.go
package mobiledeviceconfigurationprofiles

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest ResourceMobileDeviceConfigurationProfile
// information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceMobileDeviceConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the general section
	generalData := map[string]interface{}{
		"id":                                resource.General.ID,
		"name":                              resource.General.Name,
		"description":                       resource.General.Description,
		"site":                              map[string]interface{}{"id": resource.General.Site.ID, "name": resource.General.Site.Name},
		"category":                          map[string]interface{}{"id": resource.General.Category.ID, "name": resource.General.Category.Name},
		"uuid":                              resource.General.UUID,
		"deployment_method":                 resource.General.DeploymentMethod,
		"redeploy_on_update":                resource.General.RedeployOnUpdate,
		"redeploy_days_before_cert_expires": resource.General.RedeployDaysBeforeCertExpires,
		"payloads":                          resource.General.Payloads,
	}
	if err := d.Set("general", []interface{}{generalData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update scope information
	scopeData := map[string]interface{}{
		"all_mobile_devices":   resource.Scope.AllMobileDevices,
		"all_jss_users":        resource.Scope.AllJSSUsers,
		"mobile_devices":       common.ConvertToInterfaceSlice(resource.Scope.MobileDevices),
		"buildings":            common.ConvertToInterfaceSlice(resource.Scope.Buildings),
		"departments":          common.ConvertToInterfaceSlice(resource.Scope.Departments),
		"mobile_device_groups": common.ConvertToInterfaceSlice(resource.Scope.MobileDeviceGroups),
		"jss_users":            common.ConvertToInterfaceSlice(resource.Scope.JSSUsers),
		"jss_user_groups":      common.ConvertToInterfaceSlice(resource.Scope.JSSUserGroups),
		"limitations":          []interface{}{convertLimitationsToInterface(resource.Scope.Limitations)},
		"exclusions":           []interface{}{convertExclusionsToInterface(resource.Scope.Exclusions)},
	}
	if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update self-service information
	selfServiceData := map[string]interface{}{
		"self_service_description": resource.SelfService.SelfServiceDescription,
		"security_name":            map[string]interface{}{"removal_disallowed": resource.SelfService.SecurityName.RemovalDisallowed},
		"self_service_icon":        map[string]interface{}{"filename": resource.SelfService.SelfServiceIcon.Filename, "uri": resource.SelfService.SelfServiceIcon.URI},
		"feature_on_main_page":     resource.SelfService.FeatureOnMainPage,
		"self_service_categories":  common.ConvertToInterfaceSlice(resource.SelfService.SelfServiceCategories),
	}
	if err := d.Set("self_service", []interface{}{selfServiceData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// convertLimitationsToInterface takes a MobileDeviceConfigurationProfileSubsetLimitation object and converts it to a map[string]interface{} for Terraform state.
func convertLimitationsToInterface(limitations jamfpro.MobileDeviceConfigurationProfileSubsetLimitation) map[string]interface{} {
	return map[string]interface{}{
		"users":            common.ConvertToInterfaceSlice(limitations.Users),
		"user_groups":      common.ConvertToInterfaceSlice(limitations.UserGroups),
		"network_segments": common.ConvertToInterfaceSlice(limitations.NetworkSegments),
		"ibeacons":         common.ConvertToInterfaceSlice(limitations.Ibeacons),
	}
}

// convertExclusionsToInterface takes a MobileDeviceConfigurationProfileSubsetExclusion object and converts it to a map[string]interface{} for Terraform state.
func convertExclusionsToInterface(exclusions jamfpro.MobileDeviceConfigurationProfileSubsetExclusion) map[string]interface{} {
	return map[string]interface{}{
		"mobile_devices":       common.ConvertToInterfaceSlice(exclusions.MobileDevices),
		"mobile_device_groups": common.ConvertToInterfaceSlice(exclusions.MobileDeviceGroups),
		"users":                common.ConvertToInterfaceSlice(exclusions.Users),
		"user_groups":          common.ConvertToInterfaceSlice(exclusions.UserGroups),
		"buildings":            common.ConvertToInterfaceSlice(exclusions.Buildings),
		"departments":          common.ConvertToInterfaceSlice(exclusions.Departments),
		"network_segments":     common.ConvertToInterfaceSlice(exclusions.NetworkSegments),
		"ibeacons":             common.ConvertToInterfaceSlice(exclusions.IBeacons),
		"jss_users":            common.ConvertToInterfaceSlice(exclusions.JSSUsers),
		"jss_user_groups":      common.ConvertToInterfaceSlice(exclusions.JSSUserGroups),
	}
}
