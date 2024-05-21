package configurationprofiles

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"

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

// ConfigurationFilePlistToStructFromFile takes filepath of MacOS Configuration Profile .plist file and returns &ConfigurationProfile
func ConfigurationFilePlistToStructFromFile(filepath string) (*ConfigurationProfile, error) {
	log.Printf("[INFO] Reading plist file from path: %s", filepath)
	plistFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer plistFile.Close()

	xmlData, err := io.ReadAll(plistFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read plist/xml file: %v", err)
	}

	log.Printf("[DEBUG] Read plist file data: %s", string(xmlData))
	return plistDataToStruct(xmlData)
}

// ConfigurationProfilePlistToStructFromString takes xml of MacOS Configuration Profile .plist file and returns &ConfigurationProfile
func ConfigurationProfilePlistToStructFromString(plistData string) (*ConfigurationProfile, error) {
	log.Printf("[INFO] Parsing plist data from string")
	return plistDataToStruct([]byte(plistData))
}

// plistDataToStruct takes xml .plist bytes data and returns ConfigurationProfile
func plistDataToStruct(plistBytes []byte) (*ConfigurationProfile, error) {
	log.Printf("[DEBUG] Parsing plist data to struct: %s", string(plistBytes))
	var unmarshalledPlist map[string]interface{}
	_, err := plist.Unmarshal(plistBytes, &unmarshalledPlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist/xml: %v", err)
	}

	log.Printf("[DEBUG] Unmarshalled plist data: %v", unmarshalledPlist)
	var out ConfigurationProfile
	err = mapstructure.Decode(unmarshalledPlist, &out)
	if err != nil {
		return nil, fmt.Errorf("(mapstructure) failed to map unmarshaled configuration profile to struct: %v", err)
	}

	log.Printf("[INFO] Parsed ConfigurationProfile: %+v", out)
	return &out, nil
}

// NormalizePayload takes a plist string and returns a normalized version of the payload
func NormalizePayload(plistData string, keysToRemove []string) (string, error) {
	log.Printf("[INFO] Normalizing payload: %s", plistData)
	plistBytes := []byte(plistData)
	cleanedMap, err := decodePlistToMap(plistBytes, keysToRemove)
	if err != nil {
		return "", err
	}

	normalizedPlist, err := encodeMapToPlist(cleanedMap)
	if err != nil {
		return "", err
	}

	log.Printf("[INFO] Normalized payload: %s", normalizedPlist)
	return normalizedPlist, nil
}

// ComparePayloads compares two sets of payload-specific fields and returns true if they are equal
func ComparePayloads(payloads1, payloads2 string, keysToRemove []string) (bool, error) {
	log.Printf("[INFO] Comparing payloads: %s, %s", payloads1, payloads2)
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

	areEqual := compareMapsIgnoringOrder(map1, map2)
	log.Printf("[INFO] Payload comparison result: %v", areEqual)
	return areEqual, nil
}

// DecodePlistToMap decodes plist bytes into a map and cleans the data
func decodePlistToMap(plistBytes []byte, keysToRemove []string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Decoding plist bytes: %s", string(plistBytes))
	var unmarshalledPlist map[string]interface{}
	_, err := plist.Unmarshal(plistBytes, &unmarshalledPlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %v", err)
	}

	log.Printf("[DEBUG] Unmarshalled plist: %v", unmarshalledPlist)
	cleanedPlist := cleanMap(unmarshalledPlist, keysToRemove)
	return cleanedPlist, nil
}

// EncodeMapToPlist converts a map to a plist string
func encodeMapToPlist(data map[string]interface{}) (string, error) {
	log.Printf("[DEBUG] Encoding map to plist: %v", data)
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("    ")
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to encode plist: %v", err)
	}

	log.Printf("[DEBUG] Encoded plist: %s", buf.String())
	return buf.String(), nil
}

// CleanMap removes specified keys from the map.
func cleanMap(data map[string]interface{}, keysToRemove []string) map[string]interface{} {
	log.Printf("[DEBUG] Cleaning map: %v, keys to remove: %v", data, keysToRemove)
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

	log.Printf("[DEBUG] Cleaned map: %v", data)
	return data
}

// SortMapKeys sorts the keys of a map to ensure consistent ordering
func sortMapKeys(data map[string]interface{}) map[string]interface{} {
	log.Printf("[DEBUG] Sorting map keys for: %v", data)
	sortedMap := make(map[string]interface{})
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sortedMap[key] = data[key]
		if child, ok := data[key].(map[string]interface{}); ok {
			sortedMap[key] = sortMapKeys(child)
		}
	}
	log.Printf("[DEBUG] Sorted map: %v", sortedMap)
	return sortedMap
}

// CompareMapsIgnoringOrder compares two maps, ignoring the order of keys
func compareMapsIgnoringOrder(map1, map2 map[string]interface{}) bool {
	log.Printf("[DEBUG] Comparing maps ignoring order: %v, %v", map1, map2)
	map1 = sortMapKeys(map1)
	map2 = sortMapKeys(map2)
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
