package mobile_device_extension_attribute

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// ValidateInputType ensures that the appropriate fields are set based on the input type
func validateInputType(resource *jamfpro.ResourceMobileDeviceExtensionAttribute) error {
	switch resource.InputType {
	case "POPUP":
		if len(resource.PopupMenuChoices) == 0 {
			return fmt.Errorf("popup_menu_choices must be set when input_type is 'POPUP'")
		}
		if resource.LDAPAttributeMapping != "" {
			return fmt.Errorf("ldap_attribute_mapping should not be set when input_type is 'POPUP'")
		}
	case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
		if resource.LDAPAttributeMapping == "" {
			return fmt.Errorf("ldap_attribute_mapping must be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
		}
		if len(resource.PopupMenuChoices) > 0 {
			return fmt.Errorf("popup_menu_choices should not be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
		}
	case "TEXT":
		if len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {
			return fmt.Errorf("popup_menu_choices, and ldap_attribute_mapping should not be set when input_type is 'TEXT'")
		}
	default:
		return fmt.Errorf("invalid input_type: %s", resource.InputType)
	}

	return nil
}
