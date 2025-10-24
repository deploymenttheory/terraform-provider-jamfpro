// mobiledeviceextensionattributes_data_source.go
package mobile_device_extension_attribute

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMobileDeviceDeviceExtensionAttributes provides information about a specific mobiledevice extension attribute by its ID or Name.
func DataSourceJamfProMobileDeviceDeviceExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the mobiledevice extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro mobiledevice extension attribute.",
			},
		},
	}
}
