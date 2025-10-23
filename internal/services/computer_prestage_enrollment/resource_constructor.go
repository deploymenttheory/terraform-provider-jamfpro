package computer_prestage_enrollment

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData, isUpdate bool) (*jamfpro.ResourceComputerPrestage, error) {
	versionLock := constructors.HandleVersionLock(d.Get("version_lock"), isUpdate)

	resource := &jamfpro.ResourceComputerPrestage{
		VersionLock:                        d.Get("version_lock").(int), // for some reason, this is not incremented. weird.
		DisplayName:                        d.Get("display_name").(string),
		Mandatory:                          jamfpro.BoolPtr(d.Get("mandatory").(bool)),
		MDMRemovable:                       jamfpro.BoolPtr(d.Get("mdm_removable").(bool)),
		SupportPhoneNumber:                 d.Get("support_phone_number").(string),
		SupportEmailAddress:                d.Get("support_email_address").(string),
		Department:                         d.Get("department").(string),
		DefaultPrestage:                    jamfpro.BoolPtr(d.Get("default_prestage").(bool)),
		EnrollmentSiteId:                   constructors.GetHCLStringOrDefaultInteger(d, "enrollment_site_id"),
		KeepExistingSiteMembership:         jamfpro.BoolPtr(d.Get("keep_existing_site_membership").(bool)),
		KeepExistingLocationInformation:    jamfpro.BoolPtr(d.Get("keep_existing_location_information").(bool)),
		RequireAuthentication:              jamfpro.BoolPtr(d.Get("require_authentication").(bool)),
		AuthenticationPrompt:               d.Get("authentication_prompt").(string),
		PreventActivationLock:              jamfpro.BoolPtr(d.Get("prevent_activation_lock").(bool)),
		EnableDeviceBasedActivationLock:    jamfpro.BoolPtr(d.Get("enable_device_based_activation_lock").(bool)),
		DeviceEnrollmentProgramInstanceId:  d.Get("device_enrollment_program_instance_id").(string),
		EnableRecoveryLock:                 jamfpro.BoolPtr(d.Get("enable_recovery_lock").(bool)),
		RecoveryLockPasswordType:           d.Get("recovery_lock_password_type").(string),
		RecoveryLockPassword:               d.Get("recovery_lock_password").(string),
		RotateRecoveryLockPassword:         jamfpro.BoolPtr(d.Get("rotate_recovery_lock_password").(bool)),
		PrestageMinimumOsTargetVersionType: d.Get("prestage_minimum_os_target_version_type").(string),
		MinimumOsSpecificVersion:           d.Get("minimum_os_specific_version").(string),
		ProfileUuid:                        d.Get("profile_uuid").(string),
		SiteId:                             d.Get("site_id").(string),
		CustomPackageDistributionPointId:   constructors.GetHCLStringOrDefaultInteger(d, "custom_package_distribution_point_id"),
		EnrollmentCustomizationId:          d.Get("enrollment_customization_id").(string),
		Language:                           d.Get("language").(string),
		Region:                             d.Get("region").(string),
		AutoAdvanceSetup:                   jamfpro.BoolPtr(d.Get("auto_advance_setup").(bool)),
		InstallProfilesDuringSetup:         jamfpro.BoolPtr(d.Get("install_profiles_during_setup").(bool)),
		PssoEnabled:                        jamfpro.BoolPtr(d.Get("platform_sso_enabled").(bool)),
		PlatformSsoAppBundleId:             d.Get("platform_sso_app_bundle_id").(string),
	}

	if v, ok := d.GetOk("skip_setup_items"); ok && len(v.([]interface{})) > 0 {
		skipSetupItemsMap := v.([]interface{})[0].(map[string]interface{})
		resource.SkipSetupItems = constructSkipSetupItems(skipSetupItemsMap)
	}

	if v, ok := d.GetOk("location_information"); ok && len(v.([]interface{})) > 0 {
		locationData := v.([]interface{})[0].(map[string]interface{})
		resource.LocationInformation = constructLocationInformation(locationData, isUpdate, versionLock)
	}

	if v, ok := d.GetOk("purchasing_information"); ok && len(v.([]interface{})) > 0 {
		purchasingData := v.([]interface{})[0].(map[string]interface{})
		resource.PurchasingInformation = constructPurchasingInformation(purchasingData, isUpdate, versionLock)
	}

	if v, ok := d.GetOk("account_settings"); ok && len(v.([]interface{})) > 0 {
		accountData := v.([]interface{})[0].(map[string]interface{})
		resource.AccountSettings = constructAccountSettings(accountData, isUpdate, versionLock)
	}

	if v, ok := d.GetOk("anchor_certificates"); ok {
		anchorCertificates := make([]string, len(v.([]interface{})))
		for i, cert := range v.([]interface{}) {
			anchorCertificates[i] = cert.(string)
		}
		resource.AnchorCertificates = anchorCertificates
	}

	resource.PrestageInstalledProfileIds = make([]string, 0)
	if v, ok := d.GetOk("prestage_installed_profile_ids"); ok {
		profileSet := v.(*schema.Set)
		for _, id := range profileSet.List() {
			resource.PrestageInstalledProfileIds = append(resource.PrestageInstalledProfileIds, id.(string))
		}
	}

	resource.CustomPackageIds = make([]string, 0)
	if v, ok := d.GetOk("custom_package_ids"); ok {
		packageSet := v.(*schema.Set)
		for _, id := range packageSet.List() {
			resource.CustomPackageIds = append(resource.CustomPackageIds, id.(string))
		}
	}

	if v, ok := d.GetOk("onboarding_items"); ok {
		onboardingItems := make([]jamfpro.OnboardingItem, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			onboardingItems[i] = jamfpro.OnboardingItem{
				SelfServiceEntityType: itemMap["self_service_entity_type"].(string),
				ID:                    itemMap["id"].(string),
				EntityId:              itemMap["entity_id"].(string),
				Priority:              itemMap["priority"].(int),
			}
		}
		resource.OnboardingItems = onboardingItems
	}

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceJSON, err := common.SerializeAndRedactJSON(resource, []string{"AccountSettings.AdminPassword", "AccountSettings.AdminUsername"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Prestage to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Prestage JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructSkipSetupItems constructs the SkipSetupItems subset of a Computer Prestage resource.
func constructSkipSetupItems(data map[string]interface{}) jamfpro.ComputerPrestageSubsetSkipSetupItems {
	return jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:                 jamfpro.BoolPtr(data["biometric"].(bool)),
		TermsOfAddress:            jamfpro.BoolPtr(data["terms_of_address"].(bool)),
		FileVault:                 jamfpro.BoolPtr(data["file_vault"].(bool)),
		ICloudDiagnostics:         jamfpro.BoolPtr(data["icloud_diagnostics"].(bool)),
		Diagnostics:               jamfpro.BoolPtr(data["diagnostics"].(bool)),
		Accessibility:             jamfpro.BoolPtr(data["accessibility"].(bool)),
		AppleID:                   jamfpro.BoolPtr(data["apple_id"].(bool)),
		ScreenTime:                jamfpro.BoolPtr(data["screen_time"].(bool)),
		Siri:                      jamfpro.BoolPtr(data["siri"].(bool)),
		DisplayTone:               jamfpro.BoolPtr(data["display_tone"].(bool)),
		Restore:                   jamfpro.BoolPtr(data["restore"].(bool)),
		Appearance:                jamfpro.BoolPtr(data["appearance"].(bool)),
		Privacy:                   jamfpro.BoolPtr(data["privacy"].(bool)),
		Payment:                   jamfpro.BoolPtr(data["payment"].(bool)),
		Registration:              jamfpro.BoolPtr(data["registration"].(bool)),
		TOS:                       jamfpro.BoolPtr(data["tos"].(bool)),
		ICloudStorage:             jamfpro.BoolPtr(data["icloud_storage"].(bool)),
		Location:                  jamfpro.BoolPtr(data["location"].(bool)),
		Intelligence:              jamfpro.BoolPtr(data["intelligence"].(bool)),
		EnableLockdownMode:        jamfpro.BoolPtr(data["enable_lockdown_mode"].(bool)),
		Welcome:                   jamfpro.BoolPtr(data["welcome"].(bool)),
		Wallpaper:                 jamfpro.BoolPtr(data["wallpaper"].(bool)),
		SoftwareUpdate:            jamfpro.BoolPtr(data["software_update"].(bool)),
		AdditionalPrivacySettings: jamfpro.BoolPtr(data["additional_privacy_settings"].(bool)),
	}
}

// constructLocationInformation constructs the LocationInformation subset of a Computer Prestage resource.
func constructLocationInformation(data map[string]interface{}, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetLocationInformation {
	d := &schema.ResourceData{}
	for k, v := range data {
		d.Set(k, v)
	}

	return jamfpro.ComputerPrestageSubsetLocationInformation{
		ID:           "-1",
		VersionLock:  versionLock,
		Username:     data["username"].(string),
		Realname:     data["realname"].(string),
		Phone:        data["phone"].(string),
		Email:        data["email"].(string),
		Room:         data["room"].(string),
		Position:     data["position"].(string),
		DepartmentId: constructors.GetHCLStringOrDefaultInteger(d, "department_id"),
		BuildingId:   constructors.GetHCLStringOrDefaultInteger(d, "building_id"),
	}
}

// constructPurchasingInformation constructs the PurchasingInformation subset of a Computer Prestage resource.
func constructPurchasingInformation(data map[string]interface{}, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetPurchasingInformation {
	d := &schema.ResourceData{}
	for k, v := range data {
		d.Set(k, v)
	}

	return jamfpro.ComputerPrestageSubsetPurchasingInformation{
		ID:                "-1",
		VersionLock:       versionLock,
		Leased:            jamfpro.BoolPtr(data["leased"].(bool)),
		Purchased:         jamfpro.BoolPtr(data["purchased"].(bool)),
		AppleCareId:       data["apple_care_id"].(string),
		PONumber:          data["po_number"].(string),
		Vendor:            data["vendor"].(string),
		PurchasePrice:     data["purchase_price"].(string),
		LifeExpectancy:    data["life_expectancy"].(int),
		PurchasingAccount: data["purchasing_account"].(string),
		PurchasingContact: data["purchasing_contact"].(string),
		LeaseDate:         constructors.GetDateOrDefaultDate(d, "lease_date"),
		PODate:            constructors.GetDateOrDefaultDate(d, "po_date"),
		WarrantyDate:      constructors.GetDateOrDefaultDate(d, "warranty_date"),
	}
}

// constructAccountSettings constructs the AccountSettings subset of a Computer Prestage resource.
func constructAccountSettings(data map[string]interface{}, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetAccountSettings {
	return jamfpro.ComputerPrestageSubsetAccountSettings{
		ID:                                      "-1",
		VersionLock:                             versionLock,
		PayloadConfigured:                       jamfpro.BoolPtr(data["payload_configured"].(bool)),
		LocalAdminAccountEnabled:                jamfpro.BoolPtr(data["local_admin_account_enabled"].(bool)),
		AdminUsername:                           data["admin_username"].(string),
		AdminPassword:                           data["admin_password"].(string),
		HiddenAdminAccount:                      jamfpro.BoolPtr(data["hidden_admin_account"].(bool)),
		LocalUserManaged:                        jamfpro.BoolPtr(data["local_user_managed"].(bool)),
		UserAccountType:                         data["user_account_type"].(string),
		PrefillPrimaryAccountInfoFeatureEnabled: jamfpro.BoolPtr(data["prefill_primary_account_info_feature_enabled"].(bool)),
		PrefillType:                             data["prefill_type"].(string),
		PrefillAccountFullName:                  data["prefill_account_full_name"].(string),
		PrefillAccountUserName:                  data["prefill_account_user_name"].(string),
		PreventPrefillInfoFromModification:      jamfpro.BoolPtr(data["prevent_prefill_info_from_modification"].(bool)),
	}
}
