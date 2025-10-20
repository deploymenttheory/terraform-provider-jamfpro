package engage_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProEngageSettings defines the schema and CRUD operations for managing Jamf Pro Engage settings configuration in Terraform.
func ResourceEngageSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the Engage settings are enabled or not.",
			},
		},
	}
}
