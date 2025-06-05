package mobile_device_extension_attribute

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// ValidateInputType ensures that the appropriate fields are set based on the input type
func ValidateInputType(resource *jamfpro.ResourceMobileExtensionAttribute) error {
	// Check if InputType is empty
	if resource.InputType.Type == "" {
		return fmt.Errorf("input_type.type must be set")
	}

	switch resource.InputType.Type {
	case "Text Field":
		if len(resource.InputType.PopupChoices.Choice) > 0 {
			return fmt.Errorf("popup_choices should not be set when input_type is 'Text Field'")
		}
	case "Pop-up Menu":
		if len(resource.InputType.PopupChoices.Choice) == 0 {
			return fmt.Errorf("popup_choices must be set when input_type is 'Pop-up Menu'")
		}
	default:
		return fmt.Errorf("invalid input_type: %s", resource.InputType.Type)
	}

	return nil
}
