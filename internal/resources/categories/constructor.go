// categories_data_object.go
package categories

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProCategory constructs a Jamf Pro Category struct from Terraform resource data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceCategory, error) {
	var resource *jamfpro.ResourceCategory

	resource = &jamfpro.ResourceCategory{
		Name:     d.Get("name").(string),
		Priority: d.Get("priority").(int),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Category '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Category XML:\n%s\n", string(resourceXML))

	return resource, nil
}
