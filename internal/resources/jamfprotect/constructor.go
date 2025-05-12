// jamf_protect_constructor.go
package jamfprotect

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs the Jamf Protect registration request from the Terraform schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceJamfProtectRegisterRequest, error) {
	resource := &jamfpro.ResourceJamfProtectRegisterRequest{
		ProtectURL: d.Get("protect_url").(string),
		ClientID:   d.Get("client_id").(string),
		Password:   d.Get("password").(string),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Protect registration request to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Protect registration request JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
