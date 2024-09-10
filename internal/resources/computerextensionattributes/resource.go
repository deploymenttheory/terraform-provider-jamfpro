// computerextensionattributes_resource.go
package computerextensionattributes

import (
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
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceComputerExtensionAttributeV0().CoreConfigSchema().ImpliedType(),
				Upgrade: upgradeComputerExtensionAttributeV0toV1,
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
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the computer extension attribute.",
			},
			"data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "String",
				Description:  "Data type of the computer extension attribute. Can be String, Integer, or Date.",
				ValidateFunc: validation.StringInSlice([]string{"STRING", "INTEGER", "DATE"}, false),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enabled by default, but for inputType Script we can disable it as well.Possible values are: false or true.",
			},
			"inventory_display_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Extension Attributes",
				Description:  "Category in which to display the extension attribute in Jamf Pro.",
				ValidateFunc: validation.StringInSlice([]string{"GENERAL", "HARDWARE", "OPERATING_SYSTEM", "USER_AND_LOCATION", "PURCHASING", "EXTENSION_ATTRIBUTES"}, false),
			},
			"input_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Extension attributes collect inventory data by using an input type.The type of the Input used to populate the extension attribute.",
				ValidateFunc: validation.StringInSlice([]string{"SCRIPT", "TEXT", "POPUP", "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"}, false),
			},
			"script_contents": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When we run this script it returns a data value each time a computer submits inventory to Jamf Pro. Provide scriptContents only when inputType is 'SCRIPT'.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return normalizeScript(old) == normalizeScript(new)
				},
				StateFunc: func(v interface{}) string {
					return normalizeScript(v.(string))
				},
			},
			"popup_menu_choices": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "When added with list of choices while creating computer extension attributes these Pop-up menu can be displayed in inventory information. User can choose a value from the pop-up menu list when enrolling a computer any time using Jamf Pro. Provide popupMenuChoices only when inputType is 'POPUP'.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ldap_attribute_mapping": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Directory Service attribute use to populate the extension attribute.Required when inputType is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'.",
			},
			"ldap_extension_attribute_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Collect multiple values for this extension attribute. ldapExtensionAttributeAllowed is disabled by default, only for inputType 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING' it can be enabled. It's value cannot be modified during edit operation.Possible values are:true or false.",
			},
		},
	}
}
