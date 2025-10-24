// deviceenrollments_data_source.go
package device_enrollments

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDeviceEnrollments provides information about device enrollments in Jamf Pro.
func DataSourceJamfProDeviceEnrollments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the device enrollment.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the device enrollment.",
			},
			"supervision_identity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The supervision identity ID for the device enrollment.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The site ID associated with the device enrollment.",
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
				Description: "The administrator ID associated with the device enrollment.",
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
				Description: "The expiration date of the device enrollment token.",
			},
		},
	}
}
