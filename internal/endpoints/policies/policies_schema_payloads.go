package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicySchemaPayloads() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"override_default_settings": { // UI > payloads > software update settings
				Type:        schema.TypeList,
				Required:    true,
				Description: "Settings to override default configurations.",
				Elem:        getPolicySchemaNetworkLimitations(),
			},
			"network_requirements": { // NOT IN THE UI
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network requirements for the policy.",
				Default:     "Any",
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
				Elem:        getPolicySchemaPrinter(),
			},
			"dock_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dock items settings of the policy.",
				Elem:        getPolicySchemaDockItems(),
			},
			"account_maintenance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Account maintenance settings of the policy. Use this section to create and delete local accounts, and to reset local account passwords. Also use this section to disable an existing local account for FileVault 2.",
				Elem:        getPolicySchemaAccountMaintenance(),
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
				Elem:        getPolicySchemaFilesProcesses(),
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
				Elem:        getSharedSchemaDiskEncryption(),
			},
		},
	}

	return out
}
