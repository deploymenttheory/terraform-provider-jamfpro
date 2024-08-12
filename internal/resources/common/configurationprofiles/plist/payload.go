// common/configurationprofiles/plist/payload.go
// Description: This file contains the ConfigurationProfile and PayloadContent structs, as well as functions for unmarshalling, marshalling, and validating plist payloads.
package plist

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"howett.net/plist"
)

// ConfigurationProfile represents a root level MacOS configuration profile.
type ConfigurationProfile struct {
	// Standard / Expected
	PayloadDescription       string           `mapstructure:"PayloadDescription"`
	PayloadDisplayName       string           `mapstructure:"PayloadDisplayName" validate:"required"`
	PayloadEnabled           bool             `mapstructure:"PayloadEnabled" validate:"required"`
	PayloadIdentifier        string           `mapstructure:"PayloadIdentifier" validate:"required"`
	PayloadOrganization      string           `mapstructure:"PayloadOrganization" validate:"required"`
	PayloadRemovalDisallowed bool             `mapstructure:"PayloadRemovalDisallowed" validate:"required"`
	PayloadScope             string           `mapstructure:"PayloadScope" validate:"required,oneof=System User Computer"`
	PayloadType              string           `mapstructure:"PayloadType" validate:"required,eq=Configuration"`
	PayloadUUID              string           `mapstructure:"PayloadUUID" validate:"required"`
	PayloadVersion           int              `mapstructure:"PayloadVersion" validate:"required"`
	PayloadContent           []PayloadContent `mapstructure:"PayloadContent"`
	// Catch all for unexpected fields
	Unexpected map[string]interface{} `mapstructure:",remain"`
}

// ConfigurationPayload represents a nested MacOS configuration profile.
type PayloadContent struct {
	// Standard / Expected
	PayloadDescription  string `mapstructure:"PayloadDescription"`
	PayloadDisplayName  string `mapstructure:"PayloadDisplayName"`
	PayloadEnabled      bool   `mapstructure:"PayloadEnabled"`
	PayloadIdentifier   string `mapstructure:"PayloadIdentifier"`
	PayloadOrganization string `mapstructure:"PayloadOrganization"`
	PayloadType         string `mapstructure:"PayloadType"`
	PayloadUUID         string `mapstructure:"PayloadUUID"`
	PayloadVersion      int    `mapstructure:"PayloadVersion"`
	// Variable
	ConfigurationItems map[string]interface{} `mapstructure:",remain"`
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

	return string(xml)
}

// UnmarshalPayload unmarshals a plist payload into a ConfigurationProfile struct using mapstructure.
func UnmarshalPayload(payload string) (*ConfigurationProfile, error) {
	var profile map[string]interface{}
	_, err := plist.Unmarshal([]byte(payload), &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	var out ConfigurationProfile
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &out,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}
	err = decoder.Decode(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	return &out, nil
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

// MarshalPayload marshals a ConfigurationProfile struct into a plist payload using mapstructure.
func MarshalPayload(profile *ConfigurationProfile) (string, error) {
	mergedPayload := MergeConfigurationProfileFieldsIntoMap(profile)
	xml, err := plist.MarshalIndent(mergedPayload, plist.XMLFormat, "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	return string(xml), nil
}

// MergeConfigurationProfileFieldsIntoMap merges the fields of a ConfigurationProfile struct into a map.
func MergeConfigurationProfileFieldsIntoMap(profile *ConfigurationProfile) map[string]interface{} {
	merged := make(map[string]interface{}, len(profile.Unexpected))
	for k, v := range profile.Unexpected {
		merged[k] = v
	}

	val := reflect.ValueOf(profile).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		mapKey := field.Tag.Get("mapstructure")
		if mapKey != "" && mapKey != ",remain" {
			merged[mapKey] = val.Field(i).Interface()
		}
	}

	mergedPayloads := make([]map[string]interface{}, len(profile.PayloadContent))
	for k, v := range profile.PayloadContent {
		mergedPayloads[k] = MergeConfigurationPayloadFieldsIntoMap(&v)
	}

	merged["PayloadContent"] = mergedPayloads

	return merged
}

// MergeConfigurationPayloadFieldsIntoMap merges the fields of a ConfigurationPayload struct into a map.
func MergeConfigurationPayloadFieldsIntoMap(payload *PayloadContent) map[string]interface{} {
	merged := make(map[string]interface{}, len(payload.ConfigurationItems))
	for k, v := range payload.ConfigurationItems {
		merged[k] = v
	}

	val := reflect.ValueOf(payload).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		mapKey := field.Tag.Get("mapstructure")
		if mapKey != "" && mapKey != ",remain" && !strings.Contains(mapKey, ",-") {
			merged[mapKey] = val.Field(i).Interface()
		}
	}

	return merged
}
