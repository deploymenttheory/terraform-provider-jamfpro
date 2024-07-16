// computergroup_state.go
package staticcomputergroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Group information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resp.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.Set("site_id", resp.Site.ID)

	var assignments []interface{}
	if resp.Computers != nil {
		for _, comp := range *resp.Computers {
			assignments = append(assignments, comp.ID)
		}

		if err := d.Set("assigned_computer_ids", assignments); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
