// categories_data_object.go
package categories

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProCategory constructs a Jamf Pro Category struct from Terraform resource data.
func constructJamfProCategory(d *schema.ResourceData) (*jamfpro.ResourceCategory, error) {
	// Assuming ResourceDepartment struct now also includes a Priority field
	department := &jamfpro.ResourceCategory{
		Name:     d.Get("name").(string),
		Priority: d.Get("priority").(int),
	}

	resourceXML, err := xml.MarshalIndent(department, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Department '%s' to XML: %v", department.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Department XML:\n%s\n", string(resourceXML))

	return department, nil
}
