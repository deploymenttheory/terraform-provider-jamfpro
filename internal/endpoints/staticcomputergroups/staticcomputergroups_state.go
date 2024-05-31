// staticcomputergroup_state.go
package staticcomputergroups

import (
	"sort"

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

	// Handle Computers
	if resource.Computers != nil && resource.Computers.Computers != nil {
		computerIDs := flattenAndSortComputerIds(*resource.Computers.Computers)
		if err := d.Set("computer_ids", computerIDs); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("computer_ids", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// flattenAndSortComputerIds converts a slice of ComputerGroupSubsetComputer into a sorted slice of integers.
func flattenAndSortComputerIds(computers []jamfpro.ComputerGroupSubsetComputer) []int {
	var ids []int
	for _, computer := range computers {
		if computer.ID != 0 {
			ids = append(ids, computer.ID)
		}
	}
	sort.Ints(ids)
	return ids
}
