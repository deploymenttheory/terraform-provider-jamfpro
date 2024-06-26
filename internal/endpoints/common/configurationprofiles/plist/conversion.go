// common/configurationprofiles/plist/conversion.go contains the functions to convert HCL to plist and plist to HCL.
package plist

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConvertHCLToPlist converts the payloads list to a map and generates the plist XML.
func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
	payloadsList := d.Get("payloads").([]interface{})
	var configurationProfile ConfigurationProfile

	configurationProfile.PayloadContent = make([]ConfigurationPayload, 0)

	for _, payload := range payloadsList {
		payloadData := payload.(map[string]interface{})
		var configurationPayload ConfigurationPayload

		if payloadContent, ok := payloadData["payload_content"].([]interface{}); ok {
			configurationPayload.AdditionalFields = make(map[string]interface{})
			for _, content := range payloadContent {
				contentData := content.(map[string]interface{})
				key := contentData["key"].(string)
				value := contentData["value"]

				configurationPayload.AdditionalFields[key] = value
			}
		}
		// Add the configuration payload to the profile's payload content
		configurationProfile.PayloadContent = append(configurationProfile.PayloadContent, configurationPayload)
	}

	// Retrieve and set the root-level fields from HCL input
	fields := map[string]interface{}{
		"payload_description":        &configurationProfile.PayloadDescription,
		"payload_display_name":       &configurationProfile.PayloadDisplayName,
		"payload_enabled":            &configurationProfile.PayloadEnabled,
		"payload_identifier":         &configurationProfile.PayloadIdentifier,
		"payload_organization":       &configurationProfile.PayloadOrganization,
		"payload_removal_disallowed": &configurationProfile.PayloadRemovalDisallowed,
		"payload_scope":              &configurationProfile.PayloadScope,
		"payload_type":               &configurationProfile.PayloadType,
		"payload_uuid":               &configurationProfile.PayloadUUID,
		"payload_version":            &configurationProfile.PayloadVersion,
	}

	for field, fieldPtr := range fields {
		if v, ok := d.GetOk(fmt.Sprintf("payloads.0.%s", field)); ok {
			setField(fieldPtr, v)
		}
	}

	// Marshal the ConfigurationProfile to plist XML
	payloadsXML, err := MarshalPayload(&configurationProfile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	log.Printf("[DEBUG] Constructed plist XML from HCL:\n%s\n", payloadsXML)

	return payloadsXML, nil
}

// setField sets the value of the field based on its type
func setField(fieldPtr interface{}, value interface{}) {
	v := reflect.ValueOf(fieldPtr).Elem()
	v.Set(reflect.ValueOf(value))
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
