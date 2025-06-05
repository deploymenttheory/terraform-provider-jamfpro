// activationcode_object.go
package activation_code

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProActivationCode constructs a ResourceActivationCode object from the provided schema data and logs its XML representation.
func construct(d *schema.ResourceData) (*jamfpro.ResourceActivationCode, error) {
	resource := &jamfpro.ResourceActivationCode{
		OrganizationName: d.Get("organization_name").(string),
		Code:             d.Get("code").(string),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Activation Code to XML: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Activation Code XML:\n%s\n", string(resourceXML))

	return resource, nil
}
