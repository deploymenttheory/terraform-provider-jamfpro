package computer_prestage_enrollment

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Prestage Enrollment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceComputerPrestage) diag.Diagnostics {
	var diags diag.Diagnostics

	prestageAttributes := map[string]interface{}{
		"id":                                      resp.ID,
		"version_lock":                            resp.VersionLock,
		"display_name":                            resp.DisplayName,
		"mandatory":                               resp.Mandatory,
		"mdm_removable":                           resp.MDMRemovable,
		"support_phone_number":                    resp.SupportPhoneNumber,
		"support_email_address":                   resp.SupportEmailAddress,
		"department":                              resp.Department,
		"default_prestage":                        resp.DefaultPrestage,
		"enrollment_site_id":                      resp.EnrollmentSiteId,
		"keep_existing_site_membership":           resp.KeepExistingSiteMembership,
		"keep_existing_location_information":      resp.KeepExistingLocationInformation,
		"require_authentication":                  resp.RequireAuthentication,
		"authentication_prompt":                   resp.AuthenticationPrompt,
		"prevent_activation_lock":                 resp.PreventActivationLock,
		"enable_device_based_activation_lock":     resp.EnableDeviceBasedActivationLock,
		"device_enrollment_program_instance_id":   resp.DeviceEnrollmentProgramInstanceId,
		"skip_setup_items":                        []interface{}{skipSetupItems(resp.SkipSetupItems)},
		"anchor_certificates":                     resp.AnchorCertificates,
		"enrollment_customization_id":             resp.EnrollmentCustomizationId,
		"language":                                resp.Language,
		"region":                                  resp.Region,
		"auto_advance_setup":                      resp.AutoAdvanceSetup,
		"install_profiles_during_setup":           resp.InstallProfilesDuringSetup,
		"prestage_installed_profile_ids":          resp.PrestageInstalledProfileIds,
		"custom_package_ids":                      resp.CustomPackageIds,
		"custom_package_distribution_point_id":    resp.CustomPackageDistributionPointId,
		"enable_recovery_lock":                    resp.EnableRecoveryLock,
		"recovery_lock_password_type":             resp.RecoveryLockPasswordType,
		"recovery_lock_password":                  getHCLValue(d, "recovery_lock_password"),
		"rotate_recovery_lock_password":           resp.RotateRecoveryLockPassword,
		"prestage_minimum_os_target_version_type": resp.PrestageMinimumOsTargetVersionType,
		"minimum_os_specific_version":             resp.MinimumOsSpecificVersion,
		"profile_uuid":                            resp.ProfileUuid,
		"site_id":                                 resp.SiteId,
		// "enabled":                                 resp.Enabled,
		// "sso_for_enrollment_enabled":              resp.SsoForEnrollmentEnabled,
		// "sso_bypass_allowed":                      resp.SsoBypassAllowed,
		// "sso_enabled":                             resp.SsoEnabled,
		// "sso_for_mac_os_self_service_enabled":     resp.SsoForMacOsSelfServiceEnabled,
		// "token_expiration_disabled":               resp.TokenExpirationDisabled,
		// "user_attribute_enabled":                  resp.UserAttributeEnabled,
		// "user_attribute_name":                     resp.UserAttributeName,
		// "user_mapping":                            resp.UserMapping,
		// "enrollment_sso_for_account_driven_enrollment_enabled": resp.EnrollmentSsoForAccountDrivenEnrollmentEnabled,
		// "group_enrollment_access_enabled":                      resp.GroupEnrollmentAccessEnabled,
		// "group_attribute_name":                                 resp.GroupAttributeName,
		// "group_rdn_key":                                        resp.GroupRdnKey,
		// "group_enrollment_access_name":                         resp.GroupEnrollmentAccessName,
		// "idp_provider_type":                                    resp.IdpProviderType,
		// "other_provider_type_name":                             resp.OtherProviderTypeName,
		// "metadata_source":                                      resp.MetadataSource,
		// "session_timeout":                                      resp.SessionTimeout,
		// "device_type":                                          resp.DeviceType,
	}

	if locationInformation := resp.LocationInformation; locationInformation != (jamfpro.ComputerPrestageSubsetLocationInformation{}) {
		prestageAttributes["location_information"] = []interface{}{
			map[string]interface{}{
				"id":            locationInformation.ID,
				"username":      locationInformation.Username,
				"realname":      locationInformation.Realname,
				"phone":         locationInformation.Phone,
				"email":         locationInformation.Email,
				"room":          locationInformation.Room,
				"position":      locationInformation.Position,
				"department_id": locationInformation.DepartmentId,
				"building_id":   locationInformation.BuildingId,
				"version_lock":  locationInformation.VersionLock,
			},
		}
	}

	if purchasingInformation := resp.PurchasingInformation; purchasingInformation != (jamfpro.ComputerPrestageSubsetPurchasingInformation{}) {
		prestageAttributes["purchasing_information"] = []interface{}{
			map[string]interface{}{
				"id":                 purchasingInformation.ID,
				"leased":             purchasingInformation.Leased,
				"purchased":          purchasingInformation.Purchased,
				"apple_care_id":      purchasingInformation.AppleCareId,
				"po_number":          purchasingInformation.PONumber,
				"vendor":             purchasingInformation.Vendor,
				"purchase_price":     purchasingInformation.PurchasePrice,
				"life_expectancy":    purchasingInformation.LifeExpectancy,
				"purchasing_account": purchasingInformation.PurchasingAccount,
				"purchasing_contact": purchasingInformation.PurchasingContact,
				"lease_date":         purchasingInformation.LeaseDate,
				"po_date":            purchasingInformation.PODate,
				"warranty_date":      purchasingInformation.WarrantyDate,
				"version_lock":       purchasingInformation.VersionLock,
			},
		}
	}

	if accountSettings := resp.AccountSettings; accountSettings != (jamfpro.ComputerPrestageSubsetAccountSettings{}) {
		prestageAttributes["account_settings"] = []interface{}{
			map[string]interface{}{
				"id":                          accountSettings.ID,
				"payload_configured":          accountSettings.PayloadConfigured,
				"local_admin_account_enabled": accountSettings.LocalAdminAccountEnabled,
				"admin_username":              accountSettings.AdminUsername,
				"admin_password":              getHCLValue(d, "account_settings.0.admin_password"), // Use HCL value as server never returns the value
				"hidden_admin_account":        accountSettings.HiddenAdminAccount,
				"local_user_managed":          accountSettings.LocalUserManaged,
				"user_account_type":           accountSettings.UserAccountType,
				"version_lock":                accountSettings.VersionLock,
				"prefill_primary_account_info_feature_enabled": accountSettings.PrefillPrimaryAccountInfoFeatureEnabled,
				"prefill_type":                           accountSettings.PrefillType,
				"prefill_account_full_name":              accountSettings.PrefillAccountFullName,
				"prefill_account_user_name":              accountSettings.PrefillAccountUserName,
				"prevent_prefill_info_from_modification": accountSettings.PreventPrefillInfoFromModification,
			},
		}
	}

	if resp.OnboardingItems != nil {
		onboardingItems := make([]interface{}, len(resp.OnboardingItems))
		for i, item := range resp.OnboardingItems {
			onboardingItems[i] = map[string]interface{}{
				"self_service_entity_type": item.SelfServiceEntityType,
				"id":                       item.ID,
				"entity_id":                item.EntityId,
				"priority":                 item.Priority,
			}
		}
		prestageAttributes["onboarding_items"] = onboardingItems
	}

	for key, val := range prestageAttributes {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

// skipSetupItems converts the ComputerPrestageSubsetSkipSetupItems struct to a map
func skipSetupItems(skipSetupItems jamfpro.ComputerPrestageSubsetSkipSetupItems) map[string]interface{} {
	return map[string]interface{}{
		"biometric":                   *skipSetupItems.Biometric,
		"terms_of_address":            *skipSetupItems.TermsOfAddress,
		"file_vault":                  *skipSetupItems.FileVault,
		"icloud_diagnostics":          *skipSetupItems.ICloudDiagnostics,
		"diagnostics":                 *skipSetupItems.Diagnostics,
		"accessibility":               *skipSetupItems.Accessibility,
		"apple_id":                    *skipSetupItems.AppleID,
		"screen_time":                 *skipSetupItems.ScreenTime,
		"siri":                        *skipSetupItems.Siri,
		"display_tone":                *skipSetupItems.DisplayTone,
		"restore":                     *skipSetupItems.Restore,
		"appearance":                  *skipSetupItems.Appearance,
		"privacy":                     *skipSetupItems.Privacy,
		"payment":                     *skipSetupItems.Payment,
		"registration":                *skipSetupItems.Registration,
		"tos":                         *skipSetupItems.TOS,
		"icloud_storage":              *skipSetupItems.ICloudStorage,
		"location":                    *skipSetupItems.Location,
		"intelligence":                *skipSetupItems.Intelligence,
		"enable_lockdown_mode":        *skipSetupItems.EnableLockdownMode,
		"welcome":                     *skipSetupItems.Welcome,
		"wallpaper":                   *skipSetupItems.Wallpaper,
		"software_update":             *skipSetupItems.SoftwareUpdate,
		"additional_privacy_settings": *skipSetupItems.AdditionalPrivacySettings,
	}
}

// getHCLValue gets the value of a key from the ResourceData, either from the current state or the config.
func getHCLValue(d *schema.ResourceData, key string) interface{} {
	value, exists := d.GetOk(key)
	if !exists {
		value = d.Get(key)
	}
	return value
}
