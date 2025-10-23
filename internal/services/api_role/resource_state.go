// apiroles_state.go
package api_role

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest API Role information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceAPIRole) diag.Diagnostics {
	var diags diag.Diagnostics

	apiRoleData := map[string]any{
		"id":           resp.ID,
		"display_name": resp.DisplayName,
		"privileges":   resp.Privileges,
	}

	for key, val := range apiRoleData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set '%s': %v", key, err))...)
		}
	}

	return diags
}
