package plist

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	// Extracting payloads from the HCL
	payloads := d.Get("payloads").([]interface{})
	if len(payloads) == 0 {
		return "", fmt.Errorf("no payloads found in the provided HCL")
	}

	payloadData := payloads[0].(map[string]interface{})

	// Generate UUID if not provided
	uuidStr := GenerateUUID()

	// Extracting payload root
	payloadRootData := payloadData["payload_root"].([]interface{})[0].(map[string]interface{})

	// Extracting payload content
	payloadContentData := payloadData["payload_content"].([]interface{})
	payloadContent := make([]PayloadContent, len(payloadContentData))

	for i, pc := range payloadContentData {
		pcMap := pc.(map[string]interface{})
		configurations := pcMap["configuration"].([]interface{})
		additionalFields := make(map[string]interface{})
		for _, config := range configurations {
			configMap := config.(map[string]interface{})
			key := configMap["key"].(string)
			value := configMap["value"]
			additionalFields[key] = value
		}
		payloadContent[i] = PayloadContent{
			AdditionalFields:    additionalFields,
			PayloadDescription:  pcMap["payload_description"].(string),
			PayloadDisplayName:  pcMap["payload_display_name"].(string),
			PayloadEnabled:      pcMap["payload_enabled"].(bool),
			PayloadIdentifier:   pcMap["payload_identifier"].(string),
			PayloadOrganization: pcMap["payload_organization"].(string),
			PayloadType:         pcMap["payload_type"].(string),
			PayloadUUID:         pcMap["payload_uuid"].(string),
			PayloadVersion:      pcMap["payload_version"].(int),
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

	// Marshaling the ConfigurationProfile struct to a plist string
	plistXML, err := MarshalPayload(profile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile to plist: %v", err)
	}

	// Pretty-printing the plist XML for DEBUG logging
	prettyPlistXML, err := plist.MarshalIndent(plistXML, plist.XMLFormat, "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile to pretty plist: %v", err)
	}

	log.Printf("[DEBUG] Pretty printed plist XML:\n%s\n", string(prettyPlistXML))

	return plistXML, nil
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

// ConvertPlistToHCL converts a plist XML to the payloads list that can be set in the Terraform state.
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
			configurations = append(configurations, map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}

		payloadMap := map[string]interface{}{
			"payload_description":  configurationPayload.PayloadDescription,
			"payload_display_name": configurationPayload.PayloadDisplayName,
			"payload_enabled":      configurationPayload.PayloadEnabled,
			"payload_identifier":   configurationPayload.PayloadIdentifier,
			"payload_organization": configurationPayload.PayloadOrganization,
			"payload_type":         configurationPayload.PayloadType,
			"payload_uuid":         configurationPayload.PayloadUUID,
			"payload_version":      configurationPayload.PayloadVersion,
			"configuration":        configurations,
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
