// apiintegrations_data_object.go
package api_integration

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProApiIntegration constructs a ResourceApiIntegration object from the provided schema data and serializes it to JSON.
func construct(d *schema.ResourceData) (*jamfpro.ResourceApiIntegration, error) {
	resource := &jamfpro.ResourceApiIntegration{
		DisplayName:                d.Get("display_name").(string),
		Enabled:                    d.Get("enabled").(bool),
		AccessTokenLifetimeSeconds: d.Get("access_token_lifetime_seconds").(int),
	}

	if v, ok := d.GetOk("authorization_scopes"); ok {
		scopesList := v.(*schema.Set).List()
		authorizationScopes := make([]string, len(scopesList))
		for i, scope := range scopesList {
			scopeStr, ok := scope.(string)
			if !ok {
				return nil, fmt.Errorf("failed to assert authorization scope to string")
			}
			authorizationScopes[i] = scopeStr
		}
		resource.AuthorizationScopes = authorizationScopes
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Api Integration '%s' to JSON: %v", resource.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Api Integration JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
