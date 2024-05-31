// staticcomputergroup_object.go
package staticcomputergroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProStaticComputerGroup constructs a ResourceComputerGroup object from the schema.ResourceData
func constructJamfProStaticComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {

	// Initialize the ResourceComputerGroup object
	resource := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false, // Static computer groups are not smart
	}
	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		site := constructobject.ConstructSharedResourceSite(v.([]interface{}))
		resource.Site = &site
	} else {
		// Set default values if 'site' data is not provided
		site := constructobject.ConstructSharedResourceSite([]interface{}{})
		resource.Site = &site
	}

	// Extract the assignments information if provided
	if v, ok := d.GetOk("assignments"); ok {
		assignmentsList := v.([]interface{})
		if len(assignmentsList) > 0 {
			assignmentsData := assignmentsList[0].(map[string]interface{})
			computerIDs := assignmentsData["computer_ids"].([]interface{})

			computers := make([]jamfpro.ComputerGroupSubsetComputer, len(computerIDs))
			for i, id := range computerIDs {
				computers[i] = jamfpro.ComputerGroupSubsetComputer{
					ID: id.(int),
				}
			}
			resource.Computers = &jamfpro.ComputerGroupSubsetComputersContainer{
				Size:      len(computers),
				Computers: &computers,
			}
		}
	}
	// Serialize and pretty-print the Computer Group object as XML for logging
	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Static Computer Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Static Computer Group XML:\n%s\n", string(resourceXML))

	return resource, nil
}
