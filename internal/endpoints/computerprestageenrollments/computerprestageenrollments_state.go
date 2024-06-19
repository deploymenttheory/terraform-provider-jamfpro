// computerprestageenrollments_state.go
package computerprestageenrollments

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Prestage Enrollment information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerPrestage) diag.Diagnostics {
	var diags diag.Diagnostics

	if resource != nil {
		prestageAttributes := map[string]interface{}{
			"display_name":                          resource.DisplayName,
			"mandatory":                             resource.Mandatory,
			"mdm_removable":                         resource.MDMRemovable,
			"support_phone_number":                  resource.SupportPhoneNumber,
			"support_email_address":                 resource.SupportEmailAddress,
			"department":                            resource.Department,
			"default_prestage":                      resource.DefaultPrestage,
			"enrollment_site_id":                    resource.EnrollmentSiteId,
			"keep_existing_site_membership":         resource.KeepExistingSiteMembership,
			"keep_existing_location_information":    resource.KeepExistingLocationInformation,
			"authentication_prompt":                 resource.AuthenticationPrompt,
			"prevent_activation_lock":               resource.PreventActivationLock,
			"enable_device_based_activation_lock":   resource.EnableDeviceBasedActivationLock,
			"device_enrollment_program_instance_id": resource.DeviceEnrollmentProgramInstanceId,
			"skip_setup_items":                      resource.SkipSetupItems,
			"location_information":                  resource.LocationInformation,
			"purchasing_information":                resource.PurchasingInformation,
			"anchor_certificates":                   resource.AnchorCertificates,
			"enrollment_customization_id":           resource.EnrollmentCustomizationId,
			"language":                              resource.Language,
			"region":                                resource.Region,
			"auto_advance_setup":                    resource.AutoAdvanceSetup,
			"install_profiles_during_setup":         resource.InstallProfilesDuringSetup,
			"prestage_installed_profile_ids":        resource.PrestageInstalledProfileIds,
			"custom_package_ids":                    resource.CustomPackageIds,
			"custom_package_distribution_point_id":  resource.CustomPackageDistributionPointId,
			"enable_recovery_lock":                  resource.EnableRecoveryLock,
			"recovery_lock_password_type":           resource.RecoveryLockPasswordType,
			"recovery_lock_password":                resource.RecoveryLockPassword,
			"rotate_recovery_lock_password":         resource.RotateRecoveryLockPassword,
			"profile_uuid":                          resource.ProfileUuid,
			"site_id":                               resource.SiteId,
			"version_lock":                          resource.VersionLock,
			"account_settings":                      resource.AccountSettings,
		}

		if locationInformation := resource.LocationInformation; locationInformation != (jamfpro.ComputerPrestageSubsetLocationInformation{}) {
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
		if purchasingInformation := resource.PurchasingInformation; purchasingInformation != (jamfpro.ComputerPrestageSubsetPurchasingInformation{}) {
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
		prestageAttributes["anchor_certificates"] = resource.AnchorCertificates
		prestageAttributes["enrollment_customization_id"] = resource.EnrollmentCustomizationId
		prestageAttributes["language"] = resource.Language
		prestageAttributes["region"] = resource.Region
		prestageAttributes["auto_advance_setup"] = resource.AutoAdvanceSetup
		prestageAttributes["install_profiles_during_setup"] = resource.InstallProfilesDuringSetup
		prestageAttributes["prestage_installed_profile_ids"] = resource.PrestageInstalledProfileIds
		prestageAttributes["custom_package_ids"] = resource.CustomPackageIds
		prestageAttributes["custom_package_distribution_point_id"] = resource.CustomPackageDistributionPointId
		prestageAttributes["enable_recovery_lock"] = resource.EnableRecoveryLock
		prestageAttributes["recovery_lock_password_type"] = resource.RecoveryLockPasswordType
		prestageAttributes["recovery_lock_password"] = resource.RecoveryLockPassword
		prestageAttributes["rotate_recovery_lock_password"] = resource.RotateRecoveryLockPassword
		prestageAttributes["profile_uuid"] = resource.ProfileUuid
		prestageAttributes["site_id"] = resource.SiteId
		prestageAttributes["version_lock"] = resource.VersionLock
		if accountSettings := resource.AccountSettings; accountSettings != (jamfpro.ComputerPrestageSubsetAccountSettings{}) {
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
	}
	return diags
}
