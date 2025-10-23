package device_enrollments_public_key

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDeviceEnrollmentsPublicKey provides the public key for device enrollments in Jamf Pro.
func DataSourceJamfProDeviceEnrollmentsPublicKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key used for device enrollments.",
			},
		},
	}
}
