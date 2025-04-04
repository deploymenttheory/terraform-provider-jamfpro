package userinitiatedenrollment

import (
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// getLanguageCodesMap fetches language codes from the Jamf Pro API
// and creates a case-insensitive map of language names to their codes
func getLanguageCodesMap(client *jamfpro.Client) (map[string]string, error) {
	// Get language codes from the API
	languageCodes, err := client.GetEnrollmentLanguageCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch language codes from Jamf Pro API: %v", err)
	}

	// Create a case-insensitive map for language name to code mapping
	codeMap := make(map[string]string)

	// Process each language code returned by the API
	for _, code := range languageCodes {
		if code.Name != "" && code.Value != "" {
			// Add lowercase version for case-insensitive matching
			codeMap[strings.ToLower(strings.TrimSpace(code.Name))] = code.Value
		}
	}

	log.Printf("[DEBUG] Loaded %d language codes from Jamf Pro API", len(languageCodes))
	return codeMap, nil
}
