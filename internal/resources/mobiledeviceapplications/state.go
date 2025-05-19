package mobiledeviceapplications

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceMobileDeviceApplication
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDeviceApplication) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("name", resp.General.Name)
	d.Set("display_name", resp.General.DisplayName)
	d.Set("description", normalizeWhitespace(resp.General.Description))
	d.Set("bundle_id", resp.General.BundleID)
	d.Set("version", resp.General.Version)
	d.Set("internal_app", resp.General.InternalApp)
	d.Set("itunes_store_url", resp.General.ITunesStoreURL)
	d.Set("make_available_after_install", resp.General.MakeAvailableAfterInstall)
	d.Set("itunes_country_region", resp.General.ITunesCountryRegion)
	d.Set("itunes_sync_time", resp.General.ITunesSyncTime)
	d.Set("deployment_type", resp.General.DeploymentType)
	d.Set("deploy_automatically", resp.General.DeployAutomatically)
	d.Set("deploy_as_managed_app", resp.General.DeployAsManagedApp)
	d.Set("remove_app_when_mdm_profile_is_removed", resp.General.RemoveAppWhenMDMProfileIsRemoved)
	d.Set("prevent_backup_of_app_data", resp.General.PreventBackupOfAppData)
	d.Set("keep_description_and_icon_up_to_date", resp.General.KeepDescriptionAndIconUpToDate)
	d.Set("keep_app_updated_on_devices", resp.General.KeepAppUpdatedOnDevices)
	d.Set("free", resp.General.Free)
	d.Set("take_over_management", resp.General.TakeOverManagement)
	d.Set("host_externally", resp.General.HostExternally)
	d.Set("external_url", resp.General.ExternalURL)
	d.Set("mobile_device_provisioning_profile", resp.General.ProvisioningProfile)

	if resp.General.Category != nil {
		category := []map[string]interface{}{
			{
				"id":   resp.General.Category.ID,
				"name": resp.General.Category.Name,
			},
		}
		d.Set("category", category)
	}

	if resp.General.IPA.Name != "" || resp.General.IPA.URI != "" || resp.General.IPA.Data != "" {
		ipa := []map[string]interface{}{
			{
				"name": resp.General.IPA.Name,
				"uri":  resp.General.IPA.URI,
				"data": resp.General.IPA.Data,
			},
		}
		d.Set("ipa", ipa)
	}

	if resp.General.Icon.ID != 0 || resp.General.Icon.Name != "" || resp.General.Icon.URI != "" {
		icon := []map[string]interface{}{
			{
				"id":   resp.General.Icon.ID,
				"name": resp.General.Icon.Name,
				"uri":  resp.General.Icon.URI,
			},
		}
		d.Set("icon", icon)
	}

	if resp.General.Site != nil {
		site := []map[string]interface{}{
			{
				"id":   resp.General.Site.ID,
				"name": resp.General.Site.Name,
			},
		}
		d.Set("site", site)
	}

	if resp.SelfService.SelfServiceDescription != "" || resp.SelfService.NotificationMessage != "" {
		selfService := []map[string]interface{}{
			{
				"self_service_description": normalizeWhitespace(resp.SelfService.SelfServiceDescription),
				"feature_on_main_page":     resp.SelfService.FeatureOnMainPage,
				"notification":             resp.SelfService.Notification,
				"notification_subject":     resp.SelfService.NotificationSubject,
				"notification_message":     resp.SelfService.NotificationMessage,
			},
		}

		if resp.SelfService.SelfServiceIcon.ID != 0 || resp.SelfService.SelfServiceIcon.Name != "" || resp.SelfService.SelfServiceIcon.URI != "" {
			selfService[0]["self_service_icon"] = []map[string]interface{}{
				{
					"id":       resp.SelfService.SelfServiceIcon.ID,
					"filename": resp.SelfService.SelfServiceIcon.Name,
					"uri":      resp.SelfService.SelfServiceIcon.URI,
				},
			}
		}

		d.Set("self_service", selfService)
	}

	if resp.VPP.VPPAdminAccountID != 0 {
		vpp := []map[string]interface{}{
			{
				"assign_vpp_device_based_licenses": resp.VPP.AssignVPPDeviceBasedLicenses,
				"vpp_admin_account_id":             resp.VPP.VPPAdminAccountID,
			},
		}
		d.Set("vpp", vpp)
	}

	if resp.AppConfiguration.Preferences != "" {
		appConfig := []map[string]interface{}{
			{
				"preferences": normalizeWhitespace(resp.AppConfiguration.Preferences),
			},
		}
		d.Set("app_configuration", appConfig)
	}

	if scope := buildScopeState(&resp.Scope); scope != nil {
		d.Set("scope", []interface{}{scope})
	}

	return diags
}

