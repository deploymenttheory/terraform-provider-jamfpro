// dockitems_data_object.go
package dockitems

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDockItem constructs a ResourceDockItem object from the provided schema data.
func constructJamfProDockItem(d *schema.ResourceData) (*jamfpro.ResourceDockItem, error) {
	dockItem := &jamfpro.ResourceDockItem{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Path:     d.Get("path").(string),
		Contents: (d.Get("contents").(string)),
	}

	// Serialize and pretty-print the site object as XML
	resourceXML, err := xml.MarshalIndent(dockItem, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Dock Item '%s' to XML: %v", dockItem.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Dock Item XML:\n%s\n", string(resourceXML))

	return dockItem, nil
}
