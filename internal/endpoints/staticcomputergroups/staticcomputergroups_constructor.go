// staticcomputergroup_object.go
package staticcomputergroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProStaticComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProStaticComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	var resource *jamfpro.ResourceComputerGroup

	resource = &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	if v, ok := d.GetOk("assignments"); ok {
		assignments := v.([]interface{})
		if len(assignments) > 0 {
			computerIDs := assignments[0].(map[string]interface{})["computer_ids"].([]interface{})
			resource.Computers = constructGroupComputers(computerIDs)
		}
	}
	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Group XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// Helper functions for nested structures

// constructGroupComputers constructs a slice of ComputerGroupSubsetComputer from the provided schema data.
func constructGroupComputers(computerIDs []interface{}) *[]jamfpro.ComputerGroupSubsetComputer {
	var computers []jamfpro.ComputerGroupSubsetComputer
	for _, id := range computerIDs {
		computer := jamfpro.ComputerGroupSubsetComputer{
			ID: id.(int),
		}
		computers = append(computers, computer)
	}

	return &computers
}
