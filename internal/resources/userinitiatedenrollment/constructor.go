package userinitiatedenrollment

import (
	"fmt"
	"log"
	"strconv"

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

	return resource, nil
}

// constructEnrollmentMessaging builds enrollment language messages from schema data
func constructEnrollmentMessaging(d *schema.ResourceData) ([]jamfpro.ResourceEnrollmentLanguage, error) {
	var messages []jamfpro.ResourceEnrollmentLanguage

	if v, ok := d.GetOk("messaging"); ok {
		messagingSet := v.(*schema.Set).List()

		// Get language codes for mapping
		langCodes, err := getLanguageCodesMap()
		if err != nil {
			return nil, fmt.Errorf("failed to get language codes: %v", err)
		}

		for _, messaging := range messagingSet {
			msg := messaging.(map[string]interface{})

			// Get language code for this language name
			langName := msg["language_name"].(string)
			langCode, ok := langCodes[langName]
			if !ok {
				log.Printf("[WARN] No language code found for language name '%s', skipping", langName)
				continue
			}

			// Create language message
			message := jamfpro.ResourceEnrollmentLanguage{
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

			messages = append(messages, message)
		}
	}

	return messages, nil
}

// constructDirectoryServiceGroupSettings builds directory service group settings
func constructDirectoryServiceGroupSettings(d *schema.ResourceData) ([]*jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, error) {
	var groups []*jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup

	if v, ok := d.GetOk("directory_service_group_enrollment_settings"); ok {
		groupsSet := v.(*schema.Set).List()

		for _, groupData := range groupsSet {
			groupMap := groupData.(map[string]interface{})

			siteID := ""
			if v, ok := groupMap["site_id"]; ok && v.(int) > 0 {
				siteID = strconv.Itoa(v.(int))
			}

			group := &jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup{
				GroupID:                            groupMap["directory_service_group_id"].(string),
				LdapServerID:                       groupMap["ldap_server_id"].(string),
				Name:                               groupMap["directory_service_group_name"].(string),
				SiteID:                             siteID,
				EnterpriseEnrollmentEnabled:        groupMap["allow_group_to_enroll_institutionally_owned_devices"].(bool),
				PersonalEnrollmentEnabled:          groupMap["allow_group_to_enroll_personally_owned_devices"].(bool),
				AccountDrivenUserEnrollmentEnabled: groupMap["allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade"].(bool),
				RequireEula:                        groupMap["require_eula"].(bool),
			}

			groups = append(groups, group)
		}
	}

	return groups, nil
}

// updateEnrollmentState updates the schema with enrollment data from the API
func updateEnrollmentState(d *schema.ResourceData, enrollment *jamfpro.ResourceEnrollment) error {
	// General settings
	if err := d.Set("skip_certificate_installation_during_enrollment", !enrollment.InstallSingleProfile); err != nil {
		return err
	}
	if err := d.Set("restrict_reenrollment_to_authorized_users_only", enrollment.RestrictReenrollment); err != nil {
		return err
	}
	if err := d.Set("flush_location_information", enrollment.FlushLocationInformation); err != nil {
		return err
	}
	if err := d.Set("flush_location_history_information", enrollment.FlushLocationHistoryInformation); err != nil {
		return err
	}
	if err := d.Set("flush_policy_history", enrollment.FlushPolicyHistory); err != nil {
		return err
	}
	if err := d.Set("flush_extension_attributes", enrollment.FlushExtensionAttributes); err != nil {
		return err
	}
	if err := d.Set("flush_mdm_commands_on_reenroll", enrollment.FlushMdmCommandsOnReenroll); err != nil {
		return err
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
			return err
		}
	}

	// Third-party signing certificate
	thirdPartyCert := []map[string]interface{}{}
	if enrollment.SigningMdmProfileEnabled && enrollment.MdmSigningCertificate != nil {
		thirdPartyCert = append(thirdPartyCert, map[string]interface{}{
			"enabled":           enrollment.SigningMdmProfileEnabled,
			"filename":          enrollment.MdmSigningCertificate.Filename,
			"identity_keystore": enrollment.MdmSigningCertificate.IdentityKeystore,
			"keystore_password": enrollment.MdmSigningCertificate.KeystorePassword,
		})
	}

	if len(thirdPartyCert) > 0 {
		if err := d.Set("third_party_signing_certificate", thirdPartyCert); err != nil {
			return err
		}
	}

	// Computer enrollment settings
	computerEnrollment := []map[string]interface{}{}

	// Only set if computer enrollment is configured
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

		computerEnrollment = append(computerEnrollment, computerSettings)
	}

	if len(computerEnrollment) > 0 {
		if err := d.Set("user_initiated_enrollment_for_computers", computerEnrollment); err != nil {
			return err
		}
	}

	// Mobile device enrollment settings
	deviceEnrollment := []map[string]interface{}{}

	// Only set if any mobile device enrollment is configured
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

		deviceEnrollment = append(deviceEnrollment, deviceSettings)
	}

	if len(deviceEnrollment) > 0 {
		if err := d.Set("user_initiated_enrollment_for_devices", deviceEnrollment); err != nil {
			return err
		}
	}

	return nil
}

