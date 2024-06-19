// apiintegrations_state.go
package apiintegrations

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest API Integration information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceApiIntegration) diag.Diagnostics {
	var diags diag.Diagnostics

	apiIntegrationData := map[string]interface{}{
		"display_name":                  resource.DisplayName,
		"enabled":                       resource.Enabled,
		"access_token_lifetime_seconds": resource.AccessTokenLifetimeSeconds,
		"app_type":                      resource.AppType,
		"authorization_scopes":          resource.AuthorizationScopes,
		"client_id":                     resource.ClientID,
	}

	for key, val := range apiIntegrationData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set '%s': %v", key, err))...)
		}
	}

	return diags

}
