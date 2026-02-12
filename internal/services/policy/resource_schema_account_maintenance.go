package policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getPolicySchemaAccountMaintenance() *schema.Resource {
	out := &schema.Resource{Schema: map[string]*schema.Schema{
		"local_accounts": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Local user account configurations",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Details of each account configuration.",
						Elem:        getPolicySchemaAccount(),
					},
				},
			},
		},
		"directory_bindings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Directory binding settings for the policy. Use this section to bind computers to a directory service",
			Elem:        getPolicySchemaDirectoryBinding(),
		},
		"management_account": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Management account settings for the policy. Use this section to change or reset the management account password.",
			Elem:        getPolicySchemaManagementAccount(),
		},
		"open_firmware_efi_password": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Open Firmware/EFI password settings for the policy. Use this section to set or remove an Open Firmware/EFI password on computers with Intel-based processors.",
			Elem:        getPolicySchemaEfiFirmwarePassword(),
		},
	}}

	return out
}

func getPolicySchemaAccount() *schema.Resource {
	out := &schema.Resource{
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
				//Sensitive:   true,
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
	}

	return out
}

func getPolicySchemaDirectoryBinding() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"binding": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Details of the directory binding.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the binding.",
						},
					},
				},
			},
		},
	}

	return out
}

func getPolicySchemaManagementAccount() *schema.Resource {
	out := &schema.Resource{
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
	}

	return out
}

func getPolicySchemaEfiFirmwarePassword() *schema.Resource {
	out := &schema.Resource{
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
	}

	return out
}
