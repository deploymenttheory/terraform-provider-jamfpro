package computerextensionattributes

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputerExtensionAttributeV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id":                {Type: schema.TypeString, Computed: true},
			"name":              {Type: schema.TypeString, Required: true},
			"description":       {Type: schema.TypeString, Optional: true},
			"data_type":         {Type: schema.TypeString, Optional: true},
			"enabled":           {Type: schema.TypeBool, Required: true},
			"input_type":        {Type: schema.TypeString, Required: true},
			"input_popup":       {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"input_script":      {Type: schema.TypeString, Optional: true},
			"inventory_display": {Type: schema.TypeString, Optional: true},
			"recon_display":     {Type: schema.TypeString, Optional: true},
		},
	}
}

func upgradeComputerExtensionAttributeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	newState := make(map[string]interface{})

	// Map existing fields
	newState["id"] = rawState["id"]
	newState["name"] = rawState["name"]
	newState["description"] = rawState["description"]
	newState["enabled"] = rawState["enabled"]

	// Update data_type to use proper capitalization
	if dataType, ok := rawState["data_type"]; ok {
		newState["data_type"] = strings.Title(dataType.(string))
	} else {
		newState["data_type"] = "String" // Default value
	}

	// Map inventory_display to inventory_display_type
	if inv, ok := rawState["inventory_display"]; ok {
		newState["inventory_display_type"] = inv
	} else {
		newState["inventory_display_type"] = "Extension Attributes" // Default value
	}

	// Handle input_type and related fields
	if inputType, ok := rawState["input_type"]; ok {
		switch inputType {
		case "script":
			newState["input_type"] = "Script"
			if script, ok := rawState["input_script"]; ok {
				newState["script_contents"] = normalizeScript(script.(string))
			}
		case "Pop-up Menu":
			newState["input_type"] = "Pop-up Menu"
			if popup, ok := rawState["input_popup"]; ok {
				newState["popup_menu_choices"] = popup
			}
		case "Text Field":
			newState["input_type"] = "Text Field"
		default:
			newState["input_type"] = inputType
		}
	}

	// Initialize new fields with default values
	newState["ldap_attribute_mapping"] = ""
	newState["ldap_extension_attribute_allowed"] = false

	return newState, nil
}
