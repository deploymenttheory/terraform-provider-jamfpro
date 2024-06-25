package plistgenerator

import (
	"bytes"
	"fmt"

	"howett.net/plist"
)

// GeneratePlistFromPayloads generates a plist XML string from the given payloads map.
func GeneratePlistFromPayloads(payloads map[string]interface{}) (string, error) {
	// Convert the payloads map into a structure suitable for plist encoding
	convertedMap := convertMap(payloads)

	// Marshal the structure to plist XML
	var buffer bytes.Buffer
	encoder := plist.NewEncoder(&buffer)
	encoder.Indent("\t")
	err := encoder.Encode(convertedMap)
	if err != nil {
		return "", fmt.Errorf("failed to encode plist: %w", err)
	}
	return buffer.String(), nil
}

// convertMap recursively converts a map to a structure suitable for plist encoding
func convertMap(data map[string]interface{}) interface{} {
	converted := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			converted[key] = convertMap(v)
		case []interface{}:
			converted[key] = convertArray(v)
		default:
			converted[key] = value
		}
	}
	return converted
}

// convertArray recursively converts an array to a structure suitable for plist encoding
func convertArray(data []interface{}) interface{} {
	converted := make([]interface{}, len(data))
	for i, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			converted[i] = convertMap(v)
		case []interface{}:
			converted[i] = convertArray(v)
		default:
			converted[i] = value
		}
	}
	return converted
}
