// computergroup_state.go
package computergroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Prestage Enrollment information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	if resource != nil {
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_smart", resource.IsSmart); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		site := map[string]interface{}{
			"id":   resource.Site.ID,
			"name": resource.Site.Name,
		}
		if err := d.Set("site", []interface{}{site}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		// Set the criteria
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

		// Set the computers only if the group is not smart
		if !resource.IsSmart {
			computersList := make([]interface{}, len(resource.Computers))
			for i, comp := range resource.Computers {
				computerMap := map[string]interface{}{
					"id":              comp.ID,
					"name":            comp.Name,
					"mac_address":     comp.MacAddress,
					"alt_mac_address": comp.AltMacAddress,
					"serial_number":   comp.SerialNumber,
				}
				computersList[i] = computerMap
			}
			if err := d.Set("computers", computersList); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}
	return diags
}
