// computerextensionattributes_object.go
package computerextensionattributes

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerExtensionAttribute constructs a ResourceComputerExtensionAttribute object from the provided schema data.
func constructJamfProComputerExtensionAttribute(d *schema.ResourceData) (*jamfpro.ResourceComputerExtensionAttribute, error) {
	attribute := &jamfpro.ResourceComputerExtensionAttribute{
		Name:             d.Get("name").(string),
		Enabled:          d.Get("enabled").(bool),
		Description:      d.Get("description").(string),
		DataType:         d.Get("data_type").(string),
		InventoryDisplay: d.Get("inventory_display").(string),
		ReconDisplay:     d.Get("recon_display").(string),
	}

	// Handle nested "input_type" field
	if v, ok := d.GetOk("input_type"); ok && len(v.([]interface{})) > 0 {
		inputTypeData := v.([]interface{})[0].(map[string]interface{})
		inputType := jamfpro.ComputerExtensionAttributeSubsetInputType{
			Type:     inputTypeData["type"].(string),
			Platform: inputTypeData["platform"].(string),
			Script:   inputTypeData["script"].(string),
		}

		// Handle "choices" within "input_type"
		if choices, exists := inputTypeData["choices"]; exists {
			for _, choice := range choices.([]interface{}) {
				inputType.Choices = append(inputType.Choices, choice.(string))
			}
		}

		attribute.InputType = inputType
	}

	// Serialize and pretty-print the attribute object as XML
	resourceXML, err := xml.MarshalIndent(attribute, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Extension Attribute '%s' to XML: %v", attribute.Name, err)
	}
	log.Printf("Constructed Jamf Pro Computer Extension Attribute XML:\n%s\n", string(resourceXML))

	return attribute, nil
}
