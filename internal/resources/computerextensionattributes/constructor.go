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
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerExtensionAttribute, error) {
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
	inputTypeEnabledPopUp := len(d.Get("input_popup").([]interface{})) != 0
	inputTypeEnabledScript := d.Get("input_script") != ""
	// inputTypeEnabledDirectory := d.Get("input_directory_mapping") != ""

	log.Printf("Text: %v", d.Get("input_type").(string))
	log.Printf("Popup: %v", d.Get("input_popup"))
	log.Printf("script: %v", d.Get("input_script"))
	log.Printf("Dir: %v", d.Get("input_directory_mapping"))

	inputList := []bool{inputTypeEnabledText, inputTypeEnabledPopUp, inputTypeEnabledScript}
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

	if inputType == "Pop-up Menu" {
		choices := d.Get("input_popup").([]interface{})
		for _, v := range choices {
			resource.InputType.Choices = append(resource.InputType.Choices, v.(string))
		}
	} else if inputType == "script" {
		resource.InputType.Platform = "Mac"
		resource.InputType.Script = d.Get("input_script").(string)
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Extension Attribute '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Extension Attribute XML:\n%s\n", string(resourceXML))

	return resource, nil
}
