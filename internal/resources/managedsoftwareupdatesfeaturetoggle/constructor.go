package managedsoftwareupdatesfeaturetoggle

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructManagedSoftwareUpdatesFeatureToggle constructs a Managed Software Updates Feature Toggle object from the provided schema data
func constrconstructManagedSoftwareUpdatesFeatureToggle(d *schema.ResourceData) (*jamfpro.ResourceManagedSoftwareUpdateFeatureToggle, error) {
	resource := &jamfpro.ResourceManagedSoftwareUpdateFeatureToggle{
		Toggle: d.Get("toggle").(bool),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Managed Software Feature Toggle to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Managed Software Feature Toggle JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
