// macosconfigurationprofiles_data_validation.go
package macosconfigurationprofiles

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressPayloadDiff suppresses diff if only ignored fields are different.
// This function is a custom diff function for Terraform that compares old and new XML payloads,
// ignoring specific fields that are known to be auto-modified by external systems (like Jamf Pro).
func suppressPayloadDiff(k, old, new string, d *schema.ResourceData) bool {
	oldPayload, err := parsePayloadXML(old)
	if err != nil {
		log.Printf("[ERROR] Failed to parse old payload: %v", err)
		return false
	}

	newPayload, err := parsePayloadXML(new)
	if err != nil {
		log.Printf("[ERROR] Failed to parse new payload: %v", err)
		return false
	}

	// Check for differences in payloads, ignoring the specified fields
	return comparePayloadsWithIgnoredFields(oldPayload, newPayload)
}

// parsePayloadXML parses an XML string and returns a map of key-value pairs, excluding specific fields.
// It reads through the XML, builds a map of the XML paths to values, and skips over paths that match ignored fields.
func parsePayloadXML(xmlString string) (map[string]string, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	payload := make(map[string]string)
	var currentPath []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			currentPath = append(currentPath, se.Name.Local) // Push to path
		case xml.CharData:
			path := strings.Join(currentPath, "/")
			if !isIgnoredField(path) {
				payload[path] = string(bytes.TrimSpace(se))
			}
		case xml.EndElement:
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1] // Pop from path
			}
		}
	}

	// Remove the ignored fields from the map
	for key := range payload {
		if isIgnoredField(key) {
			delete(payload, key)
		}
	}

	return payload, nil
}

// isIgnoredField returns true if the field should be ignored.
// This helper function checks if a given XML path ends with any of the predefined ignored field names,
// such as 'PayloadUUID', 'PayloadOrganization', or 'PayloadIdentifier'.
func isIgnoredField(field string) bool {
	ignoredFields := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	for _, ignored := range ignoredFields {
		if strings.HasSuffix(field, "/"+ignored) {
			return true
		}
	}
	return false
}

// comparePayloadsWithIgnoredFields checks if the payloads are equal, ignoring specific fields.
// This function compares two payload maps (representing old and new states) and returns true if they are equal,
// ignoring changes in specific fields like 'PayloadUUID', 'PayloadOrganization', and 'PayloadIdentifier'.
func comparePayloadsWithIgnoredFields(oldPayload, newPayload map[string]string) bool {
	for key, oldValue := range oldPayload {
		if isIgnoredField(key) {
			continue
		}
		if newValue, exists := newPayload[key]; !exists || oldValue != newValue {
			return false
		}
	}
	for key, newValue := range newPayload {
		if isIgnoredField(key) {
			continue
		}
		if oldValue, exists := oldPayload[key]; !exists || oldValue != newValue {
			return false
		}
	}
	return true
}

// formatmacOSConfigurationProfileXMLPayload prepares the xml payload for upload into Jamf Pro
func formatmacOSConfigurationProfileXMLPayload(input string) (string, error) {
	// Decode the XML data
	var buffer bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewBufferString(input))
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("  ", "    ") // Set indentation: prefix for each element, indent for each level

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break // End of file, break out of loop
			}
			return "", err // Return with error
		}

		// Write the token to the buffer in a standard format
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	// Close the encoder to flush the buffer
	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
