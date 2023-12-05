// scripts_data_handling.go
package scripts

import (
	"encoding/base64"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// encodeScriptContent encode the script content to base64
func encodeScriptContent(scriptContent string) string {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(scriptContent))
	return encodedContent
}

// suppressBase64EncodedScriptDiff is a DiffSuppressFunc used to suppress differences
// between the state and the configuration when the script contents match after base64 encoding.
func suppressBase64EncodedScriptDiff(k, old, new string, d *schema.ResourceData) bool {
	// Encode the 'new' value to base64
	newEncoded := base64.StdEncoding.EncodeToString([]byte(new))

	// Log for debugging
	log.Printf("[DEBUG] suppressBase64EncodedScriptDiff - Encoded 'new' value: %s", newEncoded)
	log.Printf("[DEBUG] suppressBase64EncodedScriptDiff - 'old' value from state: %s", old)

	// Compare the base64 encoded 'new' value with the 'old' value
	isDiffSuppressed := newEncoded == old
	log.Printf("[DEBUG] suppressBase64EncodedScriptDiff - Is diff suppressed: %t", isDiffSuppressed)

	return isDiffSuppressed
}