// flattenEnrollmentMessage converts an enrollment message to a map for the schema
func flattenEnrollmentMessage(message *jamfpro.ResourceEnrollmentLanguage) (map[string]interface{}, error) {
	if message == nil {
		return nil, fmt.Errorf("message is nil")
	}

	// If language code is empty, try to get it from the language name
	languageCode := message.LanguageCode
	if languageCode == "" {
		languageCode = GetISO639Code(message.Name)
	}

	// If still no language code, return an error with more context
	if languageCode == "" {
		return nil, fmt.Errorf("unable to determine language code for language: '%s' (LanguageCode from API: '%s')",
			message.Name,
			message.LanguageCode,
		)
	}

	return map[string]interface{}{
		"language_code": message.LanguageCode,
		"language_name": message.Name,
		"page_title":    message.Title,

		// Login page
		"login_page_text":   message.LoginDescription,
		"username_text":     message.Username,
		"password_text":     message.Password,
		"login_button_text": message.LoginButton,

		// Device ownership
		"device_ownership_page_text":                  message.DeviceClassDescription,
		"personal_device_button_name":                 message.DeviceClassPersonal,
		"personal_device_management_description":      message.DeviceClassPersonalDescription,
		"institutional_ownership_button_name":         message.DeviceClassEnterprise,
		"institutional_device_management_description": message.DeviceClassEnterpriseDescription,
		"enroll_device_button_name":                   message.DeviceClassButton,

		// EULA
		"eula_personal_devices":      message.EnterpriseEula,
		"eula_institutional_devices": message.PersonalEula,
		"accept_button_text":         message.EulaButton,

		// Site selection
		"site_selection_text": message.SiteDescription,

		// CA Certificate
		"ca_certificate_installation_text":   message.CertificateText,
		"ca_certificate_install_button_name": message.CertificateButton,
		"ca_certificate_name":                message.CertificateProfileName,
		"ca_certificate_description":         message.CertificateProfileDescription,

		// Institutional MDM profile
		"institutional_mdm_profile_installation_text":   message.EnterpriseText,
		"institutional_mdm_profile_install_button_name": message.EnterpriseButton,
		"institutional_mdm_profile_name":                message.EnterpriseProfileName,
		"institutional_mdm_profile_description":         message.EnterpriseProfileDescription,
		"institutional_mdm_profile_pending_text":        message.EnterprisePending,

		// Personal MDM profile
		"personal_mdm_profile_installation_text":   message.PersonalText,
		"personal_mdm_profile_install_button_name": message.PersonalButton,
		"personal_mdm_profile_name":                message.PersonalProfileName,
		"personal_mdm_profile_description":         message.PersonalProfileDescription,

		// User enrollment MDM profile
		"user_enrollment_mdm_profile_installation_text":   message.UserEnrollmentText,
		"user_enrollment_mdm_profile_install_button_name": message.UserEnrollmentButton,
		"user_enrollment_mdm_profile_name":                message.UserEnrollmentProfileName,
		"user_enrollment_mdm_profile_description":         message.UserEnrollmentProfileDescription,

		// QuickAdd package
		"quickadd_package_installation_text":   message.QuickAddText,
		"quickadd_package_install_button_name": message.QuickAddButton,
		"quickadd_package_name":                message.QuickAddName,
		"quickadd_package_progress_text":       message.QuickAddPending,

		// Completion
		"enrollment_complete_text":           message.CompleteMessage,
		"enrollment_failed_text":             message.FailedMessage,
		"try_again_button_name":              message.TryAgainButton,
		"view_enrollment_status_button_name": message.CheckNowButton,
		"view_enrollment_status_text":        message.CheckEnrollmentMessage,
		"log_out_button_name":                message.LogoutButton,
	}, nil
}

