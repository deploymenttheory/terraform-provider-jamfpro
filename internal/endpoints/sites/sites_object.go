// sites_data_object.go
package sites

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProSite constructs a SharedResourceSite object from the provided schema data.
func constructJamfProSite(d *schema.ResourceData) (*jamfpro.SharedResourceSite, error) {
	site := &jamfpro.SharedResourceSite{
		Name: d.Get("name").(string),
	}

	resourceXML, err := xml.MarshalIndent(site, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Site '%s' to XML: %v", site.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Site XML:\n%s\n", string(resourceXML))

	return site, nil
}
