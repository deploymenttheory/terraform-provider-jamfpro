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

	// Handle Criteria
	if resource.Criteria != nil && resource.Criteria.Criterion != nil {
		criteria := setComputerGroupSubsetContainerCriteria(resource.Criteria)
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

// setComputerGroupSubsetContainerCriteria flattens a ComputerGroupSubsetContainerCriteria object into a format suitable for Terraform state.
func setComputerGroupSubsetContainerCriteria(criteria *jamfpro.ComputerGroupSubsetContainerCriteria) []interface{} {
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