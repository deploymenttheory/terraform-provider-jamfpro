// dockitems_data_object.go
package dockitems

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDockItem constructs a ResourceDockItem object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDockItem, error) {
	var resource *jamfpro.ResourceDockItem

	resource = &jamfpro.ResourceDockItem{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Path:     d.Get("path").(string),
		Contents: (d.Get("contents").(string)),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Dock Item '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Dock Item XML:\n%s\n", string(resourceXML))

	return resource, nil
}