// flattenDirectoryServiceGroupEnrollmentSettings converts API groups to maps for the schema
func flattenDirectoryServiceGroupEnrollmentSettings(results []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) []map[string]interface{} {
	var flattenedGroups []map[string]interface{}

	for _, group := range results {
		// Convert ID to integer, defaulting to 0 if conversion fails
		id := 0
		if group.ID != "" {
			parsedID, err := strconv.Atoi(group.ID)
			if err == nil {
				id = parsedID
			}
		}

		siteID := 0
		if group.SiteID != "" {
			siteIDint, err := strconv.Atoi(group.SiteID)
			if err == nil {
				siteID = siteIDint
			}
		}

		groupMap := map[string]interface{}{
			"id":                           id,
			"ldap_server_id":               group.LdapServerID,
			"directory_service_group_name": group.Name,
			"directory_service_group_id":   group.GroupID,
			"site_id":                      siteID,
			"allow_group_to_enroll_institutionally_owned_devices":                      group.EnterpriseEnrollmentEnabled,
			"allow_group_to_enroll_personally_owned_devices":                           group.PersonalEnrollmentEnabled,
			"allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade": group.AccountDrivenUserEnrollmentEnabled,
			"require_eula": group.RequireEula,
		}

		flattenedGroups = append(flattenedGroups, groupMap)
	}

	return flattenedGroups
}

// hasEnrollmentSettingsChange checks if enrollment settings have changed
func hasEnrollmentSettingsChange(d *schema.ResourceData) bool {
	// Check general settings
	if d.HasChange("skip_certificate_installation_during_enrollment") ||
		d.HasChange("restrict_reenrollment_to_authorized_users_only") ||
		d.HasChange("flush_location_information") ||
		d.HasChange("flush_location_history_information") ||
		d.HasChange("flush_policy_history") ||
		d.HasChange("flush_extension_attributes") ||
		d.HasChange("flush_mdm_commands_on_reenroll") ||
		d.HasChange("third_party_signing_certificate") ||
		d.HasChange("user_initiated_enrollment_for_computers") ||
		d.HasChange("user_initiated_enrollment_for_devices") {
		return true
	}

	return false
}

// findLanguagesToDelete identifies language codes to delete
func findLanguagesToDelete(oldMessaging, newMessaging []interface{}) ([]string, error) {
	var languagesToDelete []string

	// Get language codes map
	langCodes, err := getLanguageCodesMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get language codes: %v", err)
	}

	// Create a map of language names in the new set
	newLanguageNames := make(map[string]bool)
	for _, messaging := range newMessaging {
		msg := messaging.(map[string]interface{})
		langName := msg["language_name"].(string)
		newLanguageNames[langName] = true
	}

	// Find languages in old but not in new
	for _, messaging := range oldMessaging {
		msg := messaging.(map[string]interface{})
		langName := msg["language_name"].(string)
		langCode := msg["language_code"].(string)

		// If language has a code and is not in the new set, add to deletion list
		if !newLanguageNames[langName] {
			// If we have the language code directly from state, use it
			if langCode != "" {
				languagesToDelete = append(languagesToDelete, langCode)
			} else if code, ok := langCodes[langName]; ok {
				// Otherwise look it up
				languagesToDelete = append(languagesToDelete, code)
			}
		}
	}

	return languagesToDelete, nil
}

