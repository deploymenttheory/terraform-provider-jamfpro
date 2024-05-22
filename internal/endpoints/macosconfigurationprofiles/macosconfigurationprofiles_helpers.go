package macosconfigurationprofiles

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO rename this func and put it somewhere else
func GetAttrsListFromHCL[NestedObjectType any, ListItemPrimitiveType any](path string, target_field string, d *schema.ResourceData, home *[]NestedObjectType) (err error) {
	getAttr, ok := d.GetOk(path)

	if len(getAttr.([]interface{})) == 0 {
		return nil
	}

	if ok {
		outList := make([]NestedObjectType, 0)
		for _, v := range getAttr.([]interface{}) {
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

// TODO rename this func and put it somewhere else too
// func FixStupidDoubleKey(resp *jamfpro.ResourceMacOSConfigurationProfile, home *[]map[string]interface{}) error {
// 	var err error
// 	var correctNotifValue bool
// 	for _, k := range resp.SelfService.Notification {
// 		if k == "true" || k == "false" {
// 			correctNotifValue, err = strconv.ParseBool(k)
// 			if err != nil {
// 				return err
// 			}
// 			(*home)[0]["notification"] = correctNotifValue
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("failed to parse value %+v", resp.SelfService)
// }

// TODO Make this work later

// func GetListOfIdsFromResp[T any](targetItem []T, targetKey string) ([]int, error) {
// 	if len(targetItem) == 0 {
// 		return nil, nil
// 	} else {
// 		var out []int
// 		for k, v := range targetItem.FieldByName(targetKey) {
// 			out = append(out, v)
// 		}
// 		return out, nil
// 	}
// }
