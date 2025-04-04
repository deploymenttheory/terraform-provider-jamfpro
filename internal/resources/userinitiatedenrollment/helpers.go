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

// findLanguagesToDelete identifies language codes to delete
func findLanguagesToDelete(oldMessaging, newMessaging []interface{}, client *jamfpro.Client) ([]string, error) {
	var languagesToDelete []string

	// Create a map of language names in the new set
	newLanguageNames := make(map[string]bool)
	for _, messaging := range newMessaging {
		msg := messaging.(map[string]interface{})
		langName := strings.ToLower(strings.TrimSpace(msg["language_name"].(string)))
		newLanguageNames[langName] = true
	}

	// Find languages in old but not in new
	for _, messaging := range oldMessaging {
		msg := messaging.(map[string]interface{})
		langName := strings.ToLower(strings.TrimSpace(msg["language_name"].(string)))

		// Skip if this language is still in the new set
		if newLanguageNames[langName] {
			continue
		}

		// Get the language code from state
		if code, ok := msg["language_code"].(string); ok && code != "" {
			// Don't delete English - it's required by Jamf Pro
			if code == "en" {
				log.Printf("[WARN] Attempted to delete required English language (en), skipping")
				continue
			}
			languagesToDelete = append(languagesToDelete, code)
		}
	}

	return languagesToDelete, nil
}
