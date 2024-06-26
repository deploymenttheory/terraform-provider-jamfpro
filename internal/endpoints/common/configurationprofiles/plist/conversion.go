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

	// Extracting payload content
	payloadContentData := payloadData["payload_content"].([]interface{})
	payloadContent := make([]PayloadContent, len(payloadContentData))

	for i, pc := range payloadContentData {
		pcMap := pc.(map[string]interface{})
		key := pcMap["key"].(string)
		value := pcMap["value"]
		payloadContent[i] = PayloadContent{
			AdditionalFields:    map[string]interface{}{key: value},
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
		PayloadDescription:       payloadData["payload_description"].(string),
		PayloadDisplayName:       payloadData["payload_display_name"].(string),
		PayloadEnabled:           payloadData["payload_enabled"].(bool),
		PayloadIdentifier:        uuidStr,
		PayloadOrganization:      payloadData["payload_organization"].(string),
		PayloadRemovalDisallowed: payloadData["payload_removal_disallowed"].(bool),
		PayloadScope:             payloadData["payload_scope"].(string),
		PayloadType:              payloadData["payload_type"].(string),
		PayloadUUID:              uuidStr,
		PayloadVersion:           payloadData["payload_version"].(int),
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
	profileMap := map[string]interface{}{
		"payload_description":        profile.PayloadDescription,
		"payload_display_name":       profile.PayloadDisplayName,
		"payload_enabled":            profile.PayloadEnabled,
		"payload_identifier":         profile.PayloadIdentifier,
		"payload_organization":       profile.PayloadOrganization,
		"payload_removal_disallowed": profile.PayloadRemovalDisallowed,
		"payload_scope":              profile.PayloadScope,
		"payload_type":               profile.PayloadType,
		"payload_uuid":               profile.PayloadUUID,
		"payload_version":            profile.PayloadVersion,
	}

	// Convert each PayloadContent to the appropriate format
	var payloadContentList []interface{}
	for _, configurationPayload := range profile.PayloadContent {
		payloadMap := map[string]interface{}{
			"payload_description":  configurationPayload.PayloadDescription,
			"payload_display_name": configurationPayload.PayloadDisplayName,
			"payload_enabled":      configurationPayload.PayloadEnabled,
			//"payload_identifier":   configurationPayload.PayloadIdentifier,
			"payload_organization": configurationPayload.PayloadOrganization,
			"payload_type":         configurationPayload.PayloadType,
			//"payload_uuid":         configurationPayload.PayloadUUID,
			"payload_version": configurationPayload.PayloadVersion,
		}

		// Convert AdditionalFields to key-value pairs within the payload content
		for key, value := range configurationPayload.AdditionalFields {
			payloadMap[key] = value
		}

		payloadContentList = append(payloadContentList, payloadMap)
	}

	profileMap["payload_content"] = payloadContentList
	payloadsList = append(payloadsList, profileMap)

	return payloadsList, nil
}
