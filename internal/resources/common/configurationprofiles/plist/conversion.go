// common/configurationprofiles/plist/conversion.go
// Description: This file contains the functions to convert the HCL data to a plist XML string and vice versa.

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

// ConvertHCLToPlist converts the HCL data and serializes it to a plist XML string.
// A UUID is generated for the payload identifier and payload UUID for each payload.
// This is required for a successful POST request to the Jamf Pro API.
func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	uuidStr := GenerateUUID()
	// Extracting HCL data
	payloads := d.Get("payloads").([]interface{})
	if len(payloads) == 0 {
		return "", fmt.Errorf("no payloads found in the provided HCL")
	}

	payloadData := payloads[0].(map[string]interface{})

	payloadRootData := payloadData["payload_root"].([]interface{})[0].(map[string]interface{})
	payloadContentData := payloadData["payload_content"].([]interface{})

	payloadContent := make([]PayloadContent, len(payloadContentData))

	for i, pc := range payloadContentData {
		pcMap := pc.(map[string]interface{})
		configurations := pcMap["configuration"].([]interface{})
		additionalFields := make(map[string]interface{})
		for _, config := range configurations {
			configMap := config.(map[string]interface{})
			key := configMap["key"].(string)
			value := GetTypedValue(configMap["value"])
			if dictionaries, ok := configMap["dictionary"]; ok {
				value = ExtractNestedDictionary(dictionaries.([]interface{}))
			}
			additionalFields[key] = value
		}
		payloadContent[i] = PayloadContent{
			AdditionalFields:    additionalFields,
			PayloadDescription:  pcMap["payload_description"].(string),
			PayloadDisplayName:  pcMap["payload_display_name"].(string),
			PayloadEnabled:      pcMap["payload_enabled"].(bool),
			PayloadIdentifier:   uuidStr,
			PayloadOrganization: pcMap["payload_organization"].(string),
			PayloadType:         pcMap["payload_type"].(string),
			PayloadUUID:         uuidStr,
			PayloadVersion:      pcMap["payload_version"].(int),
			PayloadScope:        pcMap["payload_scope"].(string),
		}
	}

	// Creating a ConfigurationProfile struct from the extracted data
	profile := &ConfigurationProfile{
		PayloadDescription:       payloadRootData["payload_description_root"].(string),
		PayloadDisplayName:       payloadRootData["payload_display_name_root"].(string),
		PayloadEnabled:           payloadRootData["payload_enabled_root"].(bool),
		PayloadIdentifier:        uuidStr,
		PayloadOrganization:      payloadRootData["payload_organization_root"].(string),
		PayloadRemovalDisallowed: payloadRootData["payload_removal_disallowed_root"].(bool),
		PayloadScope:             payloadRootData["payload_scope_root"].(string),
		PayloadType:              payloadRootData["payload_type_root"].(string),
		PayloadUUID:              uuidStr,
		PayloadVersion:           payloadRootData["payload_version_root"].(int),
		PayloadContent:           payloadContent,
	}

	plistXML, err := MarshalPayload(profile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile to plist: %v", err)
	}

	prettyPlistXML, err := plist.MarshalIndent(plistXML, plist.XMLFormat, "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile to pretty plist: %v", err)
	}
	unescapedPrettyPlistXML := html.UnescapeString(string(prettyPlistXML))

	log.Printf("[DEBUG] Constructed Plist XML from HCL serialization:\n%s\n", unescapedPrettyPlistXML)

	return plistXML, nil
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
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

// ExtractNestedDictionary recursively extracts nested dictionaries from the given data.
func ExtractNestedDictionary(data []interface{}) map[string]interface{} {
	nestedDict := make(map[string]interface{})
	for _, item := range data {
		itemMap := item.(map[string]interface{})
		key := itemMap["key"].(string)
		value := GetTypedValue(itemMap["value"])
		if subDicts, ok := itemMap["dictionary"]; ok {
			value = ExtractNestedDictionary(subDicts.([]interface{}))
		}
		nestedDict[key] = value
	}
	return nestedDict
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
		configurations := make([]interface{}, 0, len(configurationPayload.AdditionalFields))
		for key, value := range configurationPayload.AdditionalFields {
			configMap := map[string]interface{}{
				"key": key,
			}
			if nestedDict, ok := value.(map[string]interface{}); ok {
				configMap["dictionary"] = FlattenNestedDictionary(nestedDict)
			} else {
				configMap["value"] = fmt.Sprintf("%v", value)
			}
			configurations = append(configurations, configMap)
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

// FlattenNestedDictionary flattens the nested dictionary structure into a format suitable for HCL.
func FlattenNestedDictionary(data map[string]interface{}) []interface{} {
	flattened := make([]interface{}, 0, len(data))
	for key, value := range data {
		item := map[string]interface{}{
			"key": key,
		}
		if nestedDict, ok := value.(map[string]interface{}); ok {
			item["dictionary"] = FlattenNestedDictionary(nestedDict)
		} else {
			item["value"] = fmt.Sprintf("%v", value)
		}
		flattened = append(flattened, item)
	}
	return flattened
}
