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

// constructJamfProStaticComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProStaticComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	group := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	if v, ok := d.GetOk("site_id"); ok {
		site := constructobject.ConstructSharedResourceSite(v.([]interface{}))
		group.Site = site
	}

	if v, ok := d.GetOk("assignments"); ok {
		assignments := v.([]interface{})
		if len(assignments) > 0 {
			computerIDs := assignments[0].(map[string]interface{})["computer_ids"].([]interface{})
			group.Computers = constructGroupComputers(computerIDs)
		}
	}
	resourceXML, err := xml.MarshalIndent(group, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Group '%s' to XML: %v", group.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Group XML:\n%s\n", string(resourceXML))

	return group, nil
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
