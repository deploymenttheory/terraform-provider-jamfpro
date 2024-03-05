// macosconfigurationprofiles_helpers.go
package macosconfigurationprofiles

import (
	"bytes"
	"encoding/xml"
	"io"
)

// sanitizePayloadXML removes specific fields from the XML payload.
func sanitizePayloadXML(xmlString string) (string, error) {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlString))
	var buffer bytes.Buffer
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "    ")

	var skipElement bool
	var skipDepth int

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if skipElement {
				skipDepth++
				continue
			}
			if isFieldToRemove(se.Name.Local) {
				skipElement = true
				skipDepth = 1
				continue
			}
			err = encoder.EncodeToken(se)
			if err != nil {
				return "", err
			}
		case xml.EndElement:
			if skipElement {
				skipDepth--
				if skipDepth == 0 {
					skipElement = false
				}
				continue
			}
			err = encoder.EncodeToken(se)
			if err != nil {
				return "", err
			}
		default:
			if !skipElement {
				err = encoder.EncodeToken(token)
				if err != nil {
					return "", err
				}
			}
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// isFieldToRemove checks if a field should be removed from the payload. This function handles scenario's where
// configuration profiles (.mobileconfig) include jamf pro instance specific keys and values that should be
// sanitized. This ensure that in the tf state that only the config payload and not jamf pro instance specific
// fields are stored.
func isFieldToRemove(fieldName string) bool {
	fieldsToRemove := []string{"PayloadUUID", "PayloadOrganization", "PayloadIdentifier"}
	for _, field := range fieldsToRemove {
		if field == fieldName {
			return true
		}
	}
	return false
}
