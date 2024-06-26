package macosconfigurationprofilesplistgenerator

import (
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// DiffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func DiffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
	log.Printf("Suppressing diff for key: %s", k)

	processedOldPayload, err := processPayload(old, "Terraform state payload")
	if err != nil {
		log.Printf("Error processing old payload (Terraform state): %v", err)
		return false
	}

	processedNewPayload, err := processPayload(new, "Jamf Pro server payload")
	if err != nil {
		log.Printf("Error processing new payload (Jamf Pro server): %v", err)
		return false
	}

	equal := comparePayloads(processedOldPayload, processedNewPayload)

	log.Printf("Payloads equal: %v", equal)

	return equal
}

// processPayload processes the payload by removing specified fields.
func processPayload(payload string, source string) (map[string]interface{}, error) {
	log.Printf("Processing %s: %s", source, payload)
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}
	processedPayload, err := ProcessConfigurationProfileForDiffSuppression(payload, fieldsToRemove)
	if err != nil {
		return nil, err
	}
	log.Printf("Processed %s: %v", source, processedPayload)
	return processedPayload, nil
}

// comparePayloads recursively compares two payloads, ignoring specified fields.
func comparePayloads(oldPayload, newPayload map[string]interface{}) bool {
	return deepEqualIgnoreFields(oldPayload, newPayload)
}

// deepEqualIgnoreFields recursively compares two maps, ignoring specified fields.
func deepEqualIgnoreFields(oldMap, newMap map[string]interface{}) bool {
	for key, oldValue := range oldMap {
		if newValue, found := newMap[key]; found {
			if reflect.TypeOf(oldValue) != reflect.TypeOf(newValue) {
				return false
			}
			switch oldValueTyped := oldValue.(type) {
			case map[string]interface{}:
				newValueTyped := newValue.(map[string]interface{})
				if !deepEqualIgnoreFields(oldValueTyped, newValueTyped) {
					return false
				}
			case []interface{}:
				newValueTyped := newValue.([]interface{})
				if !sliceEqualIgnoreFields(oldValueTyped, newValueTyped) {
					return false
				}
			default:
				if !reflect.DeepEqual(oldValueTyped, newValue) {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

// sliceEqualIgnoreFields compares two slices, ignoring specified fields.
func sliceEqualIgnoreFields(oldSlice, newSlice []interface{}) bool {
	if len(oldSlice) != len(newSlice) {
		return false
	}
	for i := range oldSlice {
		oldValue := oldSlice[i]
		newValue := newSlice[i]
		if reflect.TypeOf(oldValue) != reflect.TypeOf(newValue) {
			return false
		}
		switch oldValueTyped := oldValue.(type) {
		case map[string]interface{}:
			newValueTyped := newValue.(map[string]interface{})
			if !deepEqualIgnoreFields(oldValueTyped, newValueTyped) {
				return false
			}
		default:
			if !reflect.DeepEqual(oldValueTyped, newValue) {
				return false
			}
		}
	}
	return true
}

// ProcessConfigurationProfileForDiffSuppression processes the plist data, removes specified fields, and returns the cleaned map.
func ProcessConfigurationProfileForDiffSuppression(plistData string, fieldsToRemove []string) (map[string]interface{}, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Decode and clean the plist data
	plistBytes := []byte(plistData)
	cleanedData, err := decodeAndCleanPlist(plistBytes, fieldsToRemove)
	if err != nil {
		log.Printf("Error decoding and cleaning plist data: %v\n", err)
		return nil, err
	}

	// Sort keys for consistent order
	sortedData := SortPlistKeys(cleanedData)

	log.Printf("Sorted plist data: %v\n", sortedData)

	return sortedData, nil
}

// decodeAndCleanPlist decodes a plist into a map and removes specified fields.
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

// RemoveFields removes specified fields from a nested map.
func RemoveFields(data map[string]interface{}, fieldsToRemove []string, path string) {
	// Create a set of fields to remove for quick lookup
	fieldsToRemoveSet := make(map[string]struct{}, len(fieldsToRemove))
	for _, field := range fieldsToRemove {
		fieldsToRemoveSet[field] = struct{}{}
	}

	// Recursively remove fields
	recursivelyRemoveFields(data, fieldsToRemoveSet, path)
}

// recursivelyRemoveFields removes specified fields from a nested map.
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

// SortPlistKeys sorts the keys of a plist map for consistent ordering.
func SortPlistKeys(data map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedData := make(map[string]interface{}, len(data))
	for _, key := range keys {
		sortedData[key] = data[key]
	}
	return sortedData
}
