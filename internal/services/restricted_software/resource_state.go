// restrictedsoftware_state.go
package restricted_software

import (
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
func flattenScope(scope jamfpro.RestrictedSoftwareSubsetScope) []any {
	scopeMap := map[string]any{
		"all_computers":      scope.AllComputers,
		"computer_ids":       schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Computers)),
		"computer_group_ids": schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.ComputerGroups)),
		"building_ids":       schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Buildings)),
		"department_ids":     schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Departments)),
	}

	// Handle Exclusions
	if len(scope.Exclusions.Computers) > 0 || len(scope.Exclusions.ComputerGroups) > 0 ||
		len(scope.Exclusions.Buildings) > 0 || len(scope.Exclusions.Departments) > 0 ||
		len(scope.Exclusions.Users) > 0 {
		scopeMap["exclusions"] = []any{
			map[string]any{
				"computer_ids":                         schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Exclusions.Computers)),
				"computer_group_ids":                   schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Exclusions.ComputerGroups)),
				"building_ids":                         schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Exclusions.Buildings)),
				"department_ids":                       schema.NewSet(schema.HashInt, flattenScopeEntityIds(scope.Exclusions.Departments)),
				"directory_service_or_local_usernames": schema.NewSet(schema.HashString, flattenScopeEntityNames(scope.Exclusions.Users)),
			},
		}
	}

	return []any{scopeMap}
}

// flattenScopeEntityIds converts a slice of RestrictedSoftwareSubsetScopeEntity into a slice of interfaces containing IDs
func flattenScopeEntityIds(entities []jamfpro.RestrictedSoftwareSubsetScopeEntity) []any {
	var ids []any
	for _, entity := range entities {
		ids = append(ids, entity.ID)
	}
	return ids
}

// flattenScopeEntityNames converts a slice of RestrictedSoftwareSubsetScopeEntity into a slice of interfaces containing names
func flattenScopeEntityNames(entities []jamfpro.RestrictedSoftwareSubsetScopeEntity) []any {
	var names []any
	for _, entity := range entities {
		names = append(names, entity.Name)
	}
	return names
}
