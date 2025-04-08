package policies

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PopulateStructSliceFromSetField transforms a Terraform schema.Set containing simple
// primitive values (like IDs or names) into a standard Go slice of structs.
// Its primary use case is converting configuration blocks like:
//
//	 attribute_ids = [101, 102, 103]
//
//	into a Go slice suitable for API interaction, such as:
//
//	 []SomeStruct{ {ID: 101}, {ID: 102}, {ID: 103} }
//
// For each value of type `ListItemPrimitiveType` found in the `schema.Set` at the
// specified `path`, it creates a new instance of `NestedObjectType` and assigns
// the primitive value to the struct's field specified by `fieldName`.
//
// The function populates the slice pointed to by `outputSlice`. It guarantees
// idempotency by *always* initializing or clearing the target slice at the start.
// If the source path is not found, the attribute is nil, the set is empty, or an
// error occurs, `outputSlice` will point to a valid, empty slice.
func PopulateStructSliceFromSetField[NestedObjectType any, ListItemPrimitiveType comparable](path string, fieldName string, d *schema.ResourceData, outputSlice *[]NestedObjectType) error {
	*outputSlice = []NestedObjectType{}

	getAttr, ok := d.GetOk(path)
	if !ok || getAttr == nil {
		return nil
	}

	attrSet, isSet := getAttr.(*schema.Set)
	if !isSet {
		return fmt.Errorf("internal error: attribute at path %s was expected to be *schema.Set, but got %T", path, getAttr)
	}

	itemsList := attrSet.List()
	if len(itemsList) == 0 {
		return nil
	}

	result := make([]NestedObjectType, 0, len(itemsList))

	for i, v := range itemsList {
		if v == nil {
			continue
		}

		value, ok := v.(ListItemPrimitiveType)
		if !ok {
			expectedTypeName := reflect.TypeOf((*ListItemPrimitiveType)(nil)).Elem().Name()
			actualTypeName := reflect.TypeOf(v).Name()
			if actualTypeName == "" {
				actualTypeName = fmt.Sprintf("%T", v)
			}
			return fmt.Errorf("type assertion error at index %d (path %s): expected %s, got %s",
				i, path, expectedTypeName, actualTypeName)
		}

		var obj NestedObjectType
		objVal := reflect.ValueOf(&obj).Elem()
		field := objVal.FieldByName(fieldName)

		if !field.IsValid() {
			return fmt.Errorf("field '%s' not found in type %T", fieldName, obj)
		}

		if !field.CanSet() {
			return fmt.Errorf("field '%s' cannot be set in type %T (unexported?)", fieldName, obj)
		}

		valueRefl := reflect.ValueOf(value)
		if !valueRefl.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("cannot assign type %T (value: '%v') to field '%s' (type %s) in %T",
				value, value, fieldName, field.Type(), obj)
		}

		field.Set(valueRefl)
		result = append(result, obj)
	}

	*outputSlice = result
	return nil
}
