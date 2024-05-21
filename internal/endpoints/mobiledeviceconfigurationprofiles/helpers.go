package mobiledeviceconfigurationprofiles

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// CalculatePayloadHash calculates the SHA-256 hash of the given payload data
func CalculatePayloadHash(payload string) (string, error) {
	if payload == "" {
		return "", errors.New("payload is empty")
	}
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:]), nil
}
