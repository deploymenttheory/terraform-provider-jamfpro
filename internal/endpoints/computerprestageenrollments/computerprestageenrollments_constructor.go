// computerprestageenrollments_object.go
package computerprestageenrollments

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerPrestageEnrollment constructs a ResourceComputerPrestage object from the provided schema data.
func constructJamfProComputerPrestageEnrollment(d *schema.ResourceData) (*jamfpro.ResourceComputerPrestage, error) {
	var resource *jamfpro.ResourceComputerPrestage

	resource = &jamfpro.ResourceComputerPrestage{
		DisplayName:                     d.Get("display_name").(string),
		Mandatory:                       d.Get("mandatory").(bool),
		MDMRemovable:                    d.Get("mdm_removable").(bool),
		SupportPhoneNumber:              d.Get("support_phone_number").(string),
		SupportEmailAddress:             d.Get("support_email_address").(string),
		Department:                      d.Get("department").(string),
		DefaultPrestage:                 d.Get("default_prestage").(bool),
		EnrollmentSiteId:                d.Get("enrollment_site_id").(string),
		KeepExistingSiteMembership:      d.Get("keep_existing_site_membership").(bool),
		KeepExistingLocationInformation: d.Get("keep_existing_location_information").(bool),
		RequireAuthentication:           d.Get("require_authentication").(bool),
		AuthenticationPrompt:            d.Get("authentication_prompt").(string),
		PreventActivationLock:           d.Get("prevent_activation_lock").(bool),
		EnableDeviceBasedActivationLock: d.Get("enable_device_based_activation_lock").(bool),
		//
		DeviceEnrollmentProgramInstanceId: d.Get("device_enrollment_program_instance_id").(string),
		EnableRecoveryLock:                d.Get("enable_recovery_lock").(bool),
		RecoveryLockPasswordType:          d.Get("recovery_lock_password_type").(string),
		RecoveryLockPassword:              d.Get("recovery_lock_password").(string),
		RotateRecoveryLockPassword:        d.Get("rotate_recovery_lock_password").(bool),
		ProfileUuid:                       d.Get("profile_uuid").(string),
		SiteId:                            d.Get("site_id").(string),
		VersionLock:                       d.Get("version_lock").(int),
		CustomPackageDistributionPointId:  d.Get("custom_package_distribution_point_id").(string),
	}

	if v, ok := d.GetOk("skip_setup_items"); ok && len(v.([]interface{})) > 0 {
		skipSetupItemsMap := v.([]interface{})[0].(map[string]interface{})
		resource.SkipSetupItems = constructSkipSetupItems(skipSetupItemsMap)
	}

	if v, ok := d.GetOk("location_information"); ok && len(v.([]interface{})) > 0 {
		locationData := v.([]interface{})[0].(map[string]interface{})
		resource.LocationInformation = constructLocationInformation(locationData)
	}

	if v, ok := d.GetOk("purchasing_information"); ok && len(v.([]interface{})) > 0 {
		purchasingData := v.([]interface{})[0].(map[string]interface{})
		resource.PurchasingInformation = constructPurchasingInformation(purchasingData)
	}

	resource.EnrollmentCustomizationId = d.Get("enrollment_customization_id").(string)
	resource.Language = d.Get("language").(string)
	resource.Region = d.Get("region").(string)
	resource.AutoAdvanceSetup = d.Get("auto_advance_setup").(bool)
	resource.InstallProfilesDuringSetup = d.Get("install_profiles_during_setup").(bool)
	resource.EnableRecoveryLock = d.Get("enable_recovery_lock").(bool)
	resource.RecoveryLockPasswordType = d.Get("recovery_lock_password_type").(string)
	resource.RecoveryLockPassword = d.Get("recovery_lock_password").(string)
	resource.RotateRecoveryLockPassword = d.Get("rotate_recovery_lock_password").(bool)
	resource.ProfileUuid = d.Get("profile_uuid").(string)
	resource.SiteId = d.Get("site_id").(string)
	resource.VersionLock = d.Get("version_lock").(int)

	if v, ok := d.GetOk("account_settings"); ok && len(v.([]interface{})) > 0 {
		accountData := v.([]interface{})[0].(map[string]interface{})
		resource.AccountSettings = constructAccountSettings(accountData)
	}

	if v, ok := d.GetOk("anchor_certificates"); ok {
		anchorCertificates := make([]string, len(v.([]interface{})))
		for i, cert := range v.([]interface{}) {
			anchorCertificates[i] = cert.(string)
		}
		resource.AnchorCertificates = anchorCertificates
	}

	if v, ok := d.GetOk("prestage_installed_profile_ids"); ok {
		profileIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			profileIDs[i] = id.(string)
		}
		resource.PrestageInstalledProfileIds = profileIDs
	}

	if v, ok := d.GetOk("custom_package_ids"); ok {
		packageIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			packageIDs[i] = id.(string)
		}
		resource.CustomPackageIds = packageIDs
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Prestage Enrollment '%s' to XML: %v", resource.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Prestage Enrollment XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// Helper functions for complex structures
func constructSkipSetupItems(data map[string]interface{}) jamfpro.ComputerPrestageSubsetSkipSetupItems {
	return jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:         data["biometric"].(bool),
		TermsOfAddress:    data["terms_of_address"].(bool),
		FileVault:         data["file_vault"].(bool),
		ICloudDiagnostics: data["icloud_diagnostics"].(bool),
		Diagnostics:       data["diagnostics"].(bool),
		Accessibility:     data["accessibility"].(bool),
		AppleID:           data["apple_id"].(bool),
		ScreenTime:        data["screen_time"].(bool),
		Siri:              data["siri"].(bool),
		DisplayTone:       data["display_tone"].(bool),
		Restore:           data["restore"].(bool),
		Appearance:        data["appearance"].(bool),
		Privacy:           data["privacy"].(bool),
		Payment:           data["payment"].(bool),
		Registration:      data["registration"].(bool),
		TOS:               data["tos"].(bool),
		ICloudStorage:     data["icloud_storage"].(bool),
		Location:          data["location"].(bool),
	}
}

