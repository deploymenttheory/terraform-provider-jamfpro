// computergroup_state.go
package staticcomputergroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Group information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if resource == nil {
		return diags
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resource.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	if resource.Site != nil && resource.Site.ID != -1 {
		site := []interface{}{
			map[string]interface{}{
				"id": resource.Site.ID,
			},
		}
		if err := d.Set("site_id", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set the 'assignments' attribute in the state
	if resource.Computers != nil {
		computerIDs := []interface{}{}
		for _, comp := range *resource.Computers {
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
