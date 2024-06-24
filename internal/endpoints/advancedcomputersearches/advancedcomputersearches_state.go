// advancedcomputersearches_state.go
package advancedcomputersearches

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Advanced Computer Search information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceAdvancedComputerSearch) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("view_as", resp.ViewAs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort1", resp.Sort1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort2", resp.Sort2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort3", resp.Sort3); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	criteriaList := make([]interface{}, len(resp.Criteria.Criterion))
	for i, crit := range resp.Criteria.Criterion {
		criteriaMap := map[string]interface{}{
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

	if len(resp.DisplayFields) == 0 || len(resp.DisplayFields[0].DisplayField) == 0 {
		if err := d.Set("display_fields", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		displayFieldsList := make([]map[string]interface{}, len(resp.DisplayFields[0].DisplayField))
		for i, displayField := range resp.DisplayFields[0].DisplayField {
			displayFieldMap := map[string]interface{}{
				"name": displayField.Name,
			}
			displayFieldsList[i] = displayFieldMap
		}
		if err := d.Set("display_fields", displayFieldsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	d.Set("site_id", resp.Site.ID)

	return diags
}
