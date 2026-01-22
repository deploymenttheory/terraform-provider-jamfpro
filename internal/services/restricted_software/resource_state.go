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

	// General/Root level
	stateGeneral(d, resp, &diags)

	// Scope
	stateScope(d, resp, &diags)

	return diags
}

// stateGeneral reads response and states general/root level items
func stateGeneral(d *schema.ResourceData, resp *jamfpro.ResourceRestrictedSoftware, diags *diag.Diagnostics) {
	var err error

	err = d.Set("name", resp.General.Name)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("process_name", resp.General.ProcessName)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("match_exact_process_name", resp.General.MatchExactProcessName)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("send_notification", resp.General.SendNotification)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("kill_process", resp.General.KillProcess)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("delete_executable", resp.General.DeleteExecutable)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("display_message", resp.General.DisplayMessage)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("site_id", resp.General.Site.ID)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// stateScope reads response and states scope items
func stateScope(d *schema.ResourceData, resp *jamfpro.ResourceRestrictedSoftware, diags *diag.Diagnostics) {
	var err error

	out_scope := make([]map[string]any, 0)
	out_scope = append(out_scope, make(map[string]any, 1))
	out_scope[0]["all_computers"] = resp.Scope.AllComputers

	// Computers
	if len(resp.Scope.Computers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_ids"] = listOfIds
	}

	// Computer Groups
	if len(resp.Scope.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_group_ids"] = listOfIds
	}

	// Buildings
	if len(resp.Scope.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["building_ids"] = listOfIds
	}

	// Departments
	if len(resp.Scope.Departments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["department_ids"] = listOfIds
	}

	// Scope Exclusions
	out_scope_exclusions := make([]map[string]any, 0)
	out_scope_exclusions = append(out_scope_exclusions, make(map[string]any))
	var exclusionsSet bool

	// Computers
	if len(resp.Scope.Exclusions.Computers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_ids"] = listOfIds
		exclusionsSet = true
	}

	// Computer Groups
	if len(resp.Scope.Exclusions.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_group_ids"] = listOfIds
		exclusionsSet = true
	}

	// Buildings
	if len(resp.Scope.Exclusions.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["building_ids"] = listOfIds
		exclusionsSet = true
	}

	// Departments
	if len(resp.Scope.Exclusions.Departments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["department_ids"] = listOfIds
		exclusionsSet = true
	}

	// Users (directory_service_or_local_usernames)
	if len(resp.Scope.Exclusions.Users) > 0 {
		var listOfNames []string
		for _, v := range resp.Scope.Exclusions.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_exclusions[0]["directory_service_or_local_usernames"] = listOfNames
		exclusionsSet = true
	}

	// Append Exclusions if they're set
	if exclusionsSet {
		out_scope[0]["exclusions"] = out_scope_exclusions
	}

	// State Scope
	err = d.Set("scope", out_scope)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
