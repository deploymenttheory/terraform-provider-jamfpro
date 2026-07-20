package plist

import (
	"encoding/xml"
	"errors"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortPlistKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]any
		want  map[string]any
	}{
		{
			name: "Simple flat dictionary",
			input: map[string]any{
				"c": 3,
				"a": 1,
				"b": 2,
			},
			want: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			name: "Nested dictionary",
			input: map[string]any{
				"z": map[string]any{
					"inner2": 2,
					"inner1": 1,
				},
				"a": 1,
			},
			want: map[string]any{
				"a": 1,
				"z": map[string]any{
					"inner1": 1,
					"inner2": 2,
				},
			},
		},
		{
			name: "Dictionary with string array",
			input: map[string]any{
				"array": []any{"c", "b", "a"},
				"key":   "value",
			},
			want: map[string]any{
				"array": []any{"a", "b", "c"},
				"key":   "value",
			},
		},
		{
			name: "Dictionary with array of dictionaries",
			input: map[string]any{
				"array": []any{
					map[string]any{"b": 2, "a": 1},
					map[string]any{"d": 4, "c": 3},
				},
			},
			want: map[string]any{
				"array": []any{
					map[string]any{"a": 1, "b": 2},
					map[string]any{"c": 3, "d": 4},
				},
			},
		},
		{
			name: "Mixed array types (should not sort non-string arrays)",
			input: map[string]any{
				"mixedArray": []any{1, "b", 3.14, "a"},
				"key":        "value",
			},
			want: map[string]any{
				"mixedArray": []any{1, "b", 3.14, "a"},
				"key":        "value",
			},
		},
		{
			name: "Empty structures",
			input: map[string]any{
				"emptyMap":   map[string]any{},
				"emptyArray": []any{},
				"key":        "value",
			},
			want: map[string]any{
				"emptyMap":   map[string]any{},
				"emptyArray": []any{},
				"key":        "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortPlistKeys(tt.input)
			assert.Equal(t, tt.want, got, "unexpected sorted output")
		})
	}
}

// assertDictKeysSorted walks an encoded plist and asserts that every <dict>
// lists its own keys in sorted order. Sortedness is per dictionary, not global:
// in document order a nested dict's keys are interleaved between its parent's,
// so a flat scan of the whole document is not expected to be sorted.
func assertDictKeysSorted(t *testing.T, encoded string) {
	t.Helper()

	dec := xml.NewDecoder(strings.NewReader(encoded))
	var stack [][]string

	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "encoded plist must be well-formed XML")

		switch t2 := tok.(type) {
		case xml.StartElement:
			switch t2.Name.Local {
			case "dict":
				stack = append(stack, nil)
			case "key":
				var name string
				require.NoError(t, dec.DecodeElement(&name, &t2))
				require.NotEmpty(t, stack, "<key> outside of a <dict>")
				stack[len(stack)-1] = append(stack[len(stack)-1], name)
			}
		case xml.EndElement:
			if t2.Name.Local == "dict" {
				keys := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				assert.True(t, sort.StringsAreSorted(keys),
					"dict keys must be sorted, got %v", keys)
			}
		}
	}
	assert.Empty(t, stack, "unbalanced <dict> elements")
}

// SortPlistKeys exists so that two structurally equal profiles serialise to the
// same bytes, which is what makes payload diff suppression work. That ordering
// is not observable on the returned map — Go randomises map iteration order, so
// asserting over `for k := range m` is a coin flip rather than a test. It
// becomes observable only once the map is encoded, so assert it there.
func TestSortPlistKeysEncodesDeterministically(t *testing.T) {
	input := map[string]any{
		"zebra": "last",
		"apple": "first",
		"nested": map[string]any{
			"inner_z": 2,
			"inner_a": 1,
		},
		"array": []any{
			map[string]any{"d": 4, "c": 3},
			map[string]any{"b": 2, "a": 1},
		},
	}

	first, err := EncodePlist(SortPlistKeys(input))
	require.NoError(t, err)
	assertDictKeysSorted(t, first)

	// Repeat: one pass can look correct by luck if the maps happen to be
	// iterated favourably. Byte-identical output across runs is the guarantee
	// diff suppression actually relies on.
	for i := range 20 {
		again, err := EncodePlist(SortPlistKeys(input))
		require.NoError(t, err)
		require.Equal(t, first, again, "encoding must be identical on run %d", i+1)
	}
}
