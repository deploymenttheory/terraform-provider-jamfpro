// common/helpers.go
package common

import (
	"reflect"
)

// ConvertToInterfaceSlice converts a slice of any type to a slice of empty interfaces.
func ConvertToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil // Optionally, you could also return an error here
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}

	return result
}
