// macosconfigurationprofilesplistgenerator_diff_suppress.go
package macos_configuration_profile_plist_generator

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressPayloads is a custom diff suppression function for the payloads attribute.
func DiffSuppressPayloads(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diff for specific keys
	if shouldSuppressKey(k) {
		log.Printf("Suppressing diff for key: %s", k)
		return true
	}
	return false
}

// shouldSuppressKey checks if a key should be suppressed in the diff.
func shouldSuppressKey(key string) bool {
	keysToSuppress := []string{
		"payload_identifier",
		"payload_uuid",
		"payload_version",
		"payload_organization",
		//"payload_enabled",
		//"payload_scope",
		"payload_display_name",
	}

	for _, suppressKey := range keysToSuppress {
		if keyContainsSuffix(key, suppressKey) {
			return true
		}
	}
	return false
}

// keyContainsSuffix checks if the key contains a specific suffix, accounting for Terraform's nested attribute syntax.
func keyContainsSuffix(key, suffix string) bool {
	return len(key) >= len(suffix) && key[len(key)-len(suffix):] == suffix
}
