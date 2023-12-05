// scripts_data_handling.go
package scripts

import "encoding/base64"

// encodeScriptContent encode the script content to base64
func encodeScriptContent(scriptContent string) string {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(scriptContent))
	return encodedContent
}

// decodeScriptContent decodes the script content from base64
func decodeScriptContent(encodedContent string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

// isBase64Encoded validates if a string in base64 encoded
func isBase64Encoded(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
