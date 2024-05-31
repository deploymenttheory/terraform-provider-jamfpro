// smartcomputergroup_state.go
package smartcomputergroups

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the provided ResourceComputerGroup object.
func updateTerraformState(d *schema.ResourceData, group *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", group.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", group.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Site
	if group.Site != nil {
		site := flattenSharedResourceSite(group.Site)
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("site", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Handle Criteria
	if group.Criteria != nil && group.Criteria.Criterion != nil {
		criteria := flattenComputerGroupSubsetContainerCriteria(group.Criteria)
		if err := d.Set("criteria", criteria); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("criteria", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Handle ID
	if err := d.Set("id", strconv.Itoa(group.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// flattenSharedResourceSite flattens a SharedResourceSite object into a format suitable for Terraform state.
func flattenSharedResourceSite(site *jamfpro.SharedResourceSite) []interface{} {
	if site == nil {
		return nil
	}

	siteMap := map[string]interface{}{
		"id":   site.ID,
		"name": site.Name,
	}

	return []interface{}{siteMap}
}

// flattenComputerGroupSubsetContainerCriteria flattens a ComputerGroupSubsetContainerCriteria object into a format suitable for Terraform state.
func flattenComputerGroupSubsetContainerCriteria(criteria *jamfpro.ComputerGroupSubsetContainerCriteria) []interface{} {
	if criteria == nil || criteria.Criterion == nil {
		return nil
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
