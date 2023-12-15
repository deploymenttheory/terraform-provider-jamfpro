// type_assertion.go
package type_assertion

// Helper function to safely get an int value from a map. Returns 0 if key is absent or nil.
func getIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok && val != nil {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0 // Return default zero value if key is not found or nil
}

// Helper function to safely get a bool value from a map. Returns false if key is absent or nil.
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok && val != nil {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false // Return default false value if key is not found or nil
}

// Helper function to safely get a string value from a map. Returns an empty string if key is absent or nil.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return "" // Return default empty string if key is not found or nil
}

// Helper function to safely convert an interface{} to a map[string]interface{}. Returns nil if the conversion is not possible.
func getMapFromInterface(value interface{}) map[string]interface{} {
	if val, ok := value.(map[string]interface{}); ok {
		return val
	}
	return nil // Return nil if conversion is not possible
}
