package mac_application

import (
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacApplications provides information about a specific Jamf Pro Mac Application
func DataSourceJamfProMacApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the mac application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the mac application.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the mac application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle identifier of the mac application.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the mac application.",
			},
			"is_free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the application is free.",
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment type of the mac application.",
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "The scope of the mac application.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
		},
	}
}
