package configurationprofiles

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"howett.net/plist"
)

// ConfigurationProfile represents the structure of the plist data.
type ConfigurationProfile struct {
	PayloadContent     []PayloadContentListItem
	PayloadDisplayName string
	PayloadIdentifier  string
	PayloadType        string
	PayloadUuid        string
	PayloadVersion     int
	UnexpectedValues   map[string]interface{} `mapstructure:",remain"`
}

// PayloadContentListItem represents an individual payload item.
type PayloadContentListItem struct {
	PayloadDisplayName    string
	PayloadIdentifier     string
	PayloadType           string
	PayloadUuid           string
	PayloadVersion        int
	PayloadSpecificValues map[string]interface{} `mapstructure:",remain"`
}

// CleanMap removes specified keys from the map.
func cleanMap(data map[string]interface{}, keysToRemove []string) map[string]interface{} {
	for _, key := range keysToRemove {
		delete(data, key)
	}

	for k, v := range data {
		switch child := v.(type) {
		case map[string]interface{}:
			data[k] = cleanMap(child, keysToRemove)
		case []interface{}:
			for i, item := range child {
				if itemMap, ok := item.(map[string]interface{}); ok {
					child[i] = cleanMap(itemMap, keysToRemove)
				}
			}
		}
	}

	return data
}

// DecodePlistToMap decodes plist bytes into a map and cleans the data
func decodePlistToMap(plistBytes []byte, keysToRemove []string) (map[string]interface{}, error) {
	var unmarshalledPlist map[string]interface{}
	_, err := plist.Unmarshal(plistBytes, &unmarshalledPlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %v", err)
	}

	cleanedPlist := cleanMap(unmarshalledPlist, keysToRemove)
	return cleanedPlist, nil
}

// EncodeMapToPlist converts a map to a plist string
func encodeMapToPlist(data map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("    ")
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to encode plist: %v", err)
	}

	return buf.String(), nil
}

// NormalizePayload takes a plist string and returns a normalized version of the payload
func NormalizePayload(plistData string, keysToRemove []string) (string, error) {
	plistBytes := []byte(plistData)
	cleanedMap, err := decodePlistToMap(plistBytes, keysToRemove)
	if err != nil {
		return "", err
	}

	return encodeMapToPlist(cleanedMap)
}

// CompareMapsIgnoringOrder compares two maps, ignoring the order of keys
func compareMapsIgnoringOrder(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, val1 := range map1 {
		val2, ok := map2[key]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(val1, val2) {
			return false
		}
	}

	return true
}

// ComparePayloads compares two sets of payload-specific fields and returns true if they are equal
func ComparePayloads(payloads1, payloads2 string, keysToRemove []string) (bool, error) {
	normalizedPayload1, err := NormalizePayload(payloads1, keysToRemove)
	if err != nil {
		return false, err
	}

	normalizedPayload2, err := NormalizePayload(payloads2, keysToRemove)
	if err != nil {
		return false, err
	}

	map1, err := decodePlistToMap([]byte(normalizedPayload1), keysToRemove)
	if err != nil {
		return false, err
	}

	map2, err := decodePlistToMap([]byte(normalizedPayload2), keysToRemove)
	if err != nil {
		return false, err
	}

	return compareMapsIgnoringOrder(map1, map2), nil
}

// ConfigurationFilePlistToStructFromFile takes filepath of MacOS Configuration Profile .plist file and returns &ConfigurationProfile
func ConfigurationFilePlistToStructFromFile(filepath string) (*ConfigurationProfile, error) {
	plistFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer plistFile.Close()

	xmlData, err := io.ReadAll(plistFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read plist/xml file: %v", err)
	}

	return plistDataToStruct(xmlData)
}

// ConfigurationProfilePlistToStructFromString takes xml of MacOS Configuration Profile .plist file and returns &ConfigurationProfile
func ConfigurationProfilePlistToStructFromString(plistData string) (*ConfigurationProfile, error) {
	return plistDataToStruct([]byte(plistData))
}

// plistDataToStruct takes xml .plist bytes data and returns ConfigurationProfile
func plistDataToStruct(plistBytes []byte) (*ConfigurationProfile, error) {
	var unmarshalledPlist map[string]interface{}
	_, err := plist.Unmarshal(plistBytes, &unmarshalledPlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist/xml: %v", err)
	}

	var out ConfigurationProfile
	err = mapstructure.Decode(unmarshalledPlist, &out)
	if err != nil {
		return nil, fmt.Errorf("(mapstructure) failed to map unmarshaled configuration profile to struct: %v", err)
	}

	return &out, nil
}
