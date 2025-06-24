// state.go
package user_initiated_enrollment_settings

import (
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest enrollment settings from the Jamf Pro API
func updateState(d *schema.ResourceData, enrollment *jamfpro.ResourceEnrollment, messages []jamfpro.ResourceEnrollmentLanguage, accessGroups []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] updateState: Starting state update from API data.")

	// --- Set General Settings ---
	generalSettings := map[string]interface{}{
		"skip_certificate_installation_during_enrollment": enrollment.InstallSingleProfile,
		"restrict_reenrollment_to_authorized_users_only":  enrollment.RestrictReenrollment,
		"flush_location_information":                      enrollment.FlushLocationInformation,
		"flush_location_history_information":              enrollment.FlushLocationHistoryInformation,
		"flush_policy_history":                            enrollment.FlushPolicyHistory,
		"flush_extension_attributes":                      enrollment.FlushExtensionAttributes,
		"flush_software_update_plans":                     enrollment.FlushSoftwareUpdatePlans,
		"flush_mdm_commands_on_reenroll":                  enrollment.FlushMdmCommandsOnReenroll,
	}
	for key, val := range generalSettings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set general setting", Detail: "Error setting " + key + ": " + err.Error()})
		}
	}

	// --- Set MDM signing certificate details ---
	if enrollment.MdmSigningCertificateDetails.Subject != "" || enrollment.MdmSigningCertificateDetails.SerialNumber != "" {
		mdmCertDetails := []map[string]interface{}{
			{
				"subject":       enrollment.MdmSigningCertificateDetails.Subject,
				"serial_number": enrollment.MdmSigningCertificateDetails.SerialNumber,
			},
		}
		if err := d.Set("mdm_signing_certificate_details", mdmCertDetails); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set MDM signing certificate details", Detail: err.Error()})
		}
	} else {
		if err := d.Set("mdm_signing_certificate_details", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear MDM signing certificate details", Detail: err.Error()})
		}
	}

	// --- Set Third-party signing certificate ---
	if enrollment.SigningMdmProfileEnabled {
		log.Printf("[DEBUG] updateState: third_party_signing_certificate is enabled in API response.")
		configFilename := ""
		configKeystore := ""
		configPassword := ""

		if v, ok := d.GetOk("third_party_signing_certificate"); ok {
			if list, okList := v.([]interface{}); okList && len(list) > 0 {
				if certMap, okMap := list[0].(map[string]interface{}); okMap {
					if filenameVal, ok := certMap["filename"].(string); ok {
						configFilename = filenameVal
					}
					if keystoreVal, ok := certMap["identity_keystore"].(string); ok {
						configKeystore = keystoreVal
					}
					if passwordVal, ok := certMap["keystore_password"].(string); ok {
						configPassword = passwordVal
					}
					log.Printf("[DEBUG] updateState: Read third_party_signing_certificate from state - filename: '%s', keystore present: %t, password present: %t", configFilename, configKeystore != "", configPassword != "")
				}
			}
		} else {
			log.Printf("[WARN] updateState: third_party_signing_certificate enabled in API, but block not found in current state. Filename/sensitive data might be lost if not reapplied.")
		}

		certMapForState := map[string]interface{}{
			"enabled":           true,           // Use API value for enablement
			"filename":          configFilename, // USE VALUE FROM STATE
			"identity_keystore": configKeystore, // USE VALUE FROM STATE
			"keystore_password": configPassword, // USE VALUE FROM STATE
		}

		if err := d.Set("third_party_signing_certificate", []interface{}{certMapForState}); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set third-party signing certificate state", Detail: err.Error()})
		}
	} else {
		log.Printf("[DEBUG] updateState: third_party_signing_certificate is disabled in API response. Clearing state.")
		if err := d.Set("third_party_signing_certificate", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear third-party signing certificate state", Detail: err.Error()})
		}
	}

	// --- Set Computer Enrollment Settings ---
	computerBlockEnabledInAPI := enrollment.MacOsEnterpriseEnrollmentEnabled ||
		enrollment.CreateManagementAccount ||
		enrollment.EnsureSshRunning ||
		enrollment.LaunchSelfService ||
		enrollment.SignQuickAdd ||
		enrollment.AccountDrivenDeviceMacosEnrollmentEnabled

	var computerEnrollmentStateValue []interface{} // Value to set for the TypeSet

	if computerBlockEnabledInAPI {
		log.Printf("[DEBUG] updateState: user_initiated_enrollment_for_computers block is considered enabled based on API flags.")
		computerSettings := map[string]interface{}{
			"enable_user_initiated_enrollment_for_computers": enrollment.MacOsEnterpriseEnrollmentEnabled,
			"ensure_ssh_is_enabled":                          enrollment.EnsureSshRunning,
			"launch_self_service_when_done":                  enrollment.LaunchSelfService,
			"account_driven_device_enrollment":               enrollment.AccountDrivenDeviceMacosEnrollmentEnabled,
		}

		// Handle managed admin account
		if enrollment.CreateManagementAccount {
			log.Printf("[DEBUG] updateState: managed_local_administrator_account is enabled in API response.")
			adminAccountMap := map[string]interface{}{
				"create_managed_local_administrator_account":                    enrollment.CreateManagementAccount,
				"management_account_username":                                   enrollment.ManagementUsername,
				"hide_managed_local_administrator_account":                      enrollment.HideManagementAccount,
				"allow_ssh_access_for_managed_local_administrator_account_only": enrollment.AllowSshOnlyManagementAccount,
			}
			computerSettings["managed_local_administrator_account"] = []interface{}{adminAccountMap} // List for TypeSet
		} else {
			log.Printf("[DEBUG] updateState: managed_local_administrator_account is disabled in API response. Clearing sub-block state.")
			computerSettings["managed_local_administrator_account"] = nil
		}

		// Handle QuickAdd package
		if enrollment.SignQuickAdd {
			log.Printf("[DEBUG] updateState: quickadd_package signing is enabled in API response.")
			qaConfigFilename := ""
			qaConfigKeystore := ""
			qaConfigPassword := ""

			// Read existing values from state/config for quickadd
			if vComp, okComp := d.GetOk("user_initiated_enrollment_for_computers"); okComp {
				if compSet, okSet := vComp.(*schema.Set); okSet {
					compList := compSet.List()
					if len(compList) > 0 {
						if compMap, okMap := compList[0].(map[string]interface{}); okMap {
							if qaList, okQaList := compMap["quickadd_package"].([]interface{}); okQaList && len(qaList) > 0 {
								if qaMap, okQaMap := qaList[0].(map[string]interface{}); okQaMap {
									if filenameVal, ok := qaMap["filename"].(string); ok {
										qaConfigFilename = filenameVal // READ FROM STATE
									}
									if keystoreVal, ok := qaMap["identity_keystore"].(string); ok {
										qaConfigKeystore = keystoreVal // READ FROM STATE
									}
									if passwordVal, ok := qaMap["keystore_password"].(string); ok {
										qaConfigPassword = passwordVal // READ FROM STATE
									}
									log.Printf("[DEBUG] updateState: Read quickadd_package from state - filename: '%s', keystore present: %t, password present: %t", qaConfigFilename, qaConfigKeystore != "", qaConfigPassword != "")
								}
							}
						}
					}
				}
			} else {
				log.Printf("[WARN] updateState: quickadd_package enabled in API, but parent block not found in current state. Filename/sensitive data might be lost.")
			}

			quickAddMapForState := map[string]interface{}{
				"sign_quickadd_package": true,             // Use API value for enablement
				"filename":              qaConfigFilename, // USE VALUE FROM STATE
				"identity_keystore":     qaConfigKeystore, // USE VALUE FROM STATE
				"keystore_password":     qaConfigPassword, // USE VALUE FROM STATE
			}
			computerSettings["quickadd_package"] = []interface{}{quickAddMapForState}
		} else {
			log.Printf("[DEBUG] updateState: quickadd_package signing is disabled in API response. Clearing sub-block state.")
			computerSettings["quickadd_package"] = nil
		}

		// Final value is a list containing the single map (for TypeSet)
		computerEnrollmentStateValue = []interface{}{computerSettings}

	} else {
		log.Printf("[DEBUG] updateState: user_initiated_enrollment_for_computers block is considered disabled based on API flags. Clearing state.")
		computerEnrollmentStateValue = nil
	}

	if err := d.Set("user_initiated_enrollment_for_computers", computerEnrollmentStateValue); err != nil {
		diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set computer enrollment settings", Detail: err.Error()})
	}

	// --- Set Mobile Device Enrollment Settings ---
	deviceSettings := map[string]interface{}{}

	profileSettingsMap := map[string]interface{}{
		"enable_for_institutionally_owned_devices": enrollment.IosEnterpriseEnrollmentEnabled,
		"enable_for_personally_owned_devices":      enrollment.IosPersonalEnrollmentEnabled,
		"personal_device_enrollment_type":          enrollment.PersonalDeviceEnrollmentType,
	}
	deviceSettings["profile_driven_enrollment_via_url"] = []interface{}{profileSettingsMap}

	userSettingsMap := map[string]interface{}{
		"enable_for_personally_owned_mobile_devices":     enrollment.AccountDrivenUserEnrollmentEnabled,
		"enable_for_personally_owned_vision_pro_devices": enrollment.AccountDrivenUserVisionosEnrollmentEnabled,
	}
	deviceSettings["account_driven_user_enrollment"] = []interface{}{userSettingsMap}

	accountDrivenSettingsMap := map[string]interface{}{
		"enable_for_institutionally_owned_mobile_devices": enrollment.AccountDrivenDeviceIosEnrollmentEnabled,
		"enable_for_personally_owned_vision_pro_devices":  enrollment.AccountDrivenDeviceVisionosEnrollmentEnabled,
	}
	deviceSettings["account_driven_device_enrollment"] = []interface{}{accountDrivenSettingsMap}

	deviceEnrollmentStateValue := []interface{}{deviceSettings}

	if err := d.Set("user_initiated_enrollment_for_devices", deviceEnrollmentStateValue); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set device enrollment settings",
			Detail:   err.Error(),
		})
	}

	// --- Set Enrollment Messaging Configurations ---
	if len(messages) > 0 {
		messagingConfigs := make([]interface{}, 0, len(messages))
		for i := range messages {
			messagingConfigMap, err := handleEnrollmentMessage(&messages[i])
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to flatten enrollment message for state",
					Detail:   fmt.Sprintf("Error processing message for language '%s' (Code: %s): %v", messages[i].Name, messages[i].LanguageCode, err),
				})
				continue
			}
			messagingConfigs = append(messagingConfigs, messagingConfigMap)
		}
		log.Printf("[DEBUG] Setting 'messaging' state with %d flattened configurations from API.", len(messagingConfigs))
		if err := d.Set("messaging", messagingConfigs); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set messaging configurations in state", Detail: err.Error()})
		}
	} else {
		log.Print("[DEBUG] No enrollment messages returned from API, clearing 'messaging' state.")
		if err := d.Set("messaging", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear messaging configurations state", Detail: err.Error()})
		}
	}

	// --- Set Directory Service Group Enrollment Settings ---
	if len(accessGroups) > 0 {
		groupSettingsList := make([]interface{}, 0, len(accessGroups))

		// First, check if directory_service_group_enrollment_settings is in the config
		hasConfiguredGroups := false
		if _, ok := d.GetOk("directory_service_group_enrollment_settings"); ok {
			hasConfiguredGroups = true
		}

		for i := range accessGroups {
			groupConfigMap, err := handleDirectoryServiceGroupEnrollmentSettings(&accessGroups[i])
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to flatten directory service group for state",
					Detail:   fmt.Sprintf("Error processing group '%s' (API ID: %s): %v", accessGroups[i].Name, accessGroups[i].ID, err),
				})
				continue
			}

			if !hasConfiguredGroups && groupConfigMap["id"] == "1" {
				groupSettingsList = append(groupSettingsList, groupConfigMap)
			} else if hasConfiguredGroups {
				groupSettingsList = append(groupSettingsList, groupConfigMap)
			}
		}

		log.Printf("[DEBUG] Setting 'directory_service_group_enrollment_settings' state with %d flattened configurations from API.", len(groupSettingsList))
		if err := d.Set("directory_service_group_enrollment_settings", groupSettingsList); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to set directory service group settings in state", Detail: err.Error()})
		}
	} else {
		log.Print("[DEBUG] No directory service groups returned from API, clearing 'directory_service_group_enrollment_settings' state.")
		if err := d.Set("directory_service_group_enrollment_settings", nil); err != nil {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: "Failed to clear directory service group settings state", Detail: err.Error()})
		}
	}

	log.Printf("[DEBUG] updateState: Finished state update.")
	return diags
}

