// utilities.go
// For utility/helper functions to support the jamf pro tf provider
package utilities

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Helper function to extract string slice from *schema.Set
func ExtractSetToStringSlice(set *schema.Set) []string {
	list := set.List() // This converts the set to a slice of interface{}
	slice := make([]string, len(list))
	for i, item := range list {
		slice[i] = item.(string) // Type assertion, since we know the set should only contain strings
	}
	return slice
}

// Helper function to convert slice of strings to slice of empty interfaces
func ConvertToStringInterface(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, d := range slice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

// Numeric Conversions

// IntToFloat64 converts an int to a float64.
func IntToFloat64(i int) float64 {
	return float64(i)
}

// Float64ToInt converts a float64 to an int, truncating the decimal part.
func Float64ToInt(f float64) int {
	return int(f)
}

// String and Numeric Conversions

// StringToInt converts a string to an int. Returns an error if the conversion is not possible.
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("conversion failed: %w", err)
	}
	return i, nil
}

// IntToString converts an int to a string.
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Byte Slices and Strings

// BytesToString converts a byte slice to a string.
func BytesToString(b []byte) string {
	return string(b)
}

// StringToBytes converts a string to a byte slice.
func StringToBytes(s string) []byte {
	return []byte(s)
}

// Slices Conversion

// SliceIntToFloat64 converts a slice of int to a slice of float64.
func SliceIntToFloat64(slice []int) []float64 {
	result := make([]float64, len(slice))
	for i, v := range slice {
		result[i] = float64(v)
	}
	return result
}

// SliceFloat64ToInt converts a slice of float64 to a slice of int, truncating the decimal part.
func SliceFloat64ToInt(slice []float64) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = int(v)
	}
	return result
}

// ConvertSlice converts a slice of one type to a slice of another type.
// It requires a conversion function that defines how individual elements are converted.
func ConvertSlice[T any, U any](s []T, convert func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = convert(v)
	}
	return result
}

// JSON Conversion

// StructToJSON converts a struct to a JSON string. Returns an error if the conversion is not possible.
func StructToJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("JSON marshaling failed: %w", err)
	}
	return string(b), nil
}

// JSONToStruct converts a JSON string to a struct. The struct type must be provided by the caller. Returns an error if the conversion is not possible.
func JSONToStruct(data string, v interface{}) error {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		return fmt.Errorf("JSON unmarshaling failed: %w", err)
	}
	return nil
}

// Map Conversions

// MapInterfaceToString converts a map with interface{} values to a map with string values.
// Non-string values will be converted to strings using fmt.Sprint.
func MapInterfaceToString(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = val
		default:
			result[k] = fmt.Sprint(v)
		}
	}
	return result
}

// MapStringToInterface converts a map with string values to a map with interface{} values.
func MapStringToInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

// ConvertMap converts a map from one key-value pair type to another.
// It requires conversion functions for keys and values respectively.
func ConvertMap[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, convertKey func(K1) K2, convertValue func(V1) V2) map[K2]V2 {
	result := make(map[K2]V2)
	for k, v := range m {
		newKey := convertKey(k)
		newValue := convertValue(v)
		result[newKey] = newValue
	}
	return result
}

// Nested Slice and Map Conversions

// SliceMapStringToInterface converts a slice of maps with string keys and string values to a slice of maps with string keys and interface{} values.
func SliceMapStringToInterface(slice []map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(slice))
	for i, m := range slice {
		result[i] = MapStringToInterface(m)
	}
	return result
}

// SliceMapInterfaceToString converts a slice of maps with string keys and interface{} values to a slice of maps with string keys and string values.
// Non-string values will be converted to strings using fmt.Sprint.
func SliceMapInterfaceToString(slice []map[string]interface{}) []map[string]string {
	result := make([]map[string]string, len(slice))
	for i, m := range slice {
		result[i] = MapInterfaceToString(m)
	}
	return result
}

// Handling Type Assertions for Maps

// GetStringValueFromMap tries to retrieve a string value from a map using the provided key.
// Returns the value and a boolean indicating whether the operation was successful.
func GetStringValueFromMap(m map[string]interface{}, key string) (string, bool) {
	if val, exists := m[key]; exists {
		strVal, ok := val.(string)
		return strVal, ok
	}
	return "", false
}
