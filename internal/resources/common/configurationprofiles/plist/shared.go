// common/configurationprofiles/plist/shared.go contains the shared functions to process configuration profiles.
package plist

import (
	"bytes"
	"log"
	"sort"
	"strings"

	"howett.net/plist"
)

// Function to decode a plist into a map without removing any fields
func DecodePlist(plistData []byte) (map[string]interface{}, error) {
	var rawData map[string]interface{}
	_, err := plist.Unmarshal(plistData, &rawData)
	if err != nil {
		log.Printf("Error unmarshalling plist data: %v\n", err)
		return nil, err
	}
	return rawData, nil
}

// EncodePlist encodes a cleaned map back to plist XML format
func EncodePlist(cleanedData map[string]interface{}) (string, error) {
	log.Printf("Encoding plist data: %v\n", cleanedData)
	var buffer bytes.Buffer
	encoder := plist.NewEncoder(&buffer)
	encoder.Indent("\t") // Optional: for pretty-printing the XML

	if err := encoder.Encode(cleanedData); err != nil {
		log.Printf("Error encoding plist data: %v\n", err)
		return "", err
	}

	encodedString := buffer.String()

	// Post-process to remove unnecessary escaped characters while keeping essential ones
	encodedString = strings.ReplaceAll(encodedString, "&#34;", "\"") // Fix double quotes

	return encodedString, nil
}

// SortPlistKeys recursively sorts the config profile xml keys of a nested map
// into alphabetical order,and sorts elements within arrays if they are strings or dictionaries.
// This function is used to prepare the xml plist keys for diff suppression and since
// there's no guranatee what keys will be present within the XML, nor their order presented,
// this function is used to ensure that the keys are in a consistent order for comparison.
func SortPlistKeys(data map[string]interface{}) map[string]interface{} {
	sortedData := make(map[string]interface{})
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	log.Printf("[DEBUG] Unsorted keys: %v\n", keys)
	sort.Strings(keys)
	log.Printf("[DEBUG] Sorted keys: %v\n", keys)
	for _, k := range keys {
		log.Printf("[DEBUG] Processing key: %s\n", k)
		switch v := data[k].(type) {
		case map[string]interface{}:
			log.Printf("[DEBUG] Key %s is a nested map, sorting nested keys...\n", k)
			sortedData[k] = SortPlistKeys(v)
		case []interface{}:
			log.Printf("[DEBUG] Key %s is an array, processing items...\n", k)
			sortedArray := make([]interface{}, len(v))

			// First check if all elements are strings
			allStrings := true
			for _, item := range v {
				if _, ok := item.(string); !ok {
					allStrings = false
					break
				}
			}

			if allStrings {
				// If all strings, create string array and sort
				stringArray := make([]string, len(v))
				for i, item := range v {
					stringArray[i] = item.(string)
				}
				sort.Strings(stringArray)
				for i, item := range stringArray {
					sortedArray[i] = item
				}
			} else {
				// Handle non-string arrays
				for i, item := range v {
					log.Printf("[DEBUG] Processing item %d of array %s\n", i, k)
					if nestedMap, ok := item.(map[string]interface{}); ok {
						sortedArray[i] = SortPlistKeys(nestedMap)
					} else {
						sortedArray[i] = item
					}
				}
			}
			sortedData[k] = sortedArray
		default:
			sortedData[k] = v
		}
	}
	return sortedData
}
