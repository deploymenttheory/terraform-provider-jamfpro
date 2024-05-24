package policies

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO this is copied from config profiles just to make this work - it'll have a centralised home
func GetAttrsListFromHCLForPointers[NestedObjectType any, ListItemPrimitiveType any](path string, target_field string, d *schema.ResourceData, home *[]NestedObjectType) (err error) {
	getAttr, ok := d.GetOk(path)
	log.Println("START")
	log.Println(getAttr)
	log.Println(ok)

	if len(getAttr.([]interface{})) == 0 {
		log.Println("MARKER-1")
		return nil
	}

	log.Println("MARKER-2")

	if ok {
		log.Println("MARKER-3")
		*home = []NestedObjectType{}
		log.Println("MARKER-4")
		outList := make([]NestedObjectType, 0)
		log.Println("MARKER-5")
		for _, v := range getAttr.([]interface{}) {
			log.Println("MARKER-6")
			var newObj NestedObjectType
			newObjReflect := reflect.ValueOf(&newObj).Elem()
			idField := newObjReflect.FieldByName(target_field)
			if idField.IsValid() && idField.CanSet() {
				idField.Set(reflect.ValueOf(v.(ListItemPrimitiveType)))
			} else {
				return fmt.Errorf("error cannot set field line 695") // TODO write this error
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
	return fmt.Errorf("no path found/no scoped items at %v", path)
}
