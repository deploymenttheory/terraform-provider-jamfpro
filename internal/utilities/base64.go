// utilities.go
// For utility/helper functions to support the jamf pro tf provider
package utils

import "encoding/base64"

// Base64EncodeIfNot encodes a string value if it isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func Base64EncodeIfNot(value string) string {
	// Check whether the value is already Base64 encoded; don't double-encode
	if base64IsEncoded(value) {
		return value
	}

	// Base64 encode the value and return
	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	return encoded
}

func base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}
