// common/configurationprofiles/state.go contains the functions to process configuration profiles.
package configurationprofiles

import (
	"log"
)

// ProcessConfigurationProfileForState processes the plist data, removes specified fields, and returns the cleaned plist XML as a string.
func ProcessConfigurationProfileForState(plistData string) (string, error) {
	log.Println("Starting ProcessConfigurationProfile")

	// Decode and clean the plist data
	plistBytes := []byte(plistData)
	log.Printf("Decoding plist data: %s\n", plistData)
	decodedPlist, err := DecodePlist(plistBytes)
	if err != nil {
		log.Printf("Error decoding plist data: %v\n", err)
		return "", err
	}

	// Sort keys for consistent order
	log.Println("Sorting keys for consistent order...")
	sortedData := SortPlistKeys(decodedPlist)

	// Encode the cleaned and sorted data back to plist XML format
	encodedPlist, err := EncodePlist(sortedData)
	if err != nil {
		log.Printf("Error encoding cleaned data to plist: %v\n", err)
		return "", err
	}

	log.Printf("Successfully processed configuration profile\n")
	return encodedPlist, nil
}
