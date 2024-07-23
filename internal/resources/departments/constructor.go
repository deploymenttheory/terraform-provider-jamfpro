// department_data_object.go
package departments

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDepartment constructs a Jamf Pro Department struct from Terraform resource data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDepartment, error) {
	resource := &jamfpro.ResourceDepartment{
		Name: d.Get("name").(string),
	}
	log.Println("LOGHERE")
	jsonData, e := json.MarshalIndent(d, " ", "		")
	if e != nil {
		return nil, e
	}
	log.Printf("HERE: %+v", string(jsonData))

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Department '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Department XML:\n%s\n", string(resourceXML))

	return resource, nil
}
