// mobiledeviceconfigurationprofilesplist_diff_suppress.go
package mobiledeviceconfigurationprofilesplist

import (
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/hash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func DiffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
	log.Printf("Suppressing diff for key: %s", k)

	processedOldPayload, err := processPayload(old)
	if err != nil {
		log.Printf("Error processing old payload: %v", err)
		return false
	}

	processedNewPayload, err := processPayload(new)
	if err != nil {
		log.Printf("Error processing new payload: %v", err)
		return false
	}

	oldHash := hash.HashString(processedOldPayload)
	newHash := hash.HashString(processedNewPayload)

	log.Printf("Old payload hash: %s, New payload hash: %s", oldHash, newHash)

	return oldHash == newHash
}

// processPayload processes the payload by comparing the old and new payloads. It removes specified fields and compares the hashes.
func processPayload(payload string) (string, error) {
	log.Printf("Processing payload: %s", payload)
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}
	processedPayload, err := plist.ProcessConfigurationProfileForDiffSuppression(payload, fieldsToRemove)
	if err != nil {
		return "", err
	}
	log.Printf("Processed payload: %s", processedPayload)
	return processedPayload, nil
}
