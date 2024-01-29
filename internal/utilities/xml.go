// utilities.go
// For utility/helper functions to support the jamf pro tf provider
package utilities

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
)

// NormalizeXml takes an XML content represented by the `xmlString` parameter of type `interface{}` and returns a normalized, consistently formatted XML string. The function performs the following steps:
// 1. Checks if the input `xmlString` is nil or an empty string, returning an empty string if true.
// 2. Attempts to unmarshal the input string (asserted to `string` from `interface{}`) into a generic interface{} type to parse the XML structure.
// 3. If unmarshaling fails (indicating invalid XML content), returns an error message string.
// 4. If unmarshaling succeeds, marshals the object back into XML using `xml.MarshalIndent` to ensure consistent formatting and indentation.
// 5. Returns the resulting normalized XML string.
// Note: The function assumes `xmlString` can be type asserted to `string` and may panic otherwise. It is intended for use with valid XML content that needs normalization for consistency or readability.
func NormalizeXml(xmlString interface{}) string {
	if xmlString == nil || xmlString == "" {
		return ""
	}
	var x interface{}

	if err := xml.Unmarshal([]byte(xmlString.(string)), &x); err != nil {
		return fmt.Sprintf("Error parsing XML: %+v", err)
	}
	b, _ := xml.MarshalIndent(x, "", "  ")
	return string(b)
}

// EncodeXmlString escapes special XML characters in a string and replaces them with their corresponding HTML entities.
func EncodeXmlString(s string) string {
	var replacer = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"'", "&apos;",
		"\"", "&quot;",
	)
	return replacer.Replace(s)
}

// DecodeXmlString unescapes HTML entities in a string back to their original XML special characters.
func DecodeXmlString(s string) string {
	return html.UnescapeString(s)
}
