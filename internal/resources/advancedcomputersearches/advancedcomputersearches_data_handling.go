// advancedcomputersearches_data_handling.go
package advancedcomputersearches

import (
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressDisplayFieldsDiff is a DiffSuppressFunc used in the Terraform schema to determine
// whether changes to the `display_fields` attribute should be ignored. This function is necessary
// because the Jamf Pro server may return display fields in a different order than they were sent,
// leading to a perceived state change in Terraform. The function compares the old and new states
// of the `display_fields` and suppresses the diff if the only change is in the order of fields.
//
// Parameters:
// k: The key of the schema element being compared. This is not used in the function.
// old: The previous state of the attribute as a comma-separated string.
// new: The new state of the attribute as a comma-separated string.
// d: The full resource data, which can be used to access other attributes. This is not used in the function.
//
// Returns:
// A boolean indicating whether the diff should be suppressed. If it returns true, Terraform ignores
// the detected changes in the `display_fields` attribute.
func suppressDisplayFieldsDiff(k, old, new string, d *schema.ResourceData) bool {
	// Extract the display field names from the old and new states
	oldFields := strings.Split(old, ",")
	newFields := strings.Split(new, ",")

	// Sort both slices for a consistent order
	sort.Strings(oldFields)
	sort.Strings(newFields)

	// Compare sorted slices
	return reflect.DeepEqual(oldFields, newFields)
}
