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
	resource := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	assignedComputers := d.Get("assigned_computer_ids").([]interface{})
	if len(assignedComputers) > 0 {
		computers := []jamfpro.ComputerGroupSubsetComputer{}
		for _, v := range assignedComputers {
			computers = append(computers, jamfpro.ComputerGroupSubsetComputer{
				ID: v.(int),
			})
		}
		resource.Computers = &computers
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Group XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// Helper functions for nested structures
