// mobiledeviceextensionattributes_resource.go
package mobiledeviceextensionattributes

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProMobileDeviceExtensionAttributes defines the schema and CRUD operations (Create, Read, Update, Delete)
// for managing Jamf Pro MobileDevice Extension Attributes in Terraform.
func ResourceJamfProMobileDeviceExtensionAttributes() *schema.Resource {
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
				Description: "The unique identifier of the mobile device extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro mobiledevice extension attribute.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the mobiledevice extension attribute.",
			},
			"data_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Data type of the mobile device extension attribute. Can be String, Integer, or Date.",
				ValidateFunc: validation.StringInSlice([]string{"String", "Integer", "Date"}, false),
			},
			"inventory_display": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Category in which to display the extension attribute in Jamf Pro.",
				ValidateFunc: validation.StringInSlice([]string{"GENERAL", "HARDWARE", "USER_AND_LOCATION", "PURCHASING", "EXTENSION_ATTRIBUTES", "General", "Hardware", "User and Location", "Purchasing", "Extension Attributes"}, false),
			},
			"input_type": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Input type for the Extension Attribute.",
							ValidateFunc: validation.StringInSlice([]string{"Text Field", "Pop-up Menu"}, false),
						},
						"popup_choices": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description:      "List of choices for Pop-up Menu input type.",
							DiffSuppressFunc: suppressPopupChoicesDiff,
						},
					},
				},
			},
		},
	}
}
