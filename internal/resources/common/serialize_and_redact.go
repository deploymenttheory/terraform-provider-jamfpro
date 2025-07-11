// This package contains shared / common hash functions
package common

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// SerializeAndRedactXML serializes a resource to XML and redacts specified fields.
func SerializeAndRedactXML(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			}
		}
	}

	if marshaledXML, err := xml.MarshalIndent(resourceCopy.Interface(), "", "  "); err != nil {
		return "", fmt.Errorf("failed to marshal %s to XML: %v", v.Elem().Type(), err)
	} else {
		return string(marshaledXML), nil
	}
}

// SerializeAndRedactJSON serializes a resource to JSON and redacts specified fields.
func SerializeAndRedactJSON(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			} else {
				log.Printf("[DEBUG] REDACTED: '%v' Zeroed in output", field)
				f.Set(reflect.Zero(f.Type()))
			}
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
