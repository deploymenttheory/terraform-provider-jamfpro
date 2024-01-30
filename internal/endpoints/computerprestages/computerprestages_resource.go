// computerprestages_resource.go
package computerprestages

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	JamfProResourceComputerPrestage = "Computer Prestage"
)

// ResourceJamfProComputerPrestage defines the schema for managing Jamf Pro Computer Prestages in Terraform.
func ResourceJamfProComputerPrestage() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerPrestageCreate,
		ReadContext:   ResourceJamfProComputerPrestageRead,
		UpdateContext: ResourceJamfProComputerPrestageUpdate,
		DeleteContext: ResourceJamfProComputerPrestageDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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
				Optional:    true,
				Description: "The Support phone number for the organization.",
			},
			"support_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Support email address for the organization.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The department the computer prestage is assigned to.",
			},
			"default_prestage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if this is the default computer prestage enrollment configuration. If yes then new devices will be automatically assigned to this PreStage enrollment",
			},
			"enrollment_site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The jamf pro Site ID that computers will be added to during enrollment. Default is -1, aka not used.",
			},
			"keep_existing_site_membership": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if existing device site membership should be retained.",
			},
			"keep_existing_location_information": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if existing device location information should be retained.",
			},
			"require_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is required to provide username and password on computers with macOS 10.10 or later.",
			},
			"authentication_prompt": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authentication prompt message displayed to the user during enrollment.",
			},
			"prevent_activation_lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if activation lock should be prevented.",
			},
			"enable_device_based_activation_lock": {
				Type:        schema.TypeBool,
				Optional:    true,
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
				Optional:    true,
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

// constructJamfProComputerPrestage constructs a ResourceComputerPrestage object from the provided schema data.
func constructJamfProComputerPrestage(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceComputerPrestage, error) {
	prestage := &jamfpro.ResourceComputerPrestage{}

	// Extract the items from the Terraform resource data
	prestage.DisplayName = util.GetStringFromInterface(d.Get("display_name"))
	prestage.Mandatory = util.GetBoolFromInterface(d.Get("mandatory"))
	prestage.MDMRemovable = util.GetBoolFromInterface(d.Get("mdm_removable"))
	prestage.SupportPhoneNumber = util.GetStringFromInterface(d.Get("support_phone_number"))
	prestage.SupportEmailAddress = util.GetStringFromInterface(d.Get("support_email_address"))
	prestage.Department = util.GetStringFromInterface(d.Get("department"))
	prestage.DefaultPrestage = util.GetBoolFromInterface(d.Get("default_prestage"))
	prestage.EnrollmentSiteId = util.GetStringFromInterface(d.Get("enrollment_site_id"))
	prestage.KeepExistingSiteMembership = util.GetBoolFromInterface(d.Get("keep_existing_site_membership"))
	prestage.KeepExistingLocationInformation = util.GetBoolFromInterface(d.Get("keep_existing_location_information"))
	prestage.RequireAuthentication = util.GetBoolFromInterface(d.Get("require_authentication"))
	prestage.AuthenticationPrompt = util.GetStringFromInterface(d.Get("authentication_prompt"))
	prestage.PreventActivationLock = util.GetBoolFromInterface(d.Get("prevent_activation_lock"))
	prestage.EnableDeviceBasedActivationLock = util.GetBoolFromInterface(d.Get("enable_device_based_activation_lock"))
	prestage.DeviceEnrollmentProgramInstanceId = util.GetStringFromInterface(d.Get("device_enrollment_program_instance_id"))

	// Extract the 'skip_setup_items' list from the Terraform resource data
	if v, ok := d.GetOk("skip_setup_items"); ok {
		skipSetupItemsList := v.([]interface{})
		if len(skipSetupItemsList) > 0 {
			skipSetupItemsMap := skipSetupItemsList[0].(map[string]interface{})

			// Construct the ComputerPrestageSubsetSkipSetupItems struct
			prestage.SkipSetupItems = jamfpro.ComputerPrestageSubsetSkipSetupItems{
				Biometric:         util.GetBoolFromMap(skipSetupItemsMap, "biometric"),
				TermsOfAddress:    util.GetBoolFromMap(skipSetupItemsMap, "terms_of_address"),
				FileVault:         util.GetBoolFromMap(skipSetupItemsMap, "file_vault"),
				ICloudDiagnostics: util.GetBoolFromMap(skipSetupItemsMap, "icloud_diagnostics"),
				Diagnostics:       util.GetBoolFromMap(skipSetupItemsMap, "diagnostics"),
				Accessibility:     util.GetBoolFromMap(skipSetupItemsMap, "accessibility"),
				AppleID:           util.GetBoolFromMap(skipSetupItemsMap, "apple_id"),
				ScreenTime:        util.GetBoolFromMap(skipSetupItemsMap, "screen_time"),
				Siri:              util.GetBoolFromMap(skipSetupItemsMap, "siri"),
				DisplayTone:       util.GetBoolFromMap(skipSetupItemsMap, "display_tone"),
				Restore:           util.GetBoolFromMap(skipSetupItemsMap, "restore"),
				Appearance:        util.GetBoolFromMap(skipSetupItemsMap, "appearance"),
				Privacy:           util.GetBoolFromMap(skipSetupItemsMap, "privacy"),
				Payment:           util.GetBoolFromMap(skipSetupItemsMap, "payment"),
				Registration:      util.GetBoolFromMap(skipSetupItemsMap, "registration"),
				TOS:               util.GetBoolFromMap(skipSetupItemsMap, "tos"),
				ICloudStorage:     util.GetBoolFromMap(skipSetupItemsMap, "icloud_storage"),
				Location:          util.GetBoolFromMap(skipSetupItemsMap, "location"),
			}
		}
	}

	// Extract location_information
	if v, ok := d.GetOk("location_information"); ok {
		locationList := v.([]interface{})
		if len(locationList) > 0 {
			locationData := locationList[0].(map[string]interface{})
			prestage.LocationInformation = jamfpro.ComputerPrestageSubsetLocationInformation{
				Username:     util.GetStringFromMap(locationData, "username"),
				Realname:     util.GetStringFromMap(locationData, "realname"),
				Phone:        util.GetStringFromMap(locationData, "phone"),
				Email:        util.GetStringFromMap(locationData, "email"),
				Room:         util.GetStringFromMap(locationData, "room"),
				Position:     util.GetStringFromMap(locationData, "position"),
				DepartmentId: util.GetStringFromMap(locationData, "department_id"),
				BuildingId:   util.GetStringFromMap(locationData, "building_id"),
				ID:           util.GetStringFromMap(locationData, "id"),
				VersionLock:  util.GetIntFromMap(locationData, "version_lock"),
			}
		}
	}

	// Extract purchasing_information
	if v, ok := d.GetOk("purchasing_information"); ok {
		purchasingList := v.([]interface{})
		if len(purchasingList) > 0 {
			purchasingData := purchasingList[0].(map[string]interface{})
			prestage.PurchasingInformation = jamfpro.ComputerPrestageSubsetPurchasingInformation{
				ID:                util.GetStringFromMap(purchasingData, "id"),
				Leased:            util.GetBoolFromMap(purchasingData, "leased"),
				Purchased:         util.GetBoolFromMap(purchasingData, "purchased"),
				AppleCareId:       util.GetStringFromMap(purchasingData, "apple_care_id"),
				PONumber:          util.GetStringFromMap(purchasingData, "po_number"),
				Vendor:            util.GetStringFromMap(purchasingData, "vendor"),
				PurchasePrice:     util.GetStringFromMap(purchasingData, "purchase_price"),
				LifeExpectancy:    util.GetIntFromMap(purchasingData, "life_expectancy"),
				PurchasingAccount: util.GetStringFromMap(purchasingData, "purchasing_account"),
				PurchasingContact: util.GetStringFromMap(purchasingData, "purchasing_contact"),
				LeaseDate:         util.GetStringFromMap(purchasingData, "lease_date"),
				PODate:            util.GetStringFromMap(purchasingData, "po_date"),
				WarrantyDate:      util.GetStringFromMap(purchasingData, "warranty_date"),
				VersionLock:       util.GetIntFromMap(purchasingData, "version_lock"),
			}
		}
	}

	// Extract anchor_certificates
	if v, ok := d.GetOk("anchor_certificates"); ok {
		anchorCertificates := make([]string, len(v.([]interface{})))
		for i, cert := range v.([]interface{}) {
			anchorCertificates[i] = cert.(string)
		}
		prestage.AnchorCertificates = anchorCertificates
	}

	prestage.EnrollmentCustomizationId = util.GetStringFromInterface(d.Get("enrollment_customization_id"))
	prestage.Language = util.GetStringFromInterface(d.Get("language"))
	prestage.Region = util.GetStringFromInterface(d.Get("region"))
	prestage.AutoAdvanceSetup = util.GetBoolFromInterface(d.Get("auto_advance_setup"))
	prestage.InstallProfilesDuringSetup = util.GetBoolFromInterface(d.Get("install_profiles_during_setup"))

	// Extract prestage_installed_profile_ids
	if v, ok := d.GetOk("prestage_installed_profile_ids"); ok {
		profileIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			profileIDs[i] = id.(string)
		}
		prestage.PrestageInstalledProfileIds = profileIDs
	}

	// Extract custom_package_ids
	if v, ok := d.GetOk("custom_package_ids"); ok {
		packageIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			packageIDs[i] = id.(string)
		}
		prestage.CustomPackageIds = packageIDs
	}

	prestage.CustomPackageDistributionPointId = util.GetStringFromInterface(d.Get("custom_package_distribution_point_id"))
	prestage.EnableRecoveryLock = util.GetBoolFromInterface(d.Get("enable_recovery_lock"))
	prestage.RecoveryLockPasswordType = util.GetStringFromInterface(d.Get("recovery_lock_password_type"))
	prestage.RecoveryLockPassword = util.GetStringFromInterface(d.Get("recovery_lock_password"))
	prestage.RotateRecoveryLockPassword = util.GetBoolFromInterface(d.Get("rotate_recovery_lock_password"))
	prestage.ProfileUuid = util.GetStringFromInterface(d.Get("profile_uuid"))
	prestage.SiteId = util.GetStringFromInterface(d.Get("site_id"))
	prestage.VersionLock = util.GetIntFromInterface(d.Get("version_lock"))

	// Extract account_settings
	if v, ok := d.GetOk("account_settings"); ok {
		accountSettingsList := v.([]interface{})
		if len(accountSettingsList) > 0 {
			accountData := accountSettingsList[0].(map[string]interface{})
			prestage.AccountSettings = jamfpro.ComputerPrestageSubsetAccountSettings{
				ID:                                      util.GetStringFromMap(accountData, "id"),
				PayloadConfigured:                       util.GetBoolFromMap(accountData, "payload_configured"),
				LocalAdminAccountEnabled:                util.GetBoolFromMap(accountData, "local_admin_account_enabled"),
				AdminUsername:                           util.GetStringFromMap(accountData, "admin_username"),
				AdminPassword:                           util.GetStringFromMap(accountData, "admin_password"),
				HiddenAdminAccount:                      util.GetBoolFromMap(accountData, "hidden_admin_account"),
				LocalUserManaged:                        util.GetBoolFromMap(accountData, "local_user_managed"),
				UserAccountType:                         util.GetStringFromMap(accountData, "user_account_type"),
				VersionLock:                             util.GetIntFromMap(accountData, "version_lock"),
				PrefillPrimaryAccountInfoFeatureEnabled: util.GetBoolFromMap(accountData, "prefill_primary_account_info_feature_enabled"),
				PrefillType:                             util.GetStringFromMap(accountData, "prefill_type"),
				PrefillAccountFullName:                  util.GetStringFromMap(accountData, "prefill_account_full_name"),
				PrefillAccountUserName:                  util.GetStringFromMap(accountData, "prefill_account_user_name"),
				PreventPrefillInfoFromModification:      util.GetBoolFromMap(accountData, "prevent_prefill_info_from_modification"),
			}
		}
	}

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	// Serialize and pretty-print the department object as JSON
	computerPrestageJSON, err := json.MarshalIndent(prestage, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceJSONMarshalFailure(subCtx, JamfProResourceComputerPrestage, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to JSON
	logging.LogTFConstructedJSONResource(subCtx, JamfProResourceComputerPrestage, string(computerPrestageJSON))

	return prestage, nil
}

// ResourceJamfProComputerPrestageCreate is responsible for creating a new computer prestage in Jamf Pro with terraform.
// The function:
// 1. Constructs the computer prestage data using the provided Terraform configuration.
// 2. Calls the API to create the computer prestage in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created computer prestage.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProComputerPrestageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var createdAttribute *jamfpro.ResponseComputerPrestageCreate
	var apiErrorCode int

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the computerPrestage object outside the retry loop to avoid reconstructing it on each retry
	computerPrestage, err := constructJamfProComputerPrestage(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceComputerPrestage, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceComputerPrestage)

	// Retry the API call to create the computerPrestage in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		createdAttribute, apiErr = conn.CreateComputerPrestage(computerPrestage)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailure(subCtx, JamfProResourceComputerPrestage, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceComputerPrestage, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the computerPrestage and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceComputerPrestage, createdAttribute.ID)
	d.SetId(createdAttribute.ID)

	// Retry reading the computerPrestage to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerPrestageRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceComputerPrestage, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the computerPrestage, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceComputerPrestage, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceComputerPrestage, d.Id())
	}

	return diags
}

