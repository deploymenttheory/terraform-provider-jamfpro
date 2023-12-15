// type_assertion.go
package type_assertion

// Helper function to safely get an int value from a map. Returns 0 if key is absent or nil.
func GetIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok && val != nil {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0 // Return default zero value if key is not found or nil
}

// GetBoolFromMap safely retrieves a bool value from a map.
// It returns false if the key is absent, nil, or not a bool.
func GetBoolFromMap(m map[string]interface{}, key string) bool {
	val, ok := m[key]
	if !ok {
		// Key not found or value is nil, return the default false value
		return false
	}

	boolVal, ok := val.(bool)
	if !ok {
		// Value is not of type bool, return the default false value
		return false
	}

	return boolVal
}

// Helper function to safely get a string value from a map. Returns an empty string if key is absent or nil.
func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return "" // Return default empty string if key is not found or nil
}

// ConvertToMapFromInterface is a helper function to safely convert an interface{} to a map[string]interface{}. Returns nil if the conversion is not possible.
func ConvertToMapFromInterface(value interface{}) map[string]interface{} {
	if val, ok := value.(map[string]interface{}); ok {
		return val
	}
	return nil // Return nil if conversion is not possible
}

// GetStringFromInterface safely retrieves a string value from an interface{}.
// It returns an empty string if the value is absent, nil, or not a string.
func GetStringFromInterface(val interface{}) string {
	if strVal, ok := val.(string); ok {
		return strVal
	}
	return "" // Return default empty string if value is not a string
}

// GetIntFromInterface safely retrieves an int value from an interface{}.
// It returns 0 if the value is absent, nil, or not an int.
func GetIntFromInterface(val interface{}) int {
	if intVal, ok := val.(int); ok {
		return intVal
	}
	return 0 // Return default zero value if value is not an int
}

// GetBoolFromInterface safely retrieves a bool value from an interface{}.
// It returns false if the value is absent, nil, or not a bool.
func GetBoolFromInterface(val interface{}) bool {
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	return false // Return default false value if value is not a bool
}

// GetString safely performs a string type assertion.
// It returns the asserted string and a boolean indicating whether the assertion was successful.
func GetString(val interface{}) (string, bool) {
	strVal, ok := val.(string)
	return strVal, ok
}