// constructEnrollmentMessagingFromSet builds enrollment messages from a schema set
func constructEnrollmentMessagingFromSet(messagingSet []interface{}) ([]jamfpro.ResourceEnrollmentLanguage, error) {
	var messages []jamfpro.ResourceEnrollmentLanguage

	// Get language codes for mapping
	langCodes, err := getLanguageCodesMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get language codes: %v", err)
	}

	for _, messaging := range messagingSet {
		msg := messaging.(map[string]interface{})

		// Get language code
		var langCode string

		// Try to get from state first
		if code, ok := msg["language_code"].(string); ok && code != "" {
			langCode = code
		} else {
			// Otherwise look up by name
			langName := msg["language_name"].(string)
			if code, ok := langCodes[langName]; ok {
				langCode = code
			} else {
				log.Printf("[WARN] No language code found for language name '%s', skipping", langName)
				continue
			}
		}

		// Create language message
		message := jamfpro.ResourceEnrollmentLanguage{
			LanguageCode: langCode,
			Name:         msg["language_name"].(string),
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

		messages = append(messages, message)
	}

	return messages, nil
}

// processDirectoryServiceGroupChanges identifies groups to add, update, and delete
func processDirectoryServiceGroupChanges(
	oldGroups, newGroups []interface{},
	currentAPIGroups []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup,
) ([]jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, error) {
	var toDelete, toUpdate, toCreate []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup

	// Map of current groups by GroupID
	currentGroupsByGroupID := make(map[string]jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup)
	for _, group := range currentAPIGroups {
		currentGroupsByGroupID[group.GroupID] = group
	}

	// Create a map of groupIDs in new set
	newGroupIDs := make(map[string]bool)
	for _, group := range newGroups {
		groupMap := group.(map[string]interface{})
		groupID := groupMap["directory_service_group_id"].(string)
		newGroupIDs[groupID] = true

		// Check if this is an existing group to update, or a new one to create
		if currentGroup, exists := currentGroupsByGroupID[groupID]; exists {
			// Convert to API format
			siteID := ""
			if v, ok := groupMap["site_id"]; ok && v.(int) > 0 {
				siteID = strconv.Itoa(v.(int))
			}

			// Create updated group with same ID as existing one
			updatedGroup := jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup{
				ID:                                 currentGroup.ID,
				GroupID:                            groupID,
				LdapServerID:                       groupMap["ldap_server_id"].(string),
				Name:                               groupMap["directory_service_group_name"].(string),
				SiteID:                             siteID,
				EnterpriseEnrollmentEnabled:        groupMap["allow_group_to_enroll_institutionally_owned_devices"].(bool),
				PersonalEnrollmentEnabled:          groupMap["allow_group_to_enroll_personally_owned_devices"].(bool),
				AccountDrivenUserEnrollmentEnabled: groupMap["allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade"].(bool),
				RequireEula:                        groupMap["require_eula"].(bool),
			}

			// Check if anything changed
			if !areGroupsEqual(currentGroup, updatedGroup) {
				toUpdate = append(toUpdate, updatedGroup)
			}
		} else {
			// New group to create
			siteID := ""
			if v, ok := groupMap["site_id"]; ok && v.(int) > 0 {
				siteID = strconv.Itoa(v.(int))
			}

			newGroup := jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup{
				GroupID:                            groupID,
				LdapServerID:                       groupMap["ldap_server_id"].(string),
				Name:                               groupMap["directory_service_group_name"].(string),
				SiteID:                             siteID,
				EnterpriseEnrollmentEnabled:        groupMap["allow_group_to_enroll_institutionally_owned_devices"].(bool),
				PersonalEnrollmentEnabled:          groupMap["allow_group_to_enroll_personally_owned_devices"].(bool),
				AccountDrivenUserEnrollmentEnabled: groupMap["allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade"].(bool),
				RequireEula:                        groupMap["require_eula"].(bool),
			}

			toCreate = append(toCreate, newGroup)
		}
	}

	// Find groups to delete (in API but not in new set)
	for _, group := range currentAPIGroups {
		if !newGroupIDs[group.GroupID] {
			toDelete = append(toDelete, group)
		}
	}

	return toDelete, toUpdate, toCreate, nil
}

