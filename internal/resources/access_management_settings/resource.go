package access_management_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceAccessManagementSettings defines the schema and CRUD operations for managing Jamf Pro Access Management settings configuration in Terraform.
func ResourceAccessManagementSettings() *schema.Resource {
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
			"automated_device_enrollment_server_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "MDM Server UUID.",
			},
		},
	}
}
