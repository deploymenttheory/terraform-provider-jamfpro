// allowedfileextensions_object.go
package allowedfileextensions

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAllowedFileExtension creates a new ResourceAllowedFileExtension instance from Terraform data and serializes it to XML.
func construct(d *schema.ResourceData) (*jamfpro.ResourceAllowedFileExtension, error) {
	resource := &jamfpro.ResourceAllowedFileExtension{
		Extension: d.Get("extension").(string),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Allowed File Extension '%s' to XML: %v", resource.Extension, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Allowed File Extension XML:\n%s\n", string(resourceXML))

	return resource, nil
}
