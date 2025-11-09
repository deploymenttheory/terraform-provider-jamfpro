package computer_prestage_enrollment

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/redact"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData, isUpdate bool) (*jamfpro.ResourceComputerPrestage, error) {
	versionLock := constructors.HandleVersionLock(d.Get("version_lock"), isUpdate)

	resource := &jamfpro.ResourceComputerPrestage{
		VersionLock:                        d.Get("version_lock").(int), // for some reason, this is not incremented. weird.
		DisplayName:                        d.Get("display_name").(string),
		Mandatory:                          d.Get("mandatory").(bool),
		MDMRemovable:                       d.Get("mdm_removable").(bool),
		SupportPhoneNumber:                 d.Get("support_phone_number").(string),
		SupportEmailAddress:                d.Get("support_email_address").(string),
		Department:                         d.Get("department").(string),
		DefaultPrestage:                    d.Get("default_prestage").(bool),
		EnrollmentSiteId:                   constructors.GetHCLStringOrDefaultInteger(d, "enrollment_site_id"),
		KeepExistingSiteMembership:         d.Get("keep_existing_site_membership").(bool),
		KeepExistingLocationInformation:    d.Get("keep_existing_location_information").(bool),
		RequireAuthentication:              d.Get("require_authentication").(bool),
		AuthenticationPrompt:               d.Get("authentication_prompt").(string),
		PreventActivationLock:              d.Get("prevent_activation_lock").(bool),
		EnableDeviceBasedActivationLock:    d.Get("enable_device_based_activation_lock").(bool),
		DeviceEnrollmentProgramInstanceId:  d.Get("device_enrollment_program_instance_id").(string),
		EnableRecoveryLock:                 d.Get("enable_recovery_lock").(bool),
		RecoveryLockPasswordType:           d.Get("recovery_lock_password_type").(string),
		RecoveryLockPassword:               d.Get("recovery_lock_password").(string),
		RotateRecoveryLockPassword:         d.Get("rotate_recovery_lock_password").(bool),
		PrestageMinimumOsTargetVersionType: d.Get("prestage_minimum_os_target_version_type").(string),
		MinimumOsSpecificVersion:           d.Get("minimum_os_specific_version").(string),
		ProfileUuid:                        d.Get("profile_uuid").(string),
		SiteId:                             d.Get("site_id").(string),
		CustomPackageDistributionPointId:   constructors.GetHCLStringOrDefaultInteger(d, "custom_package_distribution_point_id"),
		EnrollmentCustomizationId:          d.Get("enrollment_customization_id").(string),
		Language:                           d.Get("language").(string),
		Region:                             d.Get("region").(string),
		AutoAdvanceSetup:                   d.Get("auto_advance_setup").(bool),
		InstallProfilesDuringSetup:         d.Get("install_profiles_during_setup").(bool),
		PssoEnabled:                        d.Get("platform_sso_enabled").(bool),
		PlatformSsoAppBundleId:             d.Get("platform_sso_app_bundle_id").(string),
	}

	if v, ok := d.GetOk("skip_setup_items"); ok && len(v.([]any)) > 0 {
		skipSetupItemsMap := v.([]any)[0].(map[string]any)
		resource.SkipSetupItems = constructSkipSetupItems(skipSetupItemsMap)
	}

	if v, ok := d.GetOk("location_information"); ok && len(v.([]any)) > 0 {
		locationData := v.([]any)[0].(map[string]any)
		resource.LocationInformation = constructLocationInformation(locationData, isUpdate, versionLock)
	}

	if v, ok := d.GetOk("purchasing_information"); ok && len(v.([]any)) > 0 {
		purchasingData := v.([]any)[0].(map[string]any)
		resource.PurchasingInformation = constructPurchasingInformation(purchasingData, isUpdate, versionLock)
	}

	if v, ok := d.GetOk("account_settings"); ok && len(v.([]any)) > 0 {
		accountData := v.([]any)[0].(map[string]any)
		resource.AccountSettings = constructAccountSettings(accountData, isUpdate, versionLock)
	} else {
		resource.AccountSettings = jamfpro.ComputerPrestageSubsetAccountSettings{
			VersionLock:       versionLock,
			PayloadConfigured: true,
		}
	}

	if v, ok := d.GetOk("anchor_certificates"); ok {
		anchorCertificates := make([]string, len(v.([]any)))
		for i, cert := range v.([]any) {
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

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceJSON, err := redact.SerializeAndRedactJSON(resource, []string{"AccountSettings.AdminPassword", "AccountSettings.AdminUsername"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Prestage to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Prestage JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructSkipSetupItems constructs the SkipSetupItems subset of a Computer Prestage resource.
func constructSkipSetupItems(data map[string]any) jamfpro.ComputerPrestageSubsetSkipSetupItems {
	return jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:                 data["biometric"].(bool),
		TermsOfAddress:            data["terms_of_address"].(bool),
		FileVault:                 data["file_vault"].(bool),
		ICloudDiagnostics:         data["icloud_diagnostics"].(bool),
		Diagnostics:               data["diagnostics"].(bool),
		Accessibility:             data["accessibility"].(bool),
		AppleID:                   data["apple_id"].(bool),
		ScreenTime:                data["screen_time"].(bool),
		Siri:                      data["siri"].(bool),
		DisplayTone:               data["display_tone"].(bool),
		Restore:                   data["restore"].(bool),
		Appearance:                data["appearance"].(bool),
		Privacy:                   data["privacy"].(bool),
		Payment:                   data["payment"].(bool),
		Registration:              data["registration"].(bool),
		TOS:                       data["tos"].(bool),
		ICloudStorage:             data["icloud_storage"].(bool),
		Location:                  data["location"].(bool),
		Intelligence:              data["intelligence"].(bool),
		EnableLockdownMode:        data["enable_lockdown_mode"].(bool),
		Welcome:                   data["welcome"].(bool),
		Wallpaper:                 data["wallpaper"].(bool),
		SoftwareUpdate:            data["software_update"].(bool),
		AdditionalPrivacySettings: data["additional_privacy_settings"].(bool),
		OSShowcase:                data["os_showcase"].(bool),
	}
}

// constructLocationInformation constructs the LocationInformation subset of a Computer Prestage resource.
func constructLocationInformation(data map[string]any, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetLocationInformation {
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
func constructPurchasingInformation(data map[string]any, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetPurchasingInformation {
	d := &schema.ResourceData{}
	for k, v := range data {
		d.Set(k, v)
	}

	return jamfpro.ComputerPrestageSubsetPurchasingInformation{
		ID:                "-1",
		VersionLock:       versionLock,
		Leased:            data["leased"].(bool),
		Purchased:         data["purchased"].(bool),
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
func constructAccountSettings(data map[string]any, isUpdate bool, versionLock int) jamfpro.ComputerPrestageSubsetAccountSettings {
	return jamfpro.ComputerPrestageSubsetAccountSettings{
		VersionLock:                             versionLock,
		PayloadConfigured:                       true,
		LocalAdminAccountEnabled:                data["local_admin_account_enabled"].(bool),
		AdminUsername:                           data["admin_username"].(string),
		AdminPassword:                           data["admin_password"].(string),
		HiddenAdminAccount:                      data["hidden_admin_account"].(bool),
		LocalUserManaged:                        data["local_user_managed"].(bool),
		UserAccountType:                         data["user_account_type"].(string),
		PrefillPrimaryAccountInfoFeatureEnabled: data["prefill_primary_account_info_feature_enabled"].(bool),
		PrefillType:                             data["prefill_type"].(string),
		PrefillAccountFullName:                  data["prefill_account_full_name"].(string),
		PrefillAccountUserName:                  data["prefill_account_user_name"].(string),
		PreventPrefillInfoFromModification:      data["prevent_prefill_info_from_modification"].(bool),
	}
}
