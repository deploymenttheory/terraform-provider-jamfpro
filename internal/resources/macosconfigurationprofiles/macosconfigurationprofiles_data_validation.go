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
	"html"
	"io"
	"log"
	"strconv"
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
	// Convert Terraform state to regular XML format
	terraformStateXML, err := convertTerraformStateToXML(old)
	if err != nil {
		log.Printf("[ERROR] Failed to convert Terraform state to XML: %v", err)
		return false
	}

	// Convert Jamf Pro state (if necessary) to regular XML format
	jamfProStateXML, err := convertJamfProXMLResponseToRegularXML(new)
	if err != nil {
		log.Printf("[ERROR] Failed to convert Jamf Pro state to XML: %v", err)
		return false
	}

	// Parse the normalized Terraform state XML
	currentTFStateOfConfigProfilePayload, err := parsePayloadXML(terraformStateXML)
	if err != nil {
		log.Printf("[ERROR] Failed to parse current TFState Of ConfigProfile payload: %v", err)
		return false
	}

	// Parse the normalized Jamf Pro state XML
	currentJamfProStateOfConfigProfilePayload, err := parsePayloadXML(jamfProStateXML)
	if err != nil {
		log.Printf("[ERROR] Failed to parse current Jamf Pro State Of ConfigProfile payload: %v", err)
		return false
	}

	// Compare the parsed XML payloads
	result, differences := comparePayloadsWithIgnoredFields(currentTFStateOfConfigProfilePayload, currentJamfProStateOfConfigProfilePayload)
	if !result {
		for _, diff := range differences {
			log.Printf("[DIFFERENCE] %s", diff)
		}
	}

	log.Printf("[DEBUG] Apply Difference Suppression for Configuration Profile: %t", result)
	return result
}

