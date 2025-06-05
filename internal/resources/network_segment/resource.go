// networksegments_resource.go
package network_segment

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProNetworkSegments defines the schema and CRUD operations for managing Jamf Pro NetworkSegments in Terraform.
func ResourceJamfProNetworkSegments() *schema.Resource {
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
				Description: "The unique identifier of the network segment.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network segment.",
			},
			"starting_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The starting IP address of the network segment.",
			},
			"ending_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ending IP address of the network segment.",
			},
			"distribution_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution server associated with the network segment.",
			},
			"distribution_point": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution point associated with the network segment.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL associated with the network segment.",
			},
			"swu_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The software update server associated with the network segment.",
			},
			"building": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The building associated with the network segment.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The department associated with the network segment.",
			},
			"override_buildings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if building assignments are overridden for this network segment.",
			},
			"override_departments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if department assignments are overridden for this network segment.",
			},
		},
	}
}
