package plist

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortPlistKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name: "Simple flat dictionary",
			input: map[string]interface{}{
				"c": 3,
				"a": 1,
				"b": 2,
			},
			want: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			name: "Nested dictionary",
			input: map[string]interface{}{
				"z": map[string]interface{}{
					"inner2": 2,
					"inner1": 1,
				},
				"a": 1,
			},
			want: map[string]interface{}{
				"a": 1,
				"z": map[string]interface{}{
					"inner1": 1,
					"inner2": 2,
				},
			},
		},
		{
			name: "Dictionary with string array",
			input: map[string]interface{}{
				"array": []interface{}{"c", "b", "a"},
				"key":   "value",
			},
			want: map[string]interface{}{
				"array": []interface{}{"a", "b", "c"},
				"key":   "value",
			},
		},
		{
			name: "Dictionary with array of dictionaries",
			input: map[string]interface{}{
				"array": []interface{}{
					map[string]interface{}{
						"b": 2,
						"a": 1,
					},
					map[string]interface{}{
						"d": 4,
						"c": 3,
					},
				},
			},
			want: map[string]interface{}{
				"array": []interface{}{
					map[string]interface{}{
						"a": 1,
						"b": 2,
					},
					map[string]interface{}{
						"c": 3,
						"d": 4,
					},
				},
			},
		},
		{
			name: "Mixed array types (should not sort non-string arrays)",
			input: map[string]interface{}{
				"mixedArray": []interface{}{1, "b", 3.14, "a"},
				"key":        "value",
			},
			want: map[string]interface{}{
				"mixedArray": []interface{}{1, "b", 3.14, "a"},
				"key":        "value",
			},
		},
		{
			name: "Empty structures",
			input: map[string]interface{}{
				"emptyMap":   map[string]interface{}{},
				"emptyArray": []interface{}{},
				"key":        "value",
			},
			want: map[string]interface{}{
				"emptyMap":   map[string]interface{}{},
				"emptyArray": []interface{}{},
				"key":        "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortPlistKeys(tt.input)

			// Deep equality check
			assert.Equal(t, tt.want, got)

			// Helper function to get sorted keys from a map
			getKeys := func(m map[string]interface{}) []string {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				return keys
			}

			// Helper function to verify if keys are in order
			areKeysSorted := func(m map[string]interface{}) bool {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				return sort.SliceIsSorted(keys, func(i, j int) bool {
					return keys[i] < keys[j]
				})
			}

			// Verify all nested structures recursively
			var verifyNestedStructure func(v interface{}) bool
			verifyNestedStructure = func(v interface{}) bool {
				switch x := v.(type) {
				case map[string]interface{}:
					if !areKeysSorted(x) {
						return false
					}
					for _, val := range x {
						if !verifyNestedStructure(val) {
							return false
						}
					}
				case []interface{}:
					for _, item := range x {
						if m, ok := item.(map[string]interface{}); ok {
							if !verifyNestedStructure(m) {
								return false
							}
						}
					}
				}
				return true
			}

			// Verify the entire structure
			if !verifyNestedStructure(got) {
				t.Error("structure contains unsorted keys")
			}

			// Additional specific checks for particular test cases
			if tt.name == "Dictionary with array of dictionaries" {
				if arr, ok := got["array"].([]interface{}); ok {
					for _, item := range arr {
						if dict, ok := item.(map[string]interface{}); ok {
							keys := getKeys(dict)
							assert.True(t, sort.SliceIsSorted(keys, func(i, j int) bool {
								return keys[i] < keys[j]
							}), "dictionary in array should have sorted keys")
						}
					}
				}
			}
		})
	}
}
