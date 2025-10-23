// advancedusersearches_data_source.go
package advanced_user_search

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAdvancedUserSearches provides information about a specific Advanced User Search by its ID or Name.
func DataSourceJamfProAdvancedUserSearches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the advanced user search.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the advanced user search.",
			},
		},
	}
}
