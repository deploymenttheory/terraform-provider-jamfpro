package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func mgetPolicySchemaNetworkLimitations() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"minimum_network_connection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Minimum network connection required for the policy.",
				Default:     "No Minimum",
				// ValidateFunc: validation.StringInSlice([]string{"No Minimum", "Ethernet"}, false),
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
				Default:     "",
			},
		},
	}

	return out
}
