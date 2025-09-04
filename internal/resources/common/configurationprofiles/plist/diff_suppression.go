// common/configurationprofiles/plist/plistdiffsuppression.go
// contains the functions to process configuration profiles for diff suppression.
package plist

import (
	"encoding/base64"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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

	// Step 3: Clean empty array entries
	cleanedArrays := cleanEmptyArrayEntries(processedData)

	// Step 4: Normalize base64 content
	normalizedBase64 := normalizeBase64Content(cleanedArrays)

	// Step 5: Normalize XML tags
	normalizedXML := normalizeXMLTags(normalizedBase64)

	// Step 6: Normalize empty strings
	normalizedStrings := normalizeEmptyStrings(normalizedXML)

	// Step 7: normalize HTML Entities
	normalizedData := normalizeHTMLEntitiesForDiff(normalizedStrings)

	// Step 8: Sort keys
	sortedData := SortPlistKeys(normalizedData.(map[string]interface{}))

	// Step 9: Encode back to plist
	encodedPlist, err := EncodePlist(sortedData)
	if err != nil {
		log.Printf("Error encoding plist data: %v\n", err)
		return "", err
	}

	// Step 10: Remove trailing whitespace
	return trimTrailingWhitespace(encodedPlist), nil
}

// removeSpecifiedXMLFields( removes specified fields from the plist data recursively.
// useful for removing jamfpro specific unique identifiers from the plist data.
func removeSpecifiedXMLFields(data map[string]interface{}, fieldsToRemove []string, path string) map[string]interface{} {
	// Create a set of fields to remove for quick lookup (case-insensitive)
	fieldsToRemoveSet := make(map[string]struct{}, len(fieldsToRemove))
	for _, field := range fieldsToRemove {
		fieldsToRemoveSet[strings.ToLower(field)] = struct{}{}
	}

	// Collect keys to remove to avoid modifying map while iterating
	var keysToRemove []string
	for key := range data {
		if _, exists := fieldsToRemoveSet[strings.ToLower(key)]; exists {
			log.Printf("[DEBUG] Removing field: %s from path: %s\n", key, path)
			keysToRemove = append(keysToRemove, key)
		}
	}

	// Remove the identified keys
	for _, key := range keysToRemove {
		delete(data, key)
	}

	// Recursively process nested maps and arrays
	for key, value := range data {
		newPath := path + "/" + key
		switch v := value.(type) {
		case map[string]interface{}:
			log.Printf("[DEBUG] Recursively removing fields in nested map at path: %s\n", newPath)
			removeSpecifiedXMLFields(v, fieldsToRemove, newPath)
		case []interface{}:
			for i, item := range v {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					log.Printf("[DEBUG] Recursively removing fields in array at path: %s[%d]\n", newPath, i)
					removeSpecifiedXMLFields(nestedMap, fieldsToRemove, newPath+strings.ReplaceAll(key, "/", "_")+strconv.Itoa(i))
				}
			}
			// Ensure empty arrays are preserved
			data[key] = v
		}
	}

	return data
}

func normalizeBase64Content(data interface{}) interface{} {
	// Helper to check and normalize potential base64 string values
	normalizeString := func(s string) string {
		// If string has no spaces/newlines, leave it alone
		if !strings.ContainsAny(s, " \n\t\r") {
			return s
		}

		// Remove all whitespace and try to decode
		clean := strings.Join(strings.Fields(s), "")
		_, err := base64.StdEncoding.DecodeString(clean)
		if err == nil {
			return clean
		}
		return s
	}

	switch v := data.(type) {
	case string:
		return normalizeString(v)

	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = normalizeBase64Content(value)
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeBase64Content(item)
		}
		return result

	default:
		return data
	}
}

// NormalizeBase64 normalizes base64 content by removing all whitespace characters
// Returns normalized base64 string with all spacing/formatting removed for comparison
// Base64 uses characters A-Z, a-z, 0-9, +, /, and = for padding
func NormalizeBase64(input string) string {
	// First remove all whitespace
	trimmed := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // Drop ALL whitespace (spaces, tabs, newlines, etc)
		}
		return r
	}, input)

	// Check if the result is a valid base64 string
	isBase64 := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`).MatchString(trimmed)

	if !isBase64 {
		return input // Not base64, return original
	}

	return trimmed
}

// cleanEmptyArrayEntries removes empty or whitespace-only string entries from arrays
func cleanEmptyArrayEntries(data interface{}) interface{} {
	switch v := data.(type) {
	case []interface{}:
		var cleaned []interface{}
		for _, item := range v {
			switch itemVal := item.(type) {
			case string:
				if strings.TrimSpace(itemVal) != "" {
					cleaned = append(cleaned, cleanEmptyArrayEntries(itemVal))
				}
			default:
				cleaned = append(cleaned, cleanEmptyArrayEntries(itemVal))
			}
		}
		return cleaned
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = cleanEmptyArrayEntries(value)
		}
		return result
	default:
		return data
	}
}

// normalizeXMLTags standardizes XML tag formatting for malformed config profile xml
// handles the following cases:
// < true/>
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
			normalized := regexp.MustCompile(`<\s*(\w+)\s*/>`).ReplaceAllString(trimmed, "<$1/>")
			return normalized
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

// normalizeEmptyStrings standardizes empty and whitespace-only strings
func normalizeEmptyStrings(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return ""
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = normalizeEmptyStrings(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeEmptyStrings(item)
		}
		return result
	default:
		return data
	}
}

// safeHTMLEntity is a regex to detect a single-level valid entity like &lt;, &amp;, etc.
var safeHTMLEntity = regexp.MustCompile(`&[a-zA-Z]+;`)

// normalizeHTMLEntitiesForDiff applies html.UnescapeString recursively,
// but avoids double-unescaping or unescaping intentionally escaped XML entities.
func normalizeHTMLEntitiesForDiff(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		// If it's wrapped in <string> tags, strip them for evaluation, but preserve during output
		str := strings.TrimSpace(v)
		if strings.Contains(str, "&") {
			// Unescape once
			unescaped := html.UnescapeString(str)

			// If unescaping results in a valid single-level entity, don't double-unescape
			if strings.Contains(unescaped, "<") || strings.Contains(unescaped, ">") || safeHTMLEntity.MatchString(unescaped) {
				return str
			}

			// Catch common double-escape: &amp;amp; -> &amp;
			if strings.Contains(str, "&amp;") && !strings.Contains(str, "&amp;amp;") {
				return unescaped
			}
		}
		return str

	case map[string]interface{}:
		for key, val := range v {
			v[key] = normalizeHTMLEntitiesForDiff(val)
		}
	case []interface{}:
		for i, val := range v {
			v[i] = normalizeHTMLEntitiesForDiff(val)
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
