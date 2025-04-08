package policies

import (
	"strings" // Import strings package
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestStruct is used for testing PopulateStructSliceFromSetField
type TestStruct struct {
	ID              int
	Name            string
	unexportedValue bool
}

// Test schema definition helper
func testSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"int_set": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Set:      schema.HashInt,
		},
		"string_set": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
		"wrong_type_attr": {
			Type:     schema.TypeList, // Not a Set
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"set_with_wrong_items": {
			Type:     schema.TypeSet,
			Optional: true,
			// Elem defined as int, but we'll put strings in the test data
			Elem: &schema.Schema{Type: schema.TypeInt},
			Set:  schema.HashInt,
		},
	}
}

// --- Test Cases ---

func TestPopulateStructSliceFromSetField(t *testing.T) {
	testSchemaDef := testSchema()

	t.Run("Success_IntField", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{101, 102, 103},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := []TestStruct{
			{ID: 101},
			{ID: 102},
			{ID: 103},
		}

		if len(output) != len(expected) {
			t.Fatalf("Expected slice length %d, got %d", len(expected), len(output))
		}

		// Check values (order in Set isn't guaranteed, so need to check existence)
		expectedMap := make(map[int]bool)
		for _, item := range expected {
			expectedMap[item.ID] = true
		}
		foundMap := make(map[int]bool)
		for _, item := range output {
			if !expectedMap[item.ID] {
				t.Errorf("Found unexpected item ID %d in output slice", item.ID)
			}
			if foundMap[item.ID] {
				t.Errorf("Duplicate item ID %d found in output slice", item.ID) // Sets shouldn't have duplicates
			}
			foundMap[item.ID] = true
		}
		if len(foundMap) != len(expectedMap) {
			t.Errorf("Mismatch in number of unique items found. Expected %d, got %d", len(expectedMap), len(foundMap))
		}
	})

	t.Run("Success_StringField", func(t *testing.T) {
		rawData := map[string]interface{}{
			"string_set": []interface{}{"alpha", "beta", "gamma"},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		output := make([]TestStruct, 0)
		err := PopulateStructSliceFromSetField[TestStruct, string]("string_set", "Name", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := []TestStruct{
			{Name: "alpha"},
			{Name: "beta"},
			{Name: "gamma"},
		}

		if len(output) != len(expected) {
			t.Fatalf("Expected slice length %d, got %d", len(expected), len(output))
		}

		expectedMap := make(map[string]bool)
		for _, item := range expected {
			expectedMap[item.Name] = true
		}
		foundMap := make(map[string]bool)
		for _, item := range output {
			if !expectedMap[item.Name] {
				t.Errorf("Found unexpected item Name '%s' in output slice", item.Name)
			}
			if foundMap[item.Name] {
				t.Errorf("Duplicate item Name '%s' found in output slice", item.Name)
			}
			foundMap[item.Name] = true
		}
		if len(foundMap) != len(expectedMap) {
			t.Errorf("Mismatch in number of unique items found. Expected %d, got %d", len(expectedMap), len(foundMap))
		}
	})

	t.Run("EmptySet", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{}, // Empty set
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error for empty set: %v", err)
		}
		if output == nil {
			t.Fatalf("Expected non-nil empty slice, got nil")
		}
		if len(output) != 0 {
			t.Fatalf("Expected slice length 0, got %d", len(output))
		}
	})

	t.Run("EmptySet_WithExistingSlice", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		output := []TestStruct{{ID: 999}}
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error for empty set with existing slice: %v", err)
		}
		if output == nil {
			t.Fatalf("Expected non-nil empty slice, got nil")
		}
		if len(output) != 0 {
			t.Errorf("Expected slice length 0, got %d", len(output))
		}
	})

	t.Run("PathNotFound", func(t *testing.T) {
		rawData := map[string]interface{}{} // Path "int_set" does not exist
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct = []TestStruct{{ID: 888}} // Start with non-nil slice
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error for path not found: %v", err)
		}
		if output == nil {
			t.Fatalf("Expected non-nil empty slice, got nil")
		}
		if len(output) != 0 {
			t.Fatalf("Expected slice length 0, got %d", len(output))
		}
	})

	t.Run("AttributeNotASet", func(t *testing.T) {
		rawData := map[string]interface{}{
			"wrong_type_attr": []interface{}{"a", "b"},
		}

		var d *schema.ResourceData
		var errPanic interface{}
		func() {
			defer func() {
				if r := recover(); r != nil {
					errPanic = r
				}
			}()
			d = schema.TestResourceDataRaw(t, testSchemaDef, rawData)
		}()

		if errPanic != nil {
			t.Fatalf("Panic occurred during TestResourceDataRaw setup: %v", errPanic)
		}
		if d == nil {
			t.Fatal("ResourceData is nil after setup")
		}

		var output []TestStruct
		err := PopulateStructSliceFromSetField[TestStruct, string]("wrong_type_attr", "Name", d, &output)

		if err == nil {
			t.Fatalf("Expected error for wrong attribute type, got nil")
		}

		expectedErrMsgPrefix := "internal error: attribute at path wrong_type_attr was expected to be *schema.Set"
		if !strings.HasPrefix(err.Error(), expectedErrMsgPrefix) {
			t.Errorf("Expected error message starting with '%s', got '%s'", expectedErrMsgPrefix, err.Error())
		}

		if len(output) != 0 {
			t.Errorf("Expected output slice to be empty on error, got length %d", len(output))
		}
	})

	t.Run("IncorrectItemTypeInSet", func(t *testing.T) {
		localSchemaDef := map[string]*schema.Schema{
			"int_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		}

		rawData := map[string]interface{}{
			"int_set": []interface{}{"not-an-int", "123"},
		}

		d := schema.TestResourceDataRaw(t, localSchemaDef, rawData)

		var output []TestStruct
		// Now, call the function under test, explicitly expecting INT as the primitive type,
		// even though the underlying set actually contains strings.
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output) // Expecting int!

		// We expect an error because the type assertion v.(int) inside the function will fail
		// when it encounters the string "not-an-int" or "123".
		if err == nil {
			t.Fatalf("Expected error for incorrect item type assertion, got nil. Output: %v", output)
		}

		expectedSubString1 := "type assertion error"
		expectedSubString2 := "expected int, got string"
		if !strings.Contains(err.Error(), expectedSubString1) || !strings.Contains(err.Error(), expectedSubString2) {
			t.Errorf("Expected error message containing '%s' and '%s', got '%s'", expectedSubString1, expectedSubString2, err.Error())
		}
		if len(output) != 0 {
			t.Errorf("Expected output slice to be empty on error, got length %d", len(output))
		}
	})

	t.Run("FieldNotFound", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{101},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "NonExistentField", d, &output)

		if err == nil {
			t.Fatalf("Expected error for non-existent field, got nil")
		}
		expectedErrMsg := "field 'NonExistentField' not found in type policies.TestStruct"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
		if len(output) != 0 {
			t.Errorf("Expected output slice to be empty on error, got length %d", len(output))
		}
	})

	t.Run("FieldNotSettable", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{123},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct

		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "unexportedValue", d, &output)

		if err == nil {
			t.Fatalf("Expected error for unexported field, got nil")
		}

		expectedErrMsg := "field 'unexportedValue' cannot be set in type policies.TestStruct (unexported?)"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
		if len(output) != 0 {
			t.Errorf("Expected output slice to be empty on error, got length %d", len(output))
		}
	})

	t.Run("TypeMismatchFieldAssignment", func(t *testing.T) {
		rawData := map[string]interface{}{
			"string_set": []interface{}{"not-an-int"},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct
		// Trying to assign a string from the set to the 'ID' field (which is int)
		err := PopulateStructSliceFromSetField[TestStruct, string]("string_set", "ID", d, &output)

		if err == nil {
			t.Fatalf("Expected error for field type mismatch, got nil")
		}
		// Error message check (using Contains for less brittle check)
		expectedSubString := "cannot assign type string"
		expectedSubString2 := "to field 'ID' (type int)"
		if !strings.Contains(err.Error(), expectedSubString) || !strings.Contains(err.Error(), expectedSubString2) {
			t.Errorf("Expected error message containing '%s' and '%s', got '%s'", expectedSubString, expectedSubString2, err.Error())
		}
		if len(output) != 0 {
			t.Errorf("Expected output slice to be empty on error, got length %d", len(output))
		}
	})

	t.Run("NilItemInSetHandled", func(t *testing.T) {
		rawData := map[string]interface{}{
			"int_set": []interface{}{101, 103},
		}
		d := schema.TestResourceDataRaw(t, testSchemaDef, rawData)

		var output []TestStruct
		err := PopulateStructSliceFromSetField[TestStruct, int]("int_set", "ID", d, &output)

		if err != nil {
			t.Fatalf("Unexpected error during nil check simulation: %v", err)
		}

		if len(output) != 2 {
			t.Errorf("Expected 2 items after skipping hypothetical nils, got %d", len(output))
		}
		expectedMap := map[int]bool{101: true, 103: true}
		foundMap := make(map[int]bool)
		for _, item := range output {
			if !expectedMap[item.ID] {
				t.Errorf("Found unexpected item ID %d in output slice during nil check simulation", item.ID)
			}
			foundMap[item.ID] = true
		}
		if len(foundMap) != len(expectedMap) {
			t.Errorf("Mismatch in number of unique items found. Expected %d, got %d", len(expectedMap), len(foundMap))
		}
	})
}
