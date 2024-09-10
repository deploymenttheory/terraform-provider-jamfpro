package computerextensionattributes

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceComputerExtensionAttribute object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerExtensionAttribute, error) {
	resource := &jamfpro.ResourceComputerExtensionAttribute{
		ID:                   d.Get("id").(string),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		DataType:             d.Get("data_type").(string),
		Enabled:              jamfpro.BoolPtr(d.Get("enabled").(bool)),
		InventoryDisplayType: d.Get("inventory_display_type").(string),
		InputType:            d.Get("input_type").(string),
	}

	// Handle input type specific fields
	switch resource.InputType {
	case "Script":
		resource.ScriptContents = d.Get("script_contents").(string)
	case "Pop-up Menu":
		choices := d.Get("popup_menu_choices").([]interface{})
		for _, choice := range choices {
			resource.PopupMenuChoices = append(resource.PopupMenuChoices, choice.(string))
		}
	case "LDAP Mapping":
		resource.LDAPAttributeMapping = d.Get("ldap_attribute_mapping").(string)
		resource.LDAPExtensionAttributeAllowed = jamfpro.BoolPtr(d.Get("ldap_extension_attribute_allowed").(bool))
	}

	// Validation
	if err := validateInputType(resource); err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Extension Attribute: %+v\n", resource)

	return resource, nil
}

// validateInputType ensures that the appropriate fields are set based on the input type
func validateInputType(resource *jamfpro.ResourceComputerExtensionAttribute) error {
	switch resource.InputType {
	case "Script":
		if resource.ScriptContents == "" {
			return fmt.Errorf("script_contents must be set when input_type is 'Script'")
		}
	case "Pop-up Menu":
		if len(resource.PopupMenuChoices) == 0 {
			return fmt.Errorf("popup_menu_choices must be set when input_type is 'Pop-up Menu'")
		}
	case "LDAP Mapping":
		if resource.LDAPAttributeMapping == "" {
			return fmt.Errorf("ldap_attribute_mapping must be set when input_type is 'LDAP Mapping'")
		}
	case "Text Field":
		// No additional fields required for Text Field
	default:
		return fmt.Errorf("invalid input_type: %s", resource.InputType)
	}
	return nil
}
