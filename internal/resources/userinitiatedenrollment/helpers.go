package userinitiatedenrollment

import (
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// getLanguageCodesMap fetches language codes from the Jamf Pro API
// processes each language code adds lowercase version for case-insensitive matching
// and creates a map of language names to their codes
//

func getLanguageCodesMap(client *jamfpro.Client) (map[string]string, error) {
	languageCodes, err := client.GetEnrollmentLanguageCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch language codes from Jamf Pro API: %v", err)
	}

	codeMap := make(map[string]string)

	for _, code := range languageCodes {
		if code.Name != "" && code.Value != "" {
			codeMap[strings.ToLower(strings.TrimSpace(code.Name))] = code.Value
		}
	}

	log.Printf("[DEBUG] Loaded %d language codes from Jamf Pro API", len(languageCodes))
	return codeMap, nil
}