// buildScopeState constructs the state representation of the scope
func buildScopeState(scope *jamfpro.MobileDeviceApplicationSubsetScope) map[string]interface{} {
	if scope == nil {
		return nil
	}

	scopeMap := map[string]interface{}{
		"all_mobile_devices": scope.AllMobileDevices,
		"all_jss_users":      scope.AllJSSUsers,
	}

	if len(scope.MobileDevices) > 0 {
		devices := make([]map[string]interface{}, len(scope.MobileDevices))
		for i, device := range scope.MobileDevices {
			devices[i] = map[string]interface{}{
				"id":               device.ID,
				"name":             device.Name,
				"udid":             device.UDID,
				"wifi_mac_address": device.WifiMacAddress,
			}
		}
		scopeMap["mobile_devices"] = devices
	}

	if len(scope.Buildings) > 0 {
		buildings := make([]map[string]interface{}, len(scope.Buildings))
		for i, building := range scope.Buildings {
			buildings[i] = map[string]interface{}{
				"id":   building.ID,
				"name": building.Name,
			}
		}
		scopeMap["buildings"] = buildings
	}

	if len(scope.Departments) > 0 {
		departments := make([]map[string]interface{}, len(scope.Departments))
		for i, dept := range scope.Departments {
			departments[i] = map[string]interface{}{
				"id":   dept.ID,
				"name": dept.Name,
			}
		}
		scopeMap["departments"] = departments
	}

	if len(scope.MobileDeviceGroups) > 0 {
		groups := make([]map[string]interface{}, len(scope.MobileDeviceGroups))
		for i, group := range scope.MobileDeviceGroups {
			groups[i] = map[string]interface{}{
				"id":   group.ID,
				"name": group.Name,
			}
		}
		scopeMap["mobile_device_groups"] = groups
	}

	if len(scope.JSSUsers) > 0 {
		users := make([]map[string]interface{}, len(scope.JSSUsers))
		for i, user := range scope.JSSUsers {
			users[i] = map[string]interface{}{
				"id":   user.ID,
				"name": user.Name,
			}
		}
		scopeMap["jss_users"] = users
	}

	if len(scope.JSSUserGroups) > 0 {
		groups := make([]map[string]interface{}, len(scope.JSSUserGroups))
		for i, group := range scope.JSSUserGroups {
			groups[i] = map[string]interface{}{
				"id":   group.ID,
				"name": group.Name,
			}
		}
		scopeMap["jss_user_groups"] = groups
	}

	if len(scope.Limitations.Users) > 0 || len(scope.Limitations.UserGroups) > 0 || len(scope.Limitations.NetworkSegments) > 0 {
		limitations := map[string]interface{}{}

		if len(scope.Limitations.Users) > 0 {
			users := make([]map[string]interface{}, len(scope.Limitations.Users))
			for i, user := range scope.Limitations.Users {
				users[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			limitations["users"] = users
		}

		if len(scope.Limitations.UserGroups) > 0 {
			groups := make([]map[string]interface{}, len(scope.Limitations.UserGroups))
			for i, group := range scope.Limitations.UserGroups {
				groups[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			limitations["user_groups"] = groups
		}

		if len(scope.Limitations.NetworkSegments) > 0 {
			segments := make([]map[string]interface{}, len(scope.Limitations.NetworkSegments))
			for i, segment := range scope.Limitations.NetworkSegments {
				segments[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			limitations["network_segments"] = segments
		}

		scopeMap["limitations"] = []interface{}{limitations}
	}

	if hasExclusions(scope.Exclusions) {
		exclusions := map[string]interface{}{}

		if len(scope.Exclusions.MobileDevices) > 0 {
			devices := make([]map[string]interface{}, len(scope.Exclusions.MobileDevices))
			for i, device := range scope.Exclusions.MobileDevices {
				devices[i] = map[string]interface{}{
					"id":               device.ID,
					"name":             device.Name,
					"udid":             device.UDID,
					"wifi_mac_address": device.WifiMacAddress,
				}
			}
			exclusions["mobile_devices"] = devices
		}

		if len(scope.Exclusions.Buildings) > 0 {
			buildings := make([]map[string]interface{}, len(scope.Exclusions.Buildings))
			for i, building := range scope.Exclusions.Buildings {
				buildings[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			exclusions["buildings"] = buildings
		}

		if len(scope.Exclusions.Departments) > 0 {
			departments := make([]map[string]interface{}, len(scope.Exclusions.Departments))
			for i, dept := range scope.Exclusions.Departments {
				departments[i] = map[string]interface{}{
					"id":   dept.ID,
					"name": dept.Name,
				}
			}
			exclusions["departments"] = departments
		}

		if len(scope.Exclusions.MobileDeviceGroups) > 0 {
			groups := make([]map[string]interface{}, len(scope.Exclusions.MobileDeviceGroups))
			for i, group := range scope.Exclusions.MobileDeviceGroups {
				groups[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusions["mobile_device_groups"] = groups
		}

		if len(scope.Exclusions.JSSUsers) > 0 {
			users := make([]map[string]interface{}, len(scope.Exclusions.JSSUsers))
			for i, user := range scope.Exclusions.JSSUsers {
				users[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusions["jss_users"] = users
		}

		if len(scope.Exclusions.JSSUserGroups) > 0 {
			groups := make([]map[string]interface{}, len(scope.Exclusions.JSSUserGroups))
			for i, group := range scope.Exclusions.JSSUserGroups {
				groups[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusions["jss_user_groups"] = groups
		}

		scopeMap["exclusions"] = []interface{}{exclusions}
	}

	return scopeMap
}
