// macosconfigurationprofiles_helpers.go
package macosconfigurationprofiles

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExtractNestedObjectsFromSchema is a helper function to extract a list of nested objects from HCL configuration
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

// FixDuplicateNotificationKey handles the double key issue in the notification field of the self_service block.
/*
<self_service>
        <self_service_display_name>WiFi Test</self_service_display_name>
        <install_button_text>Install</install_button_text>
        <self_service_description>null</self_service_description>
        <force_users_to_view_description>false</force_users_to_view_description>
        <security>
            <removal_disallowed>Never</removal_disallowed>
        </security>
        <self_service_icon/>
        <feature_on_main_page>false</feature_on_main_page>
        <self_service_categories/>
        <notification>false</notification>				<-- This is the issue
        <notification>Self Service</notification>  <-- This is the issue
        <notification_subject/>
        <notification_message/>
    </self_service>
*/
func FixDuplicateNotificationKey(resp *jamfpro.ResourceMacOSConfigurationProfile) (bool, error) {
	for _, k := range resp.SelfService.Notification {
		strValue := fmt.Sprintf("%v", k)
		if strValue == "true" || strValue == "false" {
			correctNotifValue, err := strconv.ParseBool(strValue)
			if err != nil {
				return false, err
			}
			return correctNotifValue, nil
		} else {
			log.Printf("Ignoring non-boolean notification value: %s", strValue)
		}
	}
	// Return default value if no valid boolean value is found
	return false, nil
}

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
