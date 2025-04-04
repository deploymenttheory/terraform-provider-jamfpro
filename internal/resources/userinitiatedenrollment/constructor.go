package userinitiatedenrollment

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructEnrollmentSettings builds the ResourceEnrollment object from schema data
func constructEnrollmentSettings(d *schema.ResourceData) (*jamfpro.ResourceEnrollment, error) {
	resource := &jamfpro.ResourceEnrollment{
		// General settings
		InstallSingleProfile:            !d.Get("skip_certificate_installation_during_enrollment").(bool),
		RestrictReenrollment:            d.Get("restrict_reenrollment_to_authorized_users_only").(bool),
		FlushLocationInformation:        d.Get("flush_location_information").(bool),
		FlushLocationHistoryInformation: d.Get("flush_location_history_information").(bool),
		FlushPolicyHistory:              d.Get("flush_policy_history").(bool),
		FlushExtensionAttributes:        d.Get("flush_extension_attributes").(bool),
		FlushMdmCommandsOnReenroll:      d.Get("flush_mdm_commands_on_reenroll").(string),
	}

	// Set third-party signing certificate
	if v, ok := d.GetOk("third_party_signing_certificate"); ok {
		certList := v.(*schema.Set).List()
		if len(certList) > 0 {
			certMap := certList[0].(map[string]interface{})

			resource.SigningMdmProfileEnabled = certMap["enabled"].(bool)

			if resource.SigningMdmProfileEnabled {
				resource.MdmSigningCertificate = &jamfpro.ResourceEnrollmentCertificate{
					Filename:         certMap["filename"].(string),
					IdentityKeystore: certMap["identity_keystore"].(string),
					KeystorePassword: certMap["keystore_password"].(string),
				}
			}
		}
	}

	// Set computer enrollment settings
	if v, ok := d.GetOk("user_initiated_enrollment_for_computers"); ok {
		computerSettingsList := v.(*schema.Set).List()
		if len(computerSettingsList) > 0 {
			computerSettings := computerSettingsList[0].(map[string]interface{})

			resource.MacOsEnterpriseEnrollmentEnabled = computerSettings["enable_user_initiated_enrollment_for_computers"].(bool)
			resource.EnsureSshRunning = computerSettings["ensure_ssh_is_enabled"].(bool)
			resource.LaunchSelfService = computerSettings["launch_self_service_when_done"].(bool)
			resource.AccountDrivenDeviceMacosEnrollmentEnabled = computerSettings["account_driven_device_enrollment"].(bool)

			// Managed local admin account
			if adminAcct, ok := computerSettings["managed_local_administrator_account"]; ok && len(adminAcct.(*schema.Set).List()) > 0 {
				adminAcctMap := adminAcct.(*schema.Set).List()[0].(map[string]interface{})

				resource.CreateManagementAccount = adminAcctMap["create_managed_local_administrator_account"].(bool)
				resource.ManagementUsername = adminAcctMap["management_account_username"].(string)
				resource.HideManagementAccount = adminAcctMap["hide_managed_local_administrator_account"].(bool)
				resource.AllowSshOnlyManagementAccount = adminAcctMap["allow_ssh_access_for_managed_local_administrator_account_only"].(bool)
			}

			// QuickAdd package settings
			if quickAdd, ok := computerSettings["quickadd_package"]; ok && len(quickAdd.(*schema.Set).List()) > 0 {
				quickAddMap := quickAdd.(*schema.Set).List()[0].(map[string]interface{})

				resource.SignQuickAdd = quickAddMap["sign_quickadd_package"].(bool)

				if resource.SignQuickAdd {
					resource.DeveloperCertificateIdentity = &jamfpro.ResourceEnrollmentCertificate{
						Filename:         quickAddMap["filename"].(string),
						IdentityKeystore: quickAddMap["identity_keystore"].(string),
						KeystorePassword: quickAddMap["keystore_password"].(string),
					}
				}
			}
		}
	}

	// Set mobile device enrollment settings
	if v, ok := d.GetOk("user_initiated_enrollment_for_devices"); ok {
		deviceSettingsList := v.(*schema.Set).List()
		if len(deviceSettingsList) > 0 {
			deviceSettings := deviceSettingsList[0].(map[string]interface{})

			// Profile-driven enrollment
			if profileDriven, ok := deviceSettings["profile_driven_enrollment_via_url"]; ok && len(profileDriven.(*schema.Set).List()) > 0 {
				profileDrivenMap := profileDriven.(*schema.Set).List()[0].(map[string]interface{})

				resource.IosEnterpriseEnrollmentEnabled = profileDrivenMap["enable_for_institutionally_owned_devices"].(bool)
				resource.IosPersonalEnrollmentEnabled = profileDrivenMap["enable_for_personally_owned_devices"].(bool)
				resource.PersonalDeviceEnrollmentType = profileDrivenMap["personal_device_enrollment_type"].(string)
			}

			// Account-driven user enrollment
			if accountUser, ok := deviceSettings["account_driven_user_enrollment"]; ok && len(accountUser.(*schema.Set).List()) > 0 {
				accountUserMap := accountUser.(*schema.Set).List()[0].(map[string]interface{})

				resource.AccountDrivenUserEnrollmentEnabled = accountUserMap["enable_for_personally_owned_mobile_devices"].(bool)
				resource.AccountDrivenUserVisionosEnrollmentEnabled = accountUserMap["enable_for_personally_owned_vision_pro_devices"].(bool)
			}

			// Account-driven device enrollment
			if accountDevice, ok := deviceSettings["account_driven_device_enrollment"]; ok && len(accountDevice.(*schema.Set).List()) > 0 {
				accountDeviceMap := accountDevice.(*schema.Set).List()[0].(map[string]interface{})

				resource.AccountDrivenDeviceIosEnrollmentEnabled = accountDeviceMap["enable_for_institutionally_owned_mobile_devices"].(bool)
				resource.AccountDrivenDeviceVisionosEnrollmentEnabled = accountDeviceMap["enable_for_personally_owned_vision_pro_devices"].(bool)
			}
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro SMTP Server Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro User-initiated enrollment Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructEnrollmentMessaging builds enrollment language messages from schema data
func constructEnrollmentMessaging(d *schema.ResourceData, client *jamfpro.Client) ([]jamfpro.ResourceEnrollmentLanguage, error) {
	var messages []jamfpro.ResourceEnrollmentLanguage

	if v, ok := d.GetOk("messaging"); ok {
		messagingSet := v.(*schema.Set).List()

		// Get language codes for mapping from the API
		langCodes, err := getLanguageCodesMap(client)
		if err != nil {
			return nil, fmt.Errorf("failed to get language codes: %v", err)
		}

		for _, messaging := range messagingSet {
			msg := messaging.(map[string]interface{})

			// Get language name and normalize it for matching
			langName := msg["language_name"].(string)
			normalizedLangName := strings.ToLower(strings.TrimSpace(langName))

			// Get language code using normalized name
			langCode, ok := langCodes[normalizedLangName]
			if !ok {
				log.Printf("[WARN] No language code found for language name '%s', skipping", langName)
				continue
			}

			// Log the language mapping for debugging
			log.Printf("[DEBUG] Mapping language name '%s' to code '%s'", langName, langCode)

			// Create language message
			resource := jamfpro.ResourceEnrollmentLanguage{
				LanguageCode: langCode,
				Name:         langName,
				Title:        msg["page_title"].(string),

				// Login page
				LoginDescription: msg["login_page_text"].(string),
				Username:         msg["username_text"].(string),
				Password:         msg["password_text"].(string),
				LoginButton:      msg["login_button_text"].(string),

				// Device ownership
				DeviceClassDescription:           msg["device_ownership_page_text"].(string),
				DeviceClassPersonal:              msg["personal_device_button_name"].(string),
				DeviceClassPersonalDescription:   msg["personal_device_management_description"].(string),
				DeviceClassEnterprise:            msg["institutional_ownership_button_name"].(string),
				DeviceClassEnterpriseDescription: msg["institutional_device_management_description"].(string),
				DeviceClassButton:                msg["enroll_device_button_name"].(string),

				// EULA
				EnterpriseEula: msg["eula_personal_devices"].(string),
				PersonalEula:   msg["eula_institutional_devices"].(string),
				EulaButton:     msg["accept_button_text"].(string),

				// Site selection
				SiteDescription: msg["site_selection_text"].(string),

				// CA Certificate
				CertificateText:               msg["ca_certificate_installation_text"].(string),
				CertificateButton:             msg["ca_certificate_install_button_name"].(string),
				CertificateProfileName:        msg["ca_certificate_name"].(string),
				CertificateProfileDescription: msg["ca_certificate_description"].(string),

				// Institutional MDM profile
				EnterpriseText:               msg["institutional_mdm_profile_installation_text"].(string),
				EnterpriseButton:             msg["institutional_mdm_profile_install_button_name"].(string),
				EnterpriseProfileName:        msg["institutional_mdm_profile_name"].(string),
				EnterpriseProfileDescription: msg["institutional_mdm_profile_description"].(string),
				EnterprisePending:            msg["institutional_mdm_profile_pending_text"].(string),

				// Personal MDM profile
				PersonalText:               msg["personal_mdm_profile_installation_text"].(string),
				PersonalButton:             msg["personal_mdm_profile_install_button_name"].(string),
				PersonalProfileName:        msg["personal_mdm_profile_name"].(string),
				PersonalProfileDescription: msg["personal_mdm_profile_description"].(string),

				// User enrollment MDM profile
				UserEnrollmentText:               msg["user_enrollment_mdm_profile_installation_text"].(string),
				UserEnrollmentButton:             msg["user_enrollment_mdm_profile_install_button_name"].(string),
				UserEnrollmentProfileName:        msg["user_enrollment_mdm_profile_name"].(string),
				UserEnrollmentProfileDescription: msg["user_enrollment_mdm_profile_description"].(string),

				// QuickAdd package
				QuickAddText:    msg["quickadd_package_installation_text"].(string),
				QuickAddButton:  msg["quickadd_package_install_button_name"].(string),
				QuickAddName:    msg["quickadd_package_name"].(string),
				QuickAddPending: msg["quickadd_package_progress_text"].(string),

				// Completion
				CompleteMessage:        msg["enrollment_complete_text"].(string),
				FailedMessage:          msg["enrollment_failed_text"].(string),
				TryAgainButton:         msg["try_again_button_name"].(string),
				CheckNowButton:         msg["view_enrollment_status_button_name"].(string),
				CheckEnrollmentMessage: msg["view_enrollment_status_text"].(string),
				LogoutButton:           msg["log_out_button_name"].(string),
			}

			// Add an extra check to ensure we have a valid language code before proceeding
			if resource.LanguageCode == "" {
				log.Printf("[ERROR] Empty language code for language name '%s', skipping", langName)
				continue
			}

			messages = append(messages, resource)
		}

	}

	messagesJSON, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro enrollment messaging to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Enrollment Messaging JSON:\n%s\n", string(messagesJSON))

	return messages, nil
}

// In the constructDirectoryServiceGroupSettings function
// constructDirectoryServiceGroupSettings builds directory service group settings
func constructDirectoryServiceGroupSettings(d *schema.ResourceData) ([]*jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, error) {
	var resource []*jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup

	if v, ok := d.GetOk("directory_service_group_enrollment_settings"); ok {
		groupsSet := v.(*schema.Set).List()

		for _, groupData := range groupsSet {
			groupMap := groupData.(map[string]interface{})

			// Create the group with all values from config
			group := &jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup{
				GroupID:                            groupMap["directory_service_group_id"].(string),
				LdapServerID:                       groupMap["ldap_server_id"].(string),
				Name:                               groupMap["directory_service_group_name"].(string),
				SiteID:                             groupMap["site_id"].(string),
				EnterpriseEnrollmentEnabled:        groupMap["allow_group_to_enroll_institutionally_owned_devices"].(bool),
				PersonalEnrollmentEnabled:          groupMap["allow_group_to_enroll_personally_owned_devices"].(bool),
				AccountDrivenUserEnrollmentEnabled: groupMap["allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade"].(bool),
				RequireEula:                        groupMap["require_eula"].(bool),
			}

			// Preserve ID if it exists in state
			if id, ok := groupMap["id"].(string); ok && id != "" {
				group.ID = id
				log.Printf("[DEBUG] Preserving existing ID '%s' for directory service group '%s'", id, group.Name)
			}

			resource = append(resource, group)
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro User-initiated enrollment Directory Service Groups to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro User-initiated enrollment Directory Service Groups JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
