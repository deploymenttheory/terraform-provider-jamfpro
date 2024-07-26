// categories_state.go
package categories

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Category information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceCategory) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", resp.Id); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("priority", resp.Priority); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
