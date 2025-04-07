package constructors

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertToInt attempts robust conversion to int. Returns value and success bool.
func ConvertToInt(val interface{}, fieldName string, index int) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		intVal := int(v)
		if float64(intVal) != v {
			log.Printf("[WARN] Loss of precision converting float64 %f to int for %s ID at index %d. Using truncated value %d.", v, fieldName, index, intVal)
		}
		return intVal, true
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

// getListFromSet is a helper function to safely extract a list of interfaces
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
