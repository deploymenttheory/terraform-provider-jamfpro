package policy

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Reads response and states scope items.
// TODO reduce the cyclomatic complexity here by centralising the repeated slice-handling logic (gocyclo)
// TODO preallocate the various output/list slices to address golangci-lint prealloc warnings
func stateScope(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error

	out_scope := make([]map[string]any, 0)
	out_scope = append(out_scope, make(map[string]any, 1))
	out_scope[0]["all_computers"] = resp.Scope.AllComputers
	out_scope[0]["all_jss_users"] = resp.Scope.AllJSSUsers

	// TODO see if we can simplify/centralise the repeated logic below
	// Computers
	if resp.Scope.Computers != nil && len(*resp.Scope.Computers) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_ids"] = listOfIds
	}

	// Computer Groups
	if resp.Scope.ComputerGroups != nil && len(*resp.Scope.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_group_ids"] = listOfIds
	}

	// JSS Users
	if resp.Scope.JSSUsers != nil && len(*resp.Scope.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_ids"] = listOfIds
	}

	// JSS User Groups
	if resp.Scope.JSSUserGroups != nil && len(*resp.Scope.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_group_ids"] = listOfIds
	}

	// Buildings
	if resp.Scope.Buildings != nil && len(*resp.Scope.Buildings) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["building_ids"] = listOfIds
	}

	// Departments
	if resp.Scope.Departments != nil && len(*resp.Scope.Departments) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["department_ids"] = listOfIds
	}

	// Scope Limitations
	out_scope_limitations := make([]map[string]any, 0)
	out_scope_limitations = append(out_scope_limitations, make(map[string]any))
	var limitationsSet bool

	// Users
	if resp.Scope.Limitations.Users != nil && len(*resp.Scope.Limitations.Users) > 0 {
		var listOfNames []string
		for _, v := range *resp.Scope.Limitations.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_limitations[0]["directory_service_or_local_usernames"] = listOfNames
		limitationsSet = true
	}

	// Network Segments
	if resp.Scope.Limitations.NetworkSegments != nil && len(*resp.Scope.Limitations.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Limitations.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["network_segment_ids"] = listOfIds
		limitationsSet = true
	}

	// IBeacons
	if resp.Scope.Limitations.IBeacons != nil && len(*resp.Scope.Limitations.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Limitations.IBeacons {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["ibeacon_ids"] = listOfIds
		limitationsSet = true
	}

	// User Groups

	if resp.Scope.Limitations.UserGroups != nil && len(*resp.Scope.Limitations.UserGroups) > 0 {
		var listOfNames []string
		for _, v := range *resp.Scope.Limitations.UserGroups {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_limitations[0]["directory_service_usergroup_names"] = listOfNames
		limitationsSet = true
	}

	if limitationsSet {
		out_scope[0]["limitations"] = out_scope_limitations
	}

	// Scope Exclusions
	out_scope_exclusions := make([]map[string]any, 0)
	out_scope_exclusions = append(out_scope_exclusions, make(map[string]any))
	var exclusionsSet bool

	// Computers
	if resp.Scope.Exclusions.Computers != nil && len(*resp.Scope.Exclusions.Computers) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_ids"] = listOfIds
		exclusionsSet = true
	}

	// Computer Groups
	if resp.Scope.Exclusions.ComputerGroups != nil && len(*resp.Scope.Exclusions.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_group_ids"] = listOfIds
		exclusionsSet = true
	}

	// Buildings
	if resp.Scope.Exclusions.Buildings != nil && len(*resp.Scope.Exclusions.Buildings) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["building_ids"] = listOfIds
		exclusionsSet = true
	}

	// Departments
	if resp.Scope.Exclusions.Departments != nil && len(*resp.Scope.Exclusions.Departments) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["department_ids"] = listOfIds
		exclusionsSet = true
	}

	// Users
	if resp.Scope.Exclusions.Users != nil && len(*resp.Scope.Exclusions.Users) > 0 {
		var listOfNames []string
		for _, v := range *resp.Scope.Exclusions.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_exclusions[0]["directory_service_or_local_usernames"] = listOfNames
		exclusionsSet = true
	}

	// User Groups
	if resp.Scope.Exclusions.UserGroups != nil && len(*resp.Scope.Exclusions.UserGroups) > 0 {
		var listOfNames []string
		for _, v := range *resp.Scope.Exclusions.UserGroups {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_exclusions[0]["directory_service_usergroup_names"] = listOfNames
		exclusionsSet = true
	}

	// Network Segments
	if resp.Scope.Exclusions.NetworkSegments != nil && len(*resp.Scope.Exclusions.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["network_segment_ids"] = listOfIds
		exclusionsSet = true
	}

	// JSS Users
	if resp.Scope.Exclusions.JSSUsers != nil && len(*resp.Scope.Exclusions.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_ids"] = listOfIds
		exclusionsSet = true
	}

	// JSS User Groups
	if resp.Scope.Exclusions.JSSUserGroups != nil && len(*resp.Scope.Exclusions.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_group_ids"] = listOfIds
		exclusionsSet = true
	}

	// IBeacons
	if resp.Scope.Exclusions.IBeacons != nil && len(*resp.Scope.Exclusions.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range *resp.Scope.Exclusions.IBeacons {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["ibeacon_ids"] = listOfIds
		exclusionsSet = true
	}

	// Append Exclusions if they're set
	if exclusionsSet {
		out_scope[0]["exclusions"] = out_scope_exclusions
	} else {
		log.Println("No exclusions set") // TODO logging
	}

	// State Scope
	err = d.Set("scope", out_scope)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
