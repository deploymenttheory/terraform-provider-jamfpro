package plist

import (
	"fmt"
)

// ConvertPlistToHCL converts a plist XML to the payloads list that can be set in the Terraform state.
func ConvertPlistToHCL(plistXML string) ([]interface{}, error) {
	// Unmarshal the plist XML into a ConfigurationProfile struct
	profile, err := UnmarshalPayload(plistXML)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	// Convert the ConfigurationProfile struct to the format required by Terraform state
	var payloadsList []interface{}
	for _, configurationPayload := range profile.PayloadContent {
		payloadMap := make(map[string]interface{})

		// Convert AdditionalFields to payload_content list
		var payloadContentList []interface{}
		for key, value := range configurationPayload.AdditionalFields {
			payloadContentList = append(payloadContentList, map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
		payloadMap["payload_content"] = payloadContentList

		// Set other fields
		payloadMap["payload_description"] = configurationPayload.PayloadDescription
		payloadMap["payload_display_name"] = configurationPayload.PayloadDisplayName
		payloadMap["payload_identifier"] = configurationPayload.PayloadIdentifier
		payloadMap["payload_organization"] = configurationPayload.PayloadOrganization
		payloadMap["payload_removal_disallowed"] = configurationPayload.PayloadRemovalDisallowed
		payloadMap["payload_scope"] = configurationPayload.PayloadScope
		payloadMap["payload_type"] = configurationPayload.PayloadType
		payloadMap["payload_uuid"] = configurationPayload.PayloadUUID
		payloadMap["payload_version"] = configurationPayload.PayloadVersion

		payloadsList = append(payloadsList, payloadMap)
	}

	return payloadsList, nil
}
