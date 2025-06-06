// deviceenrollments_object.go
package device_enrollments

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceDeviceEnrollment object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDeviceEnrollmentUpdate, error) {
	resource := &jamfpro.ResourceDeviceEnrollmentUpdate{
		Name:                  d.Get("name").(string),
		SupervisionIdentityId: d.Get("supervision_identity_id").(string),
		SiteId:                d.Get("site_id").(string),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Device Enrollment '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Device Enrollment JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
