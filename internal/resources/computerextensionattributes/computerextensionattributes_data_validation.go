// computerextensionattributes_data_validation.go
package computerextensionattributes

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateDataType ensures the provided value adheres to the accepted formats for the data_type attribute.
// The accepted formats are "String", "Integer", and a date string in the "YYYY-MM-DD hh:mm:ss" format.
func validateDataType(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)

	// Regular expression to validate the date format "YYYY-MM-DD hh:mm:ss"
	datePattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`

	if value != "String" && value != "Integer" && !regexp.MustCompile(datePattern).MatchString(value) {
		errs = append(errs, fmt.Errorf("%q must be 'String', 'Integer', or 'YYYY-MM-DD hh:mm:ss' format, got: %s", key, value))
	}
	return
}

// validateJamfProRResourceComputerExtensionAttributesDataFields performs custom validation on the Resource's schema so that passed values from
// teraform resource declarations align with attibute combinations supported by the Jamf Pro api.
func validateJamfProRResourceComputerExtensionAttributesDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Extract the first item from the input_type list, which should be a map
	inputTypes, ok := diff.GetOk("input_type")
	if !ok || len(inputTypes.([]interface{})) == 0 {
		return fmt.Errorf("input_type must be provided")
	}

	inputTypeMap := inputTypes.([]interface{})[0].(map[string]interface{})

	inputType := inputTypeMap["type"].(string)
	platform := inputTypeMap["platform"].(string)
	script := inputTypeMap["script"].(string)
	choices := inputTypeMap["choices"].([]interface{})

	switch inputType {
	case "script":
		// Ensure platform is either "Mac" or "Windows"
		if platform != "Mac" && platform != "Windows" {
			return fmt.Errorf("platform must be either 'Mac' or 'Windows' when input_type is 'script'")
		}
		// Ensure "script" is populated
		if script == "" {
			return fmt.Errorf("'script' field must be populated when input_type is 'script'")
		}
		// Ensure "choices" is not populated
		if len(choices) > 0 {
			return fmt.Errorf("'choices' must not be populated when input_type is 'script'")
		}
	case "Pop-up Menu":
		// Ensure "choices" is populated
		if len(choices) == 0 {
			return fmt.Errorf("'choices' must be populated when input_type is 'Pop-up Menu'")
		}
		// Ensure platform and script are not populated
		if platform != "" {
			return fmt.Errorf("'platform' must not be populated when input_type is 'Pop-up Menu'")
		}
		if script != "" {
			return fmt.Errorf("'script' must not be populated when input_type is 'Pop-up Menu'")
		}
	case "Text Field":
		// Ensure neither "script", "platform" nor "choices" are populated
		if script != "" {
			return fmt.Errorf("'script' field must not be populated when input_type is 'Text Field'")
		}
		if len(choices) > 0 {
			return fmt.Errorf("'choices' must not be populated when input_type is 'Text Field'")
		}
		if platform != "" {
			return fmt.Errorf("'platform' must not be populated when input_type is 'Text Field'")
		}
	}

	return nil
}
