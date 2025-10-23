package sso_failover

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSSOFailover provides information about Jamf Pro's SSO Failover configuration
func DataSourceJamfProSSOFailover() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"failover_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SSO failover URL for Jamf Pro",
			},
			"generation_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp when the failover URL was generated",
			},
		},
	}
}
