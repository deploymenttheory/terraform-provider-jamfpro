// mobiledeviceprestageenrollments_resource.go
package mobiledeviceprestageenrollments

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProMobileDevicePrestageEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: validateAuthenticationPrompt,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the mobile device prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the mobile device prestage enrollment.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Make MDM Profile Mandatory.",
			},
			"mdm_removable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allow MDM Profile Removal.",
			},
			"support_phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Support phone number for the organization.",
			},
			"support_email_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Support email address for the organization.",
			},
			"department": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Department associated with the prestage.",
			},
			"default_prestage": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether this is the default prestage enrollment configuration.",
			},
			"enrollment_site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
				Description: "Site ID for device enrollment.",
			},
			"keep_existing_site_membership": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Maintain existing site membership during enrollment.",
			},
			"keep_existing_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Maintain existing location information during enrollment.",
			},
			"require_authentication": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Require authentication during enrollment.",
			},
			"authentication_prompt": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Message displayed when authentication is required.",
			},
			"prevent_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Prevent Activation Lock on the device.",
			},
			"enable_device_based_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable device-based Activation Lock.",
			},
			"device_enrollment_program_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Automated Device Enrollment instance ID.",
			},
			"skip_setup_items": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Selected items are not displayed in the Setup Assistant during mobile device setup within Apple Device Enrollment (ADE).",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Location setup during device enrollment.",
						},
						"privacy": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Privacy setup during device enrollment.",
						},
						"biometric": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Biometric setup during device enrollment.",
						},
						"software_update": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Software Update setup during device enrollment.",
						},
						"diagnostics": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Diagnostics setup during device enrollment.",
						},
						"imessage_and_facetime": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip iMessage and FaceTime setup during device enrollment.",
						},
						"intelligence": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Intelligence setup during device enrollment.",
						},
						"tv_room": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip TV Room setup during device enrollment.",
						},
						"passcode": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Passcode setup during device enrollment.",
						},
						"sim_setup": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip SIM setup during device enrollment.",
						},
						"screen_time": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Screen Time setup during device enrollment.",
						},
						"restore_completed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Restore Completed setup during device enrollment.",
						},
						"tv_provider_sign_in": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip TV Provider Sign In during device enrollment.",
						},
						"siri": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Siri setup during device enrollment.",
						},
						"restore": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Restore setup during device enrollment.",
						},
						"screen_saver": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Screen Saver setup during device enrollment.",
						},
						"home_button_sensitivity": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Home Button Sensitivity setup during device enrollment.",
						},
						"cloud_storage": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Cloud Storage setup during device enrollment.",
						},
						"action_button": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Action Button setup during device enrollment.",
						},
						"transfer_data": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Transfer Data setup during device enrollment.",
						},
						"enable_lockdown_mode": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Enable Lockdown Mode setup during device enrollment.",
						},
						"zoom": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Zoom setup during device enrollment.",
						},
						"preferred_language": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Preferred Language setup during device enrollment.",
						},
						"voice_selection": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Voice Selection setup during device enrollment.",
						},
						"tv_home_screen_sync": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip TV Home Screen Sync during device enrollment.",
						},
						"safety": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Safety setup during device enrollment.",
						},
						"terms_of_address": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Terms of Address setup during device enrollment.",
						},
						"express_language": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Express Language setup during device enrollment.",
						},
						"camera_button": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Camera Button setup during device enrollment.",
						},
						"apple_id": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Apple ID setup during device enrollment.",
						},
						"display_tone": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Display Tone setup during device enrollment.",
						},
						"watch_migration": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Watch Migration setup during device enrollment.",
						},
						"update_completed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Update Completed setup during device enrollment.",
						},
						"appearance": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Appearance setup during device enrollment.",
						},
						"android": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Android Migration setup during device enrollment.",
						},
						"payment": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Payment setup during device enrollment.",
						},
						"onboarding": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Onboarding setup during device enrollment.",
						},
						"tos": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Terms of Service setup during device enrollment.",
						},
						"welcome": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Welcome setup during device enrollment.",
						},
						"tap_to_setup": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Skip Tap to Setup during device enrollment.",
						},
					},
				},
			},
			"location_information": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Location information associated with the Jamf Pro mobile device prestage.",
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
							Optional:    true,
							Default:     "-1",
							Description: "The jamf pro department ID associated with this computer prestage. Set to -1 if not used.",
						},
						"building_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "-1",
							Description: "The building ID associated with this computer prestage. Set to -1 if not used.",
						},
					},
				},
			},
			"purchasing_information": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Purchasing information associated with the mobile device prestage.",
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
							Optional:     true,
							Default:      "1970-01-01",
							Description:  "The lease date in YYYY-MM-DD format. Use '1970-01-01' if unused.",
							ValidateFunc: validateDateFormat,
						},
						"po_date": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "1970-01-01",
							Description:  "The purchase order date in YYYY-MM-DD format. Use '1970-01-01' if unused",
							ValidateFunc: validateDateFormat,
						},
						"warranty_date": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "1970-01-01",
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
				Description: "The language setting defined for the mobile device prestage. Leverages ISO 639-1 (two-letter language codes): https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes . Ensure you define a code supported by jamf pro. Can be left blank.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region setting defined for the mobile device prestage. Leverages ISO 3166-1 alpha-2 (two-letter country codes): https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2 . Ensure you define a code supported by jamf pro. Can be left blank.",
			},
			"auto_advance_setup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if setup should auto-advance.",
			},
			"allow_pairing": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allow device pairing.",
			},
			"multi_user": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable multi-user mode.",
			},
			"supervised": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Device is supervised.",
			},
			"maximum_shared_accounts": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum number of shared accounts.",
			},
			"configure_device_before_setup_assistant": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Configure device before Setup Assistant.",
			},
			"names": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assign_names_using": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Default Names",
								"List of Names",
								"Serial Numbers",
								"Single Name",
							}, false),
							Description: "Method to use for assigning device names. Valid values are: 'Default Names', 'List of Names', 'Serial Numbers' or 'Single Name'.",
						},
						"prestage_device_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier of the device name entry.",
									},
									"device_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name to be assigned to the device.",
									},
									"used": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates if this device name has been used.",
									},
								},
							},
							Description: "List of predefined device names when using 'List of Names' assignment method.",
						},
						"device_name_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The prefix to use when naming devices with 'Serial Numbers' method.",
						},
						"device_name_suffix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The suffix to use when naming devices with 'Serial Numbers' method.",
						},
						"single_device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name to use when using 'Single Name' assignment method.",
						},
						"manage_names": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if device names should be managed by this prestage.",
						},
						"device_naming_configured": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if device naming has been configured for this prestage.",
						},
					},
				},
			},
			"send_timezone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if timezone should be sent to the device.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The timezone to be set on the device. Default is UTC",
			},
			"storage_quota_size_megabytes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The storage quota size in megabytes.",
			},
			"use_storage_quota_size": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if storage quota size should be enforced.",
			},
			"temporary_session_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the session should be temporary only.",
			},
			"enforce_temporary_session_timeout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if temporary session timeout should be enforced.",
			},
			"temporary_session_timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout duration for temporary sessions in minutes.",
			},
			"enforce_user_session_timeout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if user session timeout should be enforced.",
			},
			"user_session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout duration for user sessions in minutes.",
			},
			"profile_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile UUID of the Automated Device Enrollment instance to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
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
			"prestage_minimum_os_target_version_type_ios": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of minimum OS version enforcement for iOS devices.",
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
			"minimum_os_specific_version_ios": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specific minimum OS version required for iOS devices when using MINIMUM_OS_SPECIFIC_VERSION type.",
			},
			"prestage_minimum_os_target_version_type_ipad": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of minimum OS version enforcement for iPadOS devices.",
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
			"minimum_os_specific_version_ipad": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specific minimum OS version required for iPadOS devices when using MINIMUM_OS_SPECIFIC_VERSION type.",
			},
			"rts_config_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the RTS configuration profile.",
			},
		},
	}
}
