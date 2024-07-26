// computerprestageenrollments_state.go
package computerprestageenrollments

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Prestage Enrollment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceComputerPrestage) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO update this logic
	prestageAttributes := map[string]interface{}{
		"display_name":                          resp.DisplayName,
		"mandatory":                             resp.Mandatory,
		"mdm_removable":                         resp.MDMRemovable,
		"support_phone_number":                  resp.SupportPhoneNumber,
		"support_email_address":                 resp.SupportEmailAddress,
		"department":                            resp.Department,
		"default_prestage":                      resp.DefaultPrestage,
		"enrollment_site_id":                    resp.EnrollmentSiteId,
		"keep_existing_site_membership":         resp.KeepExistingSiteMembership,
		"keep_existing_location_information":    resp.KeepExistingLocationInformation,
		"authentication_prompt":                 resp.AuthenticationPrompt,
		"prevent_activation_lock":               resp.PreventActivationLock,
		"enable_device_based_activation_lock":   resp.EnableDeviceBasedActivationLock,
		"device_enrollment_program_instance_id": resp.DeviceEnrollmentProgramInstanceId,
		"skip_setup_items":                      resp.SkipSetupItems,
		"location_information":                  resp.LocationInformation,
		"purchasing_information":                resp.PurchasingInformation,
		"anchor_certificates":                   resp.AnchorCertificates,
		"enrollment_customization_id":           resp.EnrollmentCustomizationId,
		"language":                              resp.Language,
		"region":                                resp.Region,
		"auto_advance_setup":                    resp.AutoAdvanceSetup,
		"install_profiles_during_setup":         resp.InstallProfilesDuringSetup,
		"prestage_installed_profile_ids":        resp.PrestageInstalledProfileIds,
		"custom_package_ids":                    resp.CustomPackageIds,
		"custom_package_distribution_point_id":  resp.CustomPackageDistributionPointId,
		"enable_recovery_lock":                  resp.EnableRecoveryLock,
		"recovery_lock_password_type":           resp.RecoveryLockPasswordType,
		"recovery_lock_password":                resp.RecoveryLockPassword,
		"rotate_recovery_lock_password":         resp.RotateRecoveryLockPassword,
		"profile_uuid":                          resp.ProfileUuid,
		"site_id":                               resp.SiteId,
		"version_lock":                          resp.VersionLock,
		"account_settings":                      resp.AccountSettings,
	}

	if locationInformation := resp.LocationInformation; locationInformation != (jamfpro.ComputerPrestageSubsetLocationInformation{}) {
		prestageAttributes["location_information"] = []interface{}{
			map[string]interface{}{
				"username":      locationInformation.Username,
				"realname":      locationInformation.Realname,
				"phone":         locationInformation.Phone,
				"email":         locationInformation.Email,
				"room":          locationInformation.Room,
				"position":      locationInformation.Position,
				"department_id": locationInformation.DepartmentId,
				"building_id":   locationInformation.BuildingId,
				"id":            locationInformation.ID,
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
	prestageAttributes["anchor_certificates"] = resp.AnchorCertificates
	prestageAttributes["enrollment_customization_id"] = resp.EnrollmentCustomizationId
	prestageAttributes["language"] = resp.Language
	prestageAttributes["region"] = resp.Region
	prestageAttributes["auto_advance_setup"] = resp.AutoAdvanceSetup
	prestageAttributes["install_profiles_during_setup"] = resp.InstallProfilesDuringSetup
	prestageAttributes["prestage_installed_profile_ids"] = resp.PrestageInstalledProfileIds
	prestageAttributes["custom_package_ids"] = resp.CustomPackageIds
	prestageAttributes["custom_package_distribution_point_id"] = resp.CustomPackageDistributionPointId
	prestageAttributes["enable_recovery_lock"] = resp.EnableRecoveryLock
	prestageAttributes["recovery_lock_password_type"] = resp.RecoveryLockPasswordType
	prestageAttributes["recovery_lock_password"] = resp.RecoveryLockPassword
	prestageAttributes["rotate_recovery_lock_password"] = resp.RotateRecoveryLockPassword
	prestageAttributes["profile_uuid"] = resp.ProfileUuid
	prestageAttributes["site_id"] = resp.SiteId
	prestageAttributes["version_lock"] = resp.VersionLock
	if accountSettings := resp.AccountSettings; accountSettings != (jamfpro.ComputerPrestageSubsetAccountSettings{}) {
		prestageAttributes["account_settings"] = []interface{}{
			map[string]interface{}{
				"id":                          accountSettings.ID,
				"payload_configured":          accountSettings.PayloadConfigured,
				"local_admin_account_enabled": accountSettings.LocalAdminAccountEnabled,
				"admin_username":              accountSettings.AdminUsername,
				"admin_password":              accountSettings.AdminPassword,
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

	for key, val := range prestageAttributes {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
