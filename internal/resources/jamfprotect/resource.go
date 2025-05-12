// jamf_protect_resource.go
package jamfprotect

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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Jamf Protect integration.",
			},
			"api_client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The API client ID for the Jamf Protect integration.",
			},
			"api_client_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the API client.",
			},
			"registration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The registration ID for the Jamf Protect integration.",
			},
			"protect_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The URL of the Jamf Protect instance.",
			},
			"last_sync_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the last synchronization.",
			},
			"sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of synchronization.",
			},
			"auto_install": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to automatically install Jamf Protect on devices.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The client ID used for registration.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The password used for registration.",
			},
		},
	}
}
