package plist

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
