package cloud_idp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProCloudIdp provides information about a specific cloud identity provider in Jamf Pro.
func DataSourceJamfProCloudIdp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The jamf pro unique identifier of the cloud identity provider.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The display name of the cloud identity provider.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the cloud identity provider is enabled.",
			},
			"provider_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the cloud identity provider. e.g AZURE",
			},
		},
	}
}
