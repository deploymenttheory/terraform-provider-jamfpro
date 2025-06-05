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

// GetHCLStringOrDefaultInteger returns the string value from the ResourceData if it exists,
// otherwise it returns the default value "-1".
func GetHCLStringOrDefaultInteger(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return "-1"
}

// GetDateOrDefaultDate returns the date string if it exists and is not empty,
// otherwise it returns the default date "1970-01-01".
func GetDateOrDefaultDate(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok && v.(string) != "" {
		return v.(string)
	}
	return "1970-01-01"
}

// HandleVersionLock manages the VersionLock field for Jamf Pro Prestage resources during update operations.
//
// https://developer.jamf.com/jamf-pro/docs/optimistic-locking
//
// Parameters:
//   - currentVersionLock: The current version lock value as an interface{}.
//   - isUpdate: A boolean flag indicating whether this is an update operation.
//
// Returns:
//   - An integer representing the version lock to be used in the API request.
//     For create operations (isUpdate == false), this will be 0.
//     For update operations (isUpdate == true), this will be the incremented version lock.
//
// Behavior:
//   - Create operations (isUpdate == false):
//   - Returns 0, as version lock is not needed for create operations.
//   - Update operations (isUpdate == true):
//   - Attempts to convert the currentVersionLock to an integer and increment it by 1.
//   - If conversion fails, logs a warning and returns 0.
//
// Error Handling:
//   - If the currentVersionLock cannot be converted to an integer during an update operation,
//     the function logs a warning and returns 0.
//
// Usage:
//   - This function should be called for each structure within a resource that requires
//     version lock handling.
func HandleVersionLock(currentVersionLock interface{}, isUpdate bool) int {
	if !isUpdate {
		log.Printf("[DEBUG] Create operation: Version lock not required, using 0")
		return 0
	}

	log.Printf("[DEBUG] Update operation: Current version lock is '%v'", currentVersionLock)

	versionLock, ok := currentVersionLock.(int)
	if !ok {
		log.Printf("[WARN] Failed to convert version lock '%v' to integer. Using 0.", currentVersionLock)
		return 0
	}

	newVersionLock := versionLock + 1
	log.Printf("[DEBUG] Update operation: Incrementing version lock from '%d' to '%d'", versionLock, newVersionLock)
	return newVersionLock
}
