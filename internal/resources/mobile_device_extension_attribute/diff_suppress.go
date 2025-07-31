package mobile_device_extension_attribute

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// diffSuppressPopupMenuChoices is a custom diff suppression function for the popup_menu_choices attribute.
// This attribute looks for matched values and ignores the order returned by the server.
func diffSuppressPopupMenuChoices(k, old, new string, d *schema.ResourceData) bool {
	oldRaw, newRaw := d.GetChange("popup_menu_choices")

	oldList := make([]string, len(oldRaw.([]interface{})))
	newList := make([]string, len(newRaw.([]interface{})))

	for i, v := range oldRaw.([]interface{}) {
		oldList[i] = v.(string)
	}
	for i, v := range newRaw.([]interface{}) {
		newList[i] = v.(string)
	}

	if len(oldList) != len(newList) {
		return false
	}

	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	for _, v := range oldList {
		oldMap[v]++
	}
	for _, v := range newList {
		newMap[v]++
	}

	if len(oldMap) != len(newMap) {
		return false
	}

	for k, v := range oldMap {
		if newMap[k] != v {
			return false
		}
	}

	return true
}
