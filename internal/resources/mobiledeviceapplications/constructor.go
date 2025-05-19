package mobiledeviceapplications

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructMobileDeviceApplication constructs a MobileDeviceApplication object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceApplication, error) {
	app := &jamfpro.ResourceMobileDeviceApplication{
		General: jamfpro.MobileDeviceApplicationSubsetGeneral{
			Name:                             d.Get("name").(string),
			DisplayName:                      d.Get("display_name").(string),
			Description:                      d.Get("description").(string),
			BundleID:                         d.Get("bundle_id").(string),
			Version:                          d.Get("version").(string),
			InternalApp:                      jamfpro.BoolPtr(d.Get("internal_app").(bool)),
			ITunesStoreURL:                   d.Get("itunes_store_url").(string),
			MakeAvailableAfterInstall:        d.Get("make_available_after_install").(bool),
			ITunesCountryRegion:              d.Get("itunes_country_region").(string),
			ITunesSyncTime:                   d.Get("itunes_sync_time").(int),
			DeploymentType:                   d.Get("deployment_type").(string),
			DeployAutomatically:              jamfpro.BoolPtr(d.Get("deploy_automatically").(bool)),
			DeployAsManagedApp:               jamfpro.BoolPtr(d.Get("deploy_as_managed_app").(bool)),
			RemoveAppWhenMDMProfileIsRemoved: jamfpro.BoolPtr(d.Get("remove_app_when_mdm_profile_is_removed").(bool)),
			PreventBackupOfAppData:           jamfpro.BoolPtr(d.Get("prevent_backup_of_app_data").(bool)),
			KeepDescriptionAndIconUpToDate:   jamfpro.BoolPtr(d.Get("keep_description_and_icon_up_to_date").(bool)),
			KeepAppUpdatedOnDevices:          jamfpro.BoolPtr(d.Get("keep_app_updated_on_devices").(bool)),
			Free:                             jamfpro.BoolPtr(d.Get("free").(bool)),
			TakeOverManagement:               jamfpro.BoolPtr(d.Get("take_over_management").(bool)),
			HostExternally:                   jamfpro.BoolPtr(d.Get("host_externally").(bool)),
			ExternalURL:                      d.Get("external_url").(string),
			ProvisioningProfile:              d.Get("mobile_device_provisioning_profile").(int),
		},
	}
	if v, ok := d.GetOk("category"); ok && len(v.([]interface{})) > 0 {
		categoryMap := v.([]interface{})[0].(map[string]interface{})
		app.General.Category = &jamfpro.SharedResourceCategory{
			ID:   categoryMap["id"].(int),
			Name: categoryMap["name"].(string),
		}
	}

	if v, ok := d.GetOk("ipa"); ok && len(v.([]interface{})) > 0 {
		ipaMap := v.([]interface{})[0].(map[string]interface{})
		app.General.IPA = jamfpro.MobileDeviceApplicationSubsetGeneralIPA{
			Name: ipaMap["name"].(string),
			URI:  ipaMap["uri"].(string),
			Data: ipaMap["data"].(string),
		}
	}

	if v, ok := d.GetOk("icon"); ok && len(v.([]interface{})) > 0 {
		iconMap := v.([]interface{})[0].(map[string]interface{})
		app.General.Icon = jamfpro.MobileDeviceApplicationSubsetIcon{
			ID:   iconMap["id"].(int),
			Name: iconMap["name"].(string),
			URI:  iconMap["uri"].(string),
		}
	}

	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteMap := v.([]interface{})[0].(map[string]interface{})
		app.General.Site = &jamfpro.SharedResourceSite{
			ID:   siteMap["id"].(int),
			Name: siteMap["name"].(string),
		}
	}

	if v, ok := d.GetOk("scope"); ok && len(v.([]interface{})) > 0 {
		scopeMap := v.([]interface{})[0].(map[string]interface{})
		scope := jamfpro.MobileDeviceApplicationSubsetScope{
			AllMobileDevices: jamfpro.BoolPtr(scopeMap["all_mobile_devices"].(bool)),
			AllJSSUsers:      jamfpro.BoolPtr(scopeMap["all_jss_users"].(bool)),
		}
		if devices, ok := scopeMap["mobile_devices"].([]interface{}); ok {
			for _, device := range devices {
				deviceMap := device.(map[string]interface{})
				scope.MobileDevices = append(scope.MobileDevices, jamfpro.MobileDeviceApplicationSubsetMobileDevice{
					ID:             deviceMap["id"].(int),
					Name:           deviceMap["name"].(string),
					UDID:           deviceMap["udid"].(string),
					WifiMacAddress: deviceMap["wifi_mac_address"].(string),
				})
			}
		}

		if buildings, ok := scopeMap["buildings"].([]interface{}); ok {
			for _, building := range buildings {
				buildingMap := building.(map[string]interface{})
				scope.Buildings = append(scope.Buildings, jamfpro.MobileDeviceApplicationSubsetBuilding{
					ID:   buildingMap["id"].(int),
					Name: buildingMap["name"].(string),
				})
			}
		}

		if departments, ok := scopeMap["departments"].([]interface{}); ok {
			for _, department := range departments {
				departmentMap := department.(map[string]interface{})
				scope.Departments = append(scope.Departments, jamfpro.MobileDeviceApplicationSubsetDepartment{
					ID:   departmentMap["id"].(int),
					Name: departmentMap["name"].(string),
				})
			}
		}
		if groups, ok := scopeMap["mobile_device_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				scope.MobileDeviceGroups = append(scope.MobileDeviceGroups, jamfpro.MobileDeviceApplicationSubsetMobileDeviceGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		if users, ok := scopeMap["jss_users"].([]interface{}); ok {
			for _, user := range users {
				userMap := user.(map[string]interface{})
				scope.JSSUsers = append(scope.JSSUsers, jamfpro.MobileDeviceApplicationSubsetJSSUser{
					ID:   userMap["id"].(int),
					Name: userMap["name"].(string),
				})
			}
		}

		if groups, ok := scopeMap["jss_user_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				scope.JSSUserGroups = append(scope.JSSUserGroups, jamfpro.MobileDeviceApplicationSubsetJSSUserGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		if limitations, ok := scopeMap["limitations"].([]interface{}); ok && len(limitations) > 0 {
			limitationsMap := limitations[0].(map[string]interface{})
			if users, ok := limitationsMap["users"].([]interface{}); ok {
				for _, user := range users {
					userMap := user.(map[string]interface{})
					scope.Limitations.Users = append(scope.Limitations.Users, jamfpro.MobileDeviceApplicationSubsetUser{
						ID:   userMap["id"].(int),
						Name: userMap["name"].(string),
					})
				}
			}

			if groups, ok := limitationsMap["user_groups"].([]interface{}); ok {
				for _, group := range groups {
					groupMap := group.(map[string]interface{})
					scope.Limitations.UserGroups = append(scope.Limitations.UserGroups, jamfpro.MobileDeviceApplicationSubsetUserGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			if segments, ok := limitationsMap["network_segments"].([]interface{}); ok {
				for _, segment := range segments {
					segmentMap := segment.(map[string]interface{})
					scope.Limitations.NetworkSegments = append(scope.Limitations.NetworkSegments, jamfpro.MobileDeviceApplicationSubsetNetworkSegment{
						ID:   segmentMap["id"].(int),
						Name: segmentMap["name"].(string),
						UID:  segmentMap["uid"].(string),
					})
				}
			}
		}

		if exclusions, ok := scopeMap["exclusions"].([]interface{}); ok && len(exclusions) > 0 {
			exclusionsMap := exclusions[0].(map[string]interface{})
			if devices, ok := exclusionsMap["mobile_devices"].([]interface{}); ok {
				for _, device := range devices {
					deviceMap := device.(map[string]interface{})
					scope.Exclusions.MobileDevices = append(scope.Exclusions.MobileDevices,
						jamfpro.MobileDeviceApplicationSubsetMobileDevice{
							ID:             deviceMap["id"].(int),
							Name:           deviceMap["name"].(string),
							UDID:           deviceMap["udid"].(string),
							WifiMacAddress: deviceMap["wifi_mac_address"].(string),
						})
				}
			}

			if buildings, ok := exclusionsMap["buildings"].([]interface{}); ok {
				for _, building := range buildings {
					buildingMap := building.(map[string]interface{})
					scope.Exclusions.Buildings = append(scope.Exclusions.Buildings,
						jamfpro.MobileDeviceApplicationSubsetBuilding{
							ID:   buildingMap["id"].(int),
							Name: buildingMap["name"].(string),
						})
				}
			}

			if departments, ok := exclusionsMap["departments"].([]interface{}); ok {
				for _, department := range departments {
					departmentMap := department.(map[string]interface{})
					scope.Exclusions.Departments = append(scope.Exclusions.Departments,
						jamfpro.MobileDeviceApplicationSubsetDepartment{
							ID:   departmentMap["id"].(int),
							Name: departmentMap["name"].(string),
						})
				}
			}

			if groups, ok := exclusionsMap["mobile_device_groups"].([]interface{}); ok {
				for _, group := range groups {
					groupMap := group.(map[string]interface{})
					scope.Exclusions.MobileDeviceGroups = append(scope.Exclusions.MobileDeviceGroups,
						jamfpro.MobileDeviceApplicationSubsetMobileDeviceGroup{
							ID:   groupMap["id"].(int),
							Name: groupMap["name"].(string),
						})
				}
			}

			if users, ok := exclusionsMap["users"].([]interface{}); ok {
				for _, user := range users {
					userMap := user.(map[string]interface{})
					scope.Exclusions.JSSUsers = append(scope.Exclusions.JSSUsers,
						jamfpro.MobileDeviceApplicationSubsetJSSUser{
							ID:   userMap["id"].(int),
							Name: userMap["name"].(string),
						})
				}
			}

			if groups, ok := exclusionsMap["user_groups"].([]interface{}); ok {
				for _, group := range groups {
					groupMap := group.(map[string]interface{})
					scope.Exclusions.JSSUserGroups = append(scope.Exclusions.JSSUserGroups,
						jamfpro.MobileDeviceApplicationSubsetJSSUserGroup{
							ID:   groupMap["id"].(int),
							Name: groupMap["name"].(string),
						})
				}
			}
		}

		app.Scope = scope
	}
	if v, ok := d.GetOk("self_service"); ok && len(v.([]interface{})) > 0 {
		selfServiceMap := v.([]interface{})[0].(map[string]interface{})
		selfService := jamfpro.MobileDeviceApplicationSubsetGeneralSelfService{
			SelfServiceDescription: selfServiceMap["self_service_description"].(string),
			FeatureOnMainPage:      jamfpro.BoolPtr(selfServiceMap["feature_on_main_page"].(bool)),
			Notification:           jamfpro.BoolPtr(selfServiceMap["notification"].(bool)),
			NotificationSubject:    selfServiceMap["notification_subject"].(string),
			NotificationMessage:    selfServiceMap["notification_message"].(string),
		}
		if icon, ok := selfServiceMap["self_service_icon"].([]interface{}); ok && len(icon) > 0 {
			iconMap := icon[0].(map[string]interface{})
			selfService.SelfServiceIcon = jamfpro.MobileDeviceApplicationSubsetIcon{
				ID:   iconMap["id"].(int),
				Name: iconMap["filename"].(string),
				URI:  iconMap["uri"].(string),
			}
		}

		app.SelfService = selfService
	}

	if v, ok := d.GetOk("vpp"); ok && len(v.([]interface{})) > 0 {
		vppMap := v.([]interface{})[0].(map[string]interface{})
		app.VPP = jamfpro.MobileDeviceApplicationSubsetGeneralVPP{
			AssignVPPDeviceBasedLicenses: jamfpro.BoolPtr(vppMap["assign_vpp_device_based_licenses"].(bool)),
			VPPAdminAccountID:            vppMap["vpp_admin_account_id"].(int),
		}
	}

	if v, ok := d.GetOk("app_configuration"); ok && len(v.([]interface{})) > 0 {
		appConfigMap := v.([]interface{})[0].(map[string]interface{})
		app.AppConfiguration = jamfpro.MobileDeviceApplicationSubsetGeneralAppConfiguration{
			Preferences: appConfigMap["preferences"].(string),
		}
	}

	resourceXML, err := xml.MarshalIndent(app, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Mobile Device Application '%s' to XML: %v", app.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Mobile Device Application XML:\n%s\n", string(resourceXML))

	return app, nil
}
