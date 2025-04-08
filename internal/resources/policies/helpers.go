package policies

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MapSetToStructs transforms a schema.Set of primitive values into a slice of structs
// by mapping each value to a named field in a new struct instance.
//
// This function uses Go generics to provide type safety while remaining flexible:
// - T any: Represents the target struct type we're creating instances of
// - V comparable: Represents the primitive value type from the schema.Set
//
// By using generics, the function can handle different combinations of struct types
// and primitive values without duplicating code or sacrificing type safety.
func MapSetToStructs[T any, V comparable](path string, fieldName string, d *schema.ResourceData, outputSlice *[]T) error {
	*outputSlice = []T{}

	setVal, ok := d.GetOk(path)
	if !ok {
		return nil
	}

	set, ok := setVal.(*schema.Set)
	if !ok || set.Len() == 0 {
		return nil
	}

	// Pre-allocate slice with capacity matching input size for better performance
	result := make([]T, 0, set.Len())

	// Process each value in the set to build our struct objects
	for _, v := range set.List() {
		if v == nil {
			continue // Skip nil values to prevent reflection panics
		}

		var obj T
		objVal := reflect.ValueOf(&obj).Elem()
		field := objVal.FieldByName(fieldName)

		field.Set(reflect.ValueOf(v))
		result = append(result, obj)
	}

	*outputSlice = result
	return nil
}
