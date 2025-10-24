package computer_extension_attribute

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// diffSuppressPopupMenuChoices is a custom diff suppression function for the popup_menu_choices attribute.
// This attribute looks for matched values and ignores the order returned by the server.
func diffSuppressPopupMenuChoices(k, old, new string, d *schema.ResourceData) bool {
	// Get the old and new values as interfaces
	oldRaw, newRaw := d.GetChange("popup_menu_choices")

	// Convert interfaces to []string
	oldList := make([]string, len(oldRaw.([]any)))
	newList := make([]string, len(newRaw.([]any)))

	for i, v := range oldRaw.([]any) {
		oldList[i] = v.(string)
	}
	for i, v := range newRaw.([]any) {
		newList[i] = v.(string)
	}

	// If lengths don't match, there's definitely a difference
	if len(oldList) != len(newList) {
		return false
	}

	// Create maps to track occurrences of each value
	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	// Count occurrences in both lists
	for _, v := range oldList {
		oldMap[v]++
	}
	for _, v := range newList {
		newMap[v]++
	}

	// Compare maps
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
