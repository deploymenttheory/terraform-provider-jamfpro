package configurationprofiles

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"howett.net/plist"
)

// ConfigurationProfile represents the structure of the plist data.
type ConfigurationProfile struct {
	PayloadContent     []map[string]interface{} `plist:"PayloadContent"`
	PayloadDescription string                   `plist:"PayloadDescription"`
	PayloadDisplayName string                   `plist:"PayloadDisplayName"`
	PayloadEnabled     bool                     `plist:"PayloadEnabled"`
	PayloadScope       string                   `plist:"PayloadScope"`
	PayloadType        string                   `plist:"PayloadType"`
	PayloadVersion     int                      `plist:"PayloadVersion"`
}

// ProcessConfigurationProfile processes the plist data, removes specified fields, and returns the cleaned plist XML as a string and its hash
func ProcessConfigurationProfile(plistData string, fieldsToRemove []string) (string, string, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Decode and clean the plist data
	plistBytes := []byte(plistData)
	log.Printf("Decoding plist data: %s\n", plistData)
	cleanedData, err := decodeAndCleanPlist(plistBytes, fieldsToRemove)
	if err != nil {
		log.Printf("Error decoding and cleaning plist data: %v\n", err)
		return "", "", err
	}

	// Encode the cleaned data back to plist XML format
	encodedPlist, err := EncodePlist(cleanedData)
	if err != nil {
		log.Printf("Error encoding cleaned data to plist: %v\n", err)
		return "", "", err
	}

	// Compute the hash of the encoded plist data
	payloadHash, err := HashPlistData(encodedPlist)
	if err != nil {
		log.Printf("Error computing hash of payload: %v\n", err)
		return "", "", err
	}

	log.Printf("Successfully processed configuration profile, hash: %s\n", payloadHash)
	return encodedPlist, payloadHash, nil
}

// Function to decode a plist into a map and remove specified fields
func decodeAndCleanPlist(plistData []byte, fieldsToRemove []string) (map[string]interface{}, error) {
	var rawData map[string]interface{}
	_, err := plist.Unmarshal(plistData, &rawData)
	if err != nil {
		log.Printf("Error unmarshalling plist data: %v\n", err)
		return nil, err
	}

	log.Printf("Raw plist data: %v\n", rawData)
	removeFields(rawData, fieldsToRemove)
	log.Printf("Cleaned plist data: %v\n", rawData)

	return rawData, nil
}

// Function to remove specified fields from a nested map
func removeFields(data map[string]interface{}, fieldsToRemove []string) {
	for _, field := range fieldsToRemove {
		log.Printf("Removing field: %s\n", field)
		delete(data, field)
	}

	for _, v := range data {
		switch v := v.(type) {
		case map[string]interface{}:
			removeFields(v, fieldsToRemove)
		case []interface{}:
			for _, item := range v {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					removeFields(nestedMap, fieldsToRemove)
				}
			}
		}
	}
}

// EncodePlist encodes a cleaned map back to plist XML format
func EncodePlist(cleanedData map[string]interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := plist.NewEncoder(&buffer)
	encoder.Indent("\t") // Optional: for pretty-printing the XML
	if err := encoder.Encode(cleanedData); err != nil {
		log.Printf("Error encoding plist data: %v\n", err)
		return "", err
	}

	return buffer.String(), nil
}

// HashPlistData computes the SHA-256 hash of the given plist data
func HashPlistData(plistData string) (string, error) {
	hash := sha256.Sum256([]byte(plistData))
	return hex.EncodeToString(hash[:]), nil
}
