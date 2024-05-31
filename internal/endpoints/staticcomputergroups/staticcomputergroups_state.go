// staticcomputergroup_state.go
package staticcomputergroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the provided ResourceComputerGroup object.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resource.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Site
	site := []interface{}{}
	if resource.Site != nil && resource.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id":   resource.Site.ID,
			"name": resource.Site.Name,
		})
	}

	if err := d.Set("site", site); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the assignments attribute
	if resource.Computers != nil && resource.Computers.Computers != nil {
		computerIDs := make([]int, len(*resource.Computers.Computers))
		for i, computer := range *resource.Computers.Computers {
			computerIDs[i] = computer.ID
		}
		assignments := map[string]interface{}{
			"computer_ids": computerIDs,
		}
		if err := d.Set("assignments", []interface{}{assignments}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("assignments", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
