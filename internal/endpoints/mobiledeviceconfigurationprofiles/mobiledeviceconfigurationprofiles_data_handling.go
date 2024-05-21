// mobiledeviceconfigurationprofiles_data_handling.go
package mobiledeviceconfigurationprofiles

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// diffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func diffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
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

	oldHash := hashString(processedOldPayload)
	newHash := hashString(processedNewPayload)

	log.Printf("Old payload hash: %s, New payload hash: %s", oldHash, newHash)

	return oldHash == newHash
}

// processPayload processes the payload using the configurationprofiles.ProcessConfigurationProfile function.
func processPayload(payload string) (string, error) {
	log.Printf("Processing payload: %s", payload)
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}
	processedPayload, err := configurationprofiles.ProcessConfigurationProfile(payload, fieldsToRemove)
	if err != nil {
		return "", err
	}
	log.Printf("Processed payload: %s", processedPayload)
	return processedPayload, nil
}

// hashString calculates the SHA-256 hash of a string and returns it as a hexadecimal string.
func hashString(s string) string {
	log.Printf("Hashing string: %s", s)
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("Computed hash: %s", hash)
	return hash
}
