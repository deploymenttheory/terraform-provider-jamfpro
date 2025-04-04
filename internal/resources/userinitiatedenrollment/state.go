// state.go
package userinitiatedenrollment

import (
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest enrollment settings from the Jamf Pro API
// updateState updates the Terraform state with the latest enrollment settings from the Jamf Pro API
func updateState(d *schema.ResourceData, enrollment *jamfpro.ResourceEnrollment, messages []jamfpro.ResourceEnrollmentLanguage, accessGroups []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	// --- Set General Settings, MDM Details, Third Party Cert ---
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
	for key, val := range generalSettings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set enrollment state", Detail: "Error setting " + key + "."})
		}
	}
	// MDM signing certificate details
	if enrollment.MdmSigningCertificateDetails.Subject != "" || enrollment.MdmSigningCertificateDetails.SerialNumber != "" {
		mdmCertDetails := []map[string]interface{}{{"subject": enrollment.MdmSigningCertificateDetails.Subject, "serial_number": enrollment.MdmSigningCertificateDetails.SerialNumber}}
		if err := d.Set("mdm_signing_certificate_details", mdmCertDetails); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set MDM signing certificate details"})
		}
	} else {
		// Ensure it's cleared if details become empty
		if err := d.Set("mdm_signing_certificate_details", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear MDM signing certificate details"})
		}
	}
	// Third-party signing certificate
	if enrollment.SigningMdmProfileEnabled {
		// Start with an empty set
		thirdPartyCertSet := []map[string]interface{}{}

		// Retrieve sensitive values from prior state if API doesn't return them
		if v, ok := d.GetOk("third_party_signing_certificate"); ok {
			certList := v.(*schema.Set).List()
			if len(certList) > 0 {
				certMap := certList[0].(map[string]interface{})

				// Create new map with updated non-sensitive values and preserved sensitive values
				newCertMap := map[string]interface{}{
					"enabled": enrollment.SigningMdmProfileEnabled,
				}

				// Add the filename from API if available
				if enrollment.MdmSigningCertificate != nil {
					newCertMap["filename"] = enrollment.MdmSigningCertificate.Filename
				} else if filename, ok := certMap["filename"].(string); ok {
					// Fallback to state value if API doesn't provide it
					newCertMap["filename"] = filename
				}

				// Preserve sensitive values from state
				if keystore, ok := certMap["identity_keystore"].(string); ok {
					newCertMap["identity_keystore"] = keystore
				}

				if password, ok := certMap["keystore_password"].(string); ok {
					newCertMap["keystore_password"] = password
				}

				thirdPartyCertSet = append(thirdPartyCertSet, newCertMap)
			}
		} else if enrollment.MdmSigningCertificate != nil {
			// No prior state, but cert is enabled in API, create with just what API returns
			thirdPartyCertSet = append(thirdPartyCertSet, map[string]interface{}{
				"enabled":           enrollment.SigningMdmProfileEnabled,
				"filename":          enrollment.MdmSigningCertificate.Filename,
				"identity_keystore": "", // Can't get from API
				"keystore_password": "", // Can't get from API
			})
		}

		// Set the derived state
		if err := d.Set("third_party_signing_certificate", thirdPartyCertSet); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set third-party signing certificate state"})
		}
	} else {
		// Ensure it's cleared if disabled in API
		if err := d.Set("third_party_signing_certificate", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear third-party signing certificate state"})
		}
	}

	// --- Set Computer/Device Enrollment Settings ---
	// Computer enrollment settings
	var computerEnrollment []map[string]interface{}

	if enrollment.MacOsEnterpriseEnrollmentEnabled || enrollment.CreateManagementAccount || enrollment.EnsureSshRunning || enrollment.LaunchSelfService || enrollment.SignQuickAdd || enrollment.AccountDrivenDeviceMacosEnrollmentEnabled {
		// First, check if we have existing state to preserve sensitive values
		var computerSettings map[string]interface{}

		if v, ok := d.GetOk("user_initiated_enrollment_for_computers"); ok {
			compList := v.(*schema.Set).List()
			if len(compList) > 0 {
				origCompMap := compList[0].(map[string]interface{})

				// Create a new map with updated values from API
				computerSettings = map[string]interface{}{
					"enable_user_initiated_enrollment_for_computers": enrollment.MacOsEnterpriseEnrollmentEnabled,
					"ensure_ssh_is_enabled":                          enrollment.EnsureSshRunning,
					"launch_self_service_when_done":                  enrollment.LaunchSelfService,
					"account_driven_device_enrollment":               enrollment.AccountDrivenDeviceMacosEnrollmentEnabled,
				}

				// Handle managed admin account
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

				// Handle QuickAdd package (with sensitive fields)
				if enrollment.SignQuickAdd {
					if origQuickAddList, ok := origCompMap["quickadd_package"].(*schema.Set); ok && origQuickAddList.Len() > 0 {
						quickAddList := origQuickAddList.List()
						if len(quickAddList) > 0 {
							origQuickAdd := quickAddList[0].(map[string]interface{})

							// Set the filename based on API response
							filename := ""
							if enrollment.DeveloperCertificateIdentity != nil {
								filename = enrollment.DeveloperCertificateIdentity.Filename
							}

							quickAddSettingsSet := []map[string]interface{}{
								{
									"sign_quickadd_package": enrollment.SignQuickAdd,
									"filename":              filename,
									// Preserve these sensitive fields from original state
									"identity_keystore": origQuickAdd["identity_keystore"],
									"keystore_password": origQuickAdd["keystore_password"],
								},
							}
							computerSettings["quickadd_package"] = quickAddSettingsSet
						}
					} else if enrollment.DeveloperCertificateIdentity != nil {
						// No existing state, create empty placeholders for sensitive fields
						quickAddSettingsSet := []map[string]interface{}{
							{
								"sign_quickadd_package": enrollment.SignQuickAdd,
								"filename":              enrollment.DeveloperCertificateIdentity.Filename,
								"identity_keystore":     "",
								"keystore_password":     "",
							},
						}
						computerSettings["quickadd_package"] = quickAddSettingsSet
					}
				}
			} else {
				// No existing state item, create a new one with API values
				computerSettings = createComputerEnrollmentMapFromAPI(enrollment)
			}
		} else {
			// No existing state at all, create a new one with API values
			computerSettings = createComputerEnrollmentMapFromAPI(enrollment)
		}

		computerEnrollment = []map[string]interface{}{computerSettings}
	}

	if err := d.Set("user_initiated_enrollment_for_computers", computerEnrollment); err != nil { // Set even if empty/nil
		diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set computer enrollment settings"})
	}

	// Mobile device enrollment settings
	var deviceEnrollment []map[string]interface{}
	if enrollment.IosEnterpriseEnrollmentEnabled || enrollment.IosPersonalEnrollmentEnabled || enrollment.AccountDrivenUserEnrollmentEnabled || enrollment.AccountDrivenUserVisionosEnrollmentEnabled || enrollment.AccountDrivenDeviceIosEnrollmentEnabled || enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled {
		deviceSettings := map[string]interface{}{}
		if enrollment.IosEnterpriseEnrollmentEnabled || enrollment.IosPersonalEnrollmentEnabled {
			profileSettings := []map[string]interface{}{{"enable_for_institutionally_owned_devices": enrollment.IosEnterpriseEnrollmentEnabled, "enable_for_personally_owned_devices": enrollment.IosPersonalEnrollmentEnabled, "personal_device_enrollment_type": enrollment.PersonalDeviceEnrollmentType}}
			deviceSettings["profile_driven_enrollment_via_url"] = profileSettings
		}
		if enrollment.AccountDrivenUserEnrollmentEnabled || enrollment.AccountDrivenUserVisionosEnrollmentEnabled {
			userSettings := []map[string]interface{}{{"enable_for_personally_owned_mobile_devices": enrollment.AccountDrivenUserEnrollmentEnabled, "enable_for_personally_owned_vision_pro_devices": enrollment.AccountDrivenUserVisionosEnrollmentEnabled}}
			deviceSettings["account_driven_user_enrollment"] = userSettings
		}
		if enrollment.AccountDrivenDeviceIosEnrollmentEnabled || enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled {
			accountDrivenSettings := []map[string]interface{}{{"enable_for_institutionally_owned_mobile_devices": enrollment.AccountDrivenDeviceIosEnrollmentEnabled, "enable_for_personally_owned_vision_pro_devices": enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled}}
			deviceSettings["account_driven_device_enrollment"] = accountDrivenSettings
		}
		deviceEnrollment = []map[string]interface{}{deviceSettings}
	}
	if err := d.Set("user_initiated_enrollment_for_devices", deviceEnrollment); err != nil { // Set even if empty/nil
		diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set device enrollment settings"})
	}

	// --- Set enrollment messaging configurations using flattenEnrollmentMessage ---
	if len(messages) > 0 {
		// Use TypeSet requires hashing, order doesn't matter.
		messagingConfigs := make([]interface{}, 0, len(messages)) // Build slice of interface{} for Set

		for i := range messages {
			// Pass address of element in slice to flatten function
			messagingConfigMap, err := flattenEnrollmentMessage(&messages[i])
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning, // Log as warning, state might be incomplete for this item
					Summary:  "Failed to flatten enrollment message for state",
					Detail:   fmt.Sprintf("Error processing message for language '%s' (Code: %s): %v", messages[i].Name, messages[i].LanguageCode, err),
				})
				continue // Skip this message
			}
			messagingConfigs = append(messagingConfigs, messagingConfigMap) // Add the map to the slice
		}

		log.Printf("[DEBUG] Setting 'messaging' state with %d flattened configurations.", len(messagingConfigs))
		// Set the entire list/set at once
		if err := d.Set("messaging", messagingConfigs); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set messaging configurations in state",
				Detail:   fmt.Sprintf("Error: %v", err),
			})
		}
	} else {
		log.Print("[DEBUG] No enrollment messages returned from API, clearing 'messaging' state.")
		// Explicitly set to nil or empty slice if the API returns none, ensuring drift detection
		if err := d.Set("messaging", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear messaging configurations state"})
		}
	}

	// --- Set directory service group enrollment settings using flattenDirectoryServiceGroupEnrollmentSettings ---
	if len(accessGroups) > 0 {
		groupSettingsList := make([]interface{}, 0, len(accessGroups)) // Build slice for Set

		for i := range accessGroups {
			groupConfigMap, err := flattenDirectoryServiceGroupEnrollmentSettings(&accessGroups[i])
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning, // Log as warning
					Summary:  "Failed to flatten directory service group for state",
					Detail:   fmt.Sprintf("Error processing group '%s' (API ID: %s): %v", accessGroups[i].Name, accessGroups[i].ID, err),
				})
				continue // Skip this group
			}
			groupSettingsList = append(groupSettingsList, groupConfigMap) // Add map to slice
		}

		log.Printf("[DEBUG] Setting 'directory_service_group_enrollment_settings' state with %d flattened configurations.", len(groupSettingsList))
		// Set the entire list/set at once
		if err := d.Set("directory_service_group_enrollment_settings", groupSettingsList); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set directory service group settings in state",
				Detail:   fmt.Sprintf("Error: %v", err),
			})
		}
	} else {
		log.Print("[DEBUG] No directory service groups returned from API, clearing 'directory_service_group_enrollment_settings' state.")
		if err := d.Set("directory_service_group_enrollment_settings", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear directory service group settings state"})
		}
	}

	return diags
}

