// // common/configurationprofiles/plist/conversion.go
// // Description: This file contains the functions to convert the HCL data to a plist XML string and vice versa.

package plist

import (
	"fmt"
	"html"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// ConvertHCLToPlist builds a plist from the Terraform HCL schema data
func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	profile := mapSchemaToProfile(d)
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
func mapSchemaToProfile(d *schema.ResourceData) *ConfigurationProfile {
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
	payloadContents := d.Get("payloads.0.payload_content").([]interface{})
	for _, v := range payloadContents {
		val := v.(map[string]interface{})
		payloadContentStruct := PayloadContent{
			PayloadDescription:  val["payload_description"].(string),
			PayloadDisplayName:  val["payload_display_name"].(string),
			PayloadEnabled:      val["payload_enabled"].(bool),
			PayloadIdentifier:   val["payload_identifier"].(string),
			PayloadOrganization: val["payload_organization"].(string),
			PayloadType:         val["payload_type"].(string),
			PayloadUUID:         val["payload_uuid"].(string),
			PayloadVersion:      val["payload_version"].(int),
			PayloadScope:        val["payload_scope"].(string),
		}

		settings := val["setting"].([]interface{})
		if len(settings) == 0 {
			return out
		}

		payloadContentStruct.ConfigurationItems = make(map[string]interface{}, 0)
		for _, s := range settings {
			settingMap := s.(map[string]interface{})
			dictionary := parseNestedDictionary(settingMap["dictionary"])
			payloadContent := map[string]interface{}{
				"key":        settingMap["key"].(string),
				"value":      GetTypedValue(settingMap["value"]),
				"dictionary": dictionary,
			}
			payloadContentStruct.ConfigurationItems[settingMap["key"].(string)] = payloadContent
		}

		out.PayloadContent = append(out.PayloadContent, payloadContentStruct)
	}

	return out
}

// parseNestedDictionary recursively parses the nested dictionary structure
func parseNestedDictionary(dict interface{}) map[string]interface{} {
	if dict == nil {
		return nil
	}

	result := make(map[string]interface{})
	dictionary := dict.([]interface{})
	for _, item := range dictionary {
		entry := item.(map[string]interface{})
		key := entry["key"].(string)
		value := GetTypedValue(entry["value"])
		if nestedDict, ok := entry["dictionary"].([]interface{}); ok {
			value = parseNestedDictionary(nestedDict)
		}
		result[key] = value
	}

	return result
}

// GetTypedValue converts the value from the HCL always stored as string into the appropriate type for plist serialization.
func GetTypedValue(value interface{}) interface{} {
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
func ConvertPlistToHCL(plistXML string) ([]interface{}, error) {
	var profile ConfigurationProfile
	if _, err := plist.Unmarshal([]byte(plistXML), &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	payloadsList, err := mapProfileToSchema(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to map profile to schema: %w", err)
	}

	log.Printf("[DEBUG] Constructed HCL schema from plist:\n%+v\n", payloadsList)

	return payloadsList, nil
}

// mapProfileToSchema maps the ConfigurationProfile struct data to the Terraform schema
func mapProfileToSchema(profile *ConfigurationProfile) ([]interface{}, error) {
	payloadHeader := map[string]interface{}{
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

	payloadContentList := []interface{}{}
	for _, content := range profile.PayloadContent {
		payloadContent := map[string]interface{}{
			"payload_description":  content.PayloadDescription,
			"payload_display_name": content.PayloadDisplayName,
			"payload_enabled":      content.PayloadEnabled,
			"payload_identifier":   content.PayloadIdentifier,
			"payload_organization": content.PayloadOrganization,
			"payload_type":         content.PayloadType,
			"payload_uuid":         content.PayloadUUID,
			"payload_version":      content.PayloadVersion,
			"payload_scope":        content.PayloadScope,
		}

		settingsList := []interface{}{}
		for _, itemInterface := range content.ConfigurationItems {
			item, ok := itemInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid configuration item format")
			}
			settingMap := map[string]interface{}{
				"key":        item["key"],
				"value":      item["value"],
				"dictionary": marshalNestedDictionary(item["dictionary"]),
			}
			settingsList = append(settingsList, settingMap)
		}

		payloadContent["setting"] = settingsList
		payloadContentList = append(payloadContentList, payloadContent)
	}

	payloadHeader["payload_content"] = payloadContentList

	return []interface{}{payloadHeader}, nil
}

// marshalNestedDictionary converts the nested dictionary structure back to an appropriate format
func marshalNestedDictionary(dict interface{}) []interface{} {
	if dict == nil {
		return nil
	}

	result := []interface{}{}
	dictionary := dict.(map[string]interface{})
	for key, value := range dictionary {
		entry := map[string]interface{}{
			"key":   key,
			"value": value,
		}
		if nestedDict, ok := value.(map[string]interface{}); ok {
			entry["dictionary"] = marshalNestedDictionary(nestedDict)
		}
		result = append(result, entry)
	}

	return result
}
