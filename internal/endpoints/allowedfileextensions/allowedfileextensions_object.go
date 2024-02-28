// allowedfileextensions_object.go
package allowedfileextensions

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAllowedFileExtension creates a new ResourceAllowedFileExtension instance from Terraform data and serializes it to XML.
func constructJamfProAllowedFileExtension(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAllowedFileExtension, error) {
	allowedFileExtension := &jamfpro.ResourceAllowedFileExtension{
		Extension: d.Get("extension").(string),
	}

	// Serialize and pretty-print the allowedFileExtension object as XML
	resourceXML, err := xml.MarshalIndent(allowedFileExtension, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Allowed File Extension '%s' to XML: %v", allowedFileExtension.Extension, err)
	}
	fmt.Printf("Constructed Jamf Pro Allowed File Extension XML:\n%s\n", string(resourceXML))

	return allowedFileExtension, nil
}
