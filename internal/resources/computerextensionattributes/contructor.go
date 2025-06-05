package computerextensionattributes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceComputerExtensionAttribute object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerExtensionAttribute, error) {
	resource := &jamfpro.ResourceComputerExtensionAttribute{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		DataType:             d.Get("data_type").(string),
		Enabled:              jamfpro.BoolPtr(d.Get("enabled").(bool)),
		InventoryDisplayType: d.Get("inventory_display_type").(string),
		InputType:            d.Get("input_type").(string),
	}

	if v, ok := d.GetOk("script_contents"); ok {
		resource.ScriptContents = v.(string)
	}

	if v, ok := d.GetOk("popup_menu_choices"); ok {
		choices := v.([]interface{})
		for _, choice := range choices {
			resource.PopupMenuChoices = append(resource.PopupMenuChoices, choice.(string))
		}
	}

	if v, ok := d.GetOk("ldap_attribute_mapping"); ok {
		resource.LDAPAttributeMapping = v.(string)
	}

	if v, ok := d.GetOk("ldap_extension_attribute_allowed"); ok {
		resource.LDAPExtensionAttributeAllowed = jamfpro.BoolPtr(v.(bool))
	}

	// Validate the input type
	if err := validateInputType(resource); err != nil {

		return nil, fmt.Errorf("failed to construct : %v", err)
	}

	// Serialize and pretty-print the inventory collection object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Extension Attribute to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Extension Attribute JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// validateInputType ensures that the appropriate fields are set based on the input type
func validateInputType(resource *jamfpro.ResourceComputerExtensionAttribute) error {
	switch resource.InputType {
	case "SCRIPT":
		if resource.ScriptContents == "" {

			return fmt.Errorf("script_contents must be set when input_type is 'SCRIPT' (current value: '%s')", resource.ScriptContents)
		}
		if len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {

			return fmt.Errorf("popup_menu_choices and ldap_attribute_mapping should not be set when input_type is 'SCRIPT'")
		}
	case "POPUP":
		if len(resource.PopupMenuChoices) == 0 {

			return fmt.Errorf("popup_menu_choices must be set when input_type is 'POPUP'")
		}
		if resource.ScriptContents != "" || resource.LDAPAttributeMapping != "" {

			return fmt.Errorf("script_contents and ldap_attribute_mapping should not be set when input_type is 'POPUP'")
		}
	case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
		if resource.LDAPAttributeMapping == "" {

			return fmt.Errorf("ldap_attribute_mapping must be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
		}
		if resource.ScriptContents != "" || len(resource.PopupMenuChoices) > 0 {

			return fmt.Errorf("script_contents and popup_menu_choices should not be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
		}
		// Note: ldap_extension_attribute_allowed is handled by the schema's default value
	case "TEXT":
		if resource.ScriptContents != "" || len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {

			return fmt.Errorf("script_contents, popup_menu_choices, and ldap_attribute_mapping should not be set when input_type is 'TEXT'")
		}
	default:

		return fmt.Errorf("invalid input_type: %s", resource.InputType)
	}

	return nil
}
