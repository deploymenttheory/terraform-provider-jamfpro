// macosconfigurationprofilesplist_diff_suppress.go
package macosconfigurationprofilesplist

import (
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func DiffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
	log.Printf("Suppressing diff for key: %s", k)

	processedOldPayload, err := processPayload(old, "Terraform state payload")
	if err != nil {
		log.Printf("Error processing old payload (Terraform state): %v", err)
		return false
	}

	processedNewPayload, err := processPayload(new, "Jamf Pro server payload")
	if err != nil {
		log.Printf("Error processing new payload (Jamf Pro server): %v", err)
		return false
	}

	oldHash := helpers.HashString(processedOldPayload)
	newHash := helpers.HashString(processedNewPayload)

	log.Printf("Old payload hash (Terraform state): %s\nOld payload (processed): %s", oldHash, processedOldPayload)
	log.Printf("New payload hash (Jamf Pro server): %s\nNew payload (processed): %s", newHash, processedNewPayload)

	return oldHash == newHash
}

// processPayload processes the payload by comparing the old and new payloads. It removes specified fields and compares the hashes.
func processPayload(payload string, source string) (string, error) {
	log.Printf("Processing %s: %s", source, payload)
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}
	processedPayload, err := plist.ProcessConfigurationProfileForDiffSuppression(payload, fieldsToRemove)
	if err != nil {
		return "", err
	}
	log.Printf("Processed %s: %s", source, processedPayload)
	return processedPayload, nil
}
