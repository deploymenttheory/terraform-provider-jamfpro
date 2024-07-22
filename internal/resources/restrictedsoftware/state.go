// restrictedsoftware_state.go
package restrictedsoftware

import (
	"sort"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest RestrictedSoftware information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceRestrictedSoftware) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.General.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("process_name", resp.General.ProcessName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("match_exact_process_name", resp.General.MatchExactProcessName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("send_notification", resp.General.SendNotification); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("kill_process", resp.General.KillProcess); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("delete_executable", resp.General.DeleteExecutable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("display_message", resp.General.DisplayMessage); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.Set("site_id", resp.General.Site.ID)

	if err := d.Set("scope", flattenScope(resp.Scope)); err != nil {
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
				"computer_ids":                         flattenAndSortScopeEntityIds(scope.Exclusions.Computers),
				"computer_group_ids":                   flattenAndSortScopeEntityIds(scope.Exclusions.ComputerGroups),
				"building_ids":                         flattenAndSortScopeEntityIds(scope.Exclusions.Buildings),
				"department_ids":                       flattenAndSortScopeEntityIds(scope.Exclusions.Departments),
				"directory_service_or_local_usernames": flattenAndSortScopeEntityNames(scope.Exclusions.Users),
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
