// computerinventorycollection_resource.go
package computerinventorycollection

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProComputerInventoryCollection defines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerInventoryCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProComputerInventoryCollectionCreate,
		ReadContext:   resourceJamfProComputerInventoryCollectionRead,
		UpdateContext: resourceJamfProComputerInventoryCollectionUpdate,
		DeleteContext: resourceJamfProComputerInventoryCollectionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"local_user_accounts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on local user accounts.",
			},
			"home_directory_sizes": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on home directory sizes.",
			},
			"hidden_accounts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on hidden accounts.",
			},
			"printers": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on printers.",
			},
			"active_services": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on active services.",
			},
			"mobile_device_app_purchasing_info": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on mobile device app purchasing.",
			},
			"computer_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect computer location information.",
			},
			"package_receipts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect package receipts.",
			},
			"available_software_updates": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on available software updates.",
			},
			"include_applications": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include applications in the inventory collection.",
			},
			"include_fonts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include fonts in the inventory collection.",
			},
			"include_plugins": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include plugins in the inventory collection.",
			},
			"applications": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of applications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the application.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the application.",
						},
					},
				},
			},
			"fonts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of fonts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the font.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the font.",
						},
					},
				},
			},
			"plugins": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of plugins.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the plugin.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the plugin.",
						},
					},
				},
			},
		},
	}
}