// parsePayloadXML parses an XML string into a hierarchical structure of XMLElements.
// parsePayloadXML parses an XML string into a hierarchical structure of XMLElements.
func parsePayloadXML(xmlString string) (*XMLElement, error) {
	log.Printf("[DEBUG] Starting parsePayloadXML")
	decoder := xml.NewDecoder(strings.NewReader(xmlString))

	var root *XMLElement
	var currentElement *XMLElement
	stack := []*XMLElement{}
	var currentKeyName string // Temporary variable to hold the key name
	currentPath := ""         // Initialize current path

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
			// Build the path for the new element
			newPath := currentPath
			if newPath != "" {
				newPath += "/" // Add a separator if the path is not empty
			}
			newPath += currentKeyName

			newElement := &XMLElement{
				KeyName:  currentKeyName,
				Value:    "",
				Path:     newPath, // Assign the path to the new element
				Children: make(map[string]*XMLElement),
			}

			if root == nil {
				root = newElement
			} else {
				currentElement.Children[currentKeyName] = newElement
			}

			stack = append(stack, newElement)
			currentElement = newElement
			currentPath = newPath // Update current path
			currentKeyName = ""   // Reset key name after starting a new element

		case xml.CharData:
			trimmedData := string(bytes.TrimSpace(se))
			if currentElement != nil && currentKeyName != "" {
				currentElement.Value = trimmedData
			} else {
				currentKeyName = trimmedData
			}

		case xml.EndElement:
			if len(stack) > 1 { // Check if there are more elements in the stack
				stack = stack[:len(stack)-1] // Pop the stack
				currentElement = stack[len(stack)-1]
				currentPath = currentElement.Path // Update current path to the parent element's path
			} else {
				// Resetting at the root level
				currentPath = ""
				currentElement = nil
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

	var compareElements func(*XMLElement, *XMLElement, string)
	compareElements = func(tfElem, jamfElem *XMLElement, path string) {
		// Log the beginning of a comparison at a new path
		log.Printf("[DEBUG] Starting comparison at path: %s", path)

		if tfElem == nil || jamfElem == nil {
			// Handle nil elements
			log.Printf("[DEBUG] One or both elements are nil at path: %s", path)
			if tfElem != nil {
				differences = append(differences, fmt.Sprintf("Path: '%s', Key missing in Jamf Pro state: '%s'", path, tfElem.KeyName))
			}
			if jamfElem != nil {
				differences = append(differences, fmt.Sprintf("Path: '%s', Key missing in Terraform state: '%s'", path, jamfElem.KeyName))
			}
			return
		}

		// Log the number of children in each element for debugging
		log.Printf("[DEBUG] Number of children in Terraform element: %d, Jamf Pro element: %d at path: %s", len(tfElem.Children), len(jamfElem.Children), path)

		for key, tfChild := range tfElem.Children {
			// Check if the current key should be ignored
			if isIgnoredField(key) {
				log.Printf("[DEBUG] Ignoring field '%s' at path: %s", key, path)
				continue // Skip comparison for this ignored field
			}

			jamfChild, exists := jamfElem.Children[key]
			newPath := path + "/" + key

			// Log the comparison of each key and value
			log.Printf("[DEBUG] Comparing key: '%s' at path: %s", key, newPath)

			if exists {
				// Compare values and log
				if tfChild.Value != jamfChild.Value {
					differences = append(differences, fmt.Sprintf("Path: '%s', Value difference for key '%s': Terraform Value: '%s', Jamf Pro Value: '%s'", newPath, key, tfChild.Value, jamfChild.Value))
					log.Printf("[DEBUG] Value difference found for key '%s' at path: %s", key, newPath)
				}
				compareElements(tfChild, jamfChild, newPath) // Recursively compare children
			} else {
				differences = append(differences, fmt.Sprintf("Path: '%s', Key missing in Jamf Pro state: '%s'", newPath, key))
				log.Printf("[DEBUG] Key missing in Jamf Pro state: '%s' at path: %s", key, newPath)
			}
		}

		// Check for new keys in Jamf Pro state
		for key := range jamfElem.Children {
			if _, exists := tfElem.Children[key]; !exists {
				differences = append(differences, fmt.Sprintf("Path: '%s', New key in Jamf Pro state: '%s'", path, key))
				log.Printf("[DEBUG] New key in Jamf Pro state: '%s' at path: %s", key, path)
			}
		}
	}

	compareElements(tfElement, jamfElement, "")
	return len(differences) == 0, differences
}

// isIgnoredField checks if a key should be ignored.
func isIgnoredField(keyName string) bool {
	ignoredFields := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	return contains(ignoredFields, keyName)
}

// contains checks if a string slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// convertTerraformStateToXML attempts to unescape a JSON-encoded XML string.
// If the string is not JSON-encoded, it returns it as is.
func convertTerraformStateToXML(jsonEncodedXML string) (string, error) {
	// Log the JSON-encoded XML for debugging
	log.Printf("[DEBUG] JSON-Encoded XML: %s", jsonEncodedXML)

	// Try to unescape as if it's a JSON string
	unescapedXML, err := strconv.Unquote(`"` + jsonEncodedXML + `"`)
	if err == nil {
		return unescapedXML, nil
	}

	// If unquote fails, check if it's already a regular XML
	if strings.HasPrefix(jsonEncodedXML, "<?xml") {
		return jsonEncodedXML, nil
	}

	// If neither, return the original error
	log.Printf("[ERROR] Failed to unescape JSON-encoded XML string: %v", err)
	return "", fmt.Errorf("failed to unescape string: %w", err)
}

// formatXML formats an XML string with proper indentation and removes unnecessary spaces.
func formatXML(xmlStr string) (string, error) {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "    ") // Indent with 4 spaces

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	// Remove spaces between XML tags
	formatted := strings.ReplaceAll(buf.String(), ">\n    <", "><")

	return formatted, nil
}

// convertJamfProXMLResponseToRegularXML converts an XML string with HTML entity-encoded characters into regular XML and formats it.
func convertJamfProXMLResponseToRegularXML(escapedXML string) (string, error) {
	// Unescape HTML entities
	unescapedXML := html.UnescapeString(escapedXML)

	// Format the unescaped XML
	formattedXML, err := formatXML(unescapedXML)
	if err != nil {
		return "", fmt.Errorf("failed to format XML: %w", err)
	}

	// Debug: Print formatted XML
	log.Printf("[DEBUG] Jamf Pro Server XML Payload is: \n%s", formattedXML)

	return formattedXML, nil
}

// formatmacOSConfigurationProfileXMLPayload prepares the xml payload for upload into Jamf Pro..
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
