// hash.go
// This package contains shared / common hash functions
package common

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"log"
)

// navigateToField recursively navigates through nested struct fields using dot notation
// Returns the reflect.Value of the final field and whether it was found
func navigateToField(v reflect.Value, fieldPath string) (reflect.Value, bool) {
	if fieldPath == "" {
		return reflect.Value{}, false
	}

	parts := strings.Split(fieldPath, ".")
	return navigateToFieldRecursive(v, parts)
}

// navigateToFieldRecursive is the recursive helper function that walks the dot path
func navigateToFieldRecursive(current reflect.Value, parts []string) (reflect.Value, bool) {
	if len(parts) == 0 {
		return current, true
	}

	// Dereference pointer
	if current.Kind() == reflect.Ptr {
		if current.IsNil() {
			return reflect.Value{}, false
		}
		current = current.Elem()
	}

	if current.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	fieldName := parts[0]
	field := current.FieldByName(fieldName)

	if !field.IsValid() {
		return reflect.Value{}, false
	}

	// Recursive case: navigate to the next level with remaining parts
	return navigateToFieldRecursive(field, parts[1:])
}

// redactField redacts a single field based on its type
func redactField(field reflect.Value, fieldPath string) {
	if !field.CanSet() {
		log.Printf("[DEBUG] Cannot set field '%s' - not settable", fieldPath)
		return
	}

	if field.Kind() == reflect.String {
		field.SetString("***REDACTED***")
		log.Printf("[DEBUG] REDACTED: String field '%s' redacted", fieldPath)
	} else {
		log.Printf("[DEBUG] REDACTED: Field '%s' zeroed in output", fieldPath)
		field.Set(reflect.Zero(field.Type()))
	}
}

// SerializeAndRedactXML serializes a resource to XML and redacts specified fields.
// Supports nested field paths using dot notation (e.g Parent.Child.Grandchild)
func SerializeAndRedactXML(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, fieldPath := range redactFields {
		if field, found := navigateToField(resourceCopy, fieldPath); found {
			redactField(field, fieldPath)
		} else {
			log.Printf("[DEBUG] Field path '%s' not found or not accessible", fieldPath)
		}
	}

	if marshaledXML, err := xml.MarshalIndent(resourceCopy.Interface(), "", "  "); err != nil {
		return "", fmt.Errorf("failed to marshal %s to XML: %v", v.Elem().Type(), err)
	} else {
		return string(marshaledXML), nil
	}
}

// SerializeAndRedactJSON serializes a resource to JSON and redacts specified fields.
// Supports nested field paths using dot notation (e.g., "Kitchen.Bowl.Fruit")
func SerializeAndRedactJSON(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, fieldPath := range redactFields {
		if field, found := navigateToField(resourceCopy, fieldPath); found {
			redactField(field, fieldPath)
		} else {
			log.Printf("[DEBUG] Field path '%s' not found or not accessible", fieldPath)
		}
	}

	marshaledJSON, err := json.MarshalIndent(resourceCopy.Interface(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s to JSON: %v", v.Elem().Type(), err)
	}

	return string(marshaledJSON), nil
}

// getIDField returns the value of the ID field in a response.
func getIDField(response interface{}) (any, error) {
	v := reflect.ValueOf(response).Elem()

	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("ID field not found in response")
	}

	str, ok := idField.Interface().(string)
	if ok {
		return str, nil
	}

	integer, ok := idField.Interface().(int)
	if ok {
		return strconv.Itoa(integer), nil
	}
	return nil, fmt.Errorf("unsupported type")
}
