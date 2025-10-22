// advancedmobiledevicesearches_resource.go
package advanced_mobile_device_search

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Advanced Mobile Device Search information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceAdvancedMobileDeviceSearch) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", resp.ID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set criteria
	if len(resp.Criteria) > 0 {
		criteriaList := make([]any, len(resp.Criteria))
		for i, crit := range resp.Criteria {
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
	}

	if len(resp.DisplayFields) > 0 {
		if err := d.Set("display_fields", resp.DisplayFields); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if resp.SiteId != nil {
		if err := d.Set("site_id", *resp.SiteId); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
