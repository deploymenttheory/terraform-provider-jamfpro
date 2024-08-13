// common/configurationprofiles/plist/state.go
// contains the functions to process configuration profiles for terraform state.
package plist

import (
	"sort"
)

// Helper function to reorder plist keys based on server logic for staing purposes
// Used by plist generator resource rther than plist.
// This function reorders the configuration keys based on observed server behavior.
// The server reorders the keys primarily based on their types and secondarily based on their names.
// The following steps outline the observations and logic implemented:
//
// Steps to Determine the Reordering Rule:
//
// 1. Alphabetical Order Check:
//   - Initially, it appeared that the server ordered the keys alphabetically by their names.
//   - However, the exact sequence did not match a simple alphabetical order.
//
// 2. Reordering Pattern:
//   - Upon further examination, the new order follows a pattern where the keys are sorted by type first:
//   - Boolean keys (`true` or `false`) are prioritized.
//   - Then, numeric keys are placed.
//   - Finally, string keys are ordered.
//
// Conclusion:
// - The keys are reordered based on their types first and then alphabetically within each type.
// - Here's a step-by-step breakdown of the observed reordering pattern:
//  1. Boolean keys (true or false) come first, ordered alphabetically by their key names.
//  2. Numeric keys come next, ordered by their values (integers).
//  3. String keys come last, ordered alphabetically by their key names.
//
// This function implements this reordering logic by sorting the configurations slice.
//
// Parameters:
// - configurations: A slice of maps, each containing a "key" and a "value".
//
// Returns:
// - The reordered slice of configurations.
func reorderConfigurationKeys(configurations []interface{}) []interface{} {
	// Sort configurations by key types and names
	sort.Slice(configurations, func(i, j int) bool {
		// Extract values
		val1 := configurations[i].(map[string]interface{})["value"]
		val2 := configurations[j].(map[string]interface{})["value"]

		// Determine type order: bool < int < string
		switch v1 := val1.(type) {
		case bool:
			if _, ok := val2.(bool); ok {
				// Both are bool, sort by key name
				key1 := configurations[i].(map[string]interface{})["key"].(string)
				key2 := configurations[j].(map[string]interface{})["key"].(string)
				return key1 < key2
			}
			// bool comes before int and string
			return true
		case int:
			if _, ok := val2.(bool); ok {
				// int comes after bool
				return false
			}
			if _, ok := val2.(int); ok {
				// Both are int, sort by value
				return v1 < val2.(int)
			}
			// int comes before string
			return true
		case string:
			if _, ok := val2.(bool); ok {
				// string comes after bool
				return false
			}
			if _, ok := val2.(int); ok {
				// string comes after int
				return false
			}
			// Both are string, sort by key name
			key1 := configurations[i].(map[string]interface{})["key"].(string)
			key2 := configurations[j].(map[string]interface{})["key"].(string)
			return key1 < key2
		default:
			return false
		}
	})

	return configurations
}
