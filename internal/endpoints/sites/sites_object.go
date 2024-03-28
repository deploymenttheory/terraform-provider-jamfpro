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
	site := &jamfpro.SharedResourceSite{
		Name: d.Get("name").(string),
	}

	// Serialize and pretty-print the Site object as XML for logging
	resourceXML, err := xml.MarshalIndent(site, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Site '%s' to XML: %v", site.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Site XML:\n%s\n", string(resourceXML))

	return site, nil
}
