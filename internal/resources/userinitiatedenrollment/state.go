// state.go
package userinitiatedenrollment

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateEnrollmentState updates the Terraform state with the latest enrollment settings from the Jamf Pro API
func updateState(d *schema.ResourceData, enrollment *jamfpro.ResourceEnrollment, messages []jamfpro.ResourceEnrollmentLanguage, accessGroups []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	// General settings
	generalSettings := map[string]interface{}{
		"skip_certificate_installation_during_enrollment": !enrollment.InstallSingleProfile,
		"restrict_reenrollment_to_authorized_users_only":  enrollment.RestrictReenrollment,
		"flush_location_information":                      enrollment.FlushLocationInformation,
		"flush_location_history_information":              enrollment.FlushLocationHistoryInformation,
		"flush_policy_history":                            enrollment.FlushPolicyHistory,
		"flush_extension_attributes":                      enrollment.FlushExtensionAttributes,
		"flush_mdm_commands_on_reenroll":                  enrollment.FlushMdmCommandsOnReenroll,
	}

	// Set general settings
	for key, val := range generalSettings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set enrollment state",
				Detail:   "Error setting " + key + " in state.",
			})
		}
	}

	// MDM signing certificate details
	if enrollment.MdmSigningCertificateDetails.Subject != "" || enrollment.MdmSigningCertificateDetails.SerialNumber != "" {
		mdmCertDetails := []map[string]interface{}{
			{
				"subject":       enrollment.MdmSigningCertificateDetails.Subject,
				"serial_number": enrollment.MdmSigningCertificateDetails.SerialNumber,
			},
		}

		if err := d.Set("mdm_signing_certificate_details", mdmCertDetails); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set MDM signing certificate details",
				Detail:   "Error setting MDM signing certificate details in state.",
			})
		}
	}

	// Third-party signing certificate
	if enrollment.SigningMdmProfileEnabled && enrollment.MdmSigningCertificate != nil {
		thirdPartyCert := []map[string]interface{}{
			{
				"enabled":           enrollment.SigningMdmProfileEnabled,
				"filename":          enrollment.MdmSigningCertificate.Filename,
				"identity_keystore": enrollment.MdmSigningCertificate.IdentityKeystore,
				"keystore_password": enrollment.MdmSigningCertificate.KeystorePassword,
			},
		}

		if err := d.Set("third_party_signing_certificate", thirdPartyCert); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set third-party signing certificate",
				Detail:   "Error setting third-party signing certificate in state.",
			})
		}
	}

	// Computer enrollment settings
	if enrollment.MacOsEnterpriseEnrollmentEnabled || enrollment.CreateManagementAccount ||
		enrollment.EnsureSshRunning || enrollment.LaunchSelfService || enrollment.SignQuickAdd ||
		enrollment.AccountDrivenDeviceMacosEnrollmentEnabled {

		computerSettings := map[string]interface{}{
			"enable_user_initiated_enrollment_for_computers": enrollment.MacOsEnterpriseEnrollmentEnabled,
			"ensure_ssh_is_enabled":                          enrollment.EnsureSshRunning,
			"launch_self_service_when_done":                  enrollment.LaunchSelfService,
			"account_driven_device_enrollment":               enrollment.AccountDrivenDeviceMacosEnrollmentEnabled,
		}

		// Managed local admin account settings
		if enrollment.CreateManagementAccount || enrollment.ManagementUsername != "" {
			adminAccount := []map[string]interface{}{
				{
					"create_managed_local_administrator_account":                    enrollment.CreateManagementAccount,
					"management_account_username":                                   enrollment.ManagementUsername,
					"hide_managed_local_administrator_account":                      enrollment.HideManagementAccount,
					"allow_ssh_access_for_managed_local_administrator_account_only": enrollment.AllowSshOnlyManagementAccount,
				},
			}
			computerSettings["managed_local_administrator_account"] = adminAccount
		}

		// QuickAdd package settings
		if enrollment.SignQuickAdd && enrollment.DeveloperCertificateIdentity != nil {
			quickAddSettings := []map[string]interface{}{
				{
					"sign_quickadd_package": enrollment.SignQuickAdd,
					"filename":              enrollment.DeveloperCertificateIdentity.Filename,
					"identity_keystore":     enrollment.DeveloperCertificateIdentity.IdentityKeystore,
					"keystore_password":     enrollment.DeveloperCertificateIdentity.KeystorePassword,
				},
			}
			computerSettings["quickadd_package"] = quickAddSettings
		}

		computerEnrollment := []map[string]interface{}{computerSettings}
		if err := d.Set("user_initiated_enrollment_for_computers", computerEnrollment); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set computer enrollment settings",
				Detail:   "Error setting computer enrollment settings in state.",
			})
		}
	}

	// Mobile device enrollment settings
	if enrollment.IosEnterpriseEnrollmentEnabled || enrollment.IosPersonalEnrollmentEnabled ||
		enrollment.AccountDrivenUserEnrollmentEnabled || enrollment.AccountDrivenUserVisionosEnrollmentEnabled ||
		enrollment.AccountDrivenDeviceIosEnrollmentEnabled || enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled {

		deviceSettings := map[string]interface{}{}

		// Profile-driven enrollment settings
		if enrollment.IosEnterpriseEnrollmentEnabled || enrollment.IosPersonalEnrollmentEnabled {
			profileSettings := []map[string]interface{}{
				{
					"enable_for_institutionally_owned_devices": enrollment.IosEnterpriseEnrollmentEnabled,
					"enable_for_personally_owned_devices":      enrollment.IosPersonalEnrollmentEnabled,
					"personal_device_enrollment_type":          enrollment.PersonalDeviceEnrollmentType,
				},
			}
			deviceSettings["profile_driven_enrollment_via_url"] = profileSettings
		}

		// Account-driven user enrollment settings
		if enrollment.AccountDrivenUserEnrollmentEnabled || enrollment.AccountDrivenUserVisionosEnrollmentEnabled {
			userSettings := []map[string]interface{}{
				{
					"enable_for_personally_owned_mobile_devices":     enrollment.AccountDrivenUserEnrollmentEnabled,
					"enable_for_personally_owned_vision_pro_devices": enrollment.AccountDrivenUserVisionosEnrollmentEnabled,
				},
			}
			deviceSettings["account_driven_user_enrollment"] = userSettings
		}

		// Account-driven device enrollment settings
		if enrollment.AccountDrivenDeviceIosEnrollmentEnabled || enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled {
			accountDrivenSettings := []map[string]interface{}{
				{
					"enable_for_institutionally_owned_mobile_devices": enrollment.AccountDrivenDeviceIosEnrollmentEnabled,
					"enable_for_personally_owned_vision_pro_devices":  enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled,
				},
			}
			deviceSettings["account_driven_device_enrollment"] = accountDrivenSettings
		}

		deviceEnrollment := []map[string]interface{}{deviceSettings}
		if err := d.Set("user_initiated_enrollment_for_devices", deviceEnrollment); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set device enrollment settings",
				Detail:   "Error setting device enrollment settings in state.",
			})
		}
	}

	// Set enrollment messaging configurations
	if len(messages) > 0 {
		messagingConfigs := make([]map[string]interface{}, 0, len(messages))

		for _, message := range messages {
			messagingConfig, err := flattenEnrollmentMessage(&message)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to flatten enrollment message",
					Detail:   "Error processing message for language " + message.Name,
				})
				continue
			}
			messagingConfigs = append(messagingConfigs, messagingConfig)
		}

		if err := d.Set("messaging", messagingConfigs); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set messaging configurations",
				Detail:   "Error setting messaging configurations in state.",
			})
		}
	}

	// Set directory service group enrollment settings
	if len(accessGroups) > 0 {
		groupSettings := flattenDirectoryServiceGroupEnrollmentSettings(accessGroups)
		if err := d.Set("directory_service_group_enrollment_settings", groupSettings); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set directory service group settings",
				Detail:   "Error setting directory service group enrollment settings in state.",
			})
		}
	}

	return diags
}
