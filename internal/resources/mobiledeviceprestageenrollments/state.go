// mobiledeviceprestageenrollments_state.go
package mobiledeviceprestageenrollments

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Mobile Device Prestage information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDevicePrestage) diag.Diagnostics {
	var diags diag.Diagnostics

	prestageAttributes := map[string]interface{}{
		"id":                                           resp.ID,
		"display_name":                                 resp.DisplayName,
		"mandatory":                                    resp.Mandatory,
		"mdm_removable":                                resp.MdmRemovable,
		"support_phone_number":                         resp.SupportPhoneNumber,
		"support_email_address":                        resp.SupportEmailAddress,
		"department":                                   resp.Department,
		"default_prestage":                             resp.DefaultPrestage,
		"enrollment_site_id":                           resp.EnrollmentSiteID,
		"keep_existing_site_membership":                resp.KeepExistingSiteMembership,
		"keep_existing_location_information":           resp.KeepExistingLocationInformation,
		"require_authentication":                       resp.RequireAuthentication,
		"authentication_prompt":                        resp.AuthenticationPrompt,
		"prevent_activation_lock":                      resp.PreventActivationLock,
		"enable_device_based_activation_lock":          resp.EnableDeviceBasedActivationLock,
		"device_enrollment_program_instance_id":        resp.DeviceEnrollmentProgramInstanceID,
		"skip_setup_items":                             []interface{}{skipSetupItems(resp.SkipSetupItems)},
		"anchor_certificates":                          resp.AnchorCertificates,
		"enrollment_customization_id":                  resp.EnrollmentCustomizationID,
		"language":                                     resp.Language,
		"region":                                       resp.Region,
		"auto_advance_setup":                           resp.AutoAdvanceSetup,
		"allow_pairing":                                resp.AllowPairing,
		"multi_user":                                   resp.MultiUser,
		"supervised":                                   resp.Supervised,
		"maximum_shared_accounts":                      resp.MaximumSharedAccounts,
		"configure_device_before_setup_assistant":      resp.ConfigureDeviceBeforeSetupAssistant,
		"send_timezone":                                resp.SendTimezone,
		"timezone":                                     resp.Timezone,
		"storage_quota_size_megabytes":                 resp.StorageQuotaSizeMegabytes,
		"use_storage_quota_size":                       resp.UseStorageQuotaSize,
		"temporary_session_only":                       resp.TemporarySessionOnly,
		"enforce_temporary_session_timeout":            resp.EnforceTemporarySessionTimeout,
		"temporary_session_timeout_seconds":            resp.TemporarySessionTimeout,
		"enforce_user_session_timeout":                 resp.EnforceUserSessionTimeout,
		"user_session_timeout":                         resp.UserSessionTimeout,
		"profile_uuid":                                 resp.ProfileUuid,
		"site_id":                                      resp.SiteId,
		"version_lock":                                 resp.VersionLock,
		"prestage_minimum_os_target_version_type_ios":  resp.PrestageMinimumOsTargetVersionTypeIos,
		"minimum_os_specific_version_ios":              resp.MinimumOsSpecificVersionIos,
		"prestage_minimum_os_target_version_type_ipad": resp.PrestageMinimumOsTargetVersionTypeIpad,
		"minimum_os_specific_version_ipad":             resp.MinimumOsSpecificVersionIpad,
		"rts_enabled":                                  resp.RTSEnabled,
		"rts_config_profile_id":                        resp.RTSConfigProfileId,
	}

	if locationInformation := resp.LocationInformation; locationInformation != (jamfpro.MobileDevicePrestageSubsetLocationInformation{}) {
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

	if purchasingInformation := resp.PurchasingInformation; purchasingInformation != (jamfpro.MobileDevicePrestageSubsetPurchasingInformation{}) {
		prestageAttributes["purchasing_information"] = []interface{}{
			map[string]interface{}{
				"id":                 purchasingInformation.ID,
				"leased":             purchasingInformation.Leased,
				"purchased":          purchasingInformation.Purchased,
				"apple_care_id":      purchasingInformation.AppleCareId,
				"po_number":          purchasingInformation.PoNumber,
				"vendor":             purchasingInformation.Vendor,
				"purchase_price":     purchasingInformation.PurchasePrice,
				"life_expectancy":    purchasingInformation.LifeExpectancy,
				"purchasing_account": purchasingInformation.PurchasingAccount,
				"purchasing_contact": purchasingInformation.PurchasingContact,
				"lease_date":         purchasingInformation.LeaseDate,
				"po_date":            purchasingInformation.PoDate,
				"warranty_date":      purchasingInformation.WarrantyDate,
				"version_lock":       purchasingInformation.VersionLock,
			},
		}
	}

	if names := resp.Names; names.AssignNamesUsing != "" || names.DeviceNamePrefix != "" || names.DeviceNameSuffix != "" || names.SingleDeviceName != "" || *names.ManageNames || *names.DeviceNamingConfigured {
		namesMap := map[string]interface{}{
			"assign_names_using":       names.AssignNamesUsing,
			"device_name_prefix":       names.DeviceNamePrefix,
			"device_name_suffix":       names.DeviceNameSuffix,
			"single_device_name":       names.SingleDeviceName,
			"manage_names":             names.ManageNames,
			"device_naming_configured": names.DeviceNamingConfigured,
		}

		if len(names.PrestageDeviceNames) > 0 {
			deviceNames := make([]interface{}, len(names.PrestageDeviceNames))
			for i, name := range names.PrestageDeviceNames {
				deviceNames[i] = map[string]interface{}{
					"id":          name.ID,
					"device_name": name.DeviceName,
					"used":        name.Used,
				}
			}
			namesMap["prestage_device_names"] = deviceNames
		}

		prestageAttributes["names"] = []interface{}{namesMap}
	}

	for key, val := range prestageAttributes {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

// skipSetupItems converts the MobileDevicePrestageSubsetSkipSetupItems struct to a map
func skipSetupItems(skipSetupItems jamfpro.MobileDevicePrestageSubsetSkipSetupItems) map[string]interface{} {
	return map[string]interface{}{
		"location":                skipSetupItems.Location,
		"privacy":                 skipSetupItems.Privacy,
		"biometric":               skipSetupItems.Biometric,
		"software_update":         skipSetupItems.SoftwareUpdate,
		"diagnostics":             skipSetupItems.Diagnostics,
		"imessage_and_facetime":   skipSetupItems.IMessageAndFaceTime,
		"intelligence":            skipSetupItems.Intelligence,
		"tv_room":                 skipSetupItems.TVRoom,
		"passcode":                skipSetupItems.Passcode,
		"sim_setup":               skipSetupItems.SIMSetup,
		"screen_time":             skipSetupItems.ScreenTime,
		"restore_completed":       skipSetupItems.RestoreCompleted,
		"tv_provider_sign_in":     skipSetupItems.TVProviderSignIn,
		"siri":                    skipSetupItems.Siri,
		"restore":                 skipSetupItems.Restore,
		"screen_saver":            skipSetupItems.ScreenSaver,
		"home_button_sensitivity": skipSetupItems.HomeButtonSensitivity,
		"cloud_storage":           skipSetupItems.CloudStorage,
		"action_button":           skipSetupItems.ActionButton,
		"transfer_data":           skipSetupItems.TransferData,
		"enable_lockdown_mode":    skipSetupItems.EnableLockdownMode,
		"zoom":                    skipSetupItems.Zoom,
		"preferred_language":      skipSetupItems.PreferredLanguage,
		"voice_selection":         skipSetupItems.VoiceSelection,
		"tv_home_screen_sync":     skipSetupItems.TVHomeScreenSync,
		"safety":                  skipSetupItems.Safety,
		"terms_of_address":        skipSetupItems.TermsOfAddress,
		"express_language":        skipSetupItems.ExpressLanguage,
		"camera_button":           skipSetupItems.CameraButton,
		"apple_id":                skipSetupItems.AppleID,
		"display_tone":            skipSetupItems.DisplayTone,
		"watch_migration":         skipSetupItems.WatchMigration,
		"update_completed":        skipSetupItems.UpdateCompleted,
		"appearance":              skipSetupItems.Appearance,
		"android":                 skipSetupItems.Android,
		"payment":                 skipSetupItems.Payment,
		"onboarding":              skipSetupItems.OnBoarding,
		"tos":                     skipSetupItems.TOS,
		"welcome":                 skipSetupItems.Welcome,
		"safety_and_handling":     skipSetupItems.SafetyAndHandling,
		"tap_to_setup":            skipSetupItems.TapToSetup,
	}
}
