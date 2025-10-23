package computer_extension_attribute

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func upgradeComputerExtensionAttributeV0toV1(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	newState := make(map[string]any)

	// Map existing fields
	newState["id"] = rawState["id"]
	newState["name"] = rawState["name"]
	newState["description"] = rawState["description"]
	newState["enabled"] = rawState["enabled"]

	// Update data_type to use proper capitalization and match new schema
	if dataType, ok := rawState["data_type"]; ok {
		caser := cases.Title(language.English)
		newDataType := caser.String(strings.ToLower(dataType.(string)))
		switch newDataType {
		case "String", "Integer", "Date":
			newState["data_type"] = strings.ToUpper(newDataType)
		default:
			newState["data_type"] = "STRING" // Default value
		}
	} else {
		newState["data_type"] = "STRING" // Default value
	}

	// Map inventory_display to inventory_display_type and match new schema
	if inv, ok := rawState["inventory_display"]; ok {
		invUpper := strings.ToUpper(inv.(string))
		switch invUpper {
		case "GENERAL", "HARDWARE", "OPERATING_SYSTEM", "USER_AND_LOCATION", "PURCHASING", "EXTENSION_ATTRIBUTES":
			newState["inventory_display_type"] = invUpper
		default:
			newState["inventory_display_type"] = "EXTENSION_ATTRIBUTES" // Default value
		}
	} else {
		newState["inventory_display_type"] = "EXTENSION_ATTRIBUTES" // Default value
	}

	// Handle input_type and related fields
	if inputType, ok := rawState["input_type"]; ok {
		switch strings.ToUpper(inputType.(string)) {
		case "SCRIPT":
			newState["input_type"] = "SCRIPT"
			if script, ok := rawState["input_script"]; ok {
				newState["script_contents"] = normalizeScript(script.(string))
			}
		case "POPUP":
			newState["input_type"] = "POPUP"
			if popup, ok := rawState["input_popup"]; ok {
				newState["popup_menu_choices"] = popup
			}
		case "TEXT":
			newState["input_type"] = "TEXT"
		case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
			newState["input_type"] = "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"
		default:
			newState["input_type"] = "TEXT" // Default to TEXT if unknown
		}
	}

	// Initialize new fields with default values
	newState["ldap_attribute_mapping"] = ""
	newState["ldap_extension_attribute_allowed"] = false

	return newState, nil
}
