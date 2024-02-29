// buildings_object.go
package buildings

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProBuilding constructs a Building object from the provided schema data.
func constructJamfProBuilding(d *schema.ResourceData) (*jamfpro.ResourceBuilding, error) {
	building := &jamfpro.ResourceBuilding{
		Name:           d.Get("name").(string),
		StreetAddress1: d.Get("street_address1").(string),
		StreetAddress2: d.Get("street_address2").(string),
		City:           d.Get("city").(string),
		StateProvince:  d.Get("state_province").(string),
		ZipPostalCode:  d.Get("zip_postal_code").(string),
		Country:        d.Get("country").(string),
	}

	// Serialize and pretty-print the building object as XML for logging
	resourceXML, err := xml.MarshalIndent(building, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Building '%s' to XML: %v", building.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Building XML:\n%s\n", string(resourceXML))

	return building, nil
}