func constructLocationInformation(data map[string]interface{}) jamfpro.ComputerPrestageSubsetLocationInformation {
	return jamfpro.ComputerPrestageSubsetLocationInformation{
		Username:     data["username"].(string),
		Realname:     data["realname"].(string),
		Phone:        data["phone"].(string),
		Email:        data["email"].(string),
		Room:         data["room"].(string),
		Position:     data["position"].(string),
		DepartmentId: data["department_id"].(string),
		BuildingId:   data["building_id"].(string),
		ID:           data["id"].(string),
		VersionLock:  data["version_lock"].(int),
	}
}

func constructPurchasingInformation(data map[string]interface{}) jamfpro.ComputerPrestageSubsetPurchasingInformation {
	return jamfpro.ComputerPrestageSubsetPurchasingInformation{
		ID:                data["id"].(string),
		Leased:            data["leased"].(bool),
		Purchased:         data["purchased"].(bool),
		AppleCareId:       data["apple_care_id"].(string),
		PONumber:          data["po_number"].(string),
		Vendor:            data["vendor"].(string),
		PurchasePrice:     data["purchase_price"].(string),
		LifeExpectancy:    data["life_expectancy"].(int),
		PurchasingAccount: data["purchasing_account"].(string),
		PurchasingContact: data["purchasing_contact"].(string),
		LeaseDate:         data["lease_date"].(string),
		PODate:            data["po_date"].(string),
		WarrantyDate:      data["warranty_date"].(string),
		VersionLock:       data["version_lock"].(int),
	}
}

func constructAccountSettings(data map[string]interface{}) jamfpro.ComputerPrestageSubsetAccountSettings {
	return jamfpro.ComputerPrestageSubsetAccountSettings{
		ID:                                      data["id"].(string),
		PayloadConfigured:                       data["payload_configured"].(bool),
		LocalAdminAccountEnabled:                data["local_admin_account_enabled"].(bool),
		AdminUsername:                           data["admin_username"].(string),
		AdminPassword:                           data["admin_password"].(string),
		HiddenAdminAccount:                      data["hidden_admin_account"].(bool),
		LocalUserManaged:                        data["local_user_managed"].(bool),
		UserAccountType:                         data["user_account_type"].(string),
		VersionLock:                             data["version_lock"].(int),
		PrefillPrimaryAccountInfoFeatureEnabled: data["prefill_primary_account_info_feature_enabled"].(bool),
		PrefillType:                             data["prefill_type"].(string),
		PrefillAccountFullName:                  data["prefill_account_full_name"].(string),
		PrefillAccountUserName:                  data["prefill_account_user_name"].(string),
		PreventPrefillInfoFromModification:      data["prevent_prefill_info_from_modification"].(bool),
	}
}
