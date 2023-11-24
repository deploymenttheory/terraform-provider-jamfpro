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

// suppressPayloadDiff suppresses diff if only ignored fields are different
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

	return comparePayloads(oldPayload, newPayload)
}

// parsePayloadXML parses XML string and returns a map of key-value pairs, excluding specific fields
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

	return payload, nil
}

// comparePayloads compares two payload maps
func comparePayloads(oldPayload, newPayload map[string]string) bool {
	for key, oldValue := range oldPayload {
		if newValue, exists := newPayload[key]; !exists || oldValue != newValue {
			return false
		}
	}
	return true
}

// isIgnoredField returns true if the field should be ignored
// "PayloadUUID", "PayloadOrganization", "PayloadIdentifier" are tenant specific fields injected into Jamf Pro config profiles and should be ignored.
func isIgnoredField(field string) bool {
	ignoredFields := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	for _, ignored := range ignoredFields {
		if strings.HasSuffix(field, "/"+ignored) {
			return true
		}
	}
	return false
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
