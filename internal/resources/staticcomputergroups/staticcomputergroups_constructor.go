// staticcomputergroup_object.go
package staticcomputergroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProStaticComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	var resource *jamfpro.ResourceComputerGroup

	resource = &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	if resource.Computers != nil {
		assignedComputers := d.Get("assigned_computer_ids").([]interface{})
		if len(assignedComputers) > 0 {
			for _, v := range assignedComputers {
				*resource.Computers = append(*resource.Computers, jamfpro.ComputerGroupSubsetComputer{
					ID: v.(int),
				})
			}
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
