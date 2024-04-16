package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getPolicySchemaScript() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of the script.",
			},
			// "name": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "Name of the script.",
			// }, // Do we need this?
			"priority": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Execution priority of the script.",
				ValidateFunc: validation.StringInSlice([]string{"Before", "After"}, false),
				Default:      "After",
			},
			"parameter4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 4 for the script.",
			},
			"parameter5": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 5 for the script.",
			},
			"parameter6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 6 for the script.",
			},
			"parameter7": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 7 for the script.",
			},
			"parameter8": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 8 for the script.",
			},
			"parameter9": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 9 for the script.",
			},
			"parameter10": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 10 for the script.",
			},
			"parameter11": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameter 11 for the script.",
			},
		},
	}
	return out
}
