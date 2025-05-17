package managedsoftwareupdatesfeaturetoggle

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceManagedSoftwareUpdateFeatureToggle defines the schema and CRUD operations for managing Jamf Pro Managed Software Updates Feature Toggle configuration in Terraform.
func ResourceManagedSoftwareUpdateFeatureToggle() *schema.Resource {
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
			"toggle": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the Managed Software Updates Feature Toggle is enabled or not.",
			},
		},
	}
}
