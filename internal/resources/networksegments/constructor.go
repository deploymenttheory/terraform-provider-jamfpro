// networksegments_object.go
package networksegments

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProNetworkSegment constructs a ResourceNetworkSegment object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceNetworkSegment, error) {
	resource := &jamfpro.ResourceNetworkSegment{
		Name:                d.Get("name").(string),
		StartingAddress:     d.Get("starting_address").(string),
		EndingAddress:       d.Get("ending_address").(string),
		DistributionServer:  d.Get("distribution_server").(string),
		DistributionPoint:   d.Get("distribution_point").(string),
		URL:                 d.Get("url").(string),
		SWUServer:           d.Get("swu_server").(string),
		Building:            d.Get("building").(string),
		Department:          d.Get("department").(string),
		OverrideBuildings:   d.Get("override_buildings").(bool),
		OverrideDepartments: d.Get("override_departments").(bool),
	}

	// Serialize and pretty-print the Network Segment object as XML for logging
	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Network Segment '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Network Segment XML:\n%s\n", string(resourceXML))

	return resource, nil
}
