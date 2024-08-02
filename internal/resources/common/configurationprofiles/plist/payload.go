// common/configurationprofiles/plist/payload.go
// Description: This file contains the ConfigurationProfile and PayloadContent structs, as well as functions for unmarshalling, marshalling, and validating plist payloads.
package plist

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
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
	PayloadScope             string           `mapstructure:"PayloadScope" validate:"required,eq=System=User=Computer"`
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
	PayloadScope        string `mapstructure:"PayloadScope"`

	// Variable
	ConfigurationItems map[string]interface{} `mapstructure:",remain"`
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

	// Anticipate Jamf's default PayloadDescription.
	if profile.PayloadDescription == "" && len(profile.PayloadContent) >= 1 {
		merged["PayloadDescription"] = fmt.Sprintf("Configuration settings for the %s preference domain.", mergedPayloads[0]["PayloadType"])
	}

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

// NormalizePayloadState normalizes a payload state by unmarshalling and remarshal it.
func NormalizePayloadState(payload any) string {
	profile, err := UnmarshalPayload(payload.(string))
	if err != nil {
		return ""
	}

	xml, err := MarshalPayload(profile)
	if err != nil {
		return ""
	}

	return xml
}

// ValidatePayload validates a payload by unmarshalling it and checking for required fields.
func ValidatePayload(payload interface{}, key string) (warns []string, errs []error) {
	profile, err := UnmarshalPayload(payload.(string))
	if err != nil {
		errs = append(errs, err)
		return warns, errs
	}

	if profile.PayloadIdentifier != profile.PayloadUUID {
		warns = append(warns, "Top-level PayloadIdentifier should match top-level PayloadUUID")
	}

	// Custom validation
	errs = ValidatePayloadFields(profile)

	return warns, errs
}

// ValidatePayloadFields validates the fields of a ConfigurationProfile struct.
func ValidatePayloadFields(profile *ConfigurationProfile) []error {
	var errs []error

	// Iterate over struct fields
	val := reflect.ValueOf(profile).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")
		if tag != "" {
			// Check for required fields
			if strings.Contains(tag, "required") {
				value := val.Field(i).Interface()
				if value == "" {
					errs = append(errs, errors.New(fmt.Sprintf("Field '%s' is required", field.Name)))
				}
			}
			// Additional validation rules can be added here
			for _, tagField := range strings.Split(tag, ",") {
				if strings.Contains(tagField, "eq=") {
					value := val.Field(i).String()
					reqValues := strings.Split(tagField, "=")[1:]
					if !slices.Contains(reqValues, value) {
						errs = append(errs, errors.New(fmt.Sprintf("Field '%s' must be one of %v. Is %v.", field.Name, reqValues, value)))
					}
				}
			}
		}
	}

	return errs
}
