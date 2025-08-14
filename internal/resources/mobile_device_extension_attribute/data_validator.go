package mobile_device_extension_attribute

import (
	"errors"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

var (
	errPopupMenuChoicesRequired                              = errors.New("popup_menu_choices must be set when input_type is 'POPUP'")
	errLDAPAttributeMappingShouldNotBeSet                    = errors.New("ldap_attribute_mapping should not be set when input_type is 'POPUP'")
	errLDAPAttributeMappingRequired                          = errors.New("ldap_attribute_mapping must be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errPopupMenuChoicesShouldNotBeSet                        = errors.New("popup_menu_choices should not be set when input_type is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errPopupMenuChoicesAndLDAPAttributeMappingShouldNotBeSet = errors.New("popup_menu_choices, and ldap_attribute_mapping should not be set when input_type is 'TEXT'")
	errInvalidInputType                                      = errors.New("invalid input_type")
)

// ValidateInputType ensures that the appropriate fields are set based on the input type
func validateInputType(resource *jamfpro.ResourceMobileDeviceExtensionAttribute) error {
	switch resource.InputType {
	case "POPUP":
		if len(resource.PopupMenuChoices) == 0 {
			return fmt.Errorf("%w", errPopupMenuChoicesRequired)
		}
		if resource.LDAPAttributeMapping != "" {
			return fmt.Errorf("%w", errLDAPAttributeMappingShouldNotBeSet)
		}
	case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
		if resource.LDAPAttributeMapping == "" {
			return fmt.Errorf("%w", errLDAPAttributeMappingRequired)
		}
		if len(resource.PopupMenuChoices) > 0 {
			return fmt.Errorf("%w", errPopupMenuChoicesShouldNotBeSet)
		}
	case "TEXT":
		if len(resource.PopupMenuChoices) > 0 || resource.LDAPAttributeMapping != "" {
			return fmt.Errorf("%w", errPopupMenuChoicesAndLDAPAttributeMappingShouldNotBeSet)
		}
	default:
		return fmt.Errorf("%w: %s", errInvalidInputType, resource.InputType)
	}

	return nil
}
