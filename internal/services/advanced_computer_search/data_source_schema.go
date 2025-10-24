// advancedcomputersearches_data_source.go
package advanced_computer_search

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAdvancedComputerSearches provides information about a specific Advanced Computer Search by its ID or Name.
func DataSourceJamfProAdvancedComputerSearches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the advanced computer search.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the advanced computer search.",
			},
		},
	}
}
