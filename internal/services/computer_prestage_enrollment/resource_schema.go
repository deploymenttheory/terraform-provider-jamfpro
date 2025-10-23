package computer_prestage_enrollment

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProComputerPrestageEnrollment defines the schema for managing Jamf Pro Computer Prestages in Terraform.
func ResourceJamfProComputerPrestageEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The display name of the computer prestage enrollment.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Make MDM Profile Mandatory and require the user to apply the MDM profile. Computers with macOS 10.15 or later automatically require the user to apply the MDM profile",
			},
			"mdm_removable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allow MDM Profile Removal and allow the user to remove the MDM profile.",
			},
			"support_phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Support phone number for the organization. Can be left blank.",
			},
			"support_email_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Support email address for the organization. Can be left blank.",
			},
			"department": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The department the computer prestage is assigned to. Can be left blank.",
			},
			"default_prestage": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if this is the default computer prestage enrollment configuration. If yes then new devices will be automatically assigned to this PreStage enrollment",
			},
			"enrollment_site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The jamf pro Site ID that computers will be added to during enrollment. Should be set to -1, if not used.",
			},
			"keep_existing_site_membership": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if enrolled should use existing site membership, if applicable",
			},
			"keep_existing_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if enrolled should use existing location information, if applicable",
			},
			"require_authentication": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the user is required to provide username and password on computers with macOS 10.10 or later.",
			},
			"authentication_prompt": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication Message to display to the user. Used when Require Authentication is enabled. Can be left blank.",
			},
			"prevent_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Prevent user from enabling Activation Lock.",
			},
			"enable_device_based_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if device-based activation lock should be enabled.",
			},
			"device_enrollment_program_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Automated Device Enrollment instance ID to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment",
			},
			"platform_sso_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if platform single sign-on (SSO) is enabled.",
			},
			"platform_sso_app_bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The app bundle ID for used for platform SSO.",
			},
			"skip_setup_items": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Selected items are not displayed in the Setup Assistant during macOS device setup within Apple Device Enrollment (ADE).",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biometric": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip biometric setup.",
						},
						"terms_of_address": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip terms of address setup.",
						},
						"file_vault": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip FileVault setup.",
						},
						"icloud_diagnostics": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip iCloud diagnostics setup.",
						},
						"diagnostics": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip diagnostics setup.",
						},
						"accessibility": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip accessibility setup.",
						},
						"apple_id": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Apple ID setup.",
						},
						"screen_time": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Screen Time setup.",
						},
						"siri": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Siri setup.",
						},
						"display_tone": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Display Tone setup. (Deprecated)",
						},
						"restore": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Restore setup.",
						},
						"appearance": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Appearance setup.",
						},
						"privacy": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Privacy setup.",
						},
						"payment": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Payment setup.",
						},
						"registration": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Registration setup.",
						},
						"tos": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Terms of Service setup.",
						},
						"icloud_storage": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip iCloud Storage setup.",
						},
						"location": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Location setup.",
						},
						"intelligence": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Apple Intelligence setup.",
						},
						"enable_lockdown_mode": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip lockdown mode setup.",
						},
						"welcome": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip welcome setup.",
						},
						"wallpaper": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip wallpaper setup.",
						},
						"software_update": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip software update setup.",
						},
						"additional_privacy_settings": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip additional privacy settings setup.",
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
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the location information.",
						},
						"version_lock": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The version lock of the location information. Optimistic locking" +
								"is a mechanism that prevents concurrent operations from taking place on a given" +
								"resource. Jamf Pro does this to safeguard resources and workflows that are" +
								"sensitive to frequent updates, ensuring that one update has completed before" +
								"any additional requests can be processed. Valid request handling is managed by" +
								"the construct function.",
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username for the location information. Can be left blank.",
						},
						"realname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The real name associated with this location. Can be left blank.",
						},
						"phone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The phone number associated with this location. Can be left blank.",
						},
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The email address associated with this location. Can be left blank.",
						},
						"room": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The room associated with this location. Can be left blank.",
						},
						"position": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The position associated with this location. Can be left blank.",
						},
						"department_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The jamf pro department ID associated with this computer prestage. Set to -1 if not used.",
						},
						"building_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The building ID associated with this computer prestage. Set to -1 if not used.",
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
							Computed:    true,
							Description: "The ID of the purchasing information.",
						},
						"version_lock": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The version lock value of the purchasing_information. Optimistic locking" +
								"is a mechanism that prevents concurrent operations from taking place on a given" +
								"resource. Jamf Pro does this to safeguard resources and workflows that are" +
								"sensitive to frequent updates, ensuring that one update has completed before" +
								"any additional requests can be processed. Valid request handling is managed by" +
								"the construct function.",
						},
						"leased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is leased. Default value to false if unused.",
						},
						"purchased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is purchased. Default value to true if unused.",
						},
						"apple_care_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AppleCare ID. Can be left blank.",
						},
						"po_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase order number. Can be left blank.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The vendor name. Can be left blank.",
						},
						"purchase_price": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase price. Can be left blank.",
						},
						"life_expectancy": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The life expectancy in years. Set to 0 if unused.",
						},
						"purchasing_account": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing account. Can be left blank.",
						},
						"purchasing_contact": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing contact. Can be left blank.",
						},
						"lease_date": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The lease date in YYYY-MM-DD format. Use '1970-01-01' if unused.",
							ValidateFunc: validateDateFormat,
						},
						"po_date": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The purchase order date in YYYY-MM-DD format. Use '1970-01-01' if unused",
							ValidateFunc: validateDateFormat,
						},
						"warranty_date": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The warranty date in YYYY-MM-DD format. Use '1970-01-01' if unused",
							ValidateFunc: validateDateFormat,
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
				Required:    true,
				Description: "The enrollment customization ID. Set to 0 if unused.",
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The language setting defined for the computer prestage. Leverages ISO 639-1 (two-letter language codes): https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes . Ensure you define a code supported by jamf pro. Can be left blank.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region setting defined for the computer prestage. Leverages ISO 3166-1 alpha-2 (two-letter country codes): https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2 . Ensure you define a code supported by jamf pro. Can be left blank.",
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
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IDs of the macOS configuration profiles installed during PreStage enrollment. requires decending order of profile IDs so uses a set rather than a list. can be left blank.",
			},
			"custom_package_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Define the Enrollment Packages by their package ID to " +
					"add an enrollment package to the PreStage enrollment. Compatible packages " +
					"must be built as flat, distribution style .pkg files and be signed by a " +
					"certificate that is trusted by managed computers. Can be left blank.",
			},
			"custom_package_distribution_point_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Set the Enrollment Packages distribution point by it's ID." +
					"Valid values are: None using '-1', Cloud Distribution Point (Jamf Cloud)" +
					"by using '-2', else all other valid valid values correspond to the" +
					"ID of the distribution point.",
			},
			"enable_recovery_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Configure how the Recovery Lock password is set on computers with macOS 11.5 or later.",
			},
			"recovery_lock_password_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Method to use to set Recovery Lock password.'MANUAL' results in " +
					"user having to enter a password. (Applies to all users) 'RANDOM' results in" +
					"automatic generation of a random password being set for the device. 'MANUAL' is the default.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
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
				Description: "Generate new Recovery Lock password 60 minutes after the password is viewed in Jamf Pro. Can be left blank.",
			},
			"rotate_recovery_lock_password": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the recovery lock password should be rotated.",
			},
			"prestage_minimum_os_target_version_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enforce a minimum macOS target version type for the prestage enrollment. Required.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validTypes := map[string]bool{
						"NO_ENFORCEMENT":                  true,
						"MINIMUM_OS_LATEST_VERSION":       true,
						"MINIMUM_OS_LATEST_MAJOR_VERSION": true,
						"MINIMUM_OS_LATEST_MINOR_VERSION": true,
						"MINIMUM_OS_SPECIFIC_VERSION":     true,
					}
					if _, valid := validTypes[v]; !valid {
						errs = append(errs, fmt.Errorf("%q must be one of 'NO_ENFORCEMENT', 'MINIMUM_OS_LATEST_VERSION', 'MINIMUM_OS_LATEST_MAJOR_VERSION', 'MINIMUM_OS_LATEST_MINOR_VERSION', 'MINIMUM_OS_SPECIFIC_VERSION', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"minimum_os_specific_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The minimum macOS version to enforce for the prestage enrollment. Only used if prestate_minimum_os_target_version_type is set to MINIMUM_OS_SPECIFIC_VERSION.",
			},
			"profile_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile UUID of the Automated Device Enrollment instance to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The jamf pro site ID. Set to -1 if not used.",
			},
			"version_lock": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The version lock value of the purchasing_information. Optimistic locking" +
					"is a mechanism that prevents concurrent operations from taking place on a given" +
					"resource. Jamf Pro does this to safeguard resources and workflows that are" +
					"sensitive to frequent updates, ensuring that one update has completed before" +
					"any additional requests can be processed. Valid request handling is managed by" +
					"the construct function.",
			},
			"account_settings": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of Account Settings.",
						},
						"version_lock": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The version lock value of the account settings block. Optimistic locking" +
								"is a mechanism that prevents concurrent operations from taking place on a given" +
								"resource. Jamf Pro does this to safeguard resources and workflows that are" +
								"sensitive to frequent updates, ensuring that one update has completed before" +
								"any additional requests can be processed. Valid request handling is managed by" +
								"the construct function.",
						},
						"payload_configured": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the payload is configured.",
						},
						"local_admin_account_enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the local admin account is enabled.",
						},
						"admin_username": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The admin username. Can be left blank if not used.",
						},
						"admin_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The admin password. Can be left blank if not used.",
						},
						"hidden_admin_account": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the admin account is hidden.",
						},
						"local_user_managed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the local user is managed.",
						},
						"user_account_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of user account (ADMINISTRATOR, STANDARD, SKIP).",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
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
						"prefill_primary_account_info_feature_enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if prefilling primary account info feature is enabled.",
						},
						"prefill_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pre-fill primary account information type (CUSTOM, DEVICE_OWNER, or UNKNOWN). Set as UNKNOWN if you wish to leave it unconfigured.",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validTypes := map[string]bool{
									"CUSTOM":       true,
									"DEVICE_OWNER": true,
									"UNKNOWN":      true,
								}
								if _, valid := validTypes[v]; !valid {
									errs = append(errs, fmt.Errorf("%q must be one of 'CUSTOM', 'DEVICE_OWNER', 'UNKNOWN' got: %s", key, v))
								}
								return warns, errs
							},
						},
						"prefill_account_full_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of information to use to pre-fill the primary account full name with. Can be left blank.",
						},
						"prefill_account_user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of information to use to pre-fill the primary account user name with. Can be left blank.",
						},
						"prevent_prefill_info_from_modification": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Lock prefill primary account information from modification.",
						},
					},
				},
			},
		},
	}
}
