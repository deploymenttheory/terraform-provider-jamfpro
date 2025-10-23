// mobile_device_applications_data_source.go
package mobile_device_application

import (
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMobileDeviceApplications provides information about a specific Jamf Pro Mobile Device Application
func DataSourceJamfProMobileDeviceApplications() *schema.Resource {
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
				Description: "The unique identifier of the mobile device application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the mobile device application.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the mobile device application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle ID of the mobile device application.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the mobile device application.",
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "The scope of the mobile device application.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedMobileDeviceSchemaScope(),
			},
		},
	}
}
