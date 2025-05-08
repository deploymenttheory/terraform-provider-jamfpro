package ssofailover

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProSSOFailover defines the schema and CRUD operations for managing Jamf Pro SSO Failover configuration in Terraform.
func ResourceJamfProSSOFailover() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"regenerate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to regenerate the failover URL",
			},
		},
	}
}
