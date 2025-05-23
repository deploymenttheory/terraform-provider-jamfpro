package mobiledeviceapplications

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructMobileDeviceApplication constructs a MobileDeviceApplication object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceApplication, error) {
	resource := &jamfpro.ResourceMobileDeviceApplication{
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

	if v, ok := d.GetOk("ipa"); ok && len(v.([]interface{})) > 0 {
		ipaMap := v.([]interface{})[0].(map[string]interface{})
		resource.General.IPA = jamfpro.MobileDeviceApplicationSubsetGeneralIPA{
			Name: ipaMap["name"].(string),
			URI:  ipaMap["uri"].(string),
			Data: ipaMap["data"].(string),
		}
	}

	if v, ok := d.GetOk("icon"); ok && len(v.([]interface{})) > 0 {
		iconMap := v.([]interface{})[0].(map[string]interface{})
		resource.General.Icon = jamfpro.MobileDeviceApplicationSubsetIcon{
			ID:   iconMap["id"].(int),
			Name: iconMap["name"].(string),
			URI:  iconMap["uri"].(string),
		}
	}

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	if _, ok := d.GetOk("scope"); ok {
		resource.Scope = constructMobileDeviceApplicationSubsetScope(d)
	} else {
		log.Printf("[DEBUG] constructJamfProMobileDeviceApplication: No scope block found or it's empty.")
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

		resource.SelfService = selfService
	}

	if v, ok := d.GetOk("vpp"); ok && len(v.([]interface{})) > 0 {
		vppMap := v.([]interface{})[0].(map[string]interface{})
		resource.VPP = jamfpro.MobileDeviceApplicationSubsetGeneralVPP{
			AssignVPPDeviceBasedLicenses: jamfpro.BoolPtr(vppMap["assign_vpp_device_based_licenses"].(bool)),
			VPPAdminAccountID:            vppMap["vpp_admin_account_id"].(int),
		}
	}

	if v, ok := d.GetOk("app_configuration"); ok && len(v.([]interface{})) > 0 {
		appConfigMap := v.([]interface{})[0].(map[string]interface{})
		resource.AppConfiguration = jamfpro.MobileDeviceApplicationSubsetGeneralAppConfiguration{
			Preferences: appConfigMap["preferences"].(string),
		}
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Mobile Device Application '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Mobile Device Application XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// constructMobileDeviceApplicationSubsetScope constructs the scope from the provided schema data.
func constructMobileDeviceApplicationSubsetScope(d *schema.ResourceData) jamfpro.MobileDeviceApplicationSubsetScope {
	scope := jamfpro.MobileDeviceApplicationSubsetScope{}
	scopeData := d.Get("scope").([]interface{})[0].(map[string]interface{})

	scope.AllMobileDevices = jamfpro.BoolPtr(scopeData["all_mobile_devices"].(bool))
	scope.AllJSSUsers = jamfpro.BoolPtr(scopeData["all_jss_users"].(bool))

	var mobileDevices []jamfpro.MobileDeviceApplicationSubsetMobileDevice
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetMobileDevice, int](
		"scope.0.mobile_device_ids", "ID", d, &mobileDevices); err == nil {
		scope.MobileDevices = mobileDevices
	}

	var mobileDeviceGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.mobile_device_group_ids", "ID", d, &mobileDeviceGroups); err == nil {
		scope.MobileDeviceGroups = mobileDeviceGroups
	}

	var buildings []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.building_ids", "ID", d, &buildings); err == nil {
		scope.Buildings = buildings
	}

	var departments []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.department_ids", "ID", d, &departments); err == nil {
		scope.Departments = departments
	}

	var jssUsers []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		scope.JSSUsers = jssUsers
	}

	var jssUserGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		scope.JSSUserGroups = jssUserGroups
	}

	if _, ok := d.GetOk("scope.0.limitations"); ok {
		scope.Limitations = constructLimitations(d)
	}

	if _, ok := d.GetOk("scope.0.exclusions"); ok {
		scope.Exclusions = constructExclusions(d)
	}

	return scope
}

// constructLimitations constructs the limitations from the provided schema data.
func constructLimitations(d *schema.ResourceData) jamfpro.MobileDeviceApplicationSubsetLimitation {
	limitations := jamfpro.MobileDeviceApplicationSubsetLimitation{}

	var users []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, string](
		"scope.0.limitations.0.directory_service_or_local_usernames", "Name", d, &users); err == nil {
		limitations.Users = users
	}

	var userGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, &userGroups); err == nil {
		limitations.UserGroups = userGroups
	}

	var networkSegments []jamfpro.MobileDeviceApplicationSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetNetworkSegment, int](
		"scope.0.limitations.0.network_segment_ids", "ID", d, &networkSegments); err == nil {
		limitations.NetworkSegments = networkSegments
	}

	return limitations
}

// constructExclusions constructs the exclusions from the provided schema data.
func constructExclusions(d *schema.ResourceData) jamfpro.MobileDeviceApplicationSubsetExclusion {
	exclusions := jamfpro.MobileDeviceApplicationSubsetExclusion{}

	var mobileDevices []jamfpro.MobileDeviceApplicationSubsetMobileDevice
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetMobileDevice, int](
		"scope.0.exclusions.0.mobile_device_ids", "ID", d, &mobileDevices); err == nil {
		exclusions.MobileDevices = mobileDevices
	}

	var mobileDeviceGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.mobile_device_group_ids", "ID", d, &mobileDeviceGroups); err == nil {
		exclusions.MobileDeviceGroups = mobileDeviceGroups
	}

	var users []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, string](
		"scope.0.exclusions.0.directory_service_or_local_usernames", "Name", d, &users); err == nil {
		exclusions.Users = users
	}

	var userGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.directory_service_usergroup_ids", "ID", d, &userGroups); err == nil {
		exclusions.UserGroups = userGroups
	}

	var buildings []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.building_ids", "ID", d, &buildings); err == nil {
		exclusions.Buildings = buildings
	}

	var departments []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.department_ids", "ID", d, &departments); err == nil {
		exclusions.Departments = departments
	}

	var networkSegments []jamfpro.MobileDeviceApplicationSubsetNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetNetworkSegment, int](
		"scope.0.exclusions.0.network_segment_ids", "ID", d, &networkSegments); err == nil {
		exclusions.NetworkSegments = networkSegments
	}

	var jssUsers []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		exclusions.JSSUsers = jssUsers
	}

	var jssUserGroups []jamfpro.MobileDeviceApplicationSubsetScopeEntity
	if err := constructors.MapSetToStructs[jamfpro.MobileDeviceApplicationSubsetScopeEntity, int](
		"scope.0.exclusions.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		exclusions.JSSUserGroups = jssUserGroups
	}

	return exclusions
}
