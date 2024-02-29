// apiintegrations_data_object.go
package apiintegrations

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProApiIntegration constructs a ResourceApiIntegration object from the provided schema data and serializes it to XML.
func constructJamfProApiIntegration(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceApiIntegration, error) {
	integration := &jamfpro.ResourceApiIntegration{
		DisplayName:                d.Get("display_name").(string),
		Enabled:                    d.Get("enabled").(bool),
		AccessTokenLifetimeSeconds: d.Get("access_token_lifetime_seconds").(int),
	}

	// Handle 'authorization_scopes' field directly without helper functions
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
		integration.AuthorizationScopes = authorizationScopes
	}

	// Serialize and pretty-print the integration object as XML
	resourceXML, err := xml.MarshalIndent(integration, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Api Integration '%s' to XML: %v", integration.DisplayName, err)
	}
	fmt.Printf("Constructed Jamf Pro Api Integration XML:\n%s\n", string(resourceXML))

	return integration, nil
}
