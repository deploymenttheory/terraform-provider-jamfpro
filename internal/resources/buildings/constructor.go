// buildings_object.go
package buildings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProBuilding constructs a Building object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceBuilding, error) {
	resource := &jamfpro.ResourceBuilding{
		Name:           d.Get("name").(string),
		StreetAddress1: d.Get("street_address1").(string),
		StreetAddress2: d.Get("street_address2").(string),
		City:           d.Get("city").(string),
		StateProvince:  d.Get("state_province").(string),
		ZipPostalCode:  d.Get("zip_postal_code").(string),
		Country:        d.Get("country").(string),
	}

	// Serialize and pretty-print the Building object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Jamf Pro Building '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Building JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
