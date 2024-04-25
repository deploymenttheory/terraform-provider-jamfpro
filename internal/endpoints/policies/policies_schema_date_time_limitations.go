package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicySchemaDateTimeLimitations() *schema.Resource {
	out := &schema.Resource{
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
			// "no_execute_on": {
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type:         schema.TypeString,
			// 		ValidateFunc: validation.StringInSlice([]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, false),
			// 	},
			// 	Description: "Client-side limitations are enforced based on the settings on computers. This field sets specific days when the policy should not execute.",
			// 	Computed:    true,
			// },
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
		}}

	return out
}
