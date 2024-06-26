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
	var resource *jamfpro.ResourceComputerExtensionAttribute

	resource = &jamfpro.ResourceComputerExtensionAttribute{
		Name:             d.Get("name").(string),
		Enabled:          d.Get("enabled").(bool),
		Description:      d.Get("description").(string),
		DataType:         d.Get("data_type").(string),
		InventoryDisplay: d.Get("inventory_display").(string),
		ReconDisplay:     d.Get("recon_display").(string),
	}

	inputTypeEnabledText := d.Get("input_type").(string) == "Text Field"
	inputTypeEnabledPopUp := d.Get("input_popup").(string) != ""
	inputTypeEnabledScript := d.Get("input_script").(string) != ""
	inputTypeEnabledDirectory := d.Get("input_directory").(string) != ""

	inputList := []bool{inputTypeEnabledText, inputTypeEnabledPopUp, inputTypeEnabledScript, inputTypeEnabledDirectory}
	var boolCount int
	for _, v := range inputList {
		if v {
			boolCount += 1
		}
	}
	if boolCount > 1 {
		return nil, fmt.Errorf("multiple input types selected, please adjust your configuratuon")
	}

	inputType := d.Get("input_type").(string)
	resource.InputType.Type = inputType

	switch inputType {
	case "Pop-up Menu":
		choices := d.Get("input_popup").([]interface{})
		for _, v := range choices {
			resource.InputType.Choices = append(resource.InputType.Choices, v.(string))
		}
	case "script":
		resource.InputType.Platform = "Mac"
		resource.InputType.Script = d.Get("input_script").(string)
	default:
		return nil, fmt.Errorf("invalid input type supplie")
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Extension Attribute '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Extension Attribute XML:\n%s\n", string(resourceXML))

	return resource, nil
}
