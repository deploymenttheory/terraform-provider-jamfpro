// hash.go
// This package contains shared / common hash functions
package common

import (
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"reflect"
)

// HashString calculates the SHA-256 hash of a string and returns it as a hexadecimal string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("Computed hash: %s", hash)
	return hash
}

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

	marshaledJSON, err := json.MarshalIndent(resourceCopy.Interface(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s to JSON: %v", v.Elem().Type(), err)
	}

	return string(marshaledJSON), nil
}
