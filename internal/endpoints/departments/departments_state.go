// department_state.go
package departments

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Department information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceDepartment) diag.Diagnostics {
	var diags diag.Diagnostics

	if resource != nil {
		if err := d.Set("id", resource.ID); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
