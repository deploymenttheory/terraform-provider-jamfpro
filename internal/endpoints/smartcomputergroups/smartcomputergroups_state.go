// smartcomputergroup_state.go
package smartcomputergroups

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
	if resource.Site != nil && resource.Site.ID != -1 {
		site := []interface{}{
			map[string]interface{}{
				"id": resource.Site.ID,
			},
		}
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Handle Criteria
	if resource.Criteria != nil && resource.Criteria.Criterion != nil {
		criteria := setComputerSmartGroupSubsetContainerCriteria(resource.Criteria)
		if err := d.Set("criteria", criteria); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("criteria", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setComputerSmartGroupSubsetContainerCriteria flattens a ComputerGroupSubsetContainerCriteria object into a format suitable for Terraform state.
func setComputerSmartGroupSubsetContainerCriteria(criteria *jamfpro.ComputerGroupSubsetContainerCriteria) []interface{} {
	if criteria == nil || criteria.Criterion == nil {
		return []interface{}{}
	}

	var criteriaList []interface{}
	for _, criterion := range *criteria.Criterion {
		criterionMap := map[string]interface{}{
			"name":          criterion.Name,
			"priority":      criterion.Priority,
			"and_or":        criterion.AndOr,
			"search_type":   criterion.SearchType,
			"value":         criterion.Value,
			"opening_paren": criterion.OpeningParen,
			"closing_paren": criterion.ClosingParen,
		}
		criteriaList = append(criteriaList, criterionMap)
	}

	return criteriaList
}
