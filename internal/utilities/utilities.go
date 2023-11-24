// utilities.go
// For utility/helper functions to support the jamf pro tf provider
package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// ConvertToXMLSafeString replaces disallowed XML characters in a string with their corresponding XML entity references.
// This function is useful for preparing a string to be safely included in an XML document.
func ConvertToXMLSafeString(s string) string {
	// Define a map of disallowed characters and their XML entity equivalents.
	replacements := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"'":  "&apos;",
		"\"": "&quot;",
	}

	// Replace each disallowed character with its entity reference.
	for key, val := range replacements {
		s = strings.ReplaceAll(s, key, val)
	}

	// Return the XML-safe string.
	return s
}

// ConvertFromXMLSafeString reverses the process of ConvertToXMLSafeString.
// It replaces XML entity references in a string back to their original characters.
// This is useful when reading XML data that contains entity references and converting them back to normal characters.
func ConvertFromXMLSafeString(s string) string {
	// Define a map of XML entities and their corresponding characters.
	replacements := map[string]string{
		"&amp;":  "&",
		"&lt;":   "<",
		"&gt;":   ">",
		"&apos;": "'",
		"&quot;": "\"",
	}

	// Replace each entity reference with its corresponding character.
	for key, val := range replacements {
		s = strings.ReplaceAll(s, key, val)
	}

	// Return the original string with characters restored.
	return s
}

func Base64EncodeCertificate(certPath string) (string, error) {
	// Read the certificate file
	data, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate file: %v", err)
	}

	// Base64 encode the file's content
	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded, nil
}

// EnsureXMLSafeString checks if a string contains disallowed XML characters.
// If it does, it converts the string to an XML-safe format using ConvertToXMLSafeString.
// This function is useful for ensuring that strings are safe for inclusion in XML documents.
func EnsureXMLSafeString(s string) string {
	// Define a set of disallowed XML characters.
	disallowedChars := []string{"&", "<", ">", "'", "\""}

	// Check if the string contains any disallowed characters.
	for _, char := range disallowedChars {
		if strings.Contains(s, char) {
			// If a disallowed character is found, convert the entire string to an XML-safe format.
			return ConvertToXMLSafeString(s)
		}
	}

	// If no disallowed characters are found, return the original string.
	return s
}
