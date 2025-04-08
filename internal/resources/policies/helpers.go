package policies

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MapSetToStructs transforms a schema.Set of primitive values (typically IDs) into a slice of structs
// by mapping each value to a specified field in a new struct instance.
//
// Type parameters:
// - NestedObjectType: The target struct type (e.g., PolicySubsetNetworkSegment)
// - ListItemPrimitiveType: The primitive value type in the set (e.g., int, string)
//
// The function relies on Terraform's schema validation to ensure values are of the correct type.
// This is typically used to transform schema.Set elements like IDs to API struct objects.
func MapSetToStructs[NestedObjectType any, ListItemPrimitiveType comparable](path string, fieldName string, d *schema.ResourceData, outputSlice *[]NestedObjectType) error {
	*outputSlice = []NestedObjectType{}

	setVal, ok := d.GetOk(path)
	if !ok {
		return nil
	}

	set, ok := setVal.(*schema.Set)
	if !ok || set.Len() == 0 {
		return nil
	}

	result := make([]NestedObjectType, 0, set.Len())

	for _, rawValue := range set.List() {

		value, ok := rawValue.(ListItemPrimitiveType)
		if !ok {
			return fmt.Errorf("value in %s has incorrect type: expected %T",
				path, *new(ListItemPrimitiveType))
		}

		var obj NestedObjectType
		objVal := reflect.ValueOf(&obj).Elem()
		field := objVal.FieldByName(fieldName)

		field.Set(reflect.ValueOf(value))
		result = append(result, obj)
	}

	*outputSlice = result
	return nil
}
