package userinitiatedenrollment

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProEnrollmentConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotImplementedCreate,
		ReadContext:   resourceEnrollmentRead,
		UpdateContext: resourceEnrollmentUpdate,
		DeleteContext: resourceNotImplementedDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			// General (page 1) Enrollment restrictions and settings
			"restrict_reenrollment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Restrict re-enrollment to authorized users only by only allowing re-enrollment of mobile devices and computers if the user has the applicable privilege (“Mobile Devices” or “Computers”) or their username matches the Username field in User and Location information.",
			},
			"install_single_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Skip certificate installation during enrollment. Certificate installation step is skipped during enrollment if your environment has an SSL certificate that was obtained from an internal CA or a trusted third-party vendor.",
			},
			"signing_mdm_profile_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Ensure that the certificate signs configuration profiles sent to computers and mobile devices, and appears as verified to users during user-initiated enrollment.",
			},
			// MDM Signing Certificate
			"mdm_signing_certificate": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "MDM signing certificate configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the certificate file in .p12 format. e.g 'my_test_certificate.p12'.",
						},
						"identity_keystore": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Base64-encoded certificate in .p12 format.",
						},
						"keystore_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for the certificate keystore",
						},
					},
				},
			},
			// Computer (Page 3)
			"create_management_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to create a management account during enrollment",
			},
			"management_username": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     "jamfadmin",
				Description: "Account to be used for computers enrolled via PreStage enrollment or user-initiated enrollment. This account is created by the Jamf management framework.",
			},
			"hide_management_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to hide the managed local administrator account from users",
			},
			"allow_ssh_only_management_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to make the managed local administrator account the only account that has SSH (Remote Login) access to computers",
			},
			"launch_self_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Ensure that computers launch Self Service immediately after they are enrolled",
			},
			"ensure_ssh_running": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to ensure SSH service is running after enrollment",
			},
			"account_driven_device_macos_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether account-driven device macOS enrollment is enabled for institutionally owned computers",
			},
			"sign_quick_add": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to ensure that the QuickAdd package is signed so that it appears as verified to users when enrolling via user-initiated enrollment with a QuickAdd package.",
			},
			// Developer certificate
			"developer_certificate_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Developer certificate identity data",
			},
			"developer_certificate_identity_details": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Details of the developer certificate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the certificate file in .p12 format. e.g 'my_test_certificate.p12'.",
						},
						"identity_keystore": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Base64-encoded certificate in .p12 format.",
						},
						"keystore_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for the certificate keystore",
						},
					},
				},
			},
			// Devices (Page 4)
			"account_driven_device_ios_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether account-driven device iOS enrollment is enabled",
			},
			"account_driven_user_visionos_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether account-driven user visionOS enrollment is enabled",
			},
			"account_driven_device_visionos_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether account-driven device visionOS enrollment is enabled",
			},

			// Flush settings
			"flush_location_information": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to flush location information during re-enrollment",
			},
			"flush_location_history_information": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to flush location history information during re-enrollment",
			},
			"flush_policy_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to flush policy history during re-enrollment",
			},
			"flush_extension_attributes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to flush extension attributes during re-enrollment",
			},
			"flush_mdm_commands_on_reenroll": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED",
				Description:  "Determines which MDM commands to flush during re-enrollment",
				ValidateFunc: validation.StringInSlice([]string{"DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED", "DELETE_EVERYTHING", "DELETE_NOTHING", "DELETE_ERRORS"}, false),
			},

			// MDM Signing Certificate Details
			"mdm_signing_certificate_details": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Details of the MDM signing certificate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subject of the MDM signing certificate",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Serial number of the MDM signing certificate",
						},
					},
				},
			},

			// Platform enrollment settings
			"macos_enterprise_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether macOS enterprise enrollment is enabled",
			},

			"ios_enterprise_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether iOS enterprise enrollment is enabled",
			},
			"ios_personal_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether iOS personal enrollment is enabled",
			},
			"personal_device_enrollment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "USERENROLLMENT",
				Description:  "Type of enrollment for personal devices",
				ValidateFunc: validation.StringInSlice([]string{"USERENROLLMENT", "DEVICEENROLLMENT"}, false),
			},
			"account_driven_user_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether account-driven user enrollment is enabled",
			},
			// /api/v3/enrollment/languages/en-gb
			// Messaging Configuration
			"messaging": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Localized text configuration for enrollment pages and messages",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Language code for this messaging configuration (e.g., 'en-gb')",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display name of the language (e.g., 'English (UK)')",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Title text for the enrollment page",
						},
						"login_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description text shown on the login page",
						},
						"certificate_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for certificate button",
						},
						"certificate_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description text for certificate profile",
						},
						"certificate_profile_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the certificate profile",
						},
						"certificate_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed for certificate information",
						},
						"check_enrollment_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message displayed when checking enrollment status",
						},
						"check_now_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for the 'Check Now' button",
						},
						"complete_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message displayed when enrollment is complete",
						},
						"device_class_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for device class button",
						},
						"device_class_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description text for device class selection",
						},
						"device_class_enterprise": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for enterprise device class option",
						},
						"device_class_enterprise_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for enterprise device class",
						},
						"device_class_personal": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for personal device class option",
						},
						"device_class_personal_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for personal device class",
						},
						"enterprise_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for enterprise enrollment button",
						},
						"enterprise_eula": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "EULA text for enterprise enrollment",
						},
						"enterprise_pending": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message shown while enterprise enrollment is pending",
						},
						"enterprise_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for enterprise enrollment profile",
						},
						"enterprise_profile_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the enterprise enrollment profile",
						},
						"enterprise_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed for enterprise enrollment information",
						},
						"eula_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for EULA acceptance button",
						},
						"failed_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message shown when enrollment fails",
						},
						"login_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for login button",
						},
						"logout_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for logout button",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label for password field",
						},
						"personal_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for personal enrollment button",
						},
						"personal_eula": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "EULA text for personal enrollment",
						},
						"personal_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for personal enrollment profile",
						},
						"personal_profile_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the personal enrollment profile",
						},
						"personal_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed for personal enrollment information",
						},
						"quick_add_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for QuickAdd button",
						},
						"quick_add_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name displayed for QuickAdd package",
						},
						"quick_add_pending": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message shown while QuickAdd enrollment is pending",
						},
						"quick_add_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed for QuickAdd enrollment information",
						},
						"site_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description text for site selection",
						},
						"try_again_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for try again button",
						},
						"user_enrollment_button": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text for user enrollment button",
						},
						"user_enrollment_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for user enrollment profile",
						},
						"user_enrollment_profile_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the user enrollment profile",
						},
						"user_enrollment_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed for user enrollment information",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label for username field",
						},
					},
				},
			},
			// api/v3/enrollment/access-groups
			// Directory Service Groups
			"directory_service_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Configuration for directory service groups allowed to enroll",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_driven_user_enrollment_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether account-driven user enrollment is enabled for this group",
						},
						"enterprise_enrollment_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether enterprise enrollment is enabled for this group",
						},
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID of the directory service group",
						},
						"ldap_server_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the LDAP server associated with this group",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the directory service group",
						},
						"personal_enrollment_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether personal enrollment is enabled for this group",
						},
						"require_eula": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to require accepting an EULA during enrollment for this group",
						},
						"site_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the site associated with this directory service group",
						},
					},
				},
			},
		},
	}
}
