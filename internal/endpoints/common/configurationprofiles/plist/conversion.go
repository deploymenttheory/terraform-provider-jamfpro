package plist

import (
	"fmt"
	"log"
	"reflect"

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
	uuid := GenerateUUID()

	// Extracting payload content
	payloadContentData := payloadData["payload_content"].([]interface{})
	payloadContent := make([]ConfigurationPayload, len(payloadContentData))

	for i, pc := range payloadContentData {
		pcMap := pc.(map[string]interface{})
		key := pcMap["key"].(string)
		value := pcMap["value"]
		payloadContent[i] = ConfigurationPayload{
			AdditionalFields: map[string]interface{}{key: value},
		}
	}

	// Creating a ConfigurationProfile struct from the extracted data
	profile := &ConfigurationProfile{
		PayloadDescription: payloadData["payload_description"].(string),
		//PayloadDisplayName:       payloadData["payload_display_name"].(string),
		PayloadEnabled:           payloadData["payload_enabled"].(bool),
		PayloadIdentifier:        uuid,
		PayloadOrganization:      payloadData["payload_organization"].(string),
		PayloadRemovalDisallowed: payloadData["payload_removal_disallowed"].(bool),
		PayloadScope:             payloadData["payload_scope"].(string),
		PayloadType:              payloadData["payload_type"].(string),
		PayloadUUID:              uuid,
		PayloadVersion:           payloadData["payload_version"].(int),
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

// // ConvertHCLToPlist converts the payloads list to a map and generates the plist XML.
// func ConvertHCLToPlist(d *schema.ResourceData) (string, error) {
// 	payloadsList := d.Get("payloads").([]interface{})
// 	var configurationProfile ConfigurationProfile

// 	configurationProfile.PayloadContent = make([]ConfigurationPayload, 0)

// 	for _, payload := range payloadsList {
// 		payloadData := payload.(map[string]interface{})
// 		var configurationPayload ConfigurationPayload
// 		configurationPayload.AdditionalFields = make(map[string]interface{})

// 		if payloadContent, ok := payloadData["payload_content"].([]interface{}); ok {
// 			for _, content := range payloadContent {
// 				contentData := content.(map[string]interface{})
// 				key := contentData["key"].(string)
// 				value := contentData["value"]

// 				// Detect the type of the value and set it accordingly
// 				switch v := value.(type) {
// 				case bool:
// 					configurationPayload.AdditionalFields[key] = v
// 				case int:
// 					configurationPayload.AdditionalFields[key] = v
// 				case float64: // Terraform SDK might return float64 for numbers
// 					configurationPayload.AdditionalFields[key] = int(v)
// 				case string:
// 					if boolVal, err := strconv.ParseBool(v); err == nil {
// 						configurationPayload.AdditionalFields[key] = boolVal
// 					} else if intVal, err := strconv.Atoi(v); err == nil {
// 						configurationPayload.AdditionalFields[key] = intVal
// 					} else {
// 						configurationPayload.AdditionalFields[key] = v
// 					}
// 				default:
// 					errorMessage := fmt.Sprintf("ERROR: Got value of type %T with value %v, unable to convert", v, v)
// 					configurationPayload.AdditionalFields[key] = "ERROR"
// 					log.Println(errorMessage)
// 				}
// 			}
// 		}

// 		// Set known fields for the ConfigurationPayload
// 		if v, ok := payloadData["payload_description"]; ok {
// 			configurationPayload.PayloadDescription = v.(string)
// 		}
// 		if v, ok := payloadData["payload_display_name"]; ok {
// 			configurationPayload.PayloadDisplayName = v.(string)
// 		}
// 		if v, ok := payloadData["payload_enabled"]; ok {
// 			configurationPayload.PayloadEnabled = v.(bool)
// 		}
// 		if v, ok := payloadData["payload_identifier"]; ok {
// 			configurationPayload.PayloadIdentifier = v.(string)
// 		}
// 		if v, ok := payloadData["payload_organization"]; ok {
// 			configurationPayload.PayloadOrganization = v.(string)
// 		}
// 		if v, ok := payloadData["payload_type"]; ok {
// 			configurationPayload.PayloadType = v.(string)
// 		}
// 		if v, ok := payloadData["payload_uuid"]; ok {
// 			configurationPayload.PayloadUUID = v.(string)
// 		}
// 		if v, ok := payloadData["payload_version"]; ok {
// 			configurationPayload.PayloadVersion = v.(int)
// 		}

// 		// Add the configuration payload to the profile's payload content
// 		configurationProfile.PayloadContent = append(configurationProfile.PayloadContent, configurationPayload)
// 	}

// 	// Retrieve and set the root-level fields from HCL input
// 	fields := map[string]interface{}{
// 		"payload_description":        &configurationProfile.PayloadDescription,
// 		"payload_display_name":       &configurationProfile.PayloadDisplayName,
// 		"payload_enabled":            &configurationProfile.PayloadEnabled,
// 		"payload_identifier":         &configurationProfile.PayloadIdentifier,
// 		"payload_organization":       &configurationProfile.PayloadOrganization,
// 		"payload_removal_disallowed": &configurationProfile.PayloadRemovalDisallowed,
// 		"payload_scope":              &configurationProfile.PayloadScope,
// 		"payload_type":               &configurationProfile.PayloadType,
// 		"payload_uuid":               &configurationProfile.PayloadUUID,
// 		"payload_version":            &configurationProfile.PayloadVersion,
// 	}

// 	for field, fieldPtr := range fields {
// 		if v, ok := d.GetOk(field); ok {
// 			setField(fieldPtr, v)
// 		}
// 	}

// 	// Marshal the ConfigurationProfile to plist XML
// 	payloadsXML, err := MarshalPayload(&configurationProfile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal payload: %w", err)
// 	}

// 	log.Printf("[DEBUG] Constructed plist XML from HCL:\n%s\n", payloadsXML)

// 	return payloadsXML, nil
// }

// setField sets the value of the field based on its type
func setField(fieldPtr interface{}, value interface{}) {
	v := reflect.ValueOf(fieldPtr).Elem()

	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(value.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value.(type) {
		case int:
			v.SetInt(int64(value.(int)))
		case float64: // Terraform SDK might return float64 for numbers
			v.SetInt(int64(value.(float64)))
		}
	case reflect.String:
		v.SetString(value.(string))
	default:
		v.Set(reflect.ValueOf(value))
	}
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
		"payload_description":  profile.PayloadDescription,
		"payload_display_name": profile.PayloadDisplayName,
		"payload_enabled":      profile.PayloadEnabled,
		//"payload_identifier":         profile.PayloadIdentifier,
		"payload_organization":       profile.PayloadOrganization,
		"payload_removal_disallowed": profile.PayloadRemovalDisallowed,
		//"payload_scope":              profile.PayloadScope,
		"payload_type": profile.PayloadType,
		//"payload_uuid":               profile.PayloadUUID,
		"payload_version": profile.PayloadVersion,
	}

	// Convert each ConfigurationPayload to the appropriate format
	var payloadContentList []interface{}
	for _, configurationPayload := range profile.PayloadContent {
		payloadMap := make(map[string]interface{})

		// Convert AdditionalFields to payload_content list
		for key, value := range configurationPayload.AdditionalFields {
			payloadMap[key] = value
		}

		payloadContentList = append(payloadContentList, payloadMap)
	}

	profileMap["payload_content"] = payloadContentList
	payloadsList = append(payloadsList, profileMap)

	return payloadsList, nil
}
