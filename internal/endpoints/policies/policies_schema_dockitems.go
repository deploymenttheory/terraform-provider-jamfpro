package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func GetPolicySchemaDockItems() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the dock item.",
			},
			// "name": {
			// 	Type:        schema.TypeString,
			// 	Description: "Name of the dock item.",
			// 	Computed:    true,
			// }, // Not needed
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Action to be performed for the dock item (e.g., Add To Beginning, Add To End, Remove).",
				ValidateFunc: validation.StringInSlice([]string{"Add To Beginning", "Add To End", "Remove"}, false),
			},
		},
	}

	return out
}
