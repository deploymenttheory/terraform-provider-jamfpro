// // common/configurationprofiles/plist/conversion.go
// // Description: This file contains the functions to convert the HCL data to a plist XML string and vice versa.

package plist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"howett.net/plist"
)

// ConvertHCLToPlist builds a plist from the Terraform HCL schema data
// Used by plist generator resource to convert HCL data to plist
func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	profile, err := mapSchemaToProfile(d)
	if err != nil {
		return "", err
	}
	plistData, err := MarshalPayload(profile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal plist: %w", err)
	}

	prettyPlistXML, err := plist.MarshalIndent(plistData, plist.XMLFormat, "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile to pretty plist: %v", err)
	}
	unescapedPrettyPlistXML := html.UnescapeString(string(prettyPlistXML))

	log.Printf("[DEBUG] Constructed Plist XML from HCL serialization:\n%s\n", unescapedPrettyPlistXML)

	return string(plistData), nil
}

// mapSchemaToProfile maps the Terraform schema data to the ConfigurationProfile struct
func mapSchemaToProfile(d *schema.ResourceData) (*ConfigurationProfile, error) {
	uuidStr := uuid.New().String()

	// Root Level
	out := &ConfigurationProfile{
		PayloadDescription:       d.Get("payloads.0.payload_description_header").(string),
		PayloadDisplayName:       d.Get("payloads.0.payload_display_name_header").(string),
		PayloadEnabled:           d.Get("payloads.0.payload_enabled_header").(bool),
		PayloadIdentifier:        uuidStr,
		PayloadOrganization:      d.Get("payloads.0.payload_organization_header").(string),
		PayloadRemovalDisallowed: d.Get("payloads.0.payload_removal_disallowed_header").(bool),
		PayloadScope:             d.Get("payloads.0.payload_scope_header").(string),
		PayloadType:              d.Get("payloads.0.payload_type_header").(string),
		PayloadUUID:              uuidStr,
		PayloadVersion:           d.Get("payloads.0.payload_version_header").(int),
	}

	// Contents
	payloadContents := d.Get("payloads.0.payload_content").([]any)
	for _, v := range payloadContents {
		val := v.(map[string]any)
		payloadContentStruct := PayloadContent{
			PayloadDescription:  val["payload_description"].(string),
			PayloadDisplayName:  val["payload_display_name"].(string),
			PayloadEnabled:      val["payload_enabled"].(bool),
			PayloadIdentifier:   val["payload_identifier"].(string),
			PayloadOrganization: val["payload_organization"].(string),
			PayloadType:         val["payload_type"].(string),
			PayloadUUID:         val["payload_uuid"].(string),
			PayloadVersion:      val["payload_version"].(int),
		}

		settings := dictionaryEntries(val["setting"])
		if len(settings) == 0 {
			out.PayloadContent = append(out.PayloadContent, payloadContentStruct)
			continue
		}

		payloadContentStruct.ConfigurationItems = make(map[string]any, 0)
		for _, s := range settings {
			setting := s.(map[string]any)
			if _, exists := payloadContentStruct.ConfigurationItems[setting["key"].(string)]; exists {
				return nil, fmt.Errorf("duplicate setting key %q", setting["key"])
			}
			value, err := settingValue(setting)
			if err != nil {
				return nil, err
			}
			payloadContentStruct.ConfigurationItems[setting["key"].(string)] = value
		}

		out.PayloadContent = append(out.PayloadContent, payloadContentStruct)
	}

	return out, nil
}

// dictionaryEntries accepts legacy lists and current sets without changing arrays.
func dictionaryEntries(value any) []any {
	switch v := value.(type) {
	case *schema.Set:
		return v.List()
	case []any:
		return v
	default:
		return nil
	}
}

