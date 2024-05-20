// restrictedsoftware_state.go
package restrictedsoftware

import (
	"sort"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest RestrictedSoftware information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceRestrictedSoftware) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	if err := d.Set("id", strconv.Itoa(resource.General.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("process_name", resource.General.ProcessName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("match_exact_process_name", resource.General.MatchExactProcessName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("send_notification", resource.General.SendNotification); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("kill_process", resource.General.KillProcess); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("delete_executable", resource.General.DeleteExecutable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("display_message", resource.General.DisplayMessage); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	site := []interface{}{}
	if resource.General.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.General.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Update the scope data
	if err := d.Set("scope", flattenScope(resource.Scope)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// flattenScope converts the scope structure into a format suitable for setting in the Terraform state.
func flattenScope(scope jamfpro.RestrictedSoftwareSubsetScope) []interface{} {
	scopeMap := map[string]interface{}{
		"all_computers":      scope.AllComputers,
		"computer_ids":       flattenAndSortScopeEntityIds(scope.Computers),
		"computer_group_ids": flattenAndSortScopeEntityIds(scope.ComputerGroups),
		"building_ids":       flattenAndSortScopeEntityIds(scope.Buildings),
		"department_ids":     flattenAndSortScopeEntityIds(scope.Departments),
	}

	if len(scope.Exclusions.Computers) > 0 || len(scope.Exclusions.ComputerGroups) > 0 || len(scope.Exclusions.Buildings) > 0 || len(scope.Exclusions.Departments) > 0 || len(scope.Exclusions.Users) > 0 {
		scopeMap["exclusions"] = []interface{}{
			map[string]interface{}{
				"computer_ids":       flattenAndSortScopeEntityIds(scope.Exclusions.Computers),
				"computer_group_ids": flattenAndSortScopeEntityIds(scope.Exclusions.ComputerGroups),
				"building_ids":       flattenAndSortScopeEntityIds(scope.Exclusions.Buildings),
				"department_ids":     flattenAndSortScopeEntityIds(scope.Exclusions.Departments),
				"user_names":         flattenAndSortScopeEntityNames(scope.Exclusions.Users),
			},
		}
	}

	return []interface{}{scopeMap}
}

// flattenAndSortScopeEntityIds converts a slice of RestrictedSoftwareSubsetScopeEntity into a sorted slice of integers.
func flattenAndSortScopeEntityIds(entities []jamfpro.RestrictedSoftwareSubsetScopeEntity) []int {
	var ids []int
	for _, entity := range entities {
		ids = append(ids, entity.ID)
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortScopeEntityNames converts a slice of RestrictedSoftwareSubsetScopeEntity into a sorted slice of strings.
func flattenAndSortScopeEntityNames(entities []jamfpro.RestrictedSoftwareSubsetScopeEntity) []string {
	var names []string
	for _, entity := range entities {
		names = append(names, entity.Name)
	}
	sort.Strings(names)
	return names
}

// setScopeEntities converts a slice of jamfpro.RestrictedSoftwareSubsetScopeEntity structs into a slice of map[string]interface{} for Terraform.
func setScopeEntities(scopeEntities []jamfpro.RestrictedSoftwareSubsetScopeEntity) []interface{} {
	var tfScopeEntities []interface{}

	for _, entity := range scopeEntities {
		tfEntity := map[string]interface{}{
			"id":   entity.ID,
			"name": entity.Name,
		}
		tfScopeEntities = append(tfScopeEntities, tfEntity)
	}

	return tfScopeEntities
}