// ResourceJamfProComputerPrestageRead is responsible for reading the current state of a Building Resource from the remote system.
// The function:
// 1. Fetches the building's current state using its ID. If it fails, then obtain the building's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the building being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProComputerPrestageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var computerPrestage *jamfpro.ResourceComputerPrestage
	var apiErrorCode int // Variable to capture the API error code
	resourceID := d.Id()

	// Try fetching the computerPrestage using the ID
	computerPrestage, err := conn.GetComputerPrestageByID(resourceID)
	if err != nil {
		// Extract and log the API error code if available
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceComputerPrestage, resourceID, err.Error(), apiErrorCode)
		d.SetId("") // Remove from Terraform state
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceComputerPrestage, resourceID)
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourceComputerPrestage, resourceID)

	// Check if prestage data exists
	if computerPrestage != nil {
		// Construct a map of computer prestage attributes
		prestageAttributes := map[string]interface{}{
			"display_name":                          computerPrestage.DisplayName,
			"mandatory":                             computerPrestage.Mandatory,
			"mdm_removable":                         computerPrestage.MDMRemovable,
			"support_phone_number":                  computerPrestage.SupportPhoneNumber,
			"support_email_address":                 computerPrestage.SupportEmailAddress,
			"department":                            computerPrestage.Department,
			"default_prestage":                      computerPrestage.DefaultPrestage,
			"enrollment_site_id":                    computerPrestage.EnrollmentSiteId,
			"keep_existing_site_membership":         computerPrestage.KeepExistingSiteMembership,
			"keep_existing_location_information":    computerPrestage.KeepExistingLocationInformation,
			"authentication_prompt":                 computerPrestage.AuthenticationPrompt,
			"prevent_activation_lock":               computerPrestage.PreventActivationLock,
			"enable_device_based_activation_lock":   computerPrestage.EnableDeviceBasedActivationLock,
			"device_enrollment_program_instance_id": computerPrestage.DeviceEnrollmentProgramInstanceId,
			"skip_setup_items":                      computerPrestage.SkipSetupItems,
			"location_information":                  computerPrestage.LocationInformation,
			"purchasing_information":                computerPrestage.PurchasingInformation,
			"anchor_certificates":                   computerPrestage.AnchorCertificates,
			"enrollment_customization_id":           computerPrestage.EnrollmentCustomizationId,
			"language":                              computerPrestage.Language,
			"region":                                computerPrestage.Region,
			"auto_advance_setup":                    computerPrestage.AutoAdvanceSetup,
			"install_profiles_during_setup":         computerPrestage.InstallProfilesDuringSetup,
			"prestage_installed_profile_ids":        computerPrestage.PrestageInstalledProfileIds,
			"custom_package_ids":                    computerPrestage.CustomPackageIds,
			"custom_package_distribution_point_id":  computerPrestage.CustomPackageDistributionPointId,
			"enable_recovery_lock":                  computerPrestage.EnableRecoveryLock,
			"recovery_lock_password_type":           computerPrestage.RecoveryLockPasswordType,
			"recovery_lock_password":                computerPrestage.RecoveryLockPassword,
			"rotate_recovery_lock_password":         computerPrestage.RotateRecoveryLockPassword,
			"profile_uuid":                          computerPrestage.ProfileUuid,
			"site_id":                               computerPrestage.SiteId,
			"version_lock":                          computerPrestage.VersionLock,
			"account_settings":                      computerPrestage.AccountSettings,
		}

		// Handle nested location_information
		if locationInformation := computerPrestage.LocationInformation; locationInformation != (jamfpro.ComputerPrestageSubsetLocationInformation{}) {
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
		if purchasingInformation := computerPrestage.PurchasingInformation; purchasingInformation != (jamfpro.ComputerPrestageSubsetPurchasingInformation{}) {
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
		prestageAttributes["anchor_certificates"] = computerPrestage.AnchorCertificates
		prestageAttributes["enrollment_customization_id"] = computerPrestage.EnrollmentCustomizationId
		prestageAttributes["language"] = computerPrestage.Language
		prestageAttributes["region"] = computerPrestage.Region
		prestageAttributes["auto_advance_setup"] = computerPrestage.AutoAdvanceSetup
		prestageAttributes["install_profiles_during_setup"] = computerPrestage.InstallProfilesDuringSetup
		prestageAttributes["prestage_installed_profile_ids"] = computerPrestage.PrestageInstalledProfileIds
		prestageAttributes["custom_package_ids"] = computerPrestage.CustomPackageIds
		prestageAttributes["custom_package_distribution_point_id"] = computerPrestage.CustomPackageDistributionPointId
		prestageAttributes["enable_recovery_lock"] = computerPrestage.EnableRecoveryLock
		prestageAttributes["recovery_lock_password_type"] = computerPrestage.RecoveryLockPasswordType
		prestageAttributes["recovery_lock_password"] = computerPrestage.RecoveryLockPassword
		prestageAttributes["rotate_recovery_lock_password"] = computerPrestage.RotateRecoveryLockPassword
		prestageAttributes["profile_uuid"] = computerPrestage.ProfileUuid
		prestageAttributes["site_id"] = computerPrestage.SiteId
		prestageAttributes["version_lock"] = computerPrestage.VersionLock
		// Handle nested account_settings
		if accountSettings := computerPrestage.AccountSettings; accountSettings != (jamfpro.ComputerPrestageSubsetAccountSettings{}) {
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
		for key, value := range prestageAttributes {
			if err := d.Set(key, value); err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
		}
	} else {
		// If the prestage is not found, clear the ID from the state
		d.SetId("")
	}

	return diags
}

// ResourceJamfProComputerPrestageUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProComputerPrestageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Construct the resource object
	computerPrestage, err := constructJamfProComputerPrestage(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceComputerPrestage, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceComputerPrestage)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateComputerPrestageByID(resourceID, computerPrestage)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateComputerPrestageByName(resourceName, computerPrestage)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceComputerPrestage, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the computerPrestage to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerPrestageRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceComputerPrestage, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceComputerPrestage, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceComputerPrestage, resourceID)
	}

	return nil
}

// ResourceJamfProComputerPrestageDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProComputerPrestageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteComputerPrestageByID(resourceID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteComputerPrestageByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceComputerPrestage, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceComputerPrestage, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