func settingValue(entry map[string]any) (any, error) {
	dictionaryPresent := len(dictionaryEntries(entry["dictionary"])) > 0
	if leaf, ok := entry["dictionary"].(map[string]any); ok && len(leaf) > 0 {
		dictionaryPresent = true
	}
	scalar, _ := entry["value"].(string)
	if dictionaryPresent && scalar != "" {
		return nil, fmt.Errorf("setting %q: value cannot be combined with dictionary", entry["key"])
	}
	if raw, ok := entry["array_json"].(string); ok && raw != "" {
		if value, _ := entry["value"].(string); value != "" || dictionaryPresent {
			return nil, fmt.Errorf("setting %q: array_json cannot be combined with value or dictionary", entry["key"])
		}
		var array []any
		decoder := json.NewDecoder(bytes.NewBufferString(raw))
		decoder.UseNumber()
		if err := decoder.Decode(&array); err != nil || array == nil {
			return nil, fmt.Errorf("setting %q: array_json must contain a JSON array", entry["key"])
		}
		return normalizeJSONNumbers(array), nil
	}
	if leaf, ok := entry["dictionary"].(map[string]any); ok && len(leaf) > 0 {
		result := make(map[string]any, len(leaf))
		for key, value := range leaf {
			result[key] = GetTypedValue(value)
		}
		return result, nil
	}
	if entries := dictionaryEntries(entry["dictionary"]); len(entries) > 0 {
		result := make(map[string]any, len(entries))
		for _, item := range entries {
			nested := item.(map[string]any)
			key := nested["key"].(string)
			if _, exists := result[key]; exists {
				return nil, fmt.Errorf("duplicate dictionary key %q", key)
			}
			value, err := settingValue(nested)
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
		return result, nil
	}
	return GetTypedValue(entry["value"]), nil
}

// GetTypedValue converts the value from the HCL always stored as string into the appropriate type for plist serialization.
func GetTypedValue(value any) any {
	strValue := fmt.Sprintf("%v", value)
	if boolValue, err := strconv.ParseBool(strValue); err == nil {
		return boolValue
	}
	if intValue, err := strconv.Atoi(strValue); err == nil {
		return intValue
	}
	return strValue
}

// ConvertPlistToHCL converts a plist XML string to Terraform HCL schema data
// Used by plist generator resource
func ConvertPlistToHCL(plistXML string) ([]any, error) {
	profile, err := UnmarshalPayload(plistXML)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	// Log the entire profile to verify its contents
	profileData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profile to JSON: %w", err)
	}
	log.Printf("[DEBUG] Unmarshaled profile: %s", string(profileData))

	if profile.PayloadContent != nil {
		for i, content := range profile.PayloadContent {
			contentData, err := json.MarshalIndent(content, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal content to JSON: %w", err)
			}
			log.Printf("[DEBUG] PayloadContent %d: %s", i, string(contentData))
		}
	} else {
		log.Printf("[DEBUG] PayloadContent is nil")
	}

	payloadsList, err := mapProfileToSchema(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to map profile to schema: %w", err)
	}

	// Convert the payload list to JSON for pretty print
	jsonData, err := json.MarshalIndent(payloadsList, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	log.Printf("[DEBUG] Constructed TF state structure from plist:\n%s\n", string(jsonData))

	return payloadsList, nil
}

// mapProfileToSchema maps the ConfigurationProfile struct data to the Terraform schema
func mapProfileToSchema(profile *ConfigurationProfile) ([]any, error) {
	payloadHeader := map[string]any{
		"payload_description_header":        profile.PayloadDescription,
		"payload_display_name_header":       profile.PayloadDisplayName,
		"payload_enabled_header":            profile.PayloadEnabled,
		"payload_identifier_header":         profile.PayloadIdentifier,
		"payload_organization_header":       profile.PayloadOrganization,
		"payload_removal_disallowed_header": profile.PayloadRemovalDisallowed,
		"payload_scope_header":              profile.PayloadScope,
		"payload_type_header":               profile.PayloadType,
		"payload_uuid_header":               profile.PayloadUUID,
		"payload_version_header":            profile.PayloadVersion,
	}

	payloadContentList := []any{}
	for _, content := range profile.PayloadContent {
		payloadContent := map[string]any{
			"payload_description":  content.PayloadDescription,
			"payload_display_name": content.PayloadDisplayName,
			"payload_enabled":      content.PayloadEnabled,
			"payload_identifier":   content.PayloadIdentifier,
			"payload_organization": content.PayloadOrganization,
			"payload_type":         content.PayloadType,
			"payload_uuid":         content.PayloadUUID,
			"payload_version":      content.PayloadVersion,
		}

		log.Printf("[DEBUG] ConfigurationItems being passed: %v", content.ConfigurationItems)
		settingsList := []any{}
		if err := extractNestedConfigurationSettings(content.ConfigurationItems, &settingsList); err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Final settingsList: %v", settingsList)

		payloadContent["setting"] = settingsList
		payloadContentList = append(payloadContentList, payloadContent)
	}

	payloadHeader["payload_content"] = payloadContentList

	return []any{payloadHeader}, nil
}

// extractNestedConfigurationSettings recursively extracts key-value pairs from nested dictionaries and appends them to settingsList
func extractNestedConfigurationSettings(items map[string]any, settingsList *[]any) error {
	log.Printf("[DEBUG] Raw data being processed: %v", items)
	for key, value := range items {
		log.Printf("[DEBUG] Processing configuration item key: %s, value: %v", key, value)
		settingMap := map[string]any{
			"key": key,
		}

		switch v := value.(type) {
		case map[string]any:
			if len(v) > 0 {
				nestedSettings := []any{}
				if err := extractNestedConfigurationSettings(v, &nestedSettings); err != nil {
					return err
				}
				settingMap["dictionary"] = nestedSettings
			} else {
				settingMap["value"] = "{}"
			}
		case []any:
			// Arrays are ordered values, including duplicate items and dictionary elements.
			if err := validateJSONArrayValue(v); err != nil {
				return fmt.Errorf("setting %q: %w", key, err)
			}
			encoded, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("setting %q: %w", key, err)
			}
			settingMap["array_json"] = string(encoded)
		case bool, int, float64, string:
			settingMap["value"] = fmt.Sprintf("%v", v)
		default:
			settingMap["value"] = v
		}

		log.Printf("[DEBUG] Adding settingMap: %v", settingMap)
		*settingsList = append(*settingsList, settingMap)
	}
	return nil
}

