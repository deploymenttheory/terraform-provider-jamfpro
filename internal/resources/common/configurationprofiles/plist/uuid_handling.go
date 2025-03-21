package plist

import "fmt"

// extractUUIDs recursively extracts config profile UUIDs from a plist structure
// and stores them in a map by PayloadDisplayName.
func ExtractUUIDs(data interface{}, uuidMap map[string]string) {
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

// updateUUIDs recursively updates config profile UUIDs in a
// plist structure
func UpdateUUIDs(data interface{}, uuidMap map[string]string) {
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
}
