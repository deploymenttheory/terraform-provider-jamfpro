// common/configurationprofiles/plist/validate.go
package plist

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidatePayloadFields validates the mapstructure field tags of a ConfigurationProfile struct
// and ensures that the required struct field tags are present and correctly populated.
// It checks for specific validation rules, such as ensuring that required fields are not empty.
// Additional custom validation rules can be added within this function to enforce other constraints
// based on the `validate` tags associated with the struct fields.
func ValidatePayloadFields(profile *ConfigurationProfile) []error {
	var errs []error

	// Iterate over struct fields to validate 'required' tags
	val := reflect.ValueOf(profile).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")
		if tag != "" {
			// Check for required fields
			if strings.Contains(tag, "required") {
				value := val.Field(i).Interface()
				if value == "" {

					errs = append(errs, fmt.Errorf("plist key '%s' is required", field.Name))
				}
			}
			// Additional validation rules can be added here
		}
	}

	return errs
}
