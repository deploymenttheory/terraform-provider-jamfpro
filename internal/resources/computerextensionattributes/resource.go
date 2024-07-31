// computerextensionattributes_resource.go
package computerextensionattributes

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProComputerExtensionAttributes defines the schema and CRUD operations (Create, Read, Update, Delete)
// for managing Jamf Pro Computer Extension Attributes in Terraform.
func ResourceJamfProComputerExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the computer extension attribute is enabled.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the computer extension attribute.",
			},
			"data_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "string",
				Description: "Data type of the computer extension attribute. Can be string / integer / date (YYYY-MM-DD hh:mm:ss). Value defaults to `String`.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) == strings.ToLower(new)
				},
				ValidateFunc: validation.StringInSlice([]string{"string", "integer", "date"}, false),
			},
			"input_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Extension Attribute Input Type",
				ValidateFunc: validation.StringInSlice([]string{"script", "Text Field", "Pop-up Menu"}, true),
			},
			"input_popup": {
				Type:        schema.TypeList,
				Description: "List of popup choices",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"input_script": {
				Type:        schema.TypeString,
				Description: "Script to populate extension attribute",
				Optional:    true,
			},
			// "input_directory_mapping": {
			// 	Type:        schema.TypeString,
			// 	Description: "Script to populate extension attribute",
			// 	Optional:    true,
			// },
			"inventory_display": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "General",
				Description:  "Display details for inventory for the computer extension attribute. Value defaults to `General`.",
				ValidateFunc: validation.StringInSlice([]string{"General", "Hardware", "Operating System", "User and Location", "Purchasing", "Extension Attributes"}, false),
			},
			"recon_display": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display details for recon for the computer extension attribute.",
			},
		},
	}
}
