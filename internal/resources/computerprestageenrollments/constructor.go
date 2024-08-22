// computerprestageenrollments_object.go
package computerprestageenrollments

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData, isUpdate bool) (*jamfpro.ResourceComputerPrestage, error) {
	resource := &jamfpro.ResourceComputerPrestage{
		ID:                                d.Get("id").(string),
		DisplayName:                       d.Get("display_name").(string),
		Mandatory:                         jamfpro.BoolPtr(d.Get("mandatory").(bool)),
		MDMRemovable:                      jamfpro.BoolPtr(d.Get("mdm_removable").(bool)),
		SupportPhoneNumber:                d.Get("support_phone_number").(string),
		SupportEmailAddress:               d.Get("support_email_address").(string),
		Department:                        d.Get("department").(string),
		DefaultPrestage:                   jamfpro.BoolPtr(d.Get("default_prestage").(bool)),
		EnrollmentSiteId:                  d.Get("enrollment_site_id").(string),
		KeepExistingSiteMembership:        jamfpro.BoolPtr(d.Get("keep_existing_site_membership").(bool)),
		KeepExistingLocationInformation:   jamfpro.BoolPtr(d.Get("keep_existing_location_information").(bool)),
		RequireAuthentication:             jamfpro.BoolPtr(d.Get("require_authentication").(bool)),
		AuthenticationPrompt:              d.Get("authentication_prompt").(string),
		PreventActivationLock:             jamfpro.BoolPtr(d.Get("prevent_activation_lock").(bool)),
		EnableDeviceBasedActivationLock:   jamfpro.BoolPtr(d.Get("enable_device_based_activation_lock").(bool)),
		DeviceEnrollmentProgramInstanceId: d.Get("device_enrollment_program_instance_id").(string),
		EnableRecoveryLock:                jamfpro.BoolPtr(d.Get("enable_recovery_lock").(bool)),
		RecoveryLockPasswordType:          d.Get("recovery_lock_password_type").(string),
		RecoveryLockPassword:              d.Get("recovery_lock_password").(string),
		RotateRecoveryLockPassword:        jamfpro.BoolPtr(d.Get("rotate_recovery_lock_password").(bool)),
		ProfileUuid:                       d.Get("profile_uuid").(string),
		SiteId:                            d.Get("site_id").(string),
		VersionLock:                       d.Get("version_lock").(int),
		CustomPackageDistributionPointId:  d.Get("custom_package_distribution_point_id").(string),
		EnrollmentCustomizationId:         d.Get("enrollment_customization_id").(string),
		Language:                          d.Get("language").(string),
		Region:                            d.Get("region").(string),
		AutoAdvanceSetup:                  jamfpro.BoolPtr(d.Get("auto_advance_setup").(bool)),
		InstallProfilesDuringSetup:        jamfpro.BoolPtr(d.Get("install_profiles_during_setup").(bool)),
		Enabled:                           jamfpro.BoolPtr(d.Get("enabled").(bool)),
		SsoForEnrollmentEnabled:           jamfpro.BoolPtr(d.Get("sso_for_enrollment_enabled").(bool)),
		SsoBypassAllowed:                  jamfpro.BoolPtr(d.Get("sso_bypass_allowed").(bool)),
		SsoEnabled:                        jamfpro.BoolPtr(d.Get("sso_enabled").(bool)),
		SsoForMacOsSelfServiceEnabled:     jamfpro.BoolPtr(d.Get("sso_for_mac_os_self_service_enabled").(bool)),
		TokenExpirationDisabled:           jamfpro.BoolPtr(d.Get("token_expiration_disabled").(bool)),
		UserAttributeEnabled:              jamfpro.BoolPtr(d.Get("user_attribute_enabled").(bool)),
		UserAttributeName:                 d.Get("user_attribute_name").(string),
		UserMapping:                       d.Get("user_mapping").(string),
		EnrollmentSsoForAccountDrivenEnrollmentEnabled: jamfpro.BoolPtr(d.Get("enrollment_sso_for_account_driven_enrollment_enabled").(bool)),
		GroupEnrollmentAccessEnabled:                   jamfpro.BoolPtr(d.Get("group_enrollment_access_enabled").(bool)),
		GroupAttributeName:                             d.Get("group_attribute_name").(string),
		GroupRdnKey:                                    d.Get("group_rdn_key").(string),
		GroupEnrollmentAccessName:                      d.Get("group_enrollment_access_name").(string),
		IdpProviderType:                                d.Get("idp_provider_type").(string),
		OtherProviderTypeName:                          d.Get("other_provider_type_name").(string),
		MetadataSource:                                 d.Get("metadata_source").(string),
		SessionTimeout:                                 d.Get("session_timeout").(int),
		DeviceType:                                     d.Get("device_type").(string),
	}

	if v, ok := d.GetOk("skip_setup_items"); ok && len(v.([]interface{})) > 0 {
		skipSetupItemsMap := v.([]interface{})[0].(map[string]interface{})
		resource.SkipSetupItems = constructSkipSetupItems(skipSetupItemsMap)
	}

	if v, ok := d.GetOk("location_information"); ok && len(v.([]interface{})) > 0 {
		locationData := v.([]interface{})[0].(map[string]interface{})
		resource.LocationInformation = constructLocationInformation(locationData, isUpdate)
	}

	if v, ok := d.GetOk("purchasing_information"); ok && len(v.([]interface{})) > 0 {
		purchasingData := v.([]interface{})[0].(map[string]interface{})
		resource.PurchasingInformation = constructPurchasingInformation(purchasingData, isUpdate)
	}

	if v, ok := d.GetOk("account_settings"); ok && len(v.([]interface{})) > 0 {
		accountData := v.([]interface{})[0].(map[string]interface{})
		resource.AccountSettings = constructAccountSettings(accountData, isUpdate)
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

	if v, ok := d.GetOk("onboarding_items"); ok {
		onboardingItems := make([]jamfpro.OnboardingItem, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			onboardingItems[i] = jamfpro.OnboardingItem{
				SelfServiceEntityType: itemMap["self_service_entity_type"].(string),
				ID:                    itemMap["id"].(string),
				EntityId:              itemMap["entity_id"].(string),
				Priority:              itemMap["priority"].(int),
			}
		}
		resource.OnboardingItems = onboardingItems
	}

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceXML, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Prestage to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Computer Prestage JSON:\n%s\n", string(resourceXML))

	return resource, nil
}

func constructSkipSetupItems(data map[string]interface{}) jamfpro.ComputerPrestageSubsetSkipSetupItems {
	return jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:         jamfpro.BoolPtr(data["biometric"].(bool)),
		TermsOfAddress:    jamfpro.BoolPtr(data["terms_of_address"].(bool)),
		FileVault:         jamfpro.BoolPtr(data["file_vault"].(bool)),
		ICloudDiagnostics: jamfpro.BoolPtr(data["icloud_diagnostics"].(bool)),
		Diagnostics:       jamfpro.BoolPtr(data["diagnostics"].(bool)),
		Accessibility:     jamfpro.BoolPtr(data["accessibility"].(bool)),
		AppleID:           jamfpro.BoolPtr(data["apple_id"].(bool)),
		ScreenTime:        jamfpro.BoolPtr(data["screen_time"].(bool)),
		Siri:              jamfpro.BoolPtr(data["siri"].(bool)),
		DisplayTone:       jamfpro.BoolPtr(data["display_tone"].(bool)),
		Restore:           jamfpro.BoolPtr(data["restore"].(bool)),
		Appearance:        jamfpro.BoolPtr(data["appearance"].(bool)),
		Privacy:           jamfpro.BoolPtr(data["privacy"].(bool)),
		Payment:           jamfpro.BoolPtr(data["payment"].(bool)),
		Registration:      jamfpro.BoolPtr(data["registration"].(bool)),
		TOS:               jamfpro.BoolPtr(data["tos"].(bool)),
		ICloudStorage:     jamfpro.BoolPtr(data["icloud_storage"].(bool)),
		Location:          jamfpro.BoolPtr(data["location"].(bool)),
	}
}

