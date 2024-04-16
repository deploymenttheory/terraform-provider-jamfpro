// categories_state.go
package categories

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Category information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceCategory) diag.Diagnostics {

	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	if resource != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", resource.Id); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("priority", resource.Priority); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
