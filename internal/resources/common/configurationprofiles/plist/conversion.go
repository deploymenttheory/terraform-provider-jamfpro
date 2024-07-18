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

	// Make payloads list here
	payloadContents := d.Get("payloads.0.payload_content").([]map[string]interface{})
	for _, v := range payloadContents {

		payloadContentStruct := PayloadContent{
			PayloadDescription:  v["payload_description"].(string),
			PayloadDisplayName:  v["payload_display_name"].(string),
			PayloadEnabled:      v["payload_enabled"].(bool),
			PayloadIdentifier:   v["payload_identifier"].(string),
			PayloadOrganization: v["payload_organisation"].(string),
			PayloadType:         v["payload_type"].(string),
			PayloadUUID:         v["payload_uuid"].(string),
			PayloadVersion:      v["payload_version"].(int),
			PayloadScope:        v["payload_scope"].(string),
		}

		// Retrieve the payload contents
		settings := v["settings"].(map[string]interface{})
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

// ConvertPlistToHCL converts a plist XML string to HCL data. Used for stating the configuration profile data.
func ConvertPlistToHCL(plistXML string) ([]interface{}, error) {
	// Unmarshal the plist XML into a ConfigurationProfile struct
	profile, err := UnmarshalPayload(plistXML)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	// Convert the ConfigurationProfile struct to the format required by Terraform state
	var payloadsList []interface{}

	// Create a map for root-level fields
	profileRootMap := map[string]interface{}{
		"payload_description_root":        profile.PayloadDescription,
		"payload_display_name_root":       profile.PayloadDisplayName,
		"payload_enabled_root":            profile.PayloadEnabled,
		"payload_identifier_root":         profile.PayloadIdentifier,
		"payload_organization_root":       profile.PayloadOrganization,
		"payload_removal_disallowed_root": profile.PayloadRemovalDisallowed,
		"payload_scope_root":              profile.PayloadScope,
		"payload_type_root":               profile.PayloadType,
		"payload_uuid_root":               profile.PayloadUUID,
		"payload_version_root":            profile.PayloadVersion,
	}

	// Convert each PayloadContent to the appropriate format
	var payloadContentList []interface{}
	for _, configurationPayload := range profile.PayloadContent {
		configurations := make([]interface{}, 0, len(configurationPayload.ConfigurationItems))
		for key, value := range configurationPayload.ConfigurationItems {
			// Ensure all values are converted to strings for storage in the state
			strValue := fmt.Sprintf("%v", value)
			configurations = append(configurations, map[string]interface{}{
				"key":   key,
				"value": strValue,
			})
		}

		// Reorder configurations based on the jamf pro server logic
		reorderedConfigurations := reorderConfigurationKeys(configurations)

		payloadMap := map[string]interface{}{
			"payload_description":  configurationPayload.PayloadDescription,
			"payload_display_name": configurationPayload.PayloadDisplayName,
			"payload_enabled":      configurationPayload.PayloadEnabled,
			"payload_identifier":   configurationPayload.PayloadIdentifier,
			"payload_organization": configurationPayload.PayloadOrganization,
			"payload_type":         configurationPayload.PayloadType,
			"payload_uuid":         configurationPayload.PayloadUUID,
			"payload_version":      configurationPayload.PayloadVersion,
			"payload_scope":        configurationPayload.PayloadScope,
			"configuration":        reorderedConfigurations,
		}

		payloadContentList = append(payloadContentList, payloadMap)
	}

	// Create the full payloads map
	payloadsMap := map[string]interface{}{
		"payload_root":    []interface{}{profileRootMap},
		"payload_content": payloadContentList,
	}

	payloadsList = append(payloadsList, payloadsMap)

	return payloadsList, nil
}
