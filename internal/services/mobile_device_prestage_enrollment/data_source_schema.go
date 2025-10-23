// mobiledeviceprestageenrollments_data_source.go
package mobile_device_prestage_enrollment

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMobileDevicePrestageEnrollment provides information about a specific mobile device prestage enrollment in Jamf Pro.
func DataSourceJamfProMobileDevicePrestageEnrollment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the mobile device prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the mobile device prestage.",
			},
		},
	}
}
