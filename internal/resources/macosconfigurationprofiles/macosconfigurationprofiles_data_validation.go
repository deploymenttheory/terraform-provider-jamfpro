// macosconfigurationprofiles_data_validation.go
/*
Hierarchy and Order: The comparison is not sensitive to the order of fields, but it is essential that the hierarchy is preserved when nesting key names and their values.

Difference Suppression: Differences should be suppressed if fields at a certain level in the hierarchy match in both key and value. This includes ensuring that nested keys and values maintain their hierarchical structure.

New Keys Detection: The function needs to detect new keys in the Jamf Pro state that do not exist in the Terraform state.

Key Value Changes: It should accurately detect changes in key values across both payloads.

Comparison Logic: The comparison logic should be refined to compare key names and values within their hierarchical context, ensuring that only fields with matching key names are compared.
*/
package macosconfigurationprofiles

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// XMLElement represents an xml key and value pair and it's path.
type XMLElement struct {
	KeyName  string
	Value    string
	Path     string
	Children map[string]*XMLElement
}

// suppressPayloadDiff compares Terraform state and Jamf Pro state payloads, suppressing diffs for specified fields.
// It calls comparePayloadsWithIgnoredFields with the root elements of both XML structures.
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

	result, differences := comparePayloadsWithIgnoredFields(currentTFStateOfConfigProfilePayload, currentJamfProStateOfConfigProfilePayload)
	if !result {
		for _, diff := range differences {
			log.Printf("[DIFFERENCE] %s", diff)
		}
	}

	log.Printf("[DEBUG] suppressPayloadDiff result: %t", result)
	return result
}

// parsePayloadXML parses an XML string into a hierarchical structure of XMLElements.
func parsePayloadXML(xmlString string) (*XMLElement, error) {
	log.Printf("[DEBUG] Starting parsePayloadXML")
	decoder := xml.NewDecoder(strings.NewReader(xmlString))

	root := &XMLElement{Children: make(map[string]*XMLElement)} // Root of the XML tree
	currentElement := root
	stack := []*XMLElement{}  // Stack to keep track of parent elements
	var currentKeyName string // Variable to store the current key name

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			log.Printf("[ERROR] Error parsing XML: %v", err)
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			// Only create a new element if it's not a <key> element
			if se.Name.Local != "key" {
				newElement := &XMLElement{
					KeyName:  currentKeyName,
					Children: make(map[string]*XMLElement),
				}
				// Add the new element as a child of the current element
				currentElement.Children[currentKeyName] = newElement
				// Push the current element onto the stack and make the new element the current element
				stack = append(stack, currentElement)
				currentElement = newElement
				// Reset the currentKeyName since we've used it
				currentKeyName = ""
			}

		case xml.CharData:
			trimmedData := string(bytes.TrimSpace(se))
			if currentElement.KeyName == "" {
				// If KeyName is empty, this CharData is the key name
				currentKeyName = trimmedData
			} else {
				// Else, it's the value of the current element
				currentElement.Value = trimmedData
			}

		case xml.EndElement:
			// Pop the last element from the stack and make it the current element
			if len(stack) > 0 {
				currentElement = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
	}

	log.Printf("[DEBUG] Finished parsing XML.")
	return root, nil
}

// comparePayloadsWithIgnoredFields recursively compares XML payloads from Terraform state and Jamf Pro server state.
// This function respects the hierarchical structure of XML elements and is designed to ignore specific fields.
// Differences are identified based on the key names and their corresponding values within their hierarchical context.
//
// The comparison is hierarchy-aware and is not sensitive to the order of elements within the same level of hierarchy.
//
// The function returns false if it finds differences that should not be ignored. These differences include:
// - A key exists in one payload but not the other, and it's not an ignored field.
// - The same key exists in both payloads but with different values, and it's not an ignored field.
// - There is a new field in the Jamf Pro state that is not in the Terraform state, and it's not an ignored field.
//
// Parameters:
// - tfElement: The root element of the Terraform state XML payload.
// - jamfElement: The root element of the Jamf Pro state XML payload.
//
// Returns:
// - A boolean indicating if the payloads are considered equal when ignoring specified fields.
// - A slice of strings detailing the differences found, if any.
func comparePayloadsWithIgnoredFields(tfElement, jamfElement *XMLElement) (bool, []string) {
	var differences []string

	// Function to recursively compare elements
	var compareElements func(*XMLElement, *XMLElement)
	compareElements = func(tfElem, jamfElem *XMLElement) {
		// Skip ignored fields
		if isIgnoredField(tfElem.KeyName) {
			return
		}

		// Check for value difference at the same hierarchy level
		if tfElem.KeyName == jamfElem.KeyName && tfElem.Value != jamfElem.Value {
			differences = append(differences, fmt.Sprintf("Value difference for key '%s': Terraform Value: '%s', Jamf Pro Value: '%s'", tfElem.KeyName, tfElem.Value, jamfElem.Value))
		}

		// Recursively compare children elements
		for key, tfChild := range tfElem.Children {
			if jamfChild, exists := jamfElem.Children[key]; exists {
				compareElements(tfChild, jamfChild)
			} else {
				differences = append(differences, fmt.Sprintf("Key missing in Jamf Pro state: '%s'", key))
			}
		}

		// Check for new keys in Jamf Pro state
		for key := range jamfElem.Children {
			if _, exists := tfElem.Children[key]; !exists {
				differences = append(differences, fmt.Sprintf("New key found in Jamf Pro state: '%s'", key))
			}
		}
	}

	compareElements(tfElement, jamfElement)
	return len(differences) == 0, differences
}

// isIgnoredField checks if a key should be ignored
func isIgnoredField(keyName string) bool {
	ignoredFields := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	return contains(ignoredFields, keyName)
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
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
