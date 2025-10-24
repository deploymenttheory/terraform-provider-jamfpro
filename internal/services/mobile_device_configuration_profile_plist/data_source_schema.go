// mobiledeviceconfigurationprofilesplist_data_source.go
package mobile_device_configuration_profile_plist

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMobileDeviceConfigurationProfiles provides information about a specific mobile device configuration profile in Jamf Pro.
func DataSourceJamfProMobileDeviceConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the mobile device configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the mobile device configuration profile.",
			},
		},
	}
}
