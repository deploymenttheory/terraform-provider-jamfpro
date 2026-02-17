package smart_mobile_device_group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSmartMobileGroups provides information about a specific mobile group in Jamf Pro.
func DataSourceJamfProSmartMobileGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the Jamf Pro Smart mobile group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro Smart mobile group.",
			},
		},
	}
}
