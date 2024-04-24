// common/redact.go
package constructobject

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

// SerializeAndRedactXML serializes a resource to XML and redacts specified fields.
func SerializeAndRedactXML(resource interface{}, redactFields []string) (string, error) {
	// Ensure the resource passed is a pointer to a struct
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	// Create a deep copy of the resource to avoid modifying the original data
	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	// Apply redactions
	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			}
		}
	}

	// Serialize the redacted resource to XML
	if marshaledXML, err := xml.MarshalIndent(resourceCopy.Interface(), "", "  "); err != nil {
		return "", fmt.Errorf("failed to marshal %s to XML: %v", v.Elem().Type(), err)
	} else {
		return string(marshaledXML), nil
	}
}

// SerializeAndRedactJSON serializes a resource to JSON and redacts specified fields.
func SerializeAndRedactJSON(resource interface{}, redactFields []string) (string, error) {
	// Ensure the resource passed is a pointer to a struct
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	// Create a deep copy of the resource to avoid modifying the original data
	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	// Apply redactions
	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			}
		}
	}

	// Serialize the redacted resource to JSON
	marshaledJSON, err := json.MarshalIndent(resourceCopy.Interface(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s to JSON: %v", v.Elem().Type(), err)
	}

	return string(marshaledJSON), nil
}
