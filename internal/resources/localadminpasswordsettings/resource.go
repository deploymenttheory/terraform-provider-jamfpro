package localadminpasswordsettings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProActivationCode defines the schema and CRUD operations for managing Jamf Pro activation code configuration in Terraform.
func ResourceLocalAdminPasswordSettings() *schema.Resource {
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
			"auto_deploy_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When enabled, all appropriate computers will have the SetAutoAdminPassword command sent to them automatically.",
			},
			"password_rotation_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "The amount of time in seconds that the local admin password will be rotated after viewing.",
			},
			"auto_rotate_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When enabled, all appropriate computers will automatically have their password expired and rotated after the configured autoRotateExpirationTime.",
			},
			"auto_rotate_expiration_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7776000,
				Description: "The amount of time in seconds that the local admin password will be rotated automatically if it is never viewed.",
			},
		},
	}
}
