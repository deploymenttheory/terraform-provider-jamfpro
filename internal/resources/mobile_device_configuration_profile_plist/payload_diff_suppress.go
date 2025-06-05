// mobiledeviceconfigurationprofilesplist_diff_suppress.go
package mobiledeviceconfigurationprofilesplist

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func DiffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
	fmt.Printf("[DIFFSUPPRESS] Checking diff for key: %s\n", k)

	processedOldPayload, err := processPayload(old, "Terraform state payload")
	if err != nil {
		fmt.Printf("[DIFFSUPPRESS] Error processing old payload (Terraform state): %v\n", err)
		return false
	}

	processedNewPayload, err := processPayload(new, "Jamf Pro server payload")
	if err != nil {
		fmt.Printf("[DIFFSUPPRESS] Error processing new payload (Jamf Pro server): %v\n", err)
		return false
	}

	oldHash := common.HashString(processedOldPayload)
	newHash := common.HashString(processedNewPayload)

	fmt.Printf("[DIFFSUPPRESS] Old payload hash: %s\n", oldHash)
	fmt.Printf("[DIFFSUPPRESS] New payload hash: %s\n", newHash)

	return oldHash == newHash
}

// processPayload processes the payload by comparing the old and new payloads. It removes specified fields and compares the hashes.
func processPayload(payload string, source string) (string, error) {
	fmt.Printf("Processing %s: %s", source, payload)
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}
	processedPayload, err := plist.ProcessConfigurationProfileForDiffSuppression(payload, fieldsToRemove)
	if err != nil {
		return "", err
	}
	fmt.Printf("Processed %s: %s", source, processedPayload)
	return processedPayload, nil
}
