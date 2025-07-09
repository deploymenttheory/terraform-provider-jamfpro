package user_initiated_enrollment_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProUserInitatedEnrollmentSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: customizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			// General (page 1) Enrollment restrictions and settings
			// /api/v4/enrollment
			"restrict_reenrollment_to_authorized_users_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Maps to payload field 'restrictReenrollment'.Restrict re-enrollment to authorized users only by only allowing re-enrollment of mobile devices and computers if the user has the applicable privilege (“Mobile Devices” or “Computers”) or their username matches the Username field in User and Location information.",
			},
			"skip_certificate_installation_during_enrollment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Maps to payload field 'installSingleProfile'. Certificate installation step is skipped during enrollment if your environment has an SSL certificate that was obtained from an internal CA or a trusted third-party vendor.",
			},
			"third_party_signing_certificate": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Third-party signing certificate configuration to ensure that the certificate signs configuration profiles sent to computers and mobile devices, and appears as verified to users during user-initiated enrollment. Maps to 'mdmSigningCertificate'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Maps to payload field 'signingMdmProfileEnabled'. Whether to use a third-party signing certificate.",
						},
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
			// Messaging (page2)
			// /api/v3/enrollment/languages/ab <- Two letter ISO 639-1 Language Code
			/*
				Field Mapping Reference Table:
				+------------------------------------------+-----------------------------+----------+
				| Terraform Schema Field (GUI Name)        | API Request Field           | Required |
				+------------------------------------------+-----------------------------+----------+
				| language_code                            | (not mapped, for Terraform) | Yes      |
				| name                                     | name                        | Yes      |
				| page_title                               | title                       | Yes      |
				| login_page_text                          | loginDescription            | No       |
				| username_text                            | username                    | Yes      |
				| password_text                            | password                    | Yes      |
				| login_button_text                        | loginButton                 | Yes      |
				| device_ownership_page_text               | deviceClassDescription      | No       |
				| personal_device_button_name              | deviceClassPersonal         | Yes      |
				| institutional_ownership_button_name      | deviceClassEnterprise       | Yes      |
				| personal_device_management_description   | deviceClassPersonalDesc...  | No       |
				| institutional_device_management_desc...  | deviceClassEnterpriseDe...  | No       |
				| enroll_device_button_name                | deviceClassButton           | Yes      |
				| eula_personal_devices                    | enterpriseEula              | No       |
				| eula_institutional_devices               | personalEula                | No       |
				| accept_button_text                       | eulaButton                  | Yes      |
				| site_selection_text                      | siteDescription             | No       |
				| ca_certificate_installation_text         | certificateText             | No       |
				| ca_certificate_name                      | certificateProfileName      | Yes      |
				| ca_certificate_description               | certificateProfileDescr...  | No       |
				| ca_certificate_install_button_name       | certificateButton           | Yes      |
				| institutional_mdm_profile_installation_t | enterpriseText              | No       |
				| institutional_mdm_profile_name           | enterpriseProfileName       | Yes      |
				| institutional_mdm_profile_description    | enterpriseProfileDescri...  | No       |
				| institutional_mdm_profile_pending_text   | enterprisePending           | No       |
				| institutional_mdm_profile_install_button | enterpriseButton            | Yes      |
				| personal_mdm_profile_installation_text   | personalText                | No       |
				| personal_mdm_profile_name                | personalProfileName         | Yes      |
				| personal_mdm_profile_description         | personalProfileDescription  | No       |
				| personal_mdm_profile_install_button_name | personalButton              | Yes      |
				| user_enrollment_mdm_profile_installation | userEnrollmentText          | No       |
				| user_enrollment_mdm_profile_name         | userEnrollmentProfileName   | Yes      |
				| user_enrollment_mdm_profile_description  | userEnrollmentProfileDesc...| No       |
				| user_enrollment_mdm_profile_install_but  | userEnrollmentButton        | Yes      |
				| quickadd_package_installation_text       | quickAddText                | No       |
				| quickadd_package_name                    | quickAddName                | Yes      |
				| quickadd_package_progress_text           | quickAddPending             | No       |
				| quickadd_package_install_button_name     | quickAddButton              | Yes      |
				| enrollment_complete_text                 | completeMessage             | No       |
				| enrollment_failed_text                   | failedMessage               | No       |
				| try_again_button_name                    | tryAgainButton              | No       |
				| view_enrollment_status_button_name       | checkNowButton              | No       |
				| view_enrollment_status_text              | checkEnrollmentMessage      | No       |
				| log_out_button_name                      | logoutButton                | Yes      |
				+------------------------------------------+-----------------------------+----------+
			*/
			// Messaging Configuration Schema
			"messaging": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Localized text configuration for enrollment pages and messages",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Two letter ISO 639-1 Language 'set 1' code for this messaging configuration language (e.g., 'en') for the 'English' language. Ref. https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes. Maps to request field 'languageCode'",
						},
						"language_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ISO 639 language name for this message configuration. (e.g., 'English'). Maps to request field 'name'",
						},
						"page_title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Title to display on all enrollment pages. Maps to request field 'title'",
						},
						"login_page_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display below the title on the login page during enrollment. Maps to request field 'loginDescription'",
						},
						"username_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Text to display for the username field on the login page during enrollment. Maps to request field 'username'",
						},
						"password_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Text to display for the password field on the login page during enrollment. Maps to request field 'password'",
						},
						"login_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap/click to log in. Maps to request field 'loginButton'",
						},
						"device_ownership_page_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display during enrollment that prompts the user to specify the device ownership type. Maps to request field 'deviceClassDescription'",
						},
						"personal_device_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to enroll a personally owned device. Maps to request field 'deviceClassPersonal'",
						},
						"institutional_ownership_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to enroll an institutionally owned device. Maps to request field 'deviceClassEnterprise'",
						},
						"personal_device_management_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for personal device management when users enroll a personally owned device. Maps to request field 'deviceClassPersonalDescription'",
						},
						"institutional_device_management_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for institutional device management when users enroll an institutionally owned device. Maps to request field 'deviceClassEnterpriseDescription'",
						},
						"enroll_device_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to start enrollment. Maps to request field 'deviceClassButton'",
						},
						"eula_personal_devices": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "End User License Agreement to display during enrollment of personally owned devices. Maps to request field 'enterpriseEula'",
						},
						"eula_institutional_devices": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "End User License Agreement to display during enrollment of institutionally owned devices and computers. Maps to request field 'personalEula'",
						},
						"accept_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap/click to accept the End User License Agreement. Maps to request field 'eulaButton'",
						},
						"site_selection_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display that prompts the user to select a site if the user has more than one site to choose from during enrollment. Maps to request field 'siteDescription'",
						},
						"ca_certificate_installation_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when installing the CA certificate during enrollment. Maps to request field 'certificateText'",
						},
						"ca_certificate_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to display for the CA certificate during enrollment. Maps to request field 'certificateProfileName'",
						},
						"ca_certificate_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for the CA certificate during enrollment. Maps to request field 'certificateProfileDescription'",
						},
						"ca_certificate_install_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to install the CA certificate. Maps to request field 'certificateButton'",
						},
						"institutional_mdm_profile_installation_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when installing the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseText'",
						},
						"institutional_mdm_profile_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to display for the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseProfileName'",
						},
						"institutional_mdm_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseProfileDescription'",
						},
						"institutional_mdm_profile_pending_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when the user is installing the MDM profile on their computer. Maps to request field 'enterprisePending'",
						},
						"institutional_mdm_profile_install_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to install the MDM profile. Maps to request field 'enterpriseButton'",
						},
						"personal_mdm_profile_installation_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when installing the MDM profile during enrollment of a personally owned device. Maps to request field 'personalText'",
						},
						"personal_mdm_profile_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to display for the MDM profile during enrollment of a personally owned device. Maps to request field 'personalProfileName'",
						},
						"personal_mdm_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for the MDM profile during enrollment of a personally owned device. Maps to request field 'personalProfileDescription'",
						},
						"personal_mdm_profile_install_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to install the MDM profile. Maps to request field 'personalButton'",
						},
						"user_enrollment_mdm_profile_installation_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when prompting to install the MDM profile. Maps to request field 'userEnrollmentText'",
						},
						"user_enrollment_mdm_profile_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to display for the MDM profile. Maps to request field 'userEnrollmentProfileName'",
						},
						"user_enrollment_mdm_profile_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for the MDM profile. Maps to request field 'userEnrollmentProfileDescription'",
						},
						"user_enrollment_mdm_profile_install_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to install the MDM profile. Maps to request field 'userEnrollmentButton'",
						},
						"quickadd_package_installation_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when installing the QuickAdd package during enrollment. Maps to request field 'quickAddText'",
						},
						"quickadd_package_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to display for the QuickAdd package during enrollment. Maps to request field 'quickAddName'",
						},
						"quickadd_package_progress_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when the QuickAdd package is downloading. Maps to request field 'quickAddPending'",
						},
						"quickadd_package_install_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap to install the QuickAdd package. Maps to request field 'quickAddButton'",
						},
						"enrollment_complete_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when enrollment is complete. Maps to request field 'completeMessage'",
						},
						"enrollment_failed_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display when enrollment fails. Maps to request field 'failedMessage'",
						},
						"try_again_button_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name for the button that users tap/click to try enrolling again. Maps to request field 'tryAgainButton'",
						},
						"view_enrollment_status_button_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name for the button that users tap to view the enrollment status for the device. Maps to request field 'checkNowButton'",
						},
						"view_enrollment_status_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to display during enrollment that prompts the user to view the enrollment status for the device. Maps to request field 'checkEnrollmentMessage'",
						},
						"log_out_button_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the button that users tap/click to log out. Maps to request field 'logoutButton'",
						},
					},
				},
			},
			// User Initiated Enrollment for Computers block
			// /api/v4/enrollment
			"user_initiated_enrollment_for_computers": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for user initiated enrollment for computers",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_user_initiated_enrollment_for_computers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Allow users to enroll computers by going to https://JAMF_PRO_URL.jamfcloud.com/enroll (hosted in Jamf Cloud) or https://JAMF_PRO_URL.com:8443/enroll (hosted on-premise). Maps to request field 'macOsEnterpriseEnrollmentEnabled'",
						},
						// Managed Local Administrator Account block
						"managed_local_administrator_account": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration for the managed local administrator account",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"create_managed_local_administrator_account": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Enables managed local administrator account to be used for computers enrolled via PreStage enrollment or user-initiated enrollment. Passwords are set to a random 29 characters that automatically rotate every 1 hour. You can change the rotation settings on the https://your-instance.jamfcloud.com/computerSecurity page. Maps to request field 'createManagementAccount'",
									},
									"management_account_username": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "jamfadmin",
										Description: "Account username to be used for computers enrolled via PreStage enrollment or user-initiated enrollment. This account is created by the Jamf management framework. Maps to request field 'managementUsername'",
									},
									"hide_managed_local_administrator_account": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Hide the managed local administrator account from users. Maps to request field 'hideManagementAccount'",
									},
									"allow_ssh_access_for_managed_local_administrator_account_only": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Make the managed local administrator account the only account that has SSH (Remote Login) access to computers. Maps to request field 'allowSshOnlyManagementAccount'",
									},
								},
							},
						},
						"ensure_ssh_is_enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable SSH (Remote Login) on computers that have it disabled. Maps to request field 'ensureSshRunning'",
						},
						"launch_self_service_when_done": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Ensure that computers launch Self Service immediately after they are enrolled. Maps to request field 'launchSelfService'",
						},
						"quickadd_package": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Third-party signing certificate configuration to ensure that the certificate signs configuration profiles sent to computers and mobile devices, and appears as verified to users during user-initiated enrollment.Maps to 'developerCertificateIdentity'.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sign_quickadd_package": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Ensure that the QuickAdd package is signed and appears as verified to users when enrolling via user-initiated enrollment with a QuickAdd package. Maps to request field 'signQuickAdd'",
									},
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
						"account_driven_device_enrollment": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Account-Driven Device Enrollment enablement for institutionally owned computers. Maps to request field 'accountDrivenDeviceMacosEnrollmentEnabled'",
						},
					},
				},
			},
			// User Initiated Enrollment for Mobile Devices block
			// /api/v4/enrollment
			"user_initiated_enrollment_for_devices": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Allow users to enroll iOS/iPadOS devicesby going to https://JAMF_PRO_URL.jamfcloud.com/enroll (hosted in Jamf Cloud) or https://JAMF_PRO_URL.com:8443/enroll (hosted on-premise)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile_driven_enrollment_via_url": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration for the Profile-Driven Enrollment via URL",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_for_institutionally_owned_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'iosEnterpriseEnrollmentEnabled'",
									},
									"enable_for_personally_owned_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for iOS/iPadOS personally owned devices is enabled. Maps to request field 'iosPersonalEnrollmentEnabled'",
									},
									"personal_device_enrollment_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "USERENROLLMENT",
										Description:  "Enrollment type for user initiated device enrollment for iOS/iPadOS personally owned devices. Note: Personal data on devices enrolled via User Enrollment is retained when the MDM profile is removed. (iOS 13.1 or later, iPadOS 13.1 or later). Maps to request field 'personalDeviceEnrollmentType'",
										ValidateFunc: validation.StringInSlice([]string{"USERENROLLMENT", "DEVICEENROLLMENT"}, false),
									},
								},
							},
						},
						"account_driven_user_enrollment": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration for the Account-Driven User Enrollment",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_for_personally_owned_mobile_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for iOS/iPadOS personally owned devices is enabled. Maps to request field 'accountDrivenUserEnrollmentEnabled'",
									},
									"enable_for_personally_owned_vision_pro_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for visionOS personally owned devices is enabled. Maps to request field 'accountDrivenUserVisionosEnrollmentEnabled'",
									},
								},
							},
						},
						"account_driven_device_enrollment": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration for the Account-Driven Device Enrollment",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_for_institutionally_owned_mobile_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'accountDrivenDeviceIosEnrollmentEnabled'",
									},
									"enable_for_personally_owned_vision_pro_devices": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether user initiated device enrollment for visionOS iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'accountDrivenDeviceVisionosEnrollmentEnabled'",
									},
								},
							},
						},
					},
				},
			},

			//
	  	// Sunsetting Re-enrollment options from this resource
		  // Use jamfpro_reenrollment resource instead
		  //
			// Flush settings
			// "flush_location_information": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Whether to flush location information during re-enrollment",
			// },
			// "flush_location_history_information": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Whether to flush location history information during re-enrollment",
			// },
			// "flush_policy_history": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Whether to flush policy history during re-enrollment",
			// },
			// "flush_extension_attributes": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Whether to flush extension attributes during re-enrollment",
			// },
			// "flush_software_update_plans": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Whether to flush software update plans during re-enrollment",
			// },
			// "flush_mdm_commands_on_reenroll": {
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	Default:      "DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED",
			// 	Description:  "Determines which MDM commands to flush during re-enrollment",
			// 	ValidateFunc: validation.StringInSlice([]string{"DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED", "DELETE_EVERYTHING", "DELETE_NOTHING", "DELETE_ERRORS"}, false),
			// },

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
			// api/v3/enrollment/access-groups
			// Directory Service Groups
			"directory_service_group_enrollment_settings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Directory Service groups to configure enrollment access for",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The unique identifier of the Directory Service group enrollment setting id. Only ID '1' is valid as it represents the built-in 'All Directory Service Users' group.",
							ValidateFunc: validation.StringInSlice([]string{"1"}, false)},
						"directory_service_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "id of the Directory Service group to configure enrollment access for. Maps to request field 'groupId'",
						},
						"ldap_server_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier of the Directory Service group id.",
						},
						"directory_service_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the Directory Service group to configure enrollment access for. Maps to request field 'name'",
						},
						"site_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "ID of the Site to allow this LDAP user group to select during enrollment. Set to '-1' for no site.",
						},
						"allow_group_to_enroll_institutionally_owned_devices": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether directory group can perform user initiated device enrollment for iOS/iPadOS institutionally owned devices. Maps to request field 'enterpriseEnrollmentEnabled'",
						},
						"allow_group_to_enroll_personally_owned_devices": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether directory group can perform user initiated device enrollment for iOS/iPadOS personally owned devices. Maps to request field 'personalEnrollmentEnabled'",
						},
						"allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether directory group can perform user initiated device enrollment for iOS/iPadOS by signing in to their device using a Managed Apple ID with ade. Maps to request field 'accountDrivenUserEnrollmentEnabled'",
						},
						"require_eula": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Upon enrollment is the eula required to be accepted",
						},
					},
				},
			},
		},
	}
}
