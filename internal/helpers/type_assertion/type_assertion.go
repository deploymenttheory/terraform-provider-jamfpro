// type_assertion.go
package type_assertion

/*
Package type_assertion provides a set of helper functions designed to handle type assertions in scenarios where the value can be null (nil in Go).
These functions are tailored to return default values when encountering nil inputs or when type assertions fail, making them robust for use in various scenarios.

For Map Functions (GetIntFromMap, GetBoolFromMap, GetStringFromMap):
- These functions check if the specified key exists in a map and whether the value associated with the key is nil.
- If the key doesn't exist or the value is nil, they return default values (0 for integers, false for booleans, and "" for strings).
- If the value exists but fails the type assertion, they also return the default values.

For Interface Functions (GetStringFromInterface, GetIntFromInterface, GetBoolFromInterface):
- These functions perform a direct type assertion on the provided interface{} value.
- They return a default value if the value is nil or if the type assertion fails.

For Direct Type Assertion Functions (GetString, GetInt):
- Similar to the interface functions, these functions directly assert the type of the provided interface{} value.
- They return default values if the value is nil or the type assertion fails.
- This approach simplifies handling nil values and type assertion errors in calling code, ensuring more readable and maintainable code.

The design of these functions makes them suitable for scenarios where nil values are acceptable and should not be treated as errors,
such as when extracting data from a Terraform schema.
*/

// GetIntFromMap safely retrieves an int value from a map, returning a default value for nil or not found.
func GetIntFromMap(m map[string]interface{}, key string) int {
	val, ok := m[key]
	if !ok || val == nil {
		return 0
	}
	intVal, ok := val.(int)
	if !ok {
		return 0
	}
	return intVal
}

// GetBoolFromMap safely retrieves a bool value from a map, returning a default value for nil or not found.
func GetBoolFromMap(m map[string]interface{}, key string) bool {
	val, ok := m[key]
	if !ok || val == nil {
		return false
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false
	}
	return boolVal
}

// GetStringFromMap safely retrieves a string value from a map, returning a default value for nil or not found.
func GetStringFromMap(m map[string]interface{}, key string) string {
	val, ok := m[key]
	if !ok || val == nil {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

// ConvertToMapFromInterface safely converts an interface{} to a map[string]interface{}, handling nil values.
func ConvertToMapFromInterface(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}
	val, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	return val
}

// GetStringFromInterface safely retrieves a string value from an interface{}, handling nil values.
func GetStringFromInterface(val interface{}) string {
	if val == nil {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

// GetIntFromInterface safely retrieves an int value from an interface{}, handling nil values.
func GetIntFromInterface(val interface{}) int {
	if val == nil {
		return 0
	}
	intVal, ok := val.(int)
	if !ok {
		return 0
	}
	return intVal
}

// GetBoolFromInterface safely retrieves a bool value from an interface{}, handling nil values.
func GetBoolFromInterface(val interface{}) bool {
	if val == nil {
		return false
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false
	}
	return boolVal
}

// GetString safely performs a string type assertion, handling nil values.
func GetString(val interface{}) string {
	if val == nil {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

// GetInt safely performs an int type assertion, handling nil values.
func GetInt(val interface{}) int {
	if val == nil {
		return 0
	}
	intVal, ok := val.(int)
	if !ok {
		return 0
	}
	return intVal
}

// GetStringBoolMapFromInterface safely retrieves a map[string]bool value from an interface{}, handling nil values and various keys.
func GetStringBoolMapFromInterface(val interface{}) map[string]bool {
	if val == nil {
		return nil
	}

	resultMap := make(map[string]bool)
	mapVal, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}

	for key, v := range mapVal {
		// Assuming that non-existent or non-boolean values should default to false
		boolVal, _ := v.(bool)
		resultMap[key] = boolVal
	}
	return resultMap
}

// GetStringSliceFromInterface safely converts an interface slice to a string slice.
func GetStringSliceFromInterface(input interface{}) []string {
	var result []string
	if inputSlice, ok := input.([]interface{}); ok {
		for _, item := range inputSlice {
			strVal := GetStringFromInterface(item) // Utilizing existing function
			if strVal != "" {                      // Append only if the value is non-empty
				result = append(result, strVal)
			}
		}
	}
	return result
}

// ConvertInterfaceSliceToStringMap safely converts an interface slice to a map of string slices.
func ConvertInterfaceSliceToStringMap(input interface{}) map[string][]string {
	result := make(map[string][]string)
	if inputSlice, ok := input.([]interface{}); ok {
		for _, item := range inputSlice {
			if itemMap, ok := item.(map[string]interface{}); ok {
				for key, val := range itemMap {
					if valSlice, ok := val.([]interface{}); ok {
						for _, valItem := range valSlice {
							strVal := GetStringFromInterface(valItem)
							if strVal != "" {
								result[key] = append(result[key], strVal)
							}
						}
					}
				}
			}
		}
	}
	return result
}
