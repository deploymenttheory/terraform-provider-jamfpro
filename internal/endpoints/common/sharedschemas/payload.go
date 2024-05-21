package sharedschemas

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"howett.net/plist"
)

type ConfigurationProfile struct {
	PayloadDescription       string                 `mapstructure:"PayloadDescription"`
	PayloadDisplayName       string                 `mapstructure:"PayloadDisplayName" validate:"required"`
	PayloadIdentifier        string                 `mapstructure:"PayloadIdentifier" validate:"required"`
	PayloadOrganization      string                 `mapstructure:"PayloadOrganization" validate:"required"`
	PayloadRemovalDisallowed bool                   `mapstructe:"PayloadRemovalDisallowed" validate:"required"`
	PayloadScope             string                 `mapstructure:"PayloadScope" validate:"required,oneof=System User"`
	PayloadType              string                 `mapstructure:"PayloadType" validate:"required,eq=Configuration"`
	PayloadUUID              string                 `mapstructure:"PayloadUUID" validate:"required"`
	PayloadVersion           int                    `mapstructure:"PayloadVersion" validate:"required,eq=1"`
	AdditionalFields         map[string]interface{} `mapstructure:",remain"`
}

func GetSharedSchemaPayload() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		StateFunc:    NormalizePayloadState,
		ValidateFunc: ValidatePayload,
		Description:  "A MacOS configuration profile as a plist-formatted XML string.",
	}
}

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

func MarshalPayload(profile *ConfigurationProfile) (string, error) {
	mergedPayload := MergePayloadFieldsIntoMap(profile)
	xml, err := plist.MarshalIndent(mergedPayload, plist.XMLFormat, "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	return string(xml), nil
}

func MergePayloadFieldsIntoMap(profile *ConfigurationProfile) map[string]interface{} {
	merged := make(map[string]interface{}, len(profile.AdditionalFields))
	for k, v := range profile.AdditionalFields {
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

	return merged
}

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
	errs = validateFields(profile)

	return warns, errs
}

func validateFields(profile *ConfigurationProfile) []error {
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
		}
	}

	return errs
}

func SuppressPayloadDiff(k, old string, new string, d *schema.ResourceData) bool {
	return NormalizePayloadState(old) == NormalizePayloadState(new)
}
