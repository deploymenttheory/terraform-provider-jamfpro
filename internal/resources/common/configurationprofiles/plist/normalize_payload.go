// common/configurationprofiles/plist/payload.go
// Description: This file contains the ConfigurationProfile and PayloadContent structs, as well as functions for unmarshalling, marshalling, and validating plist payloads.
package plist

import (
	"howett.net/plist"
)

// ConfigurationProfile represents a root level MacOS configuration profile.
type ConfigurationProfile struct {
	PayloadDescription       string                 `mapstructure:"PayloadDescription"`
	PayloadDisplayName       string                 `mapstructure:"PayloadDisplayName" validate:"required"`
	PayloadEnabled           bool                   `mapstructure:"PayloadEnabled" validate:"required"`
	PayloadIdentifier        string                 `mapstructure:"PayloadIdentifier" validate:"required"`
	PayloadOrganization      string                 `mapstructure:"PayloadOrganization" validate:"required"`
	PayloadRemovalDisallowed bool                   `mapstructure:"PayloadRemovalDisallowed" validate:"required"`
	PayloadScope             string                 `mapstructure:"PayloadScope" validate:"required,oneof=System User Computer"`
	PayloadType              string                 `mapstructure:"PayloadType" validate:"required,eq=Configuration"`
	PayloadUUID              string                 `mapstructure:"PayloadUUID" validate:"required"`
	PayloadVersion           int                    `mapstructure:"PayloadVersion" validate:"required"`
	PayloadContent           []PayloadContent       `mapstructure:"PayloadContent"`
	Unexpected               map[string]interface{} `mapstructure:",remain"`
}

// ConfigurationPayload represents a nested MacOS configuration profile.
type PayloadContent struct {
	PayloadDescription  string                 `mapstructure:"PayloadDescription"`
	PayloadDisplayName  string                 `mapstructure:"PayloadDisplayName"`
	PayloadEnabled      bool                   `mapstructure:"PayloadEnabled"`
	PayloadIdentifier   string                 `mapstructure:"PayloadIdentifier"`
	PayloadOrganization string                 `mapstructure:"PayloadOrganization"`
	PayloadType         string                 `mapstructure:"PayloadType"`
	PayloadUUID         string                 `mapstructure:"PayloadUUID"`
	PayloadVersion      int                    `mapstructure:"PayloadVersion"`
	ConfigurationItems  map[string]interface{} `mapstructure:",remain"`
}

// NormalizePayloadState processes and normalizes a macOS Configuration Profile payload.
// This function is crucial for maintaining consistency in plist structures, especially
// when working with Terraform state management for Jamf Pro configuration profiles.
//
// The function performs the following steps:
//  1. Unmarshals the input payload (expected to be a plist XML string) into a generic map structure.
//  2. Normalizes the payload content using normalizePlistPayloadContent, which recursively processes
//     nested PayloadContent fields without altering the overall structure.
//  3. Marshals the normalized data back into a plist XML string.
//
// The function is designed to work with the plist library for unmarshalling and marshalling,
// avoiding the use of struct-based approaches (like mapstructure) that might inadvertently
// add or remove fields.
//
// Parameters:
//   - payload: Any type, expected to be a string containing a plist XML representation of a Configuration Profile.
//
// Returns:
//   - A string containing the normalized plist XML. If any error occurs during processing, an empty string is returned.
func NormalizePayloadState(payload any) string {
	var plistData map[string]interface{}
	_, err := plist.Unmarshal([]byte(payload.(string)), &plistData)
	if err != nil {
		return ""
	}

	normalizePlistPayloadContent(plistData)

	xml, err := plist.MarshalIndent(plistData, plist.XMLFormat, "\t")
	if err != nil {
		return ""
	}

	return trimTrailingWhitespace(string(xml))
}

// normalizePayloadContent recursively processes the PayloadContent field of a Configuration Profile plist.
// This function is crucial for maintaining consistency in plist structures when working with mapstructure.
//
// In the context of macOS Configuration Profiles, plists can have deeply nested PayloadContent fields.
// When using mapstructure to unmarshal and marshal these plists, there's a risk of losing or altering
// the structure of nested PayloadContent items. This function ensures that:
//
// 1. The structure of nested PayloadContent fields is preserved.
// 2. No additional fields are added or removed during the unmarshal/marshal process.
// 3. The integrity of the plist structure is maintained, especially for complex, multi-level profiles.
//
// Parameters:
//   - data: A map representing a portion of the plist structure, typically the root or a nested PayloadContent.
//
// The function modifies the data in place, recursively processing all levels of PayloadContent.
func normalizePlistPayloadContent(data map[string]interface{}) {
	if payloadContent, ok := data["PayloadContent"].([]interface{}); ok {
		for i, content := range payloadContent {
			if contentMap, ok := content.(map[string]interface{}); ok {
				if nestedContent, exists := contentMap["PayloadContent"]; exists {
					if nestedMap, ok := nestedContent.(map[string]interface{}); ok {
						normalizePlistPayloadContent(nestedMap)
					}
				}
				payloadContent[i] = contentMap
			}
		}
		data["PayloadContent"] = payloadContent
	}
}
