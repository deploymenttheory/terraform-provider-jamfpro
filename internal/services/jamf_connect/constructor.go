package jamf_connect

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct creates a ResourceJamfConnectConfigProfileUpdate from the provided schema data
func construct(d *schema.ResourceData) (*jamfpro.ResourceJamfConnectConfigProfileUpdate, error) {
	resource := &jamfpro.ResourceJamfConnectConfigProfileUpdate{
		JamfConnectVersion: d.Get("jamf_connect_version").(string),
		AutoDeploymentType: d.Get("auto_deployment_type").(string),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Connect Config Profile to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Connect Config Profile JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
