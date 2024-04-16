// policies_resource.go
package policies

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProPolicies defines the schema and CRUD operations for managing Jamf Pro Policy in Terraform.
func ResourceJamfProPolicies() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPoliciesCreate,
		ReadContext:   ResourceJamfProPoliciesRead,
		UpdateContext: ResourceJamfProPoliciesUpdate,
		DeleteContext: ResourceJamfProPoliciesDelete,
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
				Description: "The unique identifier of the Jamf Pro policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Define whether the policy is enabled.",
			},
			// "trigger": { // NOTE appears to be redundant when used with the below. Maybe this use to be a multiple choice option?
			// 	Type:         schema.TypeString,
			// 	Required:     true,
			// 	Description:  "Event(s) triggers to use to initiate the policy. Values can be 'USER_INITIATED' for self self trigger and 'EVENT' for an event based trigger",
			// 	ValidateFunc: validation.StringInSlice([]string{"EVENT", "USER_INITIATED"}, false),
			// },
			"trigger_checkin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro",
			},
			"trigger_enrollment_complete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when device enrollment is complete.",
			},
			"trigger_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work",
			},
			// "trigger_logout": { // NOTE appears to be redundant
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Trigger policy when a user logout.",
			// 	Default:     false,
			// },
			"trigger_network_state_changed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)",
			},
			"trigger_startup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work",
			},
			"trigger_other": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Any other trigger for the policy.",
				// TODO need a validation func here to make sure this cannot be provided as empty.
			},
			"frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Frequency of policy execution.",
				Default:     "Once per computer",
				ValidateFunc: validation.StringInSlice([]string{
					"Once per computer",
					"Once per user per computer",
					"Once per user",
					"Once every day",
					"Once every week",
					"Once every month",
					"Ongoing",
				}, false),
			},
			"retry_event": { // Retry only relevant if frequency is Once Per Computer
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Event on which to retry policy execution.",
				Default:     "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"trigger",
					"check-in",
				}, false),
			},
			"retry_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of retry attempts for the jamf pro policy. Valid values are -1 (not configured) and 1 through 10.",
				Default:     -1,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetInt(val)
					if v == -1 || (v > 0 && v <= 10) {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be -1 if not being set or between 1 and 10 if it is being set, got: %d", key, v))
					return warns, errs
				},
			},
			"notify_on_each_failed_retry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Send notifications for each failed policy retry attempt. ",
				Default:     false,
			},
			// "location_user_only": { // NOTE Can't find in GUI
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Location-based policy for user only.",
			// 	Default:     false,
			// },
			"target_drive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default",
				Default:     "/",
			},
			"offline": { // Only avaible if frequency set to continuous else not needed
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Make policy available offline by caching the policy to the macOS device to ensure it runs when Jamf Pro is unavailable. Only used when execution policy is set to 'ongoing'. ",
				Default:     false,
			},
			"category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Category to add the policy to.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The category ID assigned to the jamf pro policy. Defaults to '-1' aka not used.",
							Default:     "-1",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Category Name for assigned jamf pro policy. Value defaults to 'No category assigned' aka not used",
							Default:     "No category assigned",
						},
					},
				},
			},
			"date_time_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"activation_date": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The activation date of the policy.",
							Computed:    true,
						},
						"activation_date_epoch": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The epoch time of the activation date.",
							Computed:    true,
						},
						"activation_date_utc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The UTC time of the activation date.",
							Computed:    true,
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The expiration date of the policy.",
							Computed:    true,
						},
						"expiration_date_epoch": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The epoch time of the expiration date.",
							Computed:    true,
						},
						"expiration_date_utc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The UTC time of the expiration date.",
						},
						"no_execute_on": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, false),
							},
							Description: "Client-side limitations are enforced based on the settings on computers. This field sets specific days when the policy should not execute.",
							Computed:    true,
						},
						"no_execute_start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client-side limitations are enforced based on the settings on computers. This field sets the start time when the policy should not execute.",
						},
						"no_execute_end": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client-side limitations are enforced based on the settings on computers. This field sets the end time when the policy should not execute.",
						},
					},
				},
			},
			"network_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network limitations for the policy.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum_network_connection": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Minimum network connection required for the policy.",
							Default:      "No Minimum",
							ValidateFunc: validation.StringInSlice([]string{"No Minimum", "Ethernet"}, false),
						},
						"any_ip_address": { // NOT IN THE UI
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the policy applies to any IP address.",
							Default:     true,
						},
						"network_segments": { // surely this has been moved to scope now?
							Type:        schema.TypeString,
							Description: "Network segment limitations for the policy.",
							Optional:    true,
						},
					},
				},
			}, // END OF General UI
			"override_default_settings": { // UI > payloads > software update settings
				Type:        schema.TypeList,
				Required:    true,
				Description: "Settings to override default configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_drive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The drive on which to run the policy (e.g. '/Volumes/Restore/'). Defaults to '/' if no value is defined, which is the root of the file system.",
							Default:     "/",
						},
						"distribution_point": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Distribution point for the policy.",
							Default:     "default",
						},
						"force_afp_smb": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to force AFP/SMB.",
							Default:     false,
						},
						"sus": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Software Update Service for the policy.",
							Default:     "default",
						},
					},
				},
			},
			"network_requirements": { // NOT IN THE UI
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network requirements for the policy.",
				Default:     "Any",
			},
			"site": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Jamf Pro Site-related settings of the policy.",
				Elem:        sharedschemas.GetSharedSchemaSite(),
			},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Scope configuration for the profile.",
				Elem:        sharedschemas.GetSharedSchemaScope(),
			},
			"self_service": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Self-service settings of the policy.",
				Elem:        getPolicySchemaSelfService(),
			},
			"packages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package configuration settings of the policy.",
				Elem:        getPolicySchemaPackages(),
			},
			"scripts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Scripts settings of the policy.",
				Elem:        getPolicySchemaScript(),
			},
			"printers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Printers settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"leave_existing_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Policy for handling existing default printers.",
						},
						"printer": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the printer configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Unique identifier of the printer.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the printer.",
										Computed:    true,
									},
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to be performed for the printer (e.g., install, uninstall).",
										ValidateFunc: validation.StringInSlice([]string{"install", "uninstall"}, false),
									},
									"make_default": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to set the printer as the default.",
									},
								},
							},
						},
					},
				},
			},
			"dock_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dock items settings of the policy.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dock_item": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the dock item configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Unique identifier of the dock item.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the dock item.",
										Computed:    true,
									},
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to be performed for the dock item (e.g., Add To Beginning, Add To End, Remove).",
										ValidateFunc: validation.StringInSlice([]string{"Add To Beginning", "Add To End", "Remove"}, false),
									},
								},
							},
						},
					},
				},
			},
			"account_maintenance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Account maintenance settings of the policy. Use this section to create and delete local accounts, and to reset local account passwords. Also use this section to disable an existing local account for FileVault 2.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of account maintenance configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of each account configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Action to be performed on the account (e.g., Create, Reset, Delete, DisableFileVault).",
													ValidateFunc: validation.StringInSlice([]string{"Create", "Reset", "Delete", "DisableFileVault"}, false),
												},
												"username": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Username/short name for the account",
												},
												"realname": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Real name associated with the account.",
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Set a new account password. This does not update the account's login keychain password or FileVault 2 password.",
												},
												"archive_home_directory": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Permanently delete home directory. If set to true will archive the home directory.",
												},
												"archive_home_directory_to": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Path in which to archive the home directory to.",
												},
												"home": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Full path in which to create the home directory (e.g. /Users/username/ or /private/var/username/)",
												},
												"hint": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Hint to help the user remember the password",
												},
												"picture": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Full path to the account picture (e.g. /Library/User Pictures/Animals/Butterfly.tif )",
												},
												"admin": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the account has admin privileges.Setting this to true will set the user administrator privileges to the computer",
												},
												"filevault_enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Allow the user to unlock the FileVault 2-encrypted drive",
												},
											},
										},
									},
								},
							},
						},
						"directory_bindings": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Directory binding settings for the policy. Use this section to bind computers to a directory service",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binding": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the directory binding.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the binding.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "The name of the binding.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"management_account": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Management account settings for the policy. Use this section to change or reset the management account password.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to perform on the management account.Rotates management account password at next policy execution. Valid values are 'rotate' or 'doNotChange'.",
										ValidateFunc: validation.StringInSlice([]string{"rotate", "doNotChange"}, false),
										Default:      "doNotChange",
									},
									"managed_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Managed password for the account. Management account passwords will be automatically randomized with 29 characters by jamf pro.",
										//Default:     "",
										//Computed: true,
									},
									"managed_password_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length of the managed password. Only necessary when utilizing the random action",
										Default:     0,
									},
								},
							},
						},
						"open_firmware_efi_password": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Open Firmware/EFI password settings for the policy. Use this section to set or remove an Open Firmware/EFI password on computers with Intel-based processors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"of_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Mode for the open firmware/EFI password. Valid values are 'command' or 'none'.",
										ValidateFunc: validation.StringInSlice([]string{"command", "none"}, false),
										Default:      "none",
									},
									"of_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Password for the open firmware/EFI.",
										Default:     "",
									},
								},
							},
						},
					},
				},
			},
			"reboot": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Use this section to restart computers and specify the disk to boot them to",
				Elem:        getPolicySchemaReboot(),
			},
			"maintenance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Maintenance settings of the policy. Use this section to update inventory, reset computer names, install all cached packages, and run common maintenance tasks.",
				Elem:        getPolicySchemaMaintenance(),
			},
			"files_processes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Files and processes settings of the policy. Use this section to search for and log specific files and processes. Also use this section to execute a command.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_by_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path of the file to search for.",
						},
						"delete_file": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to delete the file found at the specified path.",
							Default:     false,
						},
						"locate_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path of the file to locate. Name of the file, including the file extension. This field is case-sensitive and returns partial matches",
						},
						"update_locate_database": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to update the locate database. Update the locate database before searching for the file",
							Default:     false,
						},
						"spotlight_search": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Search For File Using Spotlight. File to search for. This field is not case-sensitive and returns partial matches",
						},
						"search_for_process": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the process to search for. This field is case-sensitive and returns partial matches",
						},
						"kill_process": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to kill the process if found. This works with exact matches only",
							Default:     false,
						},
						"run_command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Command to execute on computers. This command is executed as the 'root' user",
						},
					},
				},
			},
			"user_interaction": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User interaction settings of the policy.",
				Elem:        getPolicySchemaUserInteraction(),
			},
			"disk_encryption": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Disk encryption settings of the policy. Use this section to enable FileVault 2 or to issue a new recovery key.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The action to perform for disk encryption (e.g., apply, remediate).",
							ValidateFunc: validation.StringInSlice([]string{"none", "apply", "remediate"}, false),
							Default:      "none",
						},
						"disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the disk encryption configuration to apply.",
							Default:     0,
						},
						"auth_restart": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to allow authentication restart.",
							Default:     false,
						},
						"remediate_key_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Type of key to use for remediation (e.g., Individual, Institutional, Individual And Institutional).",
							ValidateFunc: validation.StringInSlice([]string{"Individual", "Institutional", "Individual And Institutional"}, false),
						},
						"remediate_disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disk encryption ID to utilize for remediating institutional recovery key types.",
							Default:     0,
						},
					},
				},
			},
		},
	}
}
