package mobiledeviceextensionattributes

import (
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressPopupChoicesDiff is a diff suppression function that ignores the order of popup choices
func suppressPopupChoicesDiff(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("input_type.0.popup_choices")
	if o == nil || n == nil {
		return false
	}
	oldChoices := o.([]interface{})
	newChoices := n.([]interface{})

	if len(oldChoices) != len(newChoices) {
		return false
	}

	oldStrings := make([]string, len(oldChoices))
	newStrings := make([]string, len(newChoices))

	for i, v := range oldChoices {
		oldStrings[i] = v.(string)
	}
	for i, v := range newChoices {
		newStrings[i] = v.(string)
	}

	sort.Strings(oldStrings)
	sort.Strings(newStrings)

	return reflect.DeepEqual(oldStrings, newStrings)
}
