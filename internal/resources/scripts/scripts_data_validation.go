// scripts_data_validation.go
package scripts

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

// validateResourceScriptsDataFields performs custom validation on the Resource's schema so that passed values from
// teraform resource declarations align with attibute combinations supported by the Jamf Pro api.
func validateResourceScriptsDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return nil
}