// UnmarshalPayload unmarshals a plist payload into a ConfigurationProfile struct using mapstructure.
func UnmarshalPayload(payload string) (*ConfigurationProfile, error) {
	var profile map[string]any
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
func MergeConfigurationProfileFieldsIntoMap(profile *ConfigurationProfile) map[string]any {
	merged := make(map[string]any, len(profile.Unexpected))
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

	mergedPayloads := make([]map[string]any, len(profile.PayloadContent))
	for k, v := range profile.PayloadContent {
		mergedPayloads[k] = MergeConfigurationPayloadFieldsIntoMap(&v)
	}

	merged["PayloadContent"] = mergedPayloads

	return merged
}

// MergeConfigurationPayloadFieldsIntoMap merges the fields of a ConfigurationPayload struct into a map.
func MergeConfigurationPayloadFieldsIntoMap(payload *PayloadContent) map[string]any {
	merged := make(map[string]any, len(payload.ConfigurationItems))
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

// JSON numbers retain their integer/real distinction when written to plist.
func normalizeJSONNumbers(value any) any {
	switch v := value.(type) {
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n
		}
		if n, err := strconv.ParseUint(string(v), 10, 64); err == nil {
			return n
		}
		n, _ := v.Float64()
		return n
	case []any:
		for i, child := range v {
			v[i] = normalizeJSONNumbers(child)
		}
		return v
	case map[string]any:
		for k, child := range v {
			v[k] = normalizeJSONNumbers(child)
		}
		return v
	default:
		return v
	}
}

// Refuse unsupported plist scalar types rather than silently converting dates
// or binary data to JSON strings and changing their type on the next apply.
func validateJSONArrayValue(value any) error {
	switch v := value.(type) {
	case string, bool, int, int64, uint64, float64, float32, json.Number:
		return nil
	case []any:
		for _, child := range v {
			if err := validateJSONArrayValue(child); err != nil {
				return err
			}
		}
	case map[string]any:
		for _, child := range v {
			if err := validateJSONArrayValue(child); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("plist array contains unsupported JSON value type %T", value)
	}
	return nil
}