func constructLocationInformation(data map[string]interface{}, isUpdate bool) jamfpro.ComputerPrestageSubsetLocationInformation {
	return jamfpro.ComputerPrestageSubsetLocationInformation{
		ID:           handleID(data["id"].(string), isUpdate),
		Username:     data["username"].(string),
		Realname:     data["realname"].(string),
		Phone:        data["phone"].(string),
		Email:        data["email"].(string),
		Room:         data["room"].(string),
		Position:     data["position"].(string),
		DepartmentId: data["department_id"].(string),
		BuildingId:   data["building_id"].(string),
		VersionLock:  data["version_lock"].(int),
	}
}

func constructPurchasingInformation(data map[string]interface{}, isUpdate bool) jamfpro.ComputerPrestageSubsetPurchasingInformation {
	newID := handleID(data["id"].(string), isUpdate)
	log.Printf("[DEBUG] constructPurchasingInformation: Using ID '%s'", newID)

	return jamfpro.ComputerPrestageSubsetPurchasingInformation{
		ID:                newID,
		Leased:            jamfpro.BoolPtr(data["leased"].(bool)),
		Purchased:         jamfpro.BoolPtr(data["purchased"].(bool)),
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

func constructAccountSettings(data map[string]interface{}, isUpdate bool) jamfpro.ComputerPrestageSubsetAccountSettings {
	newID := handleID(data["id"].(string), isUpdate)
	log.Printf("[DEBUG] constructAccountSettings: Using ID '%s'", newID)

	return jamfpro.ComputerPrestageSubsetAccountSettings{
		ID:                                      newID,
		PayloadConfigured:                       jamfpro.BoolPtr(data["payload_configured"].(bool)),
		LocalAdminAccountEnabled:                jamfpro.BoolPtr(data["local_admin_account_enabled"].(bool)),
		AdminUsername:                           data["admin_username"].(string),
		AdminPassword:                           data["admin_password"].(string),
		HiddenAdminAccount:                      jamfpro.BoolPtr(data["hidden_admin_account"].(bool)),
		LocalUserManaged:                        jamfpro.BoolPtr(data["local_user_managed"].(bool)),
		UserAccountType:                         data["user_account_type"].(string),
		VersionLock:                             data["version_lock"].(int),
		PrefillPrimaryAccountInfoFeatureEnabled: jamfpro.BoolPtr(data["prefill_primary_account_info_feature_enabled"].(bool)),
		PrefillType:                             data["prefill_type"].(string),
		PrefillAccountFullName:                  data["prefill_account_full_name"].(string),
		PrefillAccountUserName:                  data["prefill_account_user_name"].(string),
		PreventPrefillInfoFromModification:      jamfpro.BoolPtr(data["prevent_prefill_info_from_modification"].(bool)),
	}
}

// handleID manages the ID field for Jamf Pro Computer Prestage resources during create and update operations.
//
// For create operations, it always returns "-1" as required by the Jamf Pro API.
// For update operations, it increments the existing ID by 1.
//
// Parameters:
//
//   - id: A string representing the current ID of the resource.
//     For create operations, this value is ignored.
//     For update operations, this should be the existing ID of the resource.
//
//   - isUpdate: A boolean flag indicating whether this is an update operation.
//     true for update operations, false for create operations.
//
// Returns:
//   - A string representing the ID to be used in the API request.
//     For create operations, this will always be "-1".
//     For update operations, this will be the incremented ID.
//
// Behavior:
//
//   - Create operations (isUpdate == false):
//
//   - Always returns "-1", regardless of the input id.
//
//   - Logs a debug message indicating the use of the default ID.
//
//   - Update operations (isUpdate == true):
//
//   - Attempts to convert the input id to an integer and increment it by 1.
//
//   - If conversion succeeds, returns the incremented ID as a string.
//
//   - If conversion fails, logs a warning and returns the original id unchanged.
//
//   - Logs debug messages showing the original and new IDs.
//
// Error Handling:
//   - If the input id cannot be converted to an integer during an update operation,
//     the function logs a warning and returns the original id unchanged.
//   - This approach ensures that the function always returns a valid string,
//     even in error scenarios.
//
// Usage:
//   - This function should be called for each nested structure within a Computer Prestage
//     resource that requires ID handling (e.g., LocationInformation, PurchasingInformation).
//   - It ensures correct ID handling for both create and update operations in accordance
//     with Jamf Pro API requirements.
func handleID(id string, isUpdate bool) string {
	if !isUpdate {
		log.Printf("[DEBUG] Create operation: Using default ID '-1'")
		return "-1"
	}

	log.Printf("[DEBUG] Update operation: Current ID is '%s'", id)

	currentID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("[WARN] Failed to convert ID '%s' to integer: %v. Using original ID.", id, err)
		return id
	}

	newID := strconv.Itoa(currentID + 1)
	log.Printf("[DEBUG] Update operation: Incrementing ID from '%s' to '%s'", id, newID)
	return newID
}

// handleVersionLock manages the VersionLock field for Jamf Pro Computer Prestage resources during update operations.
//
// Parameters:
//   - currentVersionLock: The current version lock value as an interface{}.
//   - isUpdate: A boolean flag indicating whether this is an update operation.
//
// Returns:
//   - An integer representing the version lock to be used in the API request.
//     For create operations (isUpdate == false), this will be 0.
//     For update operations (isUpdate == true), this will be the incremented version lock.
//
// Behavior:
//   - Create operations (isUpdate == false):
//   - Returns 0, as version lock is not needed for create operations.
//   - Update operations (isUpdate == true):
//   - Attempts to convert the currentVersionLock to an integer and increment it by 1.
//   - If conversion fails, logs a warning and returns 0.
//
// Error Handling:
//   - If the currentVersionLock cannot be converted to an integer during an update operation,
//     the function logs a warning and returns 0.
//
// Usage:
//   - This function should be called for each structure within a Computer Prestage
//     resource that requires version lock handling.
func handleVersionLock(currentVersionLock interface{}, isUpdate bool) int {
	if !isUpdate {
		log.Printf("[DEBUG] Create operation: Version lock not required, using 0")
		return 0
	}

	log.Printf("[DEBUG] Update operation: Current version lock is '%v'", currentVersionLock)

	versionLock, ok := currentVersionLock.(int)
	if !ok {
		log.Printf("[WARN] Failed to convert version lock '%v' to integer. Using 0.", currentVersionLock)
		return 0
	}

	newVersionLock := versionLock + 1
	log.Printf("[DEBUG] Update operation: Incrementing version lock from '%d' to '%d'", versionLock, newVersionLock)
	return newVersionLock
}
