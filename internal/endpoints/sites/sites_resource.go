// sites_resource.go
package sites

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProSite defines the schema and CRUD operations for managing Jamf Pro Sites in Terraform.
func ResourceJamfProSites() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProSitesCreate,
		ReadContext:   resourceJamfProSitesRead,
		UpdateContext: resourceJamfProSitesUpdate,
		DeleteContext: resourceJamfProSitesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the site.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro site.",
			},
		},
	}
}