// constructDefaultEnrollmentSettings creates default enrollment settings for reset
func constructDefaultEnrollmentSettings() *jamfpro.ResourceEnrollment {
	return &jamfpro.ResourceEnrollment{
		InstallSingleProfile:                         false,
		SigningMdmProfileEnabled:                     false,
		RestrictReenrollment:                         false,
		FlushLocationInformation:                     false,
		FlushLocationHistoryInformation:              false,
		FlushPolicyHistory:                           false,
		FlushExtensionAttributes:                     false,
		FlushSoftwareUpdatePlans:                     false,
		FlushMdmCommandsOnReenroll:                   "DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED",
		MacOsEnterpriseEnrollmentEnabled:             false,
		ManagementUsername:                           "jamfadmin",
		CreateManagementAccount:                      false,
		HideManagementAccount:                        false,
		AllowSshOnlyManagementAccount:                false,
		EnsureSshRunning:                             false,
		LaunchSelfService:                            false,
		SignQuickAdd:                                 false,
		IosEnterpriseEnrollmentEnabled:               false,
		IosPersonalEnrollmentEnabled:                 false,
		PersonalDeviceEnrollmentType:                 "USERENROLLMENT",
		AccountDrivenUserEnrollmentEnabled:           false,
		AccountDrivenDeviceIosEnrollmentEnabled:      false,
		AccountDrivenDeviceMacosEnrollmentEnabled:    false,
		AccountDrivenUserVisionosEnrollmentEnabled:   false,
		AccountDrivenDeviceVisionosEnrollmentEnabled: false,
	}
}

// extractLanguageCodesToDelete extracts language codes that can be deleted
func extractLanguageCodesToDelete(messagingSet []interface{}) []string {
	var languageCodesToDelete []string

	for _, messaging := range messagingSet {
		msg := messaging.(map[string]interface{})
		langCode := msg["language_code"].(string)

		// Skip English - it's required and can't be deleted
		if langCode != "en" {
			languageCodesToDelete = append(languageCodesToDelete, langCode)
		}
	}

	return languageCodesToDelete
}

// Helper functions

// getLanguageCodesMap returns a map of language names to codes
// In a real implementation, this would call the API to get the codes
func getLanguageCodesMap() (map[string]string, error) {
	// This is a sample map - in a real implementation, you would get this
	// by calling client.GetEnrollmentLanguageCodes()
	return map[string]string{
		"English":    "en",
		"Spanish":    "es",
		"French":     "fr",
		"German":     "de",
		"Japanese":   "ja",
		"Chinese":    "zh",
		"Korean":     "ko",
		"Italian":    "it",
		"Dutch":      "nl",
		"Swedish":    "sv",
		"Russian":    "ru",
		"Portuguese": "pt",
	}, nil
}

// areGroupsEqual compares two directory service groups for equality
func areGroupsEqual(a, b jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) bool {
	return a.GroupID == b.GroupID &&
		a.LdapServerID == b.LdapServerID &&
		a.Name == b.Name &&
		a.SiteID == b.SiteID &&
		a.EnterpriseEnrollmentEnabled == b.EnterpriseEnrollmentEnabled &&
		a.PersonalEnrollmentEnabled == b.PersonalEnrollmentEnabled &&
		a.AccountDrivenUserEnrollmentEnabled == b.AccountDrivenUserEnrollmentEnabled &&
		a.RequireEula == b.RequireEula
}
