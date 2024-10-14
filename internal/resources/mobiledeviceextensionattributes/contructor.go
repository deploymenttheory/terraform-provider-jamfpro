package mobiledeviceextensionattributes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceMobileExtensionAttribute object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileExtensionAttribute, error) {
	resource := &jamfpro.ResourceMobileExtensionAttribute{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		DataType:         d.Get("data_type").(string),
		InventoryDisplay: d.Get("inventory_display").(string),
	}

	// Handle the nested input_type structure
	if v, ok := d.GetOk("input_type"); ok {
		inputTypeList := v.([]interface{})
		if len(inputTypeList) > 0 {
			inputTypeMap := inputTypeList[0].(map[string]interface{})
			resource.InputType = jamfpro.MobileExtensionAttributeSubsetInputType{
				Type: inputTypeMap["type"].(string),
			}

			// Handle popup choices
			if choices, ok := inputTypeMap["popup_choices"]; ok {
				choicesList := choices.([]interface{})
				resource.InputType.PopupChoices = jamfpro.PopupChoices{
					Choice: make([]string, len(choicesList)),
				}
				for i, choice := range choicesList {
					resource.InputType.PopupChoices.Choice[i] = choice.(string)
				}
			}
		}
	}

	// Validate the input type
	if err := ValidateInputType(resource); err != nil {
		return nil, fmt.Errorf("failed to construct: %v", err)
	}

	// Serialize and pretty-print the mobile device extension attribute object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Extension Attribute to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Extension Attribute JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
