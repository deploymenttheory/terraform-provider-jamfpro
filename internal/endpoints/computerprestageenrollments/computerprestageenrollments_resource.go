// computerprestages_resource.go
package computerprestageenrollments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProComputerPrestageEnrollmentEnrollment defines the schema for managing Jamf Pro Computer Prestages in Terraform.
func ResourceJamfProComputerPrestageEnrollmentEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerPrestageEnrollmentCreate,
		ReadContext:   ResourceJamfProComputerPrestageEnrollmentRead,
		UpdateContext: ResourceJamfProComputerPrestageEnrollmentUpdate,
		DeleteContext: ResourceJamfProComputerPrestageEnrollmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the computer prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the computer prestage.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether the computer prestage is mandatory.",
			},
			"mdm_removable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the MDM profile is removable.",
			},
			"support_phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Support phone number for the organization.",
			},
			"support_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Support email address for the organization.",
			},
			"department": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The department the computer prestage is assigned to.",
			},
			"default_prestage": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if this is the default computer prestage enrollment configuration. If yes then new devices will be automatically assigned to this PreStage enrollment",
			},
			"enrollment_site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The jamf pro Site ID that computers will be added to during enrollment. Default is -1, aka not used.",
			},
			"keep_existing_site_membership": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if existing device site membership should be retained.",
			},
			"keep_existing_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if existing device location information should be retained.",
			},
			"require_authentication": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the user is required to provide username and password on computers with macOS 10.10 or later.",
			},
			"authentication_prompt": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The authentication prompt message displayed to the user during enrollment.",
			},
			"prevent_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if activation lock should be prevented.",
			},
			"enable_device_based_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if device-based activation lock should be enabled.",
			},
			"device_enrollment_program_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The device enrollment program instance ID.",
			},
			"skip_setup_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Selected items are not displayed in the Setup Assistant during macOS device setup within Apple Device Enrollment (ADE).",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biometric": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip biometric setup.",
						},
						"terms_of_address": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip terms of address setup.",
						},
						"file_vault": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip FileVault setup.",
						},
						"icloud_diagnostics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip iCloud diagnostics setup.",
						},
						"diagnostics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip diagnostics setup.",
						},
						"accessibility": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip accessibility setup.",
						},
						"apple_id": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Apple ID setup.",
						},
						"screen_time": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Screen Time setup.",
						},
						"siri": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Siri setup.",
						},
						"display_tone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Display Tone setup.",
						},
						"restore": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Restore setup.",
						},
						"appearance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Appearance setup.",
						},
						"privacy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Privacy setup.",
						},
						"payment": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Payment setup.",
						},
						"registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Registration setup.",
						},
						"tos": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Terms of Service setup.",
						},
						"icloud_storage": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip iCloud Storage setup.",
						},
						"location": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Location setup.",
						},
					},
				},
			},
			"location_information": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Location information associated with the Jamf Pro computer prestage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username for the location information.",
						},
						"realname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The real name associated with this location.",
						},
						"phone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The phone number associated with this location.",
						},
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The email address associated with this location.",
						},
						"room": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The room associated with this location.",
						},
						"position": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The position associated with this location.",
						},
						"department_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The computerPrestage ID associated with this location.",
							Default:     "-1",
						},
						"building_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The building ID associated with this location.",
							Default:     "-1",
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the location information.",
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version lock of the location information.",
						},
					},
				},
			},
			"purchasing_information": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Purchasing information associated with the computer prestage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the purchasing information.",
						},
						"leased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is leased.",
						},
						"purchased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is purchased.",
						},
						"apple_care_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AppleCare ID.",
						},
						"po_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase order number.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The vendor name.",
						},
						"purchase_price": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase price.",
						},
						"life_expectancy": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The life expectancy in years.",
						},
						"purchasing_account": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing account.",
						},
						"purchasing_contact": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing contact.",
						},
						"lease_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The lease date.",
						},
						"po_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase order date.",
						},
						"warranty_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The warranty date.",
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version lock.",
						},
					},
				},
			},
			"anchor_certificates": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Base64 encoded PEM Certificates.",
			},
			"enrollment_customization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The enrollment customization ID.",
				Default:     "0",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The language setting.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region setting.",
			},
			"auto_advance_setup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if setup should auto-advance.",
			},
			"install_profiles_during_setup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if profiles should be installed during setup.",
			},
			"prestage_installed_profile_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IDs of profiles installed during prestage.",
			},
			"custom_package_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom package IDs.",
			},
			"custom_package_distribution_point_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom package distribution point ID.",
			},
			"enable_recovery_lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if recovery lock should be enabled.",
			},
			"recovery_lock_password_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The recovery lock password type.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					validTypes := map[string]bool{
						"MANUAL": true,
						"RANDOM": true,
					}
					if _, valid := validTypes[v]; !valid {
						errs = append(errs, fmt.Errorf("%q must be one of 'MANUAL', 'RANDOM', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"recovery_lock_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The recovery lock password.",
			},
			"rotate_recovery_lock_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the recovery lock password should be rotated.",
			},
			"profile_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The profile UUID.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The site ID.",
				Default:     "-1",
			},
			"version_lock": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The version lock.",
			},
			"account_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of Account Settings.",
						},
						"payload_configured": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the payload is configured.",
						},
						"local_admin_account_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the local admin account is enabled.",
						},
						"admin_username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The admin username.",
						},
						"admin_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The admin password.",
						},
						"hidden_admin_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the admin account is hidden.",
						},
						"local_user_managed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the local user is managed.",
						},
						"user_account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of user account (ADMINISTRATOR, STANDARD, SKIP).",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := util.GetString(val)
								validTypes := map[string]bool{
									"ADMINISTRATOR": true,
									"STANDARD":      true,
									"SKIP":          true,
								}
								if _, valid := validTypes[v]; !valid {
									errs = append(errs, fmt.Errorf("%q must be one of 'ADMINISTRATOR', 'STANDARD', 'SKIP', got: %s", key, v))
								}
								return warns, errs
							},
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The version lock for account settings.",
						},
						"prefill_primary_account_info_feature_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if prefilling primary account info feature is enabled.",
						},
						"prefill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of prefill (CUSTOM, DEVICE_OWNER).",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := util.GetString(val)
								validTypes := map[string]bool{
									"CUSTOM":       true,
									"DEVICE_OWNER": true,
								}
								if _, valid := validTypes[v]; !valid {
									errs = append(errs, fmt.Errorf("%q must be one of 'CUSTOM', 'DEVICE_OWNER', got: %s", key, v))
								}
								return warns, errs
							},
						},
						"prefill_account_full_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Full name for the account to prefill.",
						},
						"prefill_account_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Username for the account to prefill.",
						},
						"prevent_prefill_info_from_modification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if prefill info is prevented from modification.",
						},
					},
				},
			},
		},
	}
}

