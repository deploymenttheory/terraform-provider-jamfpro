package apiroles

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProApiRole constructs an ResourceAPIRole object from the provided schema data.
func constructJamfProApiRole(d *schema.ResourceData) (*jamfpro.ResourceAPIRole, error) {
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

	// Serialize and pretty-print the Api Role object as JSON for logging
	resourceJSON, err := json.MarshalIndent(apiRole, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Api Role '%s' to JSON: %v", apiRole.DisplayName, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Api Role JSON:\n%s\n", string(resourceJSON))

	return apiRole, nil
}
