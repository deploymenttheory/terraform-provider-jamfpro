// type_assertion.go
package type_assertion

import "fmt"

/*
A set of helper functions designed to handle type assertion. The funcs support scenarios where the value is null (nil in Go).
When you check for val, ok := m[key], you're not only checking if the key exists in the map but also if the value associated with that key is nil.

For Map Functions (GetIntFromMap, GetBoolFromMap, GetStringFromMap):

val, ok := m[key]: This checks if the key exists in the map. If the key does not exist, ok will be false.
if !ok || val == nil: This checks both if the key is not found (!ok) and if the value is nil (val == nil). If either condition is true, the function returns a default value (0 for integers, false for booleans, "" for strings) and false for the boolean indicator.
For Interface Functions (GetStringFromInterface, GetIntFromInterface, GetBoolFromInterface):

These functions directly perform a type assertion on the provided interface{} value.
If val is nil, the type assertion will fail, and ok will be false. This is because a nil interface cannot be asserted to any concrete type.
For Direct Type Assertion Functions (GetString, GetInt):

These functions are similar to the interface functions. They perform a type assertion and return the result along with the success indicator.
Again, if the passed interface{} is nil, the assertion will fail, and ok will be false.
*/

// GetIntFromMap safely retrieves an int value from a map, returning a default value for nil or not found.
func GetIntFromMap(m map[string]interface{}, key string) (int, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return 0, nil
	}
	intVal, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("value for key '%s' is not of type int", key)
	}
	return intVal, nil
}

// GetBoolFromMap safely retrieves a bool value from a map, returning a default value for nil or not found.
func GetBoolFromMap(m map[string]interface{}, key string) (bool, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return false, nil
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value for key '%s' is not of type bool", key)
	}
	return boolVal, nil
}

// GetStringFromMap safely retrieves a string value from a map, returning a default value for nil or not found.
func GetStringFromMap(m map[string]interface{}, key string) (string, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return "", nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value for key '%s' is not of type string", key)
	}
	return strVal, nil
}

// ConvertToMapFromInterface safely converts an interface{} to a map[string]interface{}, handling nil values.
func ConvertToMapFromInterface(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, nil
	}
	val, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not a map[string]interface{}")
	}
	return val, nil
}

// GetStringFromInterface safely retrieves a string value from an interface{}, handling nil values.
func GetStringFromInterface(val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}
	return strVal, nil
}

// GetIntFromInterface safely retrieves an int value from an interface{}, handling nil values.
func GetIntFromInterface(val interface{}) (int, error) {
	if val == nil {
		return 0, nil
	}
	intVal, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("value is not an int")
	}
	return intVal, nil
}

// GetBoolFromInterface safely retrieves a bool value from an interface{}, handling nil values.
func GetBoolFromInterface(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value is not a bool")
	}
	return boolVal, nil
}

// GetString safely performs a string type assertion, handling nil values.
func GetString(val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}
	return strVal, nil
}

// GetInt safely performs an int type assertion, handling nil values.
func GetInt(val interface{}) (int, error) {
	if val == nil {
		return 0, nil
	}
	intVal, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("value is not an int")
	}
	return intVal, nil
}
