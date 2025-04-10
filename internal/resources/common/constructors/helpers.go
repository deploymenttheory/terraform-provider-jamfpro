package constructors

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ParseResourceID ensures a value can be safely used as an ID in Jamf Pro resources.
// Handles int, float64, and string inputs with appropriate logging for conversion issues.
// Returns the validated integer ID and a success boolean.
func ParseResourceID(val interface{}, fieldName string, index int) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case string:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("[WARN] Could not convert string '%s' to int for %s ID at index %d: %v. Skipping.", v, fieldName, index, err)
			return 0, false
		}
		return intVal, true
	default:
		log.Printf("[WARN] Unexpected type %T for %s ID: %v at index %d. Skipping.", val, fieldName, val, index)
		return 0, false
	}
}

// GetListFromSet is a helper function to safely extract a list of interfaces
// from a *schema.Set stored within a map[string]interface{}.
// It returns an empty slice if the key is not found or the value is not a *schema.Set.
func GetListFromSet(data map[string]interface{}, key string) []interface{} {
	v, ok := data[key]
	if !ok || v == nil {
		return []interface{}{}
	}
	set, ok := v.(*schema.Set)
	if !ok || set == nil {
		log.Printf("[DEBUG] getListFromSet: Value for key '%s' is not a *schema.Set or is nil, type is %T", key, v)
		return []interface{}{}
	}
	return set.List()
}

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
