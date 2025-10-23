package constructors

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Test structures that will be used in our tests
type TestStruct struct {
	ID   int
	Name string
}

type TestStructWithUnexportedField struct {
	id int
}

func TestMapSetToStructs(t *testing.T) {
	testSchema := &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}

	t.Run("with valid int values", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": testSchema,
		}, map[string]any{
			"test_ids": []any{1, 2, 3},
		})

		var result []TestStruct

		err := MapSetToStructs[TestStruct, int]("test_ids", "ID", rd, &result)

		assert.NoError(t, err)
		assert.Len(t, result, 3)

		// Since sets are unordered, we need to check that all expected values are present
		// without assuming any specific order
		ids := []int{result[0].ID, result[1].ID, result[2].ID}
		assert.ElementsMatch(t, []int{1, 2, 3}, ids)
	})

	t.Run("with empty set", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": testSchema,
		}, map[string]any{
			"test_ids": []any{},
		})

		var result []TestStruct
		err := MapSetToStructs[TestStruct, int]("test_ids", "ID", rd, &result)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("with non-existent path", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"other_field": testSchema,
		}, map[string]any{})

		var result []TestStruct

		err := MapSetToStructs[TestStruct, int]("test_ids", "ID", rd, &result)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("with nil values in set", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": testSchema,
		}, map[string]any{
			"test_ids": []any{1, nil, 3},
		})

		var result []TestStruct

		err := MapSetToStructs[TestStruct, int]("test_ids", "ID", rd, &result)

		assert.NoError(t, err)

		// Check that all values are present in the result
		// In this case, it appears that nil is being converted to 0 by the test framework
		// So we need to check for 0, 1, and 3
		ids := make([]int, len(result))
		for i, r := range result {
			ids[i] = r.ID
		}
		assert.ElementsMatch(t, []int{0, 1, 3}, ids, "Expected IDs to include 0 (from nil), 1, and 3")
	})

	t.Run("with incorrect type values", func(t *testing.T) {
		stringSchema := &schema.Schema{
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}

		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": stringSchema,
		}, map[string]any{
			"test_ids": []any{"1", "2", "3"},
		})

		var result []TestStruct

		err := MapSetToStructs[TestStruct, int]("test_ids", "ID", rd, &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect type")
	})

	t.Run("with non-existent field in struct", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": testSchema,
		}, map[string]any{
			"test_ids": []any{1, 2, 3},
		})

		var result []TestStruct

		assert.Panics(t, func() {
			_ = MapSetToStructs[TestStruct, int]("test_ids", "NonExistentField", rd, &result)
		})
	})

	t.Run("with unexported field in struct", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_ids": testSchema,
		}, map[string]any{
			"test_ids": []any{1, 2, 3},
		})

		var result []TestStructWithUnexportedField

		assert.Panics(t, func() {
			_ = MapSetToStructs[TestStructWithUnexportedField, int]("test_ids", "id", rd, &result)
		})
	})
}

// This test ensures MapSetToStructs properly handles structs with different field types
func TestMapSetToStructs_DifferentTypes(t *testing.T) {
	t.Run("with string values", func(t *testing.T) {
		stringSchema := &schema.Schema{
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}

		rd := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"test_names": stringSchema,
		}, map[string]any{
			"test_names": []any{"one", "two", "three"},
		})

		var result []TestStruct

		err := MapSetToStructs[TestStruct, string]("test_names", "Name", rd, &result)

		assert.NoError(t, err)
		assert.Len(t, result, 3)

		names := make([]string, len(result))
		for i, r := range result {
			names[i] = r.Name
		}
		assert.ElementsMatch(t, []string{"one", "two", "three"}, names)
	})
}
