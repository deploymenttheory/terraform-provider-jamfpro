package apiroles

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProApiRole constructs an ResourceAPIRole object from the provided schema data.
func constructJamfProApiRole(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAPIRole, error) {
	apiRole := &jamfpro.ResourceAPIRole{
		DisplayName: d.Get("display_name").(string),
	}

	// Handle 'privileges' field directly without type assertion helper functions
	if v, ok := d.GetOk("privileges"); ok {
		// Convert privileges from interface{} to []interface{} and then to []string
		privilegesInterface := v.(*schema.Set).List()
		privileges := make([]string, len(privilegesInterface))
		for i, priv := range privilegesInterface {
			privileges[i], ok = priv.(string)
			if !ok {
				return nil, fmt.Errorf("failed to assert privilege to string")
			}
		}
		apiRole.Privileges = privileges
	}

	// Serialize and pretty-print the site object as XML
	resourceXML, err := xml.MarshalIndent(apiRole, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Building '%s' to XML: %v", apiRole.DisplayName, err)
	}
	fmt.Printf("Constructed Jamf Pro Building XML:\n%s\n", string(resourceXML))

	return apiRole, nil
}