// ResourceJamfProComputerPrestageEnrollmentCreate is responsible for creating a new computer prestage in Jamf Pro with terraform.
// The function:
// 1. Constructs the computer prestage data using the provided Terraform configuration.
// 2. Calls the API to create the computer prestage in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created computer prestage.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProComputerPrestageEnrollmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProComputerPrestageEnrollment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro computer prestage enrollment: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseComputerPrestageCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateComputerPrestage(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro computer prestage enrollment '%s' after retries: %v", resource.DisplayName, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(creationResponse.ID)

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		return apiclient.Conn.GetComputerPrestageByID(id.(string))
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro computer prestage enrollment", creationResponse.ID, checkResourceExists, 10*time.Second)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerPrestageEnrollmentRead(ctx, d, meta)
	if len(readDiags) > 0 {
		return readDiags
	}

	return diags
}

// ResourceJamfProComputerPrestageEnrollmentRead is responsible for reading the current state of a Building Resource from the remote system.
// The function:
// 1. Fetches the building's current state using its ID. If it fails, then obtain the building's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the building being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProComputerPrestageEnrollmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}

	// Use the script ID from Terraform's data schema as the resource identifier
	resourceID := d.Id()

	// Attempt to fetch the resource by ID
	resource, err := apiclient.Conn.GetComputerPrestageByID(resourceID)

	if err != nil {
		// Skip resource state removal if this is a create operation
		if !d.IsNewResource() {
			// If the error is a "not found" error, remove the resource from the state
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "410") {
				d.SetId("") // Remove the resource from Terraform state
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Resource not found",
						Detail:   fmt.Sprintf("Jamf Pro computer prestage enrollment resource with ID '%s' was not found and has been removed from the Terraform state.", resourceID),
					},
				}
			}
		}
		// For other errors, or if this is a create operation, return a diagnostic error
		return diag.FromErr(err)
	}

	// Update the Terraform state with the fetched data
	if resource != nil {
		// Construct a map of computer prestage attributes
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

		// Handle nested location_information
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
		// Handle nested purchasing_information
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
		// Add other single-level attributes
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
		// Handle nested account_settings
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
		// Update the Terraform state with prestage attributes
		for key, val := range prestageAttributes {
			if err := d.Set(key, val); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

// ResourceJamfProComputerPrestageEnrollmentUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProComputerPrestageEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Construct the resource object
	resource, err := constructJamfProComputerPrestageEnrollment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Disk Computer Prestage for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateComputerPrestageByID(resourceID, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Computer Prestage '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerPrestageEnrollmentRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerPrestageEnrollmentDelete is responsible for deleting a Jamf Pro Computer Prestage.
func ResourceJamfProComputerPrestageEnrollmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteComputerPrestageByID(resourceID)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Display Name
			resourceDisplayName := d.Get("display_name").(string)
			apiErrByDisplayName := conn.DeleteComputerPrestageByName(resourceDisplayName)
			if apiErrByDisplayName != nil {
				// If deletion by display name also fails, return a retryable error
				return retry.RetryableError(apiErrByDisplayName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Computer Prestage '%s' (ID: %s) after retries: %v", d.Get("display_name").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
