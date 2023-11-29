// computergroup_helpers.go
package computergroups

import "fmt"

// assertString attempts to perform a type assertion on the given value to a string.
// If the assertion is successful, it returns the string value and no error.
// If the assertion fails, it returns an empty string and an error message indicating the failure.
func assertString(value interface{}, key string) (string, error) {
	if strVal, ok := value.(string); ok {
		return strVal, nil
	}
	return "", fmt.Errorf("type assertion to string failed for key '%s'", key)
}

// assertInt attempts to perform a type assertion on the given value to an integer.
// If the assertion is successful, it returns the integer value and no error.
// If the assertion fails, it returns 0 and an error message indicating the failure.
func assertInt(value interface{}, key string) (int, error) {
	if intVal, ok := value.(int); ok {
		return intVal, nil
	}
	return 0, fmt.Errorf("type assertion to int failed for key '%s'", key)
}

// assertBool attempts to perform a type assertion on the given value to a boolean.
// If the assertion is successful, it returns the boolean value and no error.
// If the assertion fails, it returns false and an error message indicating the failure.
func assertBool(value interface{}, key string) (bool, error) {
	if boolVal, ok := value.(bool); ok {
		return boolVal, nil
	}
	return false, fmt.Errorf("type assertion to bool failed for key '%s'", key)
}
