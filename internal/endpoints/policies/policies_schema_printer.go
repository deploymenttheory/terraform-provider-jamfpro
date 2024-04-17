package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func GetPolicySchemaPrinter() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the printer.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Action to be performed for the printer (e.g., install, uninstall).",
				ValidateFunc: validation.StringInSlice([]string{"install", "uninstall"}, false),
			},
			"make_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to set the printer as the default.",
			},
		},
	}

	return out
}
