// common/configurationprofiles/state.go contains the functions to process configuration profiles.
package configurationprofiles

import (
	"log"

	"howett.net/plist"
)

// ProcessConfigurationProfileForState processes the plist data, removes specified fields, and returns the cleaned plist XML as a string.
func ProcessConfigurationProfileForState(plistData string, fieldsToRemove []string) (string, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Decode and clean the plist data
	plistBytes := []byte(plistData)
	log.Printf("Decoding plist data: %s\n", plistData)
	decodedPlist, err := decodePlist(plistBytes)
	if err != nil {
		log.Printf("Error decoding plist data: %v\n", err)
		return "", err
	}

	// Sort keys for consistent order
	log.Println("Sorting keys for consistent order...")
	sortedData := sortKeys(decodedPlist)

	// Encode the cleaned and sorted data back to plist XML format
	encodedPlist, err := EncodePlist(sortedData)
	if err != nil {
		log.Printf("Error encoding cleaned data to plist: %v\n", err)
		return "", err
	}

	log.Printf("Successfully processed configuration profile\n")
	return encodedPlist, nil
}

// Function to decode a plist into a map without removing any fields
func decodePlist(plistData []byte) (map[string]interface{}, error) {
	var rawData map[string]interface{}
	_, err := plist.Unmarshal(plistData, &rawData)
	if err != nil {
		log.Printf("Error unmarshalling plist data: %v\n", err)
		return nil, err
	}

	log.Printf("Decoded plist data: %v\n", rawData)
	return rawData, nil
}
