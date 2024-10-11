package mobiledeviceextensionattributes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceMobileDeviceExtensionAttribute object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileExtensionAttribute, error) {
	var dataTypeChange string
	var inventoryDisplayChange string

	switch d.Get("data_type").(string) {
	case "STRING":
		dataTypeChange = "String"
	case "INTEGER":
		dataTypeChange = "Integer"
	case "DATE":
		dataTypeChange = "Date"
	case "String":
		dataTypeChange = "String"
	case "Integer":
		dataTypeChange = "Integer"
	case "Date":
		dataTypeChange = "Date"
	default:
		dataTypeChange = "String"
	}

	switch d.Get("inventory_display").(string) {
	case "GENERAL":
		inventoryDisplayChange = "General"
	case "HARDWARE":
		inventoryDisplayChange = "Hardware"
	case "USER_AND_LOCATION":
		inventoryDisplayChange = "User and Location"
	case "PURCHASING":
		inventoryDisplayChange = "Purchasing"
	case "EXTENSION_ATTRIBUTES":
		inventoryDisplayChange = "Extension Attributes"
	case "General":
		inventoryDisplayChange = "General"
	case "Hardware":
		inventoryDisplayChange = "Hardware"
	case "User and Location":
		inventoryDisplayChange = "User and Location"
	case "Purchasing":
		inventoryDisplayChange = "Purchasing"
	case "Extension Attributes":
		inventoryDisplayChange = "Extension Attributes"
	default:
		inventoryDisplayChange = "Extension Attributes"
	}

	resource := &jamfpro.ResourceMobileExtensionAttribute{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		DataType:         jamfpro.MobileExtensionAttributeSubsetDataType{Type: dataTypeChange},
		InputType:        jamfpro.MobileExtensionAttributeSubsetInputType{Type: "Text Field"},
		InventoryDisplay: jamfpro.MobileExtensionAttributeSubsetInventoryDisplay{Type: inventoryDisplayChange},
	}

	// resource.InputType = jamfpro.MobileExtensionAttributeSubsetInputType{Type: (d.Get("input_type").(string))}

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Extension Attribute to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro MobileDevice Extension Attribute JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// Validate the input type
// if err := validateInputType(resource); err != nil {
// 	return nil, fmt.Errorf("failed to construct : %v", err)
// }

// validateInputType ensures that the appropriate fields are set based on the input type
// func validateInputType(resource *jamfpro.ResourceMobileExtensionAttribute) error {
// 	switch resource.InputType {
// 	case "SCRIPT":
// 		if resource.ScriptContents == "" {
// 			return fmt.Errorf("script_contents must be set when input_type is 'SCRIPT' (current value: '%s')", resource.ScriptContents)
// 		}
// 		if len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {
// 			return fmt.Errorf("popup_menu_choices and ldap_attribute_mapping should not be set when input_type is 'SCRIPT'")
// 		}
// 	case "POPUP":
// 		if len(resource.PopupMenuChoices) == 0 {
// 			return fmt.Errorf("popup_menu_choices must be set when input_type is 'POPUP'")
// 		}
// 		if resource.ScriptContents != "" || resource.LDAPAttributeMapping != "" {
// 			return fmt.Errorf("script_contents and ldap_attribute_mapping should not be set when input_type is 'POPUP'")
// 		}
// 	case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
// 		if resource.LDAPAttributeMapping == "" {
// 			return fmt.Errorf("ldap_attribute_mapping must be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
// 		}
// 		if resource.ScriptContents != "" || len(resource.PopupMenuChoices) > 0 {
// 			return fmt.Errorf("script_contents and popup_menu_choices should not be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
// 		}
// 		// Note: ldap_extension_attribute_allowed is handled by the schema's default value
// 	case "TEXT":
// 		if resource.ScriptContents != "" || len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {
// 			return fmt.Errorf("script_contents, popup_menu_choices, and ldap_attribute_mapping should not be set when input_type is 'TEXT'")
// 		}
// 	default:
// 		return fmt.Errorf("invalid input_type: %s", resource.InputType)
// 	}

// 	return nil
// }
