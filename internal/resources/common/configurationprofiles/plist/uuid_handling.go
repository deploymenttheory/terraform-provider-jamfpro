package plist

import (
	"fmt"
	"log"
)

// ExtractUUIDs recursively traverses a plist structure represented as nested maps and slices,
// extracting all occurrences of `PayloadUUID` and associating them with their respective
// `PayloadDisplayName`. It stores these key-value pairs in the provided `uuidMap`.
// If a `PayloadDisplayName` is absent at the root level, it uses the special key "root".
// This function is typically used to map existing UUIDs from a configuration profile
// retrieved from Jamf Pro.
func ExtractUUIDs(data interface{}, uuidMap map[string]string) {
	log.Printf("[DEBUG] Extracting existing payload UUID's and PayloadDisplayName.")

	switch v := data.(type) {
	case map[string]interface{}:
		displayName, hasDisplayName := v["PayloadDisplayName"].(string)
		uuid, hasUUID := v["PayloadUUID"].(string)

		if hasDisplayName && hasUUID {
			uuidMap[displayName] = uuid
		} else if hasUUID {
			// For root level, use special key
			uuidMap["root"] = uuid
		}

		// Recursively process all values
		for _, val := range v {
			ExtractUUIDs(val, uuidMap)
		}
	case []interface{}:
		for _, item := range v {
			ExtractUUIDs(item, uuidMap)
		}
	}
}

// UpdateUUIDs recursively traverses a plist structure represented as nested maps and slices,
// updating the values of `PayloadUUID` and `PayloadIdentifier` fields using the UUIDs
// provided in `uuidMap`. It matches UUIDs based on `PayloadDisplayName`. If a `PayloadDisplayName`
// is absent at the root level, it uses the special key "root" from the map.
// This function ensures that configuration profile UUIDs remain consistent with Jamf Pro
// expectations during Terraform update operations.
func UpdateUUIDs(data interface{}, uuidMap map[string]string) {
	log.Printf("[DEBUG] Injecting Jamf Pro post creation configuration profile PayloadUUID and PayloadIdentifier.")

	switch v := data.(type) {
	case map[string]interface{}:
		displayName, hasDisplayName := v["PayloadDisplayName"].(string)
		if hasDisplayName {
			if uuid, exists := uuidMap[displayName]; exists {
				v["PayloadUUID"] = uuid
				v["PayloadIdentifier"] = uuid // Also update identifier
			}
		} else {
			// For root level
			if uuid, exists := uuidMap["root"]; exists {
				v["PayloadUUID"] = uuid
			}
		}

		// Recursively process all values
		for _, val := range v {
			UpdateUUIDs(val, uuidMap)
		}
	case []interface{}:
		for _, item := range v {
			UpdateUUIDs(item, uuidMap)
		}
	}
}

// ValidatePayloadUUIDsMatch recursively compares UUID-related fields (`PayloadUUID` and
// `PayloadIdentifier`) between two plist structures (`existingPlist` and `newPlist`) to
// confirm they match exactly. It accumulates any differences found into the provided
// `mismatches` slice, describing the exact path and mismatched values.
// This validation step ensures Terraform updates maintain consistency with Jamf Pro's
// UUID requirements and detects unintended modifications.
func ValidatePayloadUUIDsMatch(existingPlist, newPlist interface{}, path string, mismatches *[]string) {
	existingMap, existingOk := existingPlist.(map[string]interface{})
	newMap, newOk := newPlist.(map[string]interface{})

	if existingOk && newOk {
		for key, existingValue := range existingMap {
			newValue, exists := newMap[key]

			// Build the full path for clear logging
			currentPath := path + "/" + key

			if !exists {
				continue // Ignore keys that don't exist in the new payload
			}

			switch key {
			case "PayloadUUID", "PayloadIdentifier":
				existingUUID, existingIsString := existingValue.(string)
				newUUID, newIsString := newValue.(string)

				if existingIsString && newIsString && existingUUID != newUUID {
					*mismatches = append(*mismatches, fmt.Sprintf("%s (Jamf Pro: %s, Request: %s)", currentPath, existingUUID, newUUID))
				}
			default:
				ValidatePayloadUUIDsMatch(existingValue, newValue, currentPath, mismatches)
			}
		}
	} else if existingSlice, ok := existingPlist.([]interface{}); ok {
		if newSlice, newOk := newPlist.([]interface{}); newOk {
			minLen := len(existingSlice)
			if len(newSlice) < minLen {
				minLen = len(newSlice)
			}
			for i := 0; i < minLen; i++ {
				ValidatePayloadUUIDsMatch(existingSlice[i], newSlice[i], fmt.Sprintf("%s[%d]", path, i), mismatches)
			}
		}
	}

	// If this is the root level call (empty path indicates root) and no mismatches were found
	if path == "Payload" && len(*mismatches) == 0 {
		log.Printf("[DEBUG] No config profile UUID mismatches found between existing and new plist. Injection was successful.")
	}
}
