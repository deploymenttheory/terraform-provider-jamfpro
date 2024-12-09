// hash.go
// This package contains shared / common hash functions
package common

import (
	"crypto/sha256"
	"fmt"
	"log"
)

// HashString calculates the SHA-256 hash of a string and returns it as a hexadecimal string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("Computed hash: %s", hash)
	return hash
}
