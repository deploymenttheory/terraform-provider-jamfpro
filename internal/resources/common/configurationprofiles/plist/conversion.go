// // common/configurationprofiles/plist/conversion.go
// // Description: This file contains the functions to convert the HCL data to a plist XML string and vice versa.

package plist

import (
	"encoding/json"
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
		}

		log.Printf("[DEBUG] ConfigurationItems being passed: %v", content.ConfigurationItems)
		settingsList := []interface{}{}
		extractNestedConfigurationSettings(content.ConfigurationItems, &settingsList)
		log.Printf("[DEBUG] Final settingsList: %v", settingsList)

		payloadContent["setting"] = settingsList
		payloadContentList = append(payloadContentList, payloadContent)
	}

	payloadHeader["payload_content"] = payloadContentList

	return []interface{}{payloadHeader}, nil
}

// extractNestedConfigurationSettings recursively extracts key-value pairs from nested dictionaries and appends them to settingsList
func extractNestedConfigurationSettings(items map[string]interface{}, settingsList *[]interface{}) {
	log.Printf("[DEBUG] Raw data being processed: %v", items)
	for key, value := range items {
		log.Printf("[DEBUG] Processing configuration item key: %s, value: %v", key, value)
		settingMap := map[string]interface{}{
			"key": key,
		}

		switch v := value.(type) {
		case map[string]interface{}:
			if len(v) > 0 {
				nestedSettings := []interface{}{}
				extractNestedConfigurationSettings(v, &nestedSettings)
				settingMap["dictionary"] = nestedSettings
			} else {
				settingMap["value"] = "{}"
			}
		case []interface{}:
			if len(v) > 0 {
				var nestedSettings []interface{}
				for _, item := range v {
					if nestedItem, ok := item.(map[string]interface{}); ok {
						nestedSettings = append(nestedSettings, nestedItem)
					}
				}
				settingMap["dictionary"] = nestedSettings
			} else {
				settingMap["value"] = "[]"
			}
		case bool, int, float64, string:
			settingMap["value"] = fmt.Sprintf("%v", v)
		default:
			settingMap["value"] = v
		}

		log.Printf("[DEBUG] Adding settingMap: %v", settingMap)
		*settingsList = append(*settingsList, settingMap)
	}
}
