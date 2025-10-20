package sso_failover

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProSSOFailover defines the schema and CRUD operations for managing Jamf Pro SSO Failover configuration in Terraform.
func ResourceJamfProSSOFailover() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"failover_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
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
