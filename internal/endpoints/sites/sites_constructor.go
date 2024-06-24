// sites_data_object.go
package sites

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProSite constructs a SharedResourceSite object from the provided schema data.
func constructJamfProSite(d *schema.ResourceData) (*jamfpro.SharedResourceSite, error) {
	var resource *jamfpro.SharedResourceSite

	resource = &jamfpro.SharedResourceSite{
		Name: d.Get("name").(string),
	}

	// Serialize and pretty-print the Site object as XML for logging
	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Site '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Site XML:\n%s\n", string(resourceXML))

	return resource, nil
}
