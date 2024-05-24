// hash.go
// This package contains shared / common hash functions
package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HashValue takes a plaintext string and returns a SHA-256 hash of it as a hex string.
func HashValue(plaintext string) string {
	hash := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(hash[:])
}

// HashAndUpdateSensitiveField hashes the given sensitive value and updates the Terraform state if the hash is different.
func HashAndUpdateSensitiveField(d *schema.ResourceData, fieldKey string, configValue string) error {
	// Hash the sensitive value from the configuration
	hashedConfigValue := HashValue(configValue)

	// Get the current hashed value from the state
	hashedStateValue, exists := d.GetOk(fieldKey)

	// Update the state only if the new hash is different from the existing one or if the existing one doesn't exist
	if !exists || hashedConfigValue != hashedStateValue {
		if err := d.Set(fieldKey, hashedConfigValue); err != nil {
			return err
		}
	}

	return nil
}

// HashString calculates the SHA-256 hash of a string and returns it as a hexadecimal string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("Computed hash: %s", hash)
	return hash
}
