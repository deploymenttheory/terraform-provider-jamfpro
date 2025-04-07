// policies/helpers.go
package policies

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PopulateStructSliceFromSetField is a helper function that takes a path to a *schema.Set in HCL
// and populates a slice of pointers to structs, setting a specific target field in each struct.
func PopulateStructSliceFromSetField[NestedObjectType any, ListItemPrimitiveType comparable](path string, target_field string, d *schema.ResourceData, home *[]NestedObjectType) (err error) {
	getAttr, ok := d.GetOk(path)

	if !ok || getAttr == nil {
		log.Printf("[DEBUG] PopulateStructSliceFromSetField: Attribute not found or nil at path %s. Ensuring target slice is initialized.", path)
		if *home == nil {
			*home = []NestedObjectType{}
		}
		return nil
	}

	attrSet, isSet := getAttr.(*schema.Set)
	if !isSet {
		return fmt.Errorf("internal error: attribute at path %s was expected to be *schema.Set, but got %T", path, getAttr)
	}

	itemsList := attrSet.List()

	if len(itemsList) == 0 {
		log.Printf("[DEBUG] PopulateStructSliceFromSetField: Set at path %s is empty. Ensuring target slice is initialized.", path)
		if *home == nil {
			*home = []NestedObjectType{}
		}
		return nil
	}

	if *home == nil {
		*home = []NestedObjectType{}
	}

	outList := make([]NestedObjectType, 0, len(itemsList))
	expectedType := reflect.TypeOf((*ListItemPrimitiveType)(nil)).Elem() // Determine expected type once

	for i, v := range itemsList {
		if v == nil {
			log.Printf("[WARN] PopulateStructSliceFromSetField: Found nil element at index %d in set at path %s", i, path)
			continue
		}

		var primitiveValue ListItemPrimitiveType
		var conversionOk bool
		var concreteValue interface{}

		primitiveValue, conversionOk = v.(ListItemPrimitiveType)

		if !conversionOk {
			log.Printf("[DEBUG] PopulateStructSliceFromSetField: Direct type assertion failed for element at index %d (path %s). Expected %s, got %T. Attempting conversion.", i, path, expectedType, v)

			switch expectedType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if floatVal, isFloat := v.(float64); isFloat {
					convertedIntVal := reflect.ValueOf(int64(floatVal)).Convert(expectedType)
					concreteValue = convertedIntVal.Interface() // Store concrete int type
					conversionOk = true
					if float64(int64(floatVal)) != floatVal {
						log.Printf("[WARN] PopulateStructSliceFromSetField: Potential precision loss converting float64 %f to %s for element at index %d (path %s)", floatVal, expectedType, i, path)
					} else {
						log.Printf("[DEBUG] PopulateStructSliceFromSetField: Successfully converted float64 %f to %s for element at index %d (path %s)", floatVal, expectedType, i, path)
					}
				}
			case reflect.String:
				if numVal, isNum := v.(float64); isNum {
					concreteValue = fmt.Sprintf("%g", numVal)
					conversionOk = true
					log.Printf("[DEBUG] PopulateStructSliceFromSetField: Converted number %f to string for element at index %d (path %s)", numVal, i, path)
				}
			default:
				log.Printf("[WARN] PopulateStructSliceFromSetField: Unhandled expected type kind %s for conversion at index %d (path %s)", expectedType.Kind(), i, path)
			}

			if conversionOk {
				var assertOk bool
				primitiveValue, assertOk = concreteValue.(ListItemPrimitiveType)
				if !assertOk {
					return fmt.Errorf("internal error: failed to assign converted value (%T) back to generic type ListItemPrimitiveType (%s) at index %d (path %s)", concreteValue, expectedType, i, path)
				}
			} else {
				return fmt.Errorf("element type error at index %d (path %s): expected %s, got %T, and conversion failed", i, path, expectedType, v)
			}
		}

		// --- Reflection Logic (remains the same) ---
		var newObj NestedObjectType
		newObjReflect := reflect.ValueOf(&newObj).Elem()

		idField := newObjReflect.FieldByName(target_field)
		if !idField.IsValid() {
			return fmt.Errorf("reflection error: Field '%s' not found in type %T for path %s", target_field, newObj, path)
		}
		if !idField.CanSet() {
			return fmt.Errorf("reflection error: Field '%s' cannot be set in type %T (it might be unexported) for path %s", target_field, newObj, path)
		}

		convertedValueReflect := reflect.ValueOf(primitiveValue)

		if !convertedValueReflect.Type().AssignableTo(idField.Type()) {
			return fmt.Errorf("type mismatch error: cannot assign value of type %s to field '%s' (type %s) in %T for path %s", convertedValueReflect.Type(), target_field, idField.Type(), newObj, path)
		}

		idField.Set(convertedValueReflect)
		outList = append(outList, newObj)
	}

	*home = outList
	return nil
}
