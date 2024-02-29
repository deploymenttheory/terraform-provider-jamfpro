// department_data_object.go
package departments

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDepartment constructs a Jamf Pro Department struct from Terraform resource data.
func constructJamfProDepartment(d *schema.ResourceData) (*jamfpro.ResourceDepartment, error) {
	department := &jamfpro.ResourceDepartment{
		Name: d.Get("name").(string),
	}

	resourceXML, err := xml.MarshalIndent(department, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Department '%s' to XML: %v", department.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Department XML:\n%s\n", string(resourceXML))

	return department, nil
}
