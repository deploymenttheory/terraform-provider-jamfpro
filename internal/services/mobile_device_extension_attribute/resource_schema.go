// mobiledeviceextensionattributes_resource.go
package mobile_device_extension_attribute

import (
	"context"
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
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceJamfProMobileDeviceExtensionAttributesV0().CoreConfigSchema().ImpliedType(),
				Upgrade: upgradeMobileDeviceExtensionAttributesV0toV1,
				Version: 0,
			},
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
				Description:  "Data type of the mobile device extension attribute. Can be STRING, INTEGER, or DATE.",
				ValidateFunc: validation.StringInSlice([]string{"STRING", "INTEGER", "DATE"}, false),
			},
			"inventory_display_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Category in which to display the extension attribute in Jamf Pro. Can be GENERAL, HARDWARE, USER_AND_LOCATION, PURCHASING, or EXTENSION_ATTRIBUTES.",
				ValidateFunc: validation.StringInSlice([]string{"GENERAL", "HARDWARE", "USER_AND_LOCATION", "PURCHASING", "EXTENSION_ATTRIBUTES"}, false),
			},
			"input_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Extension attributes collect inventory data by using an input type.The type of the Input used to populate the extension attribute. Can be TEXT, POPUP, or DIRECTORY_SERVICE_ATTRIBUTE_MAPPING.",
				ValidateFunc: validation.StringInSlice([]string{"TEXT", "POPUP", "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"}, false),
			},
			"popup_menu_choices": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "When added with list of choices while creating mobile device extension attributes these Pop-up menu can be displayed in inventory information. User can choose a value from the pop-up menu list when enrolling a mobile device any time using Jamf Pro. Provide popupMenuChoices only when inputType is 'POPUP'.",
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

// Old schema for state upgrade
func resourceJamfProMobileDeviceExtensionAttributesV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id":                {Type: schema.TypeString, Computed: true},
			"name":              {Type: schema.TypeString, Required: true},
			"description":       {Type: schema.TypeString, Optional: true},
			"data_type":         {Type: schema.TypeString, Required: true},
			"inventory_display": {Type: schema.TypeString, Required: true},
			"input_type": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {Type: schema.TypeString, Required: true},
						"popup_choices": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// State upgrader function
func upgradeMobileDeviceExtensionAttributesV0toV1(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if v, ok := rawState["input_type"].([]any); ok && len(v) > 0 {
		block, ok := v[0].(map[string]any)
		if ok {
			if t, ok := block["type"].(string); ok {
				rawState["input_type"] = t
			}
			if choices, ok := block["popup_choices"].([]any); ok {
				rawState["popup_menu_choices"] = choices
			}
		}
	}
	if v, ok := rawState["inventory_display"]; ok {
		rawState["inventory_display_type"] = v
		newState := make(map[string]any)
		for k, val := range rawState {
			if k != "inventory_display" {
				newState[k] = val
			}
		}
		rawState = newState
	}
	return rawState, nil
}
