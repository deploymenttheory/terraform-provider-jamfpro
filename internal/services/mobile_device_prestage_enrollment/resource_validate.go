// mobiledeviceprestageenrollments_data_validator.go
package mobile_device_prestage_enrollment

import (
	"fmt"
	"regexp"
)

// validateDateFormat checks that the date is in the format YYYY-MM-DD, but only if the value is not null or empty.
func validateDateFormat(v any, k string) (ws []string, errors []error) {
	dateString, ok := v.(string)
	if !ok {
		return
	}

	if dateString == "" {
		return
	}

	datePattern := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(datePattern, dateString)

	if !match {
		errors = append(errors, fmt.Errorf("%q must be in the format YYYY-MM-DD, got: %s", k, dateString))
	}

	return
}
