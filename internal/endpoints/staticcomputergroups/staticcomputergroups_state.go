// computergroup_state.go
package staticcomputergroups

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

		// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
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

		// Set the 'assignments' attribute in the state
		assignmentsList := []interface{}{}
		if resource.Computers != nil {
			computerIDs := []interface{}{}
			for _, comp := range *resource.Computers {
				computerIDs = append(computerIDs, comp.ID)
			}
			assignments := map[string]interface{}{
				"computer_ids": computerIDs,
			}
			assignmentsList = append(assignmentsList, assignments)
		}
		if err := d.Set("assignments", assignmentsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}
