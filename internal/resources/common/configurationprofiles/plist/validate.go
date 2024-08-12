// common/configurationprofiles/plist/validate.go
// Description: This file contains the Configuration Profile validation functions.
package plist

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// ValidatePayload validates a payload by unmarshalling it and checking for required fields.
func ValidatePayload(payload interface{}, key string) (warns []string, errs []error) {
	profile, err := UnmarshalPayload(payload.(string))
	if err != nil {
		errs = append(errs, err)
		return warns, errs
	}

	if profile.PayloadIdentifier != profile.PayloadUUID {
		warns = append(warns, "Top-level PayloadIdentifier should match top-level PayloadUUID")
	}

	// Custom validation
	errs = ValidatePayloadFields(profile)

	return warns, errs
}

// ValidatePayloadFields validates the mapstructure field tags of a ConfigurationProfile struct
// and ensures that the required struct field tags are present and correctly populated.
// It checks for specific validation rules, such as ensuring that required fields are not empty.
// Additional custom validation rules can be added within this function to enforce other constraints
// based on the `validate` tags associated with the struct fields.
func ValidatePayloadFields(profile *ConfigurationProfile) []error {
	var errs []error

	// Check 1: Validate that PayloadUUID is a valid UUID
	if !isValidUUID(profile.PayloadUUID) {
		errs = append(errs, errors.New("PayloadUUID must be a valid UUID"))
	}

	// Check 2: Ensure PayloadScope is one of the expected values (System, User, Computer)
	expectedScopes := []string{"System", "User", "Computer"}
	if !contains(expectedScopes, profile.PayloadScope) {
		errs = append(errs, errors.New("PayloadScope must be one of 'System', 'User', or 'Computer'"))
	}

	// Check 3: Ensure PayloadIdentifier matches PayloadUUID
	if profile.PayloadIdentifier != profile.PayloadUUID {
		errs = append(errs, errors.New("PayloadIdentifier should match PayloadUUID"))
	}

	// Check 4: Ensure PayloadVersion is a positive integer
	if profile.PayloadVersion <= 0 {
		errs = append(errs, errors.New("PayloadVersion must be a positive integer"))
	}

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
					errs = append(errs, fmt.Errorf(fmt.Sprintf("plist key '%s' is required", field.Name)))
				}
			}
			// Additional validation rules can be added here
		}
	}

	return errs
}

// Helper function to check if a string is a valid UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// Helper function to check if a string is in a slice of strings
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
