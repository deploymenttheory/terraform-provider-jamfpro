// advancedusersearches_state.go
package advanced_user_search

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Advanced User Search information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceAdvancedUserSearch) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	criteriaList := make([]any, len(resp.Criteria.Criterion))
	for i, crit := range resp.Criteria.Criterion {
		criteriaMap := map[string]any{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        crit.AndOr,
			"search_type":   crit.SearchType,
			"value":         crit.Value,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
		criteriaList[i] = criteriaMap
	}

	if err := d.Set("criteria", criteriaList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if len(resp.DisplayFields) > 0 {
		displayFieldSet := schema.NewSet(schema.HashString, nil)
		for _, v := range resp.DisplayFields {
			displayFieldSet.Add(v.Name)
		}
		if err := d.Set("display_fields", displayFieldSet); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	d.Set("site_id", resp.Site.ID)

	return diags

}
