// volume_purchasing_locations_resource.go
package volume_purchasing_locations

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProVolumePurchasingLocations defines the schema and CRUD operations for managing Jamf Pro Volume Purchasing Locations in Terraform.
func ResourceJamfProVolumePurchasingLocations() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the volume purchasing location.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the volume purchasing location.",
			},
			"apple_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Apple ID associated with the volume purchasing location.",
			},
			"organization_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name for the volume purchasing location.",
			},
			"token_expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the service token.",
			},
			"country_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country code for the volume purchasing location.",
			},
			"location_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location name for the volume purchasing location.",
			},
			"client_context_mismatch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if there is a client context mismatch.",
			},
			"automatically_populate_purchased_content": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Automatically populate purchased content.",
			},
			"send_notification_when_no_longer_assigned": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Send notification when content is no longer assigned.",
			},
			"auto_register_managed_users": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Automatically register managed users.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
				Description: "The site ID associated with the volume purchasing location. Default is '-1'.",
			},
			"service_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The base64 encoded service token for the volume purchasing location.",
			},
			"last_sync_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last synchronization time.",
			},
			"total_purchased_licenses": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of purchased licenses.",
			},
			"total_used_licenses": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of used licenses.",
			},
			"content": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the content.",
						},
						"license_count_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of licenses for the content.",
						},
						"license_count_in_use": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of licenses currently in use.",
						},
						"license_count_reported": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of licenses reported.",
						},
						"icon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the content icon.",
						},
						"device_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of device types compatible with the content.",
						},
						"content_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the content.",
						},
						"pricing_param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pricing parameter for the content.",
						},
						"adam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ADAM ID of the content.",
						},
					},
				},
				Description: "List of content associated with the volume purchasing location.",
			},
		},
	}
}
