package self_service_plus_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrMarshalSelfServicePlus = fmt.Errorf("failed to marshal Jamf Pro Self Service Plus Settings to JSON")
)

// constructSelfServicePlus constructs a ResourceSelfServicePlus object from the provided schema data
func constructSelfServicePlusSettings(d *schema.ResourceData) (*jamfpro.ResourceSelfServicePlusSettings, error) {
	resource := &jamfpro.ResourceSelfServicePlusSettings{
		Enabled: d.Get("enabled").(bool),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMarshalSelfServicePlus, err.Error())
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Self Service Plus Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
