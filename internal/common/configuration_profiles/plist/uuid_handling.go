package plist

import (
	"fmt"
	"log"
)

// ExtractUUIDs recursively traverses a plist structure represented as nested maps and slices,
// extracting all occurrences of `PayloadUUID` and associating them with a composite key
// that combines PayloadDisplayName and PayloadType to handle duplicate display names.
// If a `PayloadDisplayName` is absent at the root level, it uses the special key "root".
// This function is typically used to map existing UUIDs from a configuration profile
// retrieved from Jamf Pro.
func ExtractUUIDs(data any, uuidMap map[string]string, isRoot bool) {
	extractUUIDsRecursive(data, uuidMap, isRoot, 0)
}

// extractUUIDsRecursive handles the recursive extraction with payload counting for unique keys
func extractUUIDsRecursive(data any, uuidMap map[string]string, isRoot bool, payloadIndex int) int {
	log.Printf("[DEBUG] Extracting existing payload UUIDs.")

	switch v := data.(type) {
	case map[string]any:
		uuid, hasUUID := v["PayloadUUID"].(string)
		displayName, hasDisplayName := v["PayloadDisplayName"].(string)
		payloadType, hasPayloadType := v["PayloadType"].(string)

		if hasUUID {
			if isRoot {
				uuidMap["root"] = uuid
				log.Printf("[DEBUG] Found root PayloadUUID: %s", uuid)
			} else {
				// Create a composite key to handle duplicate PayloadDisplayName values
				var key string
				if hasDisplayName && hasPayloadType {
					key = fmt.Sprintf("%s|%s|%d", displayName, payloadType, payloadIndex)
				} else if hasDisplayName {
					key = fmt.Sprintf("%s|%d", displayName, payloadIndex)
				} else {
					key = fmt.Sprintf("payload|%d", payloadIndex)
				}
				uuidMap[key] = uuid
				log.Printf("[DEBUG] Found inner PayloadUUID for key '%s': %s", key, uuid)
				payloadIndex++
			}
		}

		// Recurse through all values
		for _, val := range v {
			payloadIndex = extractUUIDsRecursive(val, uuidMap, false, payloadIndex)
		}

	case []any:
		for _, item := range v {
			payloadIndex = extractUUIDsRecursive(item, uuidMap, false, payloadIndex)
		}
	}

	return payloadIndex
}

// ExtractPayloadIdentifiers recursively traverses a plist structure to extract PayloadIdentifier
// values and associate them with a composite key for proper structure preservation
func ExtractPayloadIdentifiers(data any, identifierMap map[string]string, isRoot bool) {
	extractPayloadIdentifiersRecursive(data, identifierMap, isRoot, 0)
}

// extractPayloadIdentifiersRecursive handles the recursive extraction with payload counting for unique keys
func extractPayloadIdentifiersRecursive(data any, identifierMap map[string]string, isRoot bool, payloadIndex int) int {
	switch v := data.(type) {
	case map[string]any:
		identifier, hasIdentifier := v["PayloadIdentifier"].(string)
		displayName, hasDisplayName := v["PayloadDisplayName"].(string)
		payloadType, hasPayloadType := v["PayloadType"].(string)

		if hasIdentifier {
			if isRoot {
				identifierMap["root"] = identifier
				log.Printf("[DEBUG] Found root PayloadIdentifier: %s", identifier)
			} else {
				// Create a composite key to handle duplicate PayloadDisplayName values
				var key string
				if hasDisplayName && hasPayloadType {
					key = fmt.Sprintf("%s|%s|%d", displayName, payloadType, payloadIndex)
				} else if hasDisplayName {
					key = fmt.Sprintf("%s|%d", displayName, payloadIndex)
				} else {
					key = fmt.Sprintf("payload|%d", payloadIndex)
				}
				identifierMap[key] = identifier
				log.Printf("[DEBUG] Found inner PayloadIdentifier for key '%s': %s", key, identifier)
				payloadIndex++
			}
		}

		// Recurse through all values
		for _, val := range v {
			payloadIndex = extractPayloadIdentifiersRecursive(val, identifierMap, false, payloadIndex)
		}

	case []any:
		for _, item := range v {
			payloadIndex = extractPayloadIdentifiersRecursive(item, identifierMap, false, payloadIndex)
		}
	}

	return payloadIndex
}

