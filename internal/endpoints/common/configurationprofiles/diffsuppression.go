// common/configurationprofiles/sanitize.go contains the functions to process configuration profiles.
package configurationprofiles

import (
	"log"
	"strconv"
	"strings"

	"howett.net/plist"
)

// ProcessConfigurationProfileForDiffSuppression processes the plist data, removes specified fields, and returns the cleaned plist XML as a string.
func ProcessConfigurationProfileForDiffSuppression(plistData string, fieldsToRemove []string) (string, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Decode and clean the plist data
	plistBytes := []byte(plistData)
	log.Printf("Decoding plist data: %s\n", plistData)
	cleanedData, err := decodeAndCleanPlist(plistBytes, fieldsToRemove)
	if err != nil {
		log.Printf("Error decoding and cleaning plist data: %v\n", err)
		return "", err
	}

	log.Printf("Cleaned plist data: %v\n", cleanedData)

	// Sort keys for consistent order
	log.Println("Sorting keys for consistent order...")
	sortedData := SortPlistKeys(cleanedData)

	log.Printf("Sorted plist data: %v\n", sortedData)

	// Encode the cleaned and sorted data back to plist XML format
	encodedPlist, err := EncodePlist(sortedData)
	if err != nil {
		log.Printf("Error encoding cleaned data to plist: %v\n", err)
		return "", err
	}

	log.Printf("Successfully processed configuration profile\n")
	return encodedPlist, nil
}

// Function to decode a plist into a map and remove specified fields
func decodeAndCleanPlist(plistData []byte, fieldsToRemove []string) (map[string]interface{}, error) {
	var rawData map[string]interface{}
	_, err := plist.Unmarshal(plistData, &rawData)
	if err != nil {
		log.Printf("Error unmarshalling plist data: %v\n", err)
		return nil, err
	}

	log.Printf("Raw plist data: %v\n", rawData)
	RemoveFields(rawData, fieldsToRemove, "")
	log.Printf("Cleaned plist data: %v\n", rawData)

	return rawData, nil
}

// RemoveFields removes specified fields from a nested map
func RemoveFields(data map[string]interface{}, fieldsToRemove []string, path string) {
	// Create a set of fields to remove for quick lookup
	fieldsToRemoveSet := make(map[string]struct{}, len(fieldsToRemove))
	for _, field := range fieldsToRemove {
		fieldsToRemoveSet[field] = struct{}{}
	}

	// Recursively remove fields
	recursivelyRemoveFields(data, fieldsToRemoveSet, path)
}

// recursivelyRemoveFields removes specified fields from a nested map
func recursivelyRemoveFields(data map[string]interface{}, fieldsToRemoveSet map[string]struct{}, path string) {
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
			recursivelyRemoveFields(v, fieldsToRemoveSet, newPath)
		case []interface{}:
			for i, item := range v {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					log.Printf("Recursively removing fields in array at path: %s[%d]\n", newPath, i)
					recursivelyRemoveFields(nestedMap, fieldsToRemoveSet, newPath+strings.ReplaceAll(key, "/", "_")+strconv.Itoa(i))
				}
			}
			// Ensure empty arrays are preserved
			data[key] = v
		}
	}
}
