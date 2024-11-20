// common/configurationprofiles/plist/plistdiffsuppression.go
// contains the functions to process configuration profiles for diff suppression.
package plist

import (
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"

	"howett.net/plist"
)

// ProcessConfigurationProfileForDiffSuppression processes the plist data through multiple steps
// to prepare it for diff suppression. This function is used to remove specified fields, normalize base64 content,
// normalize XML tags, unescape HTML entities, sort keys, and encode back to plist.
func ProcessConfigurationProfileForDiffSuppression(plistData string, fieldsToRemove []string) (string, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Step 1: Unmarshal
	var rawData map[string]interface{}
	if _, err := plist.Unmarshal([]byte(plistData), &rawData); err != nil {
		log.Printf("Error unmarshalling plist data: %v\n", err)
		return "", err
	}

	// Step 2: Remove specified fields
	processedData := removeSpecifiedXMLFields(rawData, fieldsToRemove, "")

	// Step 3: Normalize base64 content
	normalizedBase64 := normalizeBase64Content(processedData)

	// Step 4: Normalize XML tags
	normalizedXML := normalizeXMLTags(normalizedBase64)

	// Step 5: Unescape HTML Entities
	normalizedData := unescapeHTMLEntities(normalizedXML)

	// Step 5: Sort keys
	sortedData := SortPlistKeys(normalizedData.(map[string]interface{}))

	// Step 6: Encode back to plist
	encodedPlist, err := EncodePlist(sortedData)
	if err != nil {
		log.Printf("Error encoding plist data: %v\n", err)
		return "", err
	}

	// Step 7: Remove trailing whitespace
	return trimTrailingWhitespace(encodedPlist), nil
}

// removeSpecifiedXMLFields( removes specified fields from the plist data recursively.
// useful for removing jamfpro specific unique identifiers from the plist data.
func removeSpecifiedXMLFields(data map[string]interface{}, fieldsToRemove []string, path string) map[string]interface{} {
	// Create a set of fields to remove for quick lookup
	fieldsToRemoveSet := make(map[string]struct{}, len(fieldsToRemove))
	for _, field := range fieldsToRemove {
		fieldsToRemoveSet[field] = struct{}{}
	}

	// Iterate over the map and remove fields if they exist
	for field := range fieldsToRemoveSet {
		if _, exists := data[field]; exists {
			log.Printf("Removing field: %s from path: %s\n", field, path)
			delete(data, field)
		}
	}

	// Recursively process nested maps and arrays
	for key, value := range data {
		newPath := path + "/" + key
		switch v := value.(type) {
		case map[string]interface{}:
			log.Printf("Recursively removing fields in nested map at path: %s\n", newPath)
			removeSpecifiedXMLFields(v, fieldsToRemove, newPath)
		case []interface{}:
			for i, item := range v {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					log.Printf("Recursively removing fields in array at path: %s[%d]\n", newPath, i)
					removeSpecifiedXMLFields(nestedMap, fieldsToRemove, newPath+strings.ReplaceAll(key, "/", "_")+strconv.Itoa(i))
				}
			}
			// Ensure empty arrays are preserved
			data[key] = v
		}
	}

	return data
}

// normalizeBase64Content normalizes base64 data by removing all whitespace
func normalizeBase64Content(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return NormalizeBase64(v) // Normalize Base64 string
	case map[string]interface{}:
		for key, value := range v {
			v[key] = normalizeBase64Content(value) // Recursively normalize
		}
		return v
	case []interface{}:
		for i, item := range v {
			v[i] = normalizeBase64Content(item) // Recursively normalize
		}
		return v
	default:
		return data // Unsupported types remain unchanged
	}
}

// normalizeXMLTags standardizes XML tag formatting for malformed config profile xml
// handles the following cases:
// <true />
// <true    />
// <true  \t />
// <false   />
// <string    />
func normalizeXMLTags(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		if strings.Contains(v, "/") {
			trimmed := strings.TrimSpace(v)
			return regexp.MustCompile(`\s+/>`).ReplaceAllString(trimmed, "/>")
		}
		return v
	case map[string]interface{}:
		for key, value := range v {
			v[key] = normalizeXMLTags(value)
		}
	case []interface{}:
		for i, item := range v {
			v[i] = normalizeXMLTags(item)
		}
	}
	return data
}

// unescapeHTMLEntities applies html.UnescapeString recursively
func unescapeHTMLEntities(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return html.UnescapeString(v)
	case map[string]interface{}:
		for key, value := range v {
			v[key] = unescapeHTMLEntities(value)
		}
	case []interface{}:
		for i, item := range v {
			v[i] = unescapeHTMLEntities(item)
		}
	}
	return data
}

// trimTrailingWhitespace removes trailing whitespace from each line of the plist
func trimTrailingWhitespace(plist string) string {
	lines := strings.Split(plist, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}
