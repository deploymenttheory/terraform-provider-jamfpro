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

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		site := constructobject.ConstructSharedResourceSite(v.([]interface{}))
		group.Site = &site
	} else {
		// Set default values if 'site' data is not provided
		site := constructobject.ConstructSharedResourceSite([]interface{}{})
		group.Site = &site
	}

	// Handle Computers by IDs
	if v, ok := d.GetOk("computer_ids"); ok {
		computerIDs := v.([]interface{})
		group.Computers = constructComputerGroupSubsetComputersContainerFromIDs(computerIDs)
	}

	// Serialize and pretty-print the Computer Group object as XML for logging
	resourceXML, err := xml.MarshalIndent(group, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Static Computer Group '%s' to XML: %v", group.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Static Computer Group XML:\n%s\n", string(resourceXML))

	return group, nil
}

// constructComputerGroupSubsetComputersContainerFromIDs constructs a ComputerGroupSubsetComputersContainer object from the provided list of computer IDs.
func constructComputerGroupSubsetComputersContainerFromIDs(computerIDs []interface{}) *jamfpro.ComputerGroupSubsetComputersContainer {
	computers := &jamfpro.ComputerGroupSubsetComputersContainer{
		Size:      len(computerIDs),
		Computers: &[]jamfpro.ComputerGroupSubsetComputer{},
	}

	for _, id := range computerIDs {
		computer := jamfpro.ComputerGroupSubsetComputer{
			ID: id.(int),
		}
		*computers.Computers = append(*computers.Computers, computer)
	}

	return computers
}
