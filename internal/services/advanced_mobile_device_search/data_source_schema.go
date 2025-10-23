// advancedmobiledevicesearches_data_source.go
package advanced_mobile_device_search

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAdvancedMobileDeviceSearches provides information about a specific Advanced Mobile Device Search by its ID or Name.
func DataSourceJamfProAdvancedMobileDeviceSearches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the advanced mobile device search.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the advanced mobile device search.",
			},
		},
	}
}
