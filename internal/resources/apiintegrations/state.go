// apiintegrations_state.go
package apiintegrations

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest API Integration information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceApiIntegration) diag.Diagnostics {
	var diags diag.Diagnostics

	apiIntegrationData := map[string]interface{}{
		"display_name":                  resp.DisplayName,
		"enabled":                       resp.Enabled,
		"access_token_lifetime_seconds": resp.AccessTokenLifetimeSeconds,
		"app_type":                      resp.AppType,
		"authorization_scopes":          resp.AuthorizationScopes,
		"client_id":                     resp.ClientID,
	}

	for key, val := range apiIntegrationData {
		if err := d.Set(key, val); err != nil {

			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set '%s': %v", key, err))...)
		}
	}

	return diags

}
