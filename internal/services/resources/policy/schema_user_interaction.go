package policy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicySchemaUserInteraction() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Message to display before the policy runs",
			},
			"allow_users_to_defer": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow user deferral and configure deferral type. A deferral limit must be specified for this to work.",
				Default:     false,
			},
			"allow_deferral_until_utc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Date/time at which deferrals are prohibited and the policy runs. Uses time zone settings of your hosting server. Standard environments hosted in Jamf Cloud use Coordinated Universal Time (UTC)",
			},
			"allow_deferral_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of minutes after the user was first prompted by the policy at which the policy runs and deferrals are prohibited. Must be a multiple of 1440 (minutes in day)",
				Default:     0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v%1440 != 0 {
						errs = append(errs, fmt.Errorf("%q must be a multiple of 1440 (minutes in day), got: %d", key, v))
					}
					return
				},
			},
			"message_finish": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Message to display when the policy is complete.",
			},
		},
	}

	return out
}
