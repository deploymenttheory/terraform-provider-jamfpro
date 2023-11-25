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

// suppressPayloadDiff compares Terraform state and Jamf Pro state payloads, suppressing diffs for specified fields.
func suppressPayloadDiff(k, old, new string, d *schema.ResourceData) bool {
	currentTFStateOfConfigProfilePayload, err := parsePayloadXML(old)
	if err != nil {
		log.Printf("[ERROR] Failed to parse current TFState Of ConfigProfile payload: %v", err)
		return false
	}

	currentJamfProStateOfConfigProfilePayload, err := parsePayloadXML(new)
	if err != nil {
		log.Printf("[ERROR] Failed to parse current Jamf Pro State Of ConfigProfile payload: %v", err)
		return false
	}

	// Check for differences in payloads, ignoring the specified fields
	result := comparePayloadsWithIgnoredFields(currentTFStateOfConfigProfilePayload, currentJamfProStateOfConfigProfilePayload)
	log.Printf("[DEBUG] suppressPayloadDiff result: %t", result)
	return result
}

// parsePayloadXML parses an XML string into a map of key-value pairs, excluding specific fields.
func parsePayloadXML(xmlString string) (map[string]string, error) {
	log.Printf("[DEBUG] Starting parsePayloadXML")
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	payload := make(map[string]string)
	var currentPath []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			log.Printf("[DEBUG] End of XML file reached")
			break
		}
		if err != nil {
			log.Printf("[ERROR] Error parsing XML: %v", err)
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			currentPath = append(currentPath, se.Name.Local)
			log.Printf("[DEBUG] Start Element: %s, Current Path: %s", se.Name.Local, strings.Join(currentPath, "/"))
		case xml.CharData:
			path := strings.Join(currentPath, "/")
			trimmedData := string(bytes.TrimSpace(se))
			log.Printf("[DEBUG] Char Data: %s, Path: %s, Trimmed Data: '%s'", se, path, trimmedData)
			if !isIgnoredField(path) {
				payload[path] = trimmedData
			}
		case xml.EndElement:
			log.Printf("[DEBUG] End Element: %s", se.Name.Local)
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}

	for key, value := range payload {
		if isIgnoredField(key) {
			log.Printf("[DEBUG] Ignoring field: %s", key)
			delete(payload, key)
		} else {
			log.Printf("[DEBUG] Field: %s, Value: %s", key, value)
		}
	}

	log.Printf("[DEBUG] Finished parsing XML. Total fields parsed (excluding ignored): %d", len(payload))
	return payload, nil
}

// isIgnoredField determines if a field should be ignored based on its path.
func isIgnoredField(fieldPath string) bool {
	ignoredFields := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	for _, ignored := range ignoredFields {
		if strings.Contains(fieldPath, ignored) {
			return true
		}
	}
	return false
}

// comparePayloadsWithIgnoredFields checks if the XML payloads from Terraform state and Jamf Pro server state are equal,
// while ignoring specific fields. This comparison is not sensitive to the order of XML elements.
// It first converts both payloads into sets (excluding ignored fields), then compares these sets.
// The function returns false if:
// - A key exists in one payload but not the other.
// - The same key exists in both payloads but with different values.
// - There is a new field in the Jamf Pro state that is not in the Terraform state, and it's not an ignored field.
// This approach ensures a more flexible comparison that can accurately detect meaningful differences,
// including new fields added in the Jamf Pro state, while ignoring the differences in the order of XML elements and ignored fields.
func comparePayloadsWithIgnoredFields(tfPayload, jamfPayload map[string]string) bool {
	// Convert maps to sets for comparison
	tfSet := make(map[string]struct{})
	jamfSet := make(map[string]struct{})

	// Populate the sets, ignoring the specified fields
	for key, value := range tfPayload {
		if !isIgnoredField(key) {
			tfSet[key] = struct{}{}
			log.Printf("[DEBUG] Terraform Payload - Key: %s, Value: %s", key, value)
		}
	}
	for key, value := range jamfPayload {
		if !isIgnoredField(key) {
			jamfSet[key] = struct{}{}
			log.Printf("[DEBUG] Jamf Pro Payload - Key: %s, Value: %s", key, value)
		}
	}

	// Initialize a variable to track differences
	var hasDifferences bool

	// Compare the sets and log differences
	for key := range tfSet {
		if _, exists := jamfSet[key]; !exists {
			log.Printf("[DIFFERENCE] Key missing in Jamf Pro state: %s", key)
			hasDifferences = true
			continue
		}
		if tfValue, jamfValue := tfPayload[key], jamfPayload[key]; tfValue != jamfValue {
			log.Printf("[DIFFERENCE] Value difference found for key '%s': Terraform State: '%s', Jamf Pro State: '%s'", key, tfValue, jamfValue)
			hasDifferences = true
		}
	}

	for key := range jamfSet {
		if _, exists := tfSet[key]; !exists {
			log.Printf("[DIFFERENCE] New key found in Jamf Pro state: %s", key)
			hasDifferences = true
		}
	}

	// If differences were found, log them
	if hasDifferences {
		log.Printf("[DEBUG] Differences detected between Terraform state and Jamf Pro state.")
	} else {
		log.Printf("[DEBUG] No differences found between Terraform state and Jamf Pro state.")
	}

	// Return false if differences were found
	return !hasDifferences
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
