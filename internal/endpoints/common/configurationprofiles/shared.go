// common/configurationprofiles/shared.go contains the shared functions to process configuration profiles.
package configurationprofiles

import (
	"bytes"
	"log"
	"sort"

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

	log.Printf("Decoded plist data: %v\n", rawData)
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
	return buffer.String(), nil
}

// SortPlistKeys recursively sorts the keys of a nested map in alphabetical order,
// and sorts elements within arrays if they are strings or dictionaries.
func SortPlistKeys(data map[string]interface{}) map[string]interface{} {
	sortedData := make(map[string]interface{})
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	log.Printf("Unsorted keys: %v\n", keys)
	sort.Strings(keys)
	log.Printf("Sorted keys: %v\n", keys)
	for _, k := range keys {
		log.Printf("Processing key: %s\n", k)
		switch v := data[k].(type) {
		case map[string]interface{}:
			log.Printf("Key %s is a nested map, sorting nested keys...\n", k)
			sortedData[k] = SortPlistKeys(v)
		case []interface{}:
			log.Printf("Key %s is an array, processing items...\n", k)
			sortedArray := make([]interface{}, len(v))
			for i, item := range v {
				log.Printf("Processing item %d of array %s\n", i, k)
				if nestedMap, ok := item.(map[string]interface{}); ok {
					sortedArray[i] = SortPlistKeys(nestedMap)
				} else {
					sortedArray[i] = item
				}
			}
			// Check if the array elements are strings and sort them
			if len(sortedArray) > 0 {
				switch sortedArray[0].(type) {
				case string:
					stringArray := make([]string, len(sortedArray))
					for i, item := range sortedArray {
						stringArray[i] = item.(string)
					}
					sort.Strings(stringArray)
					for i, item := range stringArray {
						sortedArray[i] = item
					}
				case map[string]interface{}:
					for i, item := range sortedArray {
						if nestedMap, ok := item.(map[string]interface{}); ok {
							sortedArray[i] = SortPlistKeys(nestedMap)
						}
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
