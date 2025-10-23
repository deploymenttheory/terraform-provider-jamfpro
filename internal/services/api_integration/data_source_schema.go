// apiintegrations_data_source.go
package api_integration

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProApiIntegrations provides information about a specific API integration by its ID or Name.
func DataSourceJamfProApiIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the API integration.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the API integration.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client secret",
			},
		},
	}
}
