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
		group.Computers = &jamfpro.ComputerGroupSubsetComputersContainer{
			Size:      len(computerIDs),
			Computers: constructComputersFromIDs(computerIDs),
		}
	}

	// Serialize and pretty-print the Computer Group object as XML for logging
	resourceXML, err := xml.MarshalIndent(group, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Static Computer Group '%s' to XML: %v", group.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Static Computer Group XML:\n%s\n", string(resourceXML))

	return group, nil
}

// constructComputersFromIDs constructs a slice of ComputerGroupSubsetComputer from the provided list of computer IDs.
func constructComputersFromIDs(ids []interface{}) *[]jamfpro.ComputerGroupSubsetComputer {
	computers := make([]jamfpro.ComputerGroupSubsetComputer, len(ids))
	for i, id := range ids {
		computers[i] = jamfpro.ComputerGroupSubsetComputer{
			ID: id.(int),
		}
	}
	return &computers
}
