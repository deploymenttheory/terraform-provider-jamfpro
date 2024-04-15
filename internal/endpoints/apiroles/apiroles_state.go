// apiroles_state.go
package apiroles

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Building information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceAPIRole) diag.Diagnostics {

	var diags diag.Diagnostics

	// Map the configuration fields from the API response to a structured map
	apiRoleData := map[string]interface{}{
		"id":           resource.ID,
		"display_name": resource.DisplayName,
		"privileges":   resource.Privileges,
	}

	// Set the structured map in the Terraform state
	for key, val := range apiRoleData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set '%s': %v", key, err))...)
		}
	}

	return diags
}