// handleDirectoryServiceGroupEnrollmentSettings flattens a directory service group enrollment setting struct into a map for Terraform state.
func handleDirectoryServiceGroupEnrollmentSettings(group *jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) (map[string]interface{}, error) {
	if group == nil {
		return nil, fmt.Errorf("cannot flatten nil access group")
	}

	flatGroup := map[string]interface{}{
		"id":                           group.ID,
		"directory_service_group_name": group.Name,
		"directory_service_group_id":   group.GroupID,
		"ldap_server_id":               group.LdapServerID,
		"allow_group_to_enroll_institutionally_owned_devices":                      group.EnterpriseEnrollmentEnabled,
		"allow_group_to_enroll_personally_owned_devices":                           group.PersonalEnrollmentEnabled,
		"allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade": group.AccountDrivenUserEnrollmentEnabled,
		"require_eula": group.RequireEula,
		"site_id":      group.SiteID,
	}

	return flatGroup, nil
}

// handleEnrollmentMessage flattens an enrollment language message struct into a map for Terraform state.
func handleEnrollmentMessage(message *jamfpro.ResourceEnrollmentLanguage) (map[string]interface{}, error) {
	if message == nil {
		return nil, fmt.Errorf("cannot flatten nil enrollment message")
	}

	flatMsg := map[string]interface{}{
		"language_code":                                   message.LanguageCode,
		"language_name":                                   strings.ToLower(message.Name), // Store consistent lowercase name
		"page_title":                                      message.Title,
		"login_page_text":                                 message.LoginDescription,
		"username_text":                                   message.Username,
		"password_text":                                   message.Password,
		"login_button_text":                               message.LoginButton,
		"device_ownership_page_text":                      message.DeviceClassDescription,
		"personal_device_button_name":                     message.DeviceClassPersonal,
		"institutional_ownership_button_name":             message.DeviceClassEnterprise,
		"personal_device_management_description":          message.DeviceClassPersonalDescription,
		"institutional_device_management_description":     message.DeviceClassEnterpriseDescription,
		"enroll_device_button_name":                       message.DeviceClassButton,
		"eula_personal_devices":                           message.EnterpriseEula,
		"eula_institutional_devices":                      message.PersonalEula,
		"accept_button_text":                              message.EulaButton,
		"site_selection_text":                             message.SiteDescription,
		"ca_certificate_installation_text":                message.CertificateText,
		"ca_certificate_name":                             message.CertificateProfileName,
		"ca_certificate_description":                      message.CertificateProfileDescription,
		"ca_certificate_install_button_name":              message.CertificateButton,
		"institutional_mdm_profile_installation_text":     message.EnterpriseText,
		"institutional_mdm_profile_name":                  message.EnterpriseProfileName,
		"institutional_mdm_profile_description":           message.EnterpriseProfileDescription,
		"institutional_mdm_profile_pending_text":          message.EnterprisePending,
		"institutional_mdm_profile_install_button_name":   message.EnterpriseButton,
		"personal_mdm_profile_installation_text":          message.PersonalText,
		"personal_mdm_profile_name":                       message.PersonalProfileName,
		"personal_mdm_profile_description":                message.PersonalProfileDescription,
		"personal_mdm_profile_install_button_name":        message.PersonalButton,
		"user_enrollment_mdm_profile_installation_text":   message.UserEnrollmentText,
		"user_enrollment_mdm_profile_name":                message.UserEnrollmentProfileName,
		"user_enrollment_mdm_profile_description":         message.UserEnrollmentProfileDescription,
		"user_enrollment_mdm_profile_install_button_name": message.UserEnrollmentButton,
		"quickadd_package_installation_text":              message.QuickAddText,
		"quickadd_package_name":                           message.QuickAddName,
		"quickadd_package_progress_text":                  message.QuickAddPending,
		"quickadd_package_install_button_name":            message.QuickAddButton,
		"enrollment_complete_text":                        message.CompleteMessage,
		"enrollment_failed_text":                          message.FailedMessage,
		"try_again_button_name":                           message.TryAgainButton,
		"view_enrollment_status_button_name":              message.CheckNowButton,
		"view_enrollment_status_text":                     message.CheckEnrollmentMessage,
		"log_out_button_name":                             message.LogoutButton,
	}

	return flatMsg, nil
}
