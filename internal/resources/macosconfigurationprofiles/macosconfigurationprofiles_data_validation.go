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

// XMLElement represents an xml key and value pair and it's path.
type XMLElement struct {
	Path  string
	Key   string
	Value string
}

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

// parsePayloadXML parses an XML string into a map of XMLElements.
func parsePayloadXML(xmlString string) (map[string]XMLElement, error) {
	log.Printf("[DEBUG] Starting parsePayloadXML")
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	payload := make(map[string]XMLElement)
	var currentPath []string
	var lastKey string

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
			lastKey = strings.Join(currentPath, "/")
		case xml.CharData:
			trimmedData := string(bytes.TrimSpace(se))
			if !isIgnoredField(lastKey) {
				payload[lastKey] = XMLElement{Key: lastKey, Value: trimmedData}
			}
		case xml.EndElement:
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
			lastKey = strings.Join(currentPath, "/")
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
// It first converts both payloads into maps of XMLElements, which include the path, key name, and value,
// then compares these maps.
// The function returns false if:
// - A key exists in one payload but not the other.
// - The same key exists in both payloads but with different values.
// - There is a new field in the Jamf Pro state that is not in the Terraform state, and it's not an ignored field.
func comparePayloadsWithIgnoredFields(tfPayload, jamfPayload map[string]XMLElement) bool {
	// Initialize a map to track differences
	differences := make(map[string][]XMLElement)

	// Compare the payloads and record differences
	for path, tfElement := range tfPayload {
		if jamfElement, exists := jamfPayload[path]; !exists {
			differences[path] = []XMLElement{tfElement, {Key: "Missing in Jamf Pro"}}
		} else if tfElement.Value != jamfElement.Value {
			differences[path] = []XMLElement{tfElement, jamfElement}
		}
	}

	for path, jamfElement := range jamfPayload {
		if _, exists := tfPayload[path]; !exists {
			differences[path] = []XMLElement{{Key: "New in Jamf Pro"}, jamfElement}
		}
	}

	// Log the differences
	if len(differences) > 0 {
		for path, elements := range differences {
			log.Printf("[DIFFERENCE] XML Path: '%s', Terraform Key Name: '%s', Terraform Value: '%s', Jamf Pro Key Name: '%s', Jamf Pro Value: '%s'",
				path, elements[0].Key, elements[0].Value, elements[1].Key, elements[1].Value)
		}
		return false
	}

	log.Printf("[DEBUG] No differences found between Configuration Profile in Terraform state and Jamf Pro state.")
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
