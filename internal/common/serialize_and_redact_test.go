package common

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

// Test structs for various scenarios
type TestStruct struct {
	Name     string         `xml:"Name"`
	Password string         `xml:"Password"`
	Age      int            `xml:"Age"`
	Nested   *NestedStruct  `xml:"Nested"`
	Direct   NestedStruct   `xml:"Direct"`
}

type NestedStruct struct {
	Secret     string            `xml:"Secret"`
	Token      string            `xml:"Token"`
	Value      int               `xml:"Value"`
	DeepNest   *DeepNestedStruct `xml:"DeepNest"`
	DeepDirect DeepNestedStruct  `xml:"DeepDirect"`
}

type DeepNestedStruct struct {
	HiddenField string `xml:"HiddenField"`
	PublicField string `xml:"PublicField"`
	Number      int    `xml:"Number"`
}

func TestNavigateToField(t *testing.T) {
	// Setup test data
	testData := &TestStruct{
		Name:     "John",
		Password: "secret123",
		Age:      30,
		Nested: &NestedStruct{
			Secret: "nested-secret",
			Token:  "auth-token",
			Value:  42,
			DeepNest: &DeepNestedStruct{
				HiddenField: "deep-secret",
				PublicField: "public-info",
				Number:      100,
			},
			DeepDirect: DeepNestedStruct{
				HiddenField: "direct-deep-secret",
				PublicField: "direct-public-info",
				Number:      200,
			},
		},
		Direct: NestedStruct{
			Secret: "direct-secret",
			Token:  "direct-token",
			Value:  99,
		},
	}

	tests := []struct {
		name        string
		fieldPath   string
		expectFound bool
		expectValue interface{}
	}{
		{
			name:        "Simple field access",
			fieldPath:   "Name",
			expectFound: true,
			expectValue: "John",
		},
		{
			name:        "Simple field access - Password",
			fieldPath:   "Password",
			expectFound: true,
			expectValue: "secret123",
		},
		{
			name:        "Nested pointer field",
			fieldPath:   "Nested.Secret",
			expectFound: true,
			expectValue: "nested-secret",
		},
		{
			name:        "Nested pointer field - Token",
			fieldPath:   "Nested.Token",
			expectFound: true,
			expectValue: "auth-token",
		},
		{
			name:        "Deep nested pointer field",
			fieldPath:   "Nested.DeepNest.HiddenField",
			expectFound: true,
			expectValue: "deep-secret",
		},
		{
			name:        "Deep nested direct field",
			fieldPath:   "Nested.DeepDirect.HiddenField",
			expectFound: true,
			expectValue: "direct-deep-secret",
		},
		{
			name:        "Direct nested field",
			fieldPath:   "Direct.Secret",
			expectFound: true,
			expectValue: "direct-secret",
		},
		{
			name:        "Non-existent field",
			fieldPath:   "NonExistent",
			expectFound: false,
			expectValue: nil,
		},
		{
			name:        "Non-existent nested field",
			fieldPath:   "Nested.NonExistent",
			expectFound: false,
			expectValue: nil,
		},
		{
			name:        "Empty field path",
			fieldPath:   "",
			expectFound: false,
			expectValue: nil,
		},
		{
			name:        "Invalid path - too deep",
			fieldPath:   "Name.Something",
			expectFound: false,
			expectValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(testData).Elem()
			field, found := navigateToField(v, tt.fieldPath)

			if found != tt.expectFound {
				t.Errorf("navigateToField() found = %v, expectFound %v", found, tt.expectFound)
				return
			}

			if tt.expectFound && field.Interface() != tt.expectValue {
				t.Errorf("navigateToField() value = %v, expectValue %v", field.Interface(), tt.expectValue)
			}
		})
	}
}

func TestNavigateToFieldWithNilPointer(t *testing.T) {
	testData := &TestStruct{
		Name:   "John",
		Nested: nil, // nil pointer
	}

	v := reflect.ValueOf(testData).Elem()
	_, found := navigateToField(v, "Nested.Secret")

	if found {
		t.Error("Expected navigateToField to return false for nil pointer navigation")
	}
}

