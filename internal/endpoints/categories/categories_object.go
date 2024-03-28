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
func constructJamfProCategory(d *schema.ResourceData) (*jamfpro.ResourceCategory, error) {

	category := &jamfpro.ResourceCategory{
		Name:     d.Get("name").(string),
		Priority: d.Get("priority").(int),
	}

	// Serialize and pretty-print the Category object as XML for logging
	resourceXML, err := xml.MarshalIndent(category, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Category '%s' to XML: %v", category.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Category XML:\n%s\n", string(resourceXML))

	return category, nil
}