// UpdateUUIDs recursively traverses a plist structure represented as nested maps and slices,
// updating the values of `PayloadUUID` and `PayloadIdentifier` fields using the UUIDs
// provided in `uuidMap` and `identifierMap`. It matches UUIDs based on composite keys.
// If a `PayloadDisplayName` is absent at the root level, it uses the special key "root" from the map.
// This function ensures that configuration profile UUIDs remain consistent with Jamf Pro
// expectations during Terraform update operations.
func UpdateUUIDs(data any, uuidMap map[string]string, identifierMap map[string]string, isRoot bool) {
	updateUUIDsRecursive(data, uuidMap, identifierMap, isRoot, 0)
}

// updateUUIDsRecursive handles the recursive update with payload counting for unique keys
func updateUUIDsRecursive(data any, uuidMap map[string]string, identifierMap map[string]string, isRoot bool, payloadIndex int) int {
	log.Printf("[DEBUG] Injecting Jamf Pro post creation configuration profile PayloadUUID and PayloadIdentifier.")

	switch v := data.(type) {
	case map[string]any:
		_, hasUUID := v["PayloadUUID"].(string)
		displayName, hasDisplayName := v["PayloadDisplayName"].(string)
		payloadType, hasPayloadType := v["PayloadType"].(string)

		// Only update root-level UUID if explicitly present in the map as "root"
		if isRoot {
			if uuid, exists := uuidMap["root"]; exists {
				v["PayloadUUID"] = uuid
			}
			if identifier, exists := identifierMap["root"]; exists {
				v["PayloadIdentifier"] = identifier
			}
		} else if hasUUID {
			// Create a composite key to match with extracted UUIDs
			var key string
			if hasDisplayName && hasPayloadType {
				key = fmt.Sprintf("%s|%s|%d", displayName, payloadType, payloadIndex)
			} else if hasDisplayName {
				key = fmt.Sprintf("%s|%d", displayName, payloadIndex)
			} else {
				key = fmt.Sprintf("payload|%d", payloadIndex)
			}

			if existingUUID, exists := uuidMap[key]; exists {
				v["PayloadUUID"] = existingUUID
				log.Printf("[DEBUG] Updated PayloadUUID for key '%s' to: %s", key, existingUUID)
			}
			if identifier, exists := identifierMap[key]; exists {
				v["PayloadIdentifier"] = identifier
				log.Printf("[DEBUG] Updated PayloadIdentifier for key '%s' to: %s", key, identifier)
			}
			payloadIndex++
		}

		// Recurse through all values
		for _, val := range v {
			payloadIndex = updateUUIDsRecursive(val, uuidMap, identifierMap, false, payloadIndex)
		}

	case []any:
		for _, item := range v {
			payloadIndex = updateUUIDsRecursive(item, uuidMap, identifierMap, false, payloadIndex)
		}
	}

	return payloadIndex
}

// ValidatePayloadUUIDsMatch recursively compares UUID-related fields (`PayloadUUID` and
// `PayloadIdentifier`) between two plist structures (`existingPlist` and `newPlist`) to
// confirm they match exactly. It accumulates any differences found into the provided
// `mismatches` slice, describing the exact path and mismatched values.
// This validation step ensures Terraform updates maintain consistency with Jamf Pro's
// UUID requirements and detects unintended modifications.
func ValidatePayloadUUIDsMatch(existingPlist, newPlist any, path string, mismatches *[]string) {
	existingMap, existingOk := existingPlist.(map[string]any)
	newMap, newOk := newPlist.(map[string]any)

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
	} else if existingSlice, ok := existingPlist.([]any); ok {
		if newSlice, newOk := newPlist.([]any); newOk {
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