// Helper function to create a new computer enrollment map from API values
func createComputerEnrollmentMapFromAPI(enrollment *jamfpro.ResourceEnrollment) map[string]interface{} {
	computerSettings := map[string]interface{}{
		"enable_user_initiated_enrollment_for_computers": enrollment.MacOsEnterpriseEnrollmentEnabled,
		"ensure_ssh_is_enabled":                          enrollment.EnsureSshRunning,
		"launch_self_service_when_done":                  enrollment.LaunchSelfService,
		"account_driven_device_enrollment":               enrollment.AccountDrivenDeviceMacosEnrollmentEnabled,
	}

	// Add managed admin account if enabled
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

	// Add QuickAdd package if signing is enabled
	if enrollment.SignQuickAdd && enrollment.DeveloperCertificateIdentity != nil {
		quickAddSettingsSet := []map[string]interface{}{
			{
				"sign_quickadd_package": enrollment.SignQuickAdd,
				"filename":              enrollment.DeveloperCertificateIdentity.Filename,
				"identity_keystore":     "", // Cannot retrieve from API
				"keystore_password":     "", // Cannot retrieve from API
			},
		}
		computerSettings["quickadd_package"] = quickAddSettingsSet
	}

	return computerSettings
}

// flattenDirectoryServiceGroupEnrollmentSettings converts an API access group struct into a map for Terraform state.
func flattenDirectoryServiceGroupEnrollmentSettings(group *jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) (map[string]interface{}, error) {
	if group == nil {
		return nil, fmt.Errorf("cannot flatten nil access group")
	}

	flatGroup := map[string]interface{}{
		// --- Core Identification ---
		"id":                           group.ID,
		"directory_service_group_name": group.Name,
		"directory_service_group_id":   group.GroupID,
		"ldap_server_id":               group.LdapServerID,

		// --- Permissions ---
		"allow_group_to_enroll_institutionally_owned_devices":                      group.EnterpriseEnrollmentEnabled,
		"allow_group_to_enroll_personally_owned_devices":                           group.PersonalEnrollmentEnabled,
		"allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade": group.AccountDrivenUserEnrollmentEnabled,

		// --- Other Settings ---
		"require_eula": group.RequireEula,
		"site_id":      group.SiteID,
	}

	return flatGroup, nil
}

