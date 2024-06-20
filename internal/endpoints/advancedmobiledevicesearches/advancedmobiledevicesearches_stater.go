// advancedmobiledevicesearches_resource.go
package advancedmobiledevicesearches

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Advanced User Search information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceAdvancedMobileDeviceSearch) diag.Diagnostics {

	var diags diag.Diagnostics
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("view_as", resource.ViewAs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort1", resource.Sort1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort2", resource.Sort2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort3", resource.Sort3); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	criteriaList := make([]interface{}, len(resource.Criteria.Criterion))
	for i, crit := range resource.Criteria.Criterion {
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

	if len(resource.DisplayFields) == 0 || len(resource.DisplayFields[0].DisplayField) == 0 {
		if err := d.Set("display_fields", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		displayFieldsList := make([]map[string]interface{}, len(resource.DisplayFields[0].DisplayField))
		for i, displayField := range resource.DisplayFields[0].DisplayField {
			displayFieldMap := map[string]interface{}{
				"name": displayField.Name,
			}
			displayFieldsList[i] = displayFieldMap
		}
		if err := d.Set("display_fields", displayFieldsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	site := []interface{}{}
	if resource.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags

}
