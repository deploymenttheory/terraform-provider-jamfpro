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

	// TODO review this.
	if resp.Computers != nil {
		computerIDs := []interface{}{}
		for _, comp := range *resp.Computers {
			computerIDs = append(computerIDs, comp.ID)
		}
		assignments := []interface{}{
			map[string]interface{}{
				"computer_ids": computerIDs,
			},
		}
		if err := d.Set("assignments", assignments); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