// flattenEnrollmentMessage converts an API language message struct into a map for Terraform state.
func flattenEnrollmentMessage(message *jamfpro.ResourceEnrollmentLanguage) (map[string]interface{}, error) {
	if message == nil {
		return nil, fmt.Errorf("cannot flatten nil enrollment message")
	}

	flatMsg := map[string]interface{}{
		// --- Core Identification ---
		"language_code": message.LanguageCode,
		"language_name": strings.ToLower(message.Name),
		"page_title":    message.Title,

		// --- Login Page ---
		"login_page_text":   message.LoginDescription,
		"username_text":     message.Username,
		"password_text":     message.Password,
		"login_button_text": message.LoginButton,

		// --- Device Ownership ---
		"device_ownership_page_text":                  message.DeviceClassDescription,
		"personal_device_button_name":                 message.DeviceClassPersonal,
		"institutional_ownership_button_name":         message.DeviceClassEnterprise,
		"personal_device_management_description":      message.DeviceClassPersonalDescription,
		"institutional_device_management_description": message.DeviceClassEnterpriseDescription,
		"enroll_device_button_name":                   message.DeviceClassButton,

		// --- EULA ---
		"eula_personal_devices":      message.EnterpriseEula,
		"eula_institutional_devices": message.PersonalEula,
		"accept_button_text":         message.EulaButton,

		// --- Site Selection ---
		"site_selection_text": message.SiteDescription,

		// --- CA Certificate ---
		"ca_certificate_installation_text":   message.CertificateText,
		"ca_certificate_name":                message.CertificateProfileName,
		"ca_certificate_description":         message.CertificateProfileDescription,
		"ca_certificate_install_button_name": message.CertificateButton,

		// --- Institutional MDM Profile ---
		"institutional_mdm_profile_installation_text":   message.EnterpriseText,
		"institutional_mdm_profile_name":                message.EnterpriseProfileName,
		"institutional_mdm_profile_description":         message.EnterpriseProfileDescription,
		"institutional_mdm_profile_pending_text":        message.EnterprisePending,
		"institutional_mdm_profile_install_button_name": message.EnterpriseButton,

		// --- Personal MDM Profile ---
		"personal_mdm_profile_installation_text":   message.PersonalText,
		"personal_mdm_profile_name":                message.PersonalProfileName,
		"personal_mdm_profile_description":         message.PersonalProfileDescription,
		"personal_mdm_profile_install_button_name": message.PersonalButton,

		// --- User Enrollment MDM Profile ---
		"user_enrollment_mdm_profile_installation_text":   message.UserEnrollmentText,
		"user_enrollment_mdm_profile_name":                message.UserEnrollmentProfileName,
		"user_enrollment_mdm_profile_description":         message.UserEnrollmentProfileDescription,
		"user_enrollment_mdm_profile_install_button_name": message.UserEnrollmentButton,

		// --- QuickAdd Package ---
		"quickadd_package_installation_text":   message.QuickAddText,
		"quickadd_package_name":                message.QuickAddName,
		"quickadd_package_progress_text":       message.QuickAddPending,
		"quickadd_package_install_button_name": message.QuickAddButton,

		// --- Completion ---
		"enrollment_complete_text":           message.CompleteMessage,
		"enrollment_failed_text":             message.FailedMessage,
		"try_again_button_name":              message.TryAgainButton,
		"view_enrollment_status_button_name": message.CheckNowButton,
		"view_enrollment_status_text":        message.CheckEnrollmentMessage,
		"log_out_button_name":                message.LogoutButton,
	}

	return flatMsg, nil
}
