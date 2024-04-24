package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getSharedSchemaDiskEncryption() *schema.Resource {
	out := &schema.Resource{
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
	}

	return out
}
