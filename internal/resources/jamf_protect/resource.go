package jamf_protect

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProtect defines the schema and CRUD operations for managing Jamf Protect integration in Terraform.
func ResourceJamfProtect() *schema.Resource {
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
		Schema: map[string]*schema.Schema{
			"protect_url": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The URL of the Jamf Protect instance",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The API client ID for Jamf Protect authentication",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The password for Jamf Protect authentication",
			},
			"auto_install": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to automatically install Jamf Protect on devices",
			},
			"sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current sync status of the Jamf Protect integration",
			},
			"api_client_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the API client used for integration",
			},
			"last_sync_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the last successful sync",
			},
			"registration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Registration ID of the Jamf Protect integration",
			},
		},
	}
}
