// common/configurationprofiles/datavalidators/helpers.go
package datavalidators

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"howett.net/plist"
)

// GetPayloadScope retrieves the 'PayloadScope' key from the decoded plist data.
func GetPayloadScope(plistData map[string]interface{}) (string, error) {
	if scope, ok := plistData["PayloadScope"].(string); ok {
		return scope, nil
	}

	return "", fmt.Errorf("'PayloadScope' key not found in plist")
}

// CheckPlistIndentationAndWhiteSpace checks the plist XML for proper indentation and whitespace.
func CheckPlistIndentationAndWhiteSpace(plistStr string) error {
	// Decode the plist XML
	var decoded interface{}

	_, err := plist.Unmarshal([]byte(plistStr), &decoded)
	if err != nil {

		return fmt.Errorf("invalid plist: %v", err)
	}

	// Re-encode the plist to check formatting
	formatted, err := FormatPlist(plistStr)
	if err != nil {

		return fmt.Errorf("error formatting plist: %v", err)
	}

	// Normalize both original and formatted plists to compare
	normalizedOriginal := NormalizeXML(plistStr)
	normalizedFormatted := NormalizeXML(formatted)

	// Split the normalized strings into lines for detailed comparison
	origLines := strings.Split(normalizedOriginal, "\n")
	formattedLines := strings.Split(normalizedFormatted, "\n")

	if len(origLines) != len(formattedLines) {
		log.Printf("[DEBUG] Line count mismatch: original has %d lines, formatted has %d lines", len(origLines), len(formattedLines))

		return fmt.Errorf("plist line count mismatch: source has %d lines, formatted has %d lines, check for trailing lines and whitespace", len(origLines), len(formattedLines))
	}

	for i := range origLines {
		if origLines[i] != formattedLines[i] {
			log.Printf("[DEBUG] Difference at line %d:\nOriginal: %s\nFormatted: %s\n", i+1, origLines[i], formattedLines[i])

			return fmt.Errorf("plist is not properly indented at line %d", i+1)
		}
	}

	return nil
}

// FormatPlist formats the plist structure to a properly indented XML string.
func FormatPlist(plistStr string) (string, error) {
	var decoded interface{}
	_, err := plist.Unmarshal([]byte(plistStr), &decoded)
	if err != nil {

		return "", fmt.Errorf("invalid plist: %v", err)
	}

	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("\t") // Indent with a single tab
	err = encoder.Encode(decoded)
	if err != nil {
		return "", err
	}

	formatted := buf.String()
	// Trim any leading/trailing whitespace for comparison
	return strings.TrimSpace(formatted), nil
}

// NormalizeXML normalizes XML string by removing extra spaces and new lines.
func NormalizeXML(xmlStr string) string {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	encoder := xml.NewEncoder(&buf)

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if err := encoder.EncodeToken(token); err != nil {
			break
		}
	}
	encoder.Flush()

	return buf.String()
}