func TestSerializeAndRedactJSON(t *testing.T) {
	testData := &TestStruct{
		Name:     "John",
		Password: "secret123",
		Age:      30,
		Nested: &NestedStruct{
			Secret: "nested-secret",
			Token:  "auth-token",
			Value:  42,
			DeepNest: &DeepNestedStruct{
				HiddenField: "deep-secret",
				PublicField: "public-info",
				Number:      100,
			},
		},
		Direct: NestedStruct{
			Secret: "direct-secret",
			Token:  "direct-token",
			Value:  99,
		},
	}

	redactFields := []string{
		"Password",
		"Nested.Secret",
		"Nested.DeepNest.HiddenField",
		"Direct.Token",
		"Age", // Non-string field to test zeroing
	}

	result, err := SerializeAndRedactJSON(testData, redactFields)
	if err != nil {
		t.Fatalf("SerializeAndRedactJSON() error = %v", err)
	}

	// Parse the result back to verify redaction
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(result), &resultMap); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Check that string fields were redacted
	if resultMap["Password"] != "***REDACTED***" {
		t.Errorf("Expected Password to be redacted, got: %v", resultMap["Password"])
	}

	// Check nested redaction
	nested := resultMap["Nested"].(map[string]interface{})
	if nested["Secret"] != "***REDACTED***" {
		t.Errorf("Expected Nested.Secret to be redacted, got: %v", nested["Secret"])
	}

	// Check that non-redacted fields remain
	if resultMap["Name"] != "John" {
		t.Errorf("Expected Name to remain unchanged, got: %v", resultMap["Name"])
	}

	// Check that non-string field was zeroed
	if resultMap["Age"] != float64(0) { // JSON unmarshals numbers as float64
		t.Errorf("Expected Age to be zeroed, got: %v", resultMap["Age"])
	}

	// Check that Token in nested field was redacted but Value remains
	if nested["Token"] != "auth-token" {
		t.Errorf("Expected Nested.Token to remain unchanged, got: %v", nested["Token"])
	}
}

func TestSerializeAndRedactXML(t *testing.T) {
	testData := &TestStruct{
		Name:     "John",
		Password: "secret123",
		Age:      30,
		Nested: &NestedStruct{
			Secret: "nested-secret",
			Token:  "auth-token",
			Value:  42,
		},
	}

	redactFields := []string{
		"Password",
		"Nested.Secret",
	}

	result, err := SerializeAndRedactXML(testData, redactFields)
	if err != nil {
		t.Fatalf("SerializeAndRedactXML() error = %v", err)
	}

	// Check that the result contains redacted values
	if !strings.Contains(result, "***REDACTED***") {
		t.Error("Expected XML to contain redacted values")
	}

	// Check that non-redacted values remain
	if !strings.Contains(result, "John") {
		t.Error("Expected XML to contain non-redacted name")
	}

	// Parse back to verify structure
	var resultStruct TestStruct
	if err := xml.Unmarshal([]byte(result), &resultStruct); err != nil {
		t.Fatalf("Failed to parse XML result: %v", err)
	}

	if resultStruct.Password != "***REDACTED***" {
		t.Errorf("Expected Password to be redacted in XML, got: %v", resultStruct.Password)
	}

	if resultStruct.Name != "John" {
		t.Errorf("Expected Name to remain unchanged in XML, got: %v", resultStruct.Name)
	}
}

func TestSerializeAndRedactJSONInvalidInput(t *testing.T) {
	// Test with non-pointer input
	testData := TestStruct{Name: "John"}
	_, err := SerializeAndRedactJSON(testData, []string{"Name"})
	if err == nil {
		t.Error("Expected error for non-pointer input")
	}

	// Test with nil pointer
	var nilData *TestStruct
	_, err = SerializeAndRedactJSON(nilData, []string{"Name"})
	if err == nil {
		t.Error("Expected error for nil pointer input")
	}

	// Test with pointer to non-struct
	stringPtr := "test"
	_, err = SerializeAndRedactJSON(&stringPtr, []string{})
	if err == nil {
		t.Error("Expected error for pointer to non-struct")
	}
}

func TestSerializeAndRedactXMLInvalidInput(t *testing.T) {
	// Test with non-pointer input
	testData := TestStruct{Name: "John"}
	_, err := SerializeAndRedactXML(testData, []string{"Name"})
	if err == nil {
		t.Error("Expected error for non-pointer input")
	}
}

func TestRedactFieldWithNonSettableField(t *testing.T) {
	// This test verifies the behavior when trying to redact a non-settable field
	// In practice, this would be logged but not cause a panic
	testData := &TestStruct{Name: "John"}
	v := reflect.ValueOf(testData).Elem()
	nameField := v.FieldByName("Name")

	// Make field non-settable by getting it through Interface() and back
	nonSettableField := reflect.ValueOf(nameField.Interface())

	// This should not panic and should log a debug message
	redactField(nonSettableField, "TestField")

	// Original field should remain unchanged since we used a copy
	if testData.Name != "John" {
		t.Error("Original field should not be modified when using non-settable copy")
	}
}

func BenchmarkNavigateToField(b *testing.B) {
	testData := &TestStruct{
		Name: "John",
		Nested: &NestedStruct{
			DeepNest: &DeepNestedStruct{
				HiddenField: "deep-secret",
			},
		},
	}

	v := reflect.ValueOf(testData).Elem()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		navigateToField(v, "Nested.DeepNest.HiddenField")
	}
}

func BenchmarkSerializeAndRedactJSON(b *testing.B) {
	testData := &TestStruct{
		Name:     "John",
		Password: "secret123",
		Nested: &NestedStruct{
			Secret: "nested-secret",
			DeepNest: &DeepNestedStruct{
				HiddenField: "deep-secret",
			},
		},
	}

	redactFields := []string{"Password", "Nested.Secret", "Nested.DeepNest.HiddenField"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SerializeAndRedactJSON(testData, redactFields)
		if err != nil {
			b.Fatalf("SerializeAndRedactJSON failed: %v", err)
		}
	}
}
