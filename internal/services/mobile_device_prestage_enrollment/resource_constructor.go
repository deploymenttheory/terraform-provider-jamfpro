// Package mobile_device_prestage_enrollment provides the schema and CRUD operations for managing Jamf Pro Mobile Device Prestage Enrollment in Terraform.
package mobile_device_prestage_enrollment

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData, isUpdate bool) (*jamfpro.ResourceMobileDevicePrestage, error) {
	versionLock := constructors.HandleVersionLock(d.Get("version_lock"), isUpdate)

	resource := &jamfpro.ResourceMobileDevicePrestage{
		DisplayName:                            d.Get("display_name").(string),
		Mandatory:                              jamfpro.BoolPtr(d.Get("mandatory").(bool)),
		MdmRemovable:                           jamfpro.BoolPtr(d.Get("mdm_removable").(bool)),
		SupportPhoneNumber:                     d.Get("support_phone_number").(string),
		SupportEmailAddress:                    d.Get("support_email_address").(string),
		Department:                             d.Get("department").(string),
		DefaultPrestage:                        jamfpro.BoolPtr(d.Get("default_prestage").(bool)),
		EnrollmentSiteID:                       d.Get("enrollment_site_id").(string),
		KeepExistingSiteMembership:             jamfpro.BoolPtr(d.Get("keep_existing_site_membership").(bool)),
		KeepExistingLocationInformation:        jamfpro.BoolPtr(d.Get("keep_existing_location_information").(bool)),
		RequireAuthentication:                  jamfpro.BoolPtr(d.Get("require_authentication").(bool)),
		AuthenticationPrompt:                   d.Get("authentication_prompt").(string),
		PreventActivationLock:                  jamfpro.BoolPtr(d.Get("prevent_activation_lock").(bool)),
		EnableDeviceBasedActivationLock:        jamfpro.BoolPtr(d.Get("enable_device_based_activation_lock").(bool)),
		DeviceEnrollmentProgramInstanceID:      d.Get("device_enrollment_program_instance_id").(string),
		EnrollmentCustomizationID:              d.Get("enrollment_customization_id").(string),
		Language:                               d.Get("language").(string),
		Region:                                 d.Get("region").(string),
		AutoAdvanceSetup:                       jamfpro.BoolPtr(d.Get("auto_advance_setup").(bool)),
		AllowPairing:                           jamfpro.BoolPtr(d.Get("allow_pairing").(bool)),
		MultiUser:                              jamfpro.BoolPtr(d.Get("multi_user").(bool)),
		Supervised:                             jamfpro.BoolPtr(d.Get("supervised").(bool)),
		MaximumSharedAccounts:                  d.Get("maximum_shared_accounts").(int),
		ConfigureDeviceBeforeSetupAssistant:    jamfpro.BoolPtr(d.Get("configure_device_before_setup_assistant").(bool)),
		SendTimezone:                           jamfpro.BoolPtr(d.Get("send_timezone").(bool)),
		Timezone:                               d.Get("timezone").(string),
		StorageQuotaSizeMegabytes:              d.Get("storage_quota_size_megabytes").(int),
		UseStorageQuotaSize:                    jamfpro.BoolPtr(d.Get("use_storage_quota_size").(bool)),
		TemporarySessionOnly:                   jamfpro.BoolPtr(d.Get("temporary_session_only").(bool)),
		EnforceTemporarySessionTimeout:         jamfpro.BoolPtr(d.Get("enforce_temporary_session_timeout").(bool)),
		TemporarySessionTimeout:                jamfpro.IntPtr(d.Get("temporary_session_timeout_seconds").(int)),
		EnforceUserSessionTimeout:              jamfpro.BoolPtr(d.Get("enforce_user_session_timeout").(bool)),
		UserSessionTimeout:                     jamfpro.IntPtr(d.Get("user_session_timeout").(int)),
		ProfileUuid:                            d.Get("profile_uuid").(string),
		SiteId:                                 d.Get("site_id").(string),
		VersionLock:                            d.Get("version_lock").(int),
		PrestageMinimumOsTargetVersionTypeIos:  d.Get("prestage_minimum_os_target_version_type_ios").(string),
		MinimumOsSpecificVersionIos:            d.Get("minimum_os_specific_version_ios").(string),
		PrestageMinimumOsTargetVersionTypeIpad: d.Get("prestage_minimum_os_target_version_type_ipad").(string),
		RTSEnabled:                             jamfpro.BoolPtr(d.Get("rts_enabled").(bool)),
		RTSConfigProfileId:                     d.Get("rts_config_profile_id").(string),
		MinimumOsSpecificVersionIpad:           d.Get("minimum_os_specific_version_ipad").(string),
		PreserveManagedApps:                    jamfpro.BoolPtr(d.Get("preserve_managed_apps").(bool)),
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

	if v, ok := d.GetOk("names"); ok && len(v.([]any)) > 0 {
		namesData := v.([]any)[0].(map[string]any)
		resource.Names = constructNames(namesData, isUpdate)
	}

	if v, ok := d.GetOk("anchor_certificates"); ok {
		anchorCertificates := make([]string, len(v.([]any)))
		for i, cert := range v.([]any) {
			anchorCertificates[i] = cert.(string)
		}
		resource.AnchorCertificates = anchorCertificates
	}

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Prestage to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Prestage JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructSkipSetupItems constructs the SkipSetupItems subset of a Mobile Device Prestage resource.
func constructSkipSetupItems(data map[string]any) jamfpro.MobileDevicePrestageSubsetSkipSetupItems {
	return jamfpro.MobileDevicePrestageSubsetSkipSetupItems{
		Location:              jamfpro.BoolPtr(data["location"].(bool)),
		Privacy:               jamfpro.BoolPtr(data["privacy"].(bool)),
		Biometric:             jamfpro.BoolPtr(data["biometric"].(bool)),
		SoftwareUpdate:        jamfpro.BoolPtr(data["software_update"].(bool)),
		Diagnostics:           jamfpro.BoolPtr(data["diagnostics"].(bool)),
		IMessageAndFaceTime:   jamfpro.BoolPtr(data["imessage_and_facetime"].(bool)),
		Intelligence:          jamfpro.BoolPtr(data["intelligence"].(bool)),
		TVRoom:                jamfpro.BoolPtr(data["tv_room"].(bool)),
		Passcode:              jamfpro.BoolPtr(data["passcode"].(bool)),
		SIMSetup:              jamfpro.BoolPtr(data["sim_setup"].(bool)),
		ScreenTime:            jamfpro.BoolPtr(data["screen_time"].(bool)),
		RestoreCompleted:      jamfpro.BoolPtr(data["restore_completed"].(bool)),
		TVProviderSignIn:      jamfpro.BoolPtr(data["tv_provider_sign_in"].(bool)),
		Siri:                  jamfpro.BoolPtr(data["siri"].(bool)),
		Restore:               jamfpro.BoolPtr(data["restore"].(bool)),
		ScreenSaver:           jamfpro.BoolPtr(data["screen_saver"].(bool)),
		HomeButtonSensitivity: jamfpro.BoolPtr(data["home_button_sensitivity"].(bool)),
		CloudStorage:          jamfpro.BoolPtr(data["cloud_storage"].(bool)),
		ActionButton:          jamfpro.BoolPtr(data["action_button"].(bool)),
		TransferData:          jamfpro.BoolPtr(data["transfer_data"].(bool)),
		EnableLockdownMode:    jamfpro.BoolPtr(data["enable_lockdown_mode"].(bool)),
		Zoom:                  jamfpro.BoolPtr(data["zoom"].(bool)),
		PreferredLanguage:     jamfpro.BoolPtr(data["preferred_language"].(bool)),
		VoiceSelection:        jamfpro.BoolPtr(data["voice_selection"].(bool)),
		TVHomeScreenSync:      jamfpro.BoolPtr(data["tv_home_screen_sync"].(bool)),
		Safety:                jamfpro.BoolPtr(data["safety"].(bool)),
		TermsOfAddress:        jamfpro.BoolPtr(data["terms_of_address"].(bool)),
		ExpressLanguage:       jamfpro.BoolPtr(data["express_language"].(bool)),
		CameraButton:          jamfpro.BoolPtr(data["camera_button"].(bool)),
		AppleID:               jamfpro.BoolPtr(data["apple_id"].(bool)),
		DisplayTone:           jamfpro.BoolPtr(data["display_tone"].(bool)),
		WatchMigration:        jamfpro.BoolPtr(data["watch_migration"].(bool)),
		UpdateCompleted:       jamfpro.BoolPtr(data["update_completed"].(bool)),
		Appearance:            jamfpro.BoolPtr(data["appearance"].(bool)),
		Android:               jamfpro.BoolPtr(data["android"].(bool)),
		Payment:               jamfpro.BoolPtr(data["payment"].(bool)),
		OnBoarding:            jamfpro.BoolPtr(data["onboarding"].(bool)),
		TOS:                   jamfpro.BoolPtr(data["tos"].(bool)),
		Welcome:               jamfpro.BoolPtr(data["welcome"].(bool)),
		SafetyAndHandling:     jamfpro.BoolPtr(data["safety_and_handling"].(bool)),
		TapToSetup:            jamfpro.BoolPtr(data["tap_to_setup"].(bool)),
		SpokenLanguage:        jamfpro.BoolPtr(data["spoken_language"].(bool)),
		Keyboard:              jamfpro.BoolPtr(data["keyboard"].(bool)),
		Multitasking:          jamfpro.BoolPtr(data["multitasking"].(bool)),
		OSShowcase:            jamfpro.BoolPtr(data["os_showcase"].(bool)),
	}
}

// constructLocationInformation constructs the LocationInformation subset of a Mobile Device Prestage resource.
func constructLocationInformation(data map[string]any, isUpdate bool, versionLock int) jamfpro.MobileDevicePrestageSubsetLocationInformation {
	return jamfpro.MobileDevicePrestageSubsetLocationInformation{
		Username:     data["username"].(string),
		Realname:     data["realname"].(string),
		Phone:        data["phone"].(string),
		Email:        data["email"].(string),
		Room:         data["room"].(string),
		Position:     data["position"].(string),
		DepartmentId: data["department_id"].(string),
		BuildingId:   data["building_id"].(string),
		ID:           "-1",
		VersionLock:  versionLock,
	}
}

// constructPurchasingInformation constructs the PurchasingInformation subset of a Mobile Device Prestage resource.
func constructPurchasingInformation(data map[string]any, isUpdate bool, versionLock int) jamfpro.MobileDevicePrestageSubsetPurchasingInformation {
	return jamfpro.MobileDevicePrestageSubsetPurchasingInformation{
		ID:                "-1",
		Leased:            jamfpro.BoolPtr(data["leased"].(bool)),
		Purchased:         jamfpro.BoolPtr(data["purchased"].(bool)),
		AppleCareId:       data["apple_care_id"].(string),
		PoNumber:          data["po_number"].(string),
		Vendor:            data["vendor"].(string),
		PurchasePrice:     data["purchase_price"].(string),
		LifeExpectancy:    data["life_expectancy"].(int),
		PurchasingAccount: data["purchasing_account"].(string),
		PurchasingContact: data["purchasing_contact"].(string),
		LeaseDate:         data["lease_date"].(string),
		PoDate:            data["po_date"].(string),
		WarrantyDate:      data["warranty_date"].(string),
		VersionLock:       versionLock,
	}
}

// constructNames constructs the Names subset of a Mobile Device Prestage resource.
func constructNames(data map[string]any, isUpdate bool) jamfpro.MobileDevicePrestageSubsetNames {
	names := jamfpro.MobileDevicePrestageSubsetNames{
		AssignNamesUsing:       data["assign_names_using"].(string),
		DeviceNamePrefix:       data["device_name_prefix"].(string),
		DeviceNameSuffix:       data["device_name_suffix"].(string),
		SingleDeviceName:       data["single_device_name"].(string),
		ManageNames:            jamfpro.BoolPtr(data["manage_names"].(bool)),
		DeviceNamingConfigured: jamfpro.BoolPtr(data["device_naming_configured"].(bool)),
	}

	if v, ok := data["prestage_device_names"]; ok {
		deviceNames := v.([]any)
		names.PrestageDeviceNames = make([]jamfpro.MobileDevicePrestageSubsetNamesName, len(deviceNames))

		for i, name := range deviceNames {
			nameMap := name.(map[string]any)
			names.PrestageDeviceNames[i] = jamfpro.MobileDevicePrestageSubsetNamesName{
				ID:         "-1",
				DeviceName: nameMap["device_name"].(string),
				Used:       jamfpro.BoolPtr(nameMap["used"].(bool)),
			}
		}
	}

	return names
}
