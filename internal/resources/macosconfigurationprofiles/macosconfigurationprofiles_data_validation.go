// macosconfigurationprofiles_data_validation.go
package macosconfigurationprofiles

import (
	"bytes"
	"encoding/xml"
)

// formatmacOSConfigurationProfileXMLPayload prepares the xml payload for upload into Jamf Pro
func formatmacOSConfigurationProfileXMLPayload(input string) (string, error) {
	// Decode the XML data
	var buffer bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewBufferString(input))
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "    ") // Set indentation here if needed

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		// Write the token to the buffer in a standard format
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	// Close the encoder to flush the buffer
	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
