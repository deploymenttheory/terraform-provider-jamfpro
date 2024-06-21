// allowedfileextensions_resource.go
package allowedfileextensions

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAllowedFileExtensions defines the schema and CRUD operations for managing AllowedFileExtentionss in Terraform.
func ResourceJamfProAllowedFileExtensions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProAllowedFileExtensionCreate,
		ReadContext:   resourceJamfProAllowedFileExtensionRead,
		UpdateContext: resourceJamfProAllowedFileExtensionUpdate,
		DeleteContext: resourceJamfProAllowedFileExtensionDelete,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
