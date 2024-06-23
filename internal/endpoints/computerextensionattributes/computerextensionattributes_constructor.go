// computerextensionattributes_object.go
package computerextensionattributes

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerExtensionAttribute constructs a ResourceComputerExtensionAttribute object from the provided schema data.
func constructJamfProComputerExtensionAttribute(d *schema.ResourceData) (*jamfpro.ResourceComputerExtensionAttribute, error) {
	var resource *jamfpro.ResourceComputerExtensionAttribute

	resource = &jamfpro.ResourceComputerExtensionAttribute{
		Name:             d.Get("name").(string),
		Enabled:          d.Get("enabled").(bool),
		Description:      d.Get("description").(string),
		DataType:         d.Get("data_type").(string),
		InventoryDisplay: d.Get("inventory_display").(string),
		ReconDisplay:     d.Get("recon_display").(string),
	}

	if v, ok := d.GetOk("input_type"); ok && len(v.([]interface{})) > 0 {
		inputTypeData := v.([]interface{})[0].(map[string]interface{})
		inputType := jamfpro.ComputerExtensionAttributeSubsetInputType{
			Type:     inputTypeData["type"].(string),
			Platform: inputTypeData["platform"].(string),
			Script:   strings.TrimSpace(inputTypeData["script"].(string)),
		}

		if choices, exists := inputTypeData["choices"]; exists {
			for _, choice := range choices.([]interface{}) {
				inputType.Choices = append(inputType.Choices, choice.(string))
			}
		}

		resource.InputType = inputType
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Extension Attribute '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Extension Attribute XML:\n%s\n", string(resourceXML))

	return resource, nil
}
