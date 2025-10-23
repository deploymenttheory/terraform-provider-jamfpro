package access_management_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrMarshalAccessManagement = fmt.Errorf("failed to marshal Jamf Pro Access Management Settings to JSON")
)

// constructAccessManagementSettings constructs a ResourceAccessManagementSettings object from the provided schema data
func constructAccessManagementSettings(d *schema.ResourceData) (*jamfpro.ResourceAccessManagementSettings, error) {
	resource := &jamfpro.ResourceAccessManagementSettings{
		AutomatedDeviceEnrollmentServerUuid: d.Get("automated_device_enrollment_server_uuid").(string),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMarshalAccessManagement, err.Error())
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Access Management Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
