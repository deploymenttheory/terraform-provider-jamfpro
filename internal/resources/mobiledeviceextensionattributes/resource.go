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
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMobileDeviceExtensionAttributeV0().CoreConfigSchema().ImpliedType(),
				Upgrade: upgradeMobileDeviceExtensionAttributeV0toV1,
				Version: 0,
			},
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
				Optional:     true,
				Default:      "String",
				Description:  "Data type of the mobiledevice extension attribute. Can be String, Integer, or Date.",
				ValidateFunc: validation.StringInSlice([]string{"String", "Integer", "Date"}, false),
			},
			"inventory_display": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Extension Attributes",
				Description:  "Category in which to display the extension attribute in Jamf Pro.",
				ValidateFunc: validation.StringInSlice([]string{"General", "Hardware", "User and Location", "Purchasing", "Extension Attributes"}, false),
			},
			"input_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Text Field",
				Description:  "Input type for the Extension Attribute.",
				ValidateFunc: validation.StringInSlice([]string{"Text Field"}, false),
			},
		},
	}
}
