// endpoints/common/sharedschemas/helpers.go
package sharedschemas

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExtractNestedObjectsFromSchema is a helper function to extract a list of nested objects from HCL configuration
/* Example usage:
err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, &out.Scope.Exclusions.IBeacons)
if err != nil {
	return nil, err
}
*/
func ExtractNestedObjectsFromSchema[NestedObjectType any, ListItemPrimitiveType any](path string, targetField string, d *schema.ResourceData, home *[]NestedObjectType) error {
	getAttr, ok := d.GetOk(path)
	if !ok {
		return nil
	}

	attrList, ok := getAttr.([]interface{})
	if !ok {
		return fmt.Errorf("failed to cast %s to []interface{}", path)
	}

	if len(attrList) == 0 {
		return nil
	}

	outList := make([]NestedObjectType, 0)
	for _, v := range attrList {
		if v == nil {
			continue
		}

		var newObj NestedObjectType
		newObjReflect := reflect.ValueOf(&newObj).Elem()
		idField := newObjReflect.FieldByName(targetField)

		if idField.IsValid() && idField.CanSet() {
			idField.Set(reflect.ValueOf(v.(ListItemPrimitiveType)))
		} else {
			return fmt.Errorf("error cannot set field %s", targetField)
		}

		outList = append(outList, newObj)
	}

	if len(outList) > 0 {
		*home = outList
	} else {
		log.Println("list is empty")
	}

	return nil
}
