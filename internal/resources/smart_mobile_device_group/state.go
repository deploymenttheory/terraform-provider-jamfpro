package smart_mobile_device_group

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the provided ResourceMobileDeviceGroup object.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDeviceGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("site_id", resp.Site.ID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resp.Criteria.Size != 0 && resp.Criteria.Criterion != nil {
		criteria := setMobileSmartGroupSubsetContainerCriteria(resp.Criteria)
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

// setMobileSmartGroupSubsetContainerCriteria flattens a MobileGroupSubsetContainerCriteria object into a format suitable for Terraform state.
func setMobileSmartGroupSubsetContainerCriteria(criteria jamfpro.SharedContainerCriteria) []interface{} {
	// TODO Review this!
	if criteria.Size == 0 || criteria.Criterion == nil {
		return []interface{}{}
	}

	var criteriaList []interface{}
	for _, criterion := range criteria.Criterion {
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
