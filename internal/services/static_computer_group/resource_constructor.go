// staticcomputergroup_object.go
package static_computer_group

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProStaticComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	resource := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() {
		if raw := rawConfig.GetAttr("assigned_computer_ids"); !raw.IsNull() {
			assignedComputers := d.Get("assigned_computer_ids").([]any)
			computers := []jamfpro.ComputerGroupSubsetComputer{}
			for _, id := range assignedComputers {
				computers = append(computers, jamfpro.ComputerGroupSubsetComputer{
					ID: id.(int),
				})
			}
			resource.Computers = &computers
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
