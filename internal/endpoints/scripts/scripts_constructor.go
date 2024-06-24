// scripts_object.go
package scripts

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProScript constructs a ResourceScript object from the provided schema data.
func constructJamfProScript(d *schema.ResourceData) (*jamfpro.ResourceScript, error) {
	var resource *jamfpro.ResourceScript

	// TODO category testing required ASAP

	resource = &jamfpro.ResourceScript{
		Name:           d.Get("name").(string),
		Info:           d.Get("info").(string),
		Notes:          d.Get("notes").(string),
		OSRequirements: d.Get("os_requirements").(string),
		Priority:       d.Get("priority").(string),
		Parameter4:     d.Get("parameter4").(string),
		Parameter5:     d.Get("parameter5").(string),
		Parameter6:     d.Get("parameter6").(string),
		Parameter7:     d.Get("parameter7").(string),
		Parameter8:     d.Get("parameter8").(string),
		Parameter9:     d.Get("parameter9").(string),
		Parameter10:    d.Get("parameter10").(string),
		Parameter11:    d.Get("parameter11").(string),
	}

	if scriptContent, ok := d.GetOk("script_contents"); ok {
		resource.ScriptContents = scriptContent.(string)
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Script '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Script JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
