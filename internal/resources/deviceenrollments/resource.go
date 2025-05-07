// deviceenrollments_resource.go
package deviceenrollments

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProDeviceEnrollments defines the schema and CRUD operations for managing Jamf Pro Device Enrollments in Terraform.
func ResourceJamfProDeviceEnrollments() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the device enrollment.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the device enrollment.",
			},
			"supervision_identity_id": {
				Type:        schema.TypeString,
				Default:     "-1",
				Optional:    true,
				Description: "The supervision identity ID associated with the device enrollment. Default is '-1'.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Default:     "-1",
				Optional:    true,
				Description: "The site ID associated with the device enrollment. Default is '-1'.",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server name for the device enrollment.",
			},
			"server_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server UUID for the device enrollment.",
			},
			"admin_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The admin ID associated with the device enrollment.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name for the device enrollment.",
			},
			"org_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization email for the device enrollment.",
			},
			"org_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization phone number for the device enrollment.",
			},
			"org_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization address for the device enrollment.",
			},
			"token_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the enrollment token.",
			},
			"token_file_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional name of the token to be saved, if no name is provided one will be auto-generated.",
			},
			"encoded_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The base64 encoded MDM server token.",
			},
		},
	}
}
