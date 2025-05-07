// deviceenrollments_object.go
package deviceenrollments

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceDeviceEnrollment object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDeviceEnrollmentUpdate, error) {
	// Create the base resource
	resource := &jamfpro.ResourceDeviceEnrollmentUpdate{
		Name:                  d.Get("name").(string),
		SupervisionIdentityId: d.Get("supervision_identity_id").(string),
		SiteId:                d.Get("site_id").(string),
	}

	// Create token upload resource if encoded_token is provided
	tokenUpload := &jamfpro.ResourceDeviceEnrollmentTokenUpload{
		TokenFileName: d.Get("token_file_name").(string),
		EncodedToken:  d.Get("encoded_token").(string),
	}

	// Serialize and pretty-print the Device Enrollment object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Device Enrollment '%s' to JSON: %v", resource.Name, err)
	}

	// Log the resource structure (excluding sensitive token data)
	log.Printf("[DEBUG] Constructed Jamf Pro Device Enrollment JSON:\n%s\n", string(resourceJSON))

	// Store token upload in the schema meta for use during create/update
	if tokenUpload.EncodedToken != "" {
		d.Set("_token_upload", tokenUpload)
	}

	return resource, nil
}
