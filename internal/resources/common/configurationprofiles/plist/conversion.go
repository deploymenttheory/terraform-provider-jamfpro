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

// ConvertHCLToPlist converts the HCL data and serializes it to a plist XML string.
// A UUID is generated for the payload identifier and payload UUID for each payload.
// This is required for a successful POST request to the Jamf Pro API.
func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	// Generate UUID for root and payloads
	uuidStr := GenerateUUID()

	// Extracting HCL data
	payloads := d.Get("payloads").([]interface{})
	if len(payloads) == 0 {
		return "", fmt.Errorf("no payloads found in the provided HCL")
	}

	payloadData := payloads[0].(map[string]interface{})

	// Extract payload content data
	payloadContentData := payloadData["payload_content"].([]interface{})

	// Convert payload content data
	payloadContent := make([]PayloadContent, len(payloadContentData))
	for i, pc := range payloadContentData {
		pcMap := pc.(map[string]interface{})
		settingData := pcMap["setting"].([]interface{})
		additionalFields := make(map[string]interface{})
		for _, setting := range settingData {
			settingMap := setting.(map[string]interface{})
			key := settingMap["key"].(string)
			value := extractAndParseValue(settingMap["value"])
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
		PayloadDescription:       d.Get("payload_description_header").(string),
		PayloadDisplayName:       d.Get("payload_display_name_header").(string),
		PayloadEnabled:           d.Get("payload_enabled_header").(bool),
		PayloadIdentifier:        uuidStr,
		PayloadOrganization:      d.Get("payload_organization_header").(string),
		PayloadRemovalDisallowed: d.Get("payload_removal_disallowed_header").(bool),
		PayloadScope:             d.Get("payload_scope_header").(string),
		PayloadType:              d.Get("payload_type_header").(string),
		PayloadUUID:              uuidStr,
		PayloadVersion:           d.Get("payload_version_header").(int),
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

// extractAndParseValue recursively retrieves nested values from a schema structure and parses them.
func extractAndParseValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			firstElem := v[0].(map[string]interface{})
			nestedValue := extractAndParseValue(firstElem["value"])
			if dictionary, ok := firstElem["dictionary"].([]interface{}); ok && len(dictionary) > 0 {
				return map[string]interface{}{
					"key":        firstElem["key"].(string),
					"value":      nestedValue,
					"dictionary": extractAndParseValue(dictionary),
				}
			}
			return map[string]interface{}{
				"key":   firstElem["key"].(string),
				"value": nestedValue,
			}
		}
		return nil
	case map[string]interface{}:
		return extractAndParseValue(v["value"])
	default:
		return GetTypedValue(value)
	}
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
