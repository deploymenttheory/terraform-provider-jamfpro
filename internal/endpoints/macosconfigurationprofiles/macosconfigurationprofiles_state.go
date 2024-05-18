// macosconfigurationprofiles_state.go
package macosconfigurationprofiles

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest MacOS Configuration Profile information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceMacOSConfigurationProfile, resourceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Stating - commented ones appear to be done automatically.

	// ID
	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Name
	if err := d.Set("name", resource.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Description
	if err := d.Set("description", resource.General.Description); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Site
	if resource.General.Site.ID != -1 && resource.General.Site.Name != "None" {
		out_site := []map[string]interface{}{
			{
				"id":   resource.General.Site.ID,
				"name": resource.General.Site.Name,
			},
		}

		if err := d.Set("site", out_site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default site response") // TODO logging
	}

	// Category
	if resource.General.Category.ID != -1 && resource.General.Category.Name != "No category assigned" {
		out_category := []map[string]interface{}{
			{
				"id":   resource.General.Category.ID,
				"name": resource.General.Category.Name,
			},
		}
		if err := d.Set("category", out_category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default category response") // TODO logging
	}

	// Payloads
	profile := sharedschemas.NormalizePayloadState(resource.General.Payloads)
	if err := d.Set("payload", profile); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Distribution Method
	if err := d.Set("distribution_method", resource.General.DistributionMethod); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// User Removeable
	if err := d.Set("user_removeable", resource.General.UserRemovable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Level
	if err := d.Set("level", resource.General.Level); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// UUID
	if err := d.Set("uuid", resource.General.UUID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Redeploy On Update - This is always "Newly Assigned" on existing profile objects
	if err := d.Set("redeploy_on_update", resource.General.RedeployOnUpdate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Scope

	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))

	out_scope[0]["all_computers"] = resource.Scope.AllComputers
	out_scope[0]["all_jss_users"] = resource.Scope.AllJSSUsers

	// Computers
	if len(resource.Scope.Computers) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_ids"] = listOfIds
	}

	// TODO make this work later. It's a replacement for the log above.
	// comps, err := GetListOfIdsFromResp[jamfpro.MacOSConfigurationProfileSubsetComputer](resource.Scope.Computers, "id")
	// out_scope[0]["computer_ids"] = comps

	// Computer Groups
	if len(resource.Scope.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_group_ids"] = listOfIds
	}

	// JSS Users
	if len(resource.Scope.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_ids"] = listOfIds
	}

	// JSS User Groups
	if len(resource.Scope.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_group_ids"] = listOfIds
	}

	// Buildings
	if len(resource.Scope.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["building_ids"] = listOfIds
	}

	// Departments
	if len(resource.Scope.Departments) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["department_ids"] = listOfIds
	}

	// Scope Limitations

	out_scope_limitations := make([]map[string]interface{}, 0)
	out_scope_limitations = append(out_scope_limitations, make(map[string]interface{}))
	var limitationsSet bool

	// Users
	if len(resource.Scope.Limitations.Users) > 0 {
		var listOfNames []string
		for _, v := range resource.Scope.Limitations.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_limitations[0]["user_names"] = listOfNames
		limitationsSet = true
	}

	// Network Segments
	if len(resource.Scope.Limitations.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Limitations.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["network_segment_ids"] = listOfIds
		limitationsSet = true
	}

	// IBeacons
	if len(resource.Scope.Limitations.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Limitations.IBeacons {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["ibeacon_ids"] = listOfIds
		limitationsSet = true
	}

	// User Groups
	if len(resource.Scope.Limitations.UserGroups) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Limitations.UserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["user_group_ids"] = listOfIds
		limitationsSet = true
	}

	if limitationsSet {
		out_scope[0]["limitations"] = out_scope_limitations
	}

	// Scope Exclusions

	out_scope_exclusions := make([]map[string]interface{}, 0)
	out_scope_exclusions = append(out_scope_exclusions, make(map[string]interface{}))
	var exclusionsSet bool

	// Computers
	if len(resource.Scope.Exclusions.Computers) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_ids"] = listOfIds
		exclusionsSet = true
	}

	// Computer Groups
	if len(resource.Scope.Exclusions.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_group_ids"] = listOfIds
		exclusionsSet = true
	}

	// Buildings
	if len(resource.Scope.Exclusions.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["building_ids"] = listOfIds
		exclusionsSet = true
	}

	// Departments
	if len(resource.Scope.Exclusions.Departments) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["department_ids"] = listOfIds
		exclusionsSet = true
	}

	// Network Segments
	if len(resource.Scope.Exclusions.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["network_segment_ids"] = listOfIds
		exclusionsSet = true
	}

	// JSS Users
	if len(resource.Scope.Exclusions.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_ids"] = listOfIds
		exclusionsSet = true
	}

	// JSS User Groups
	if len(resource.Scope.Exclusions.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_group_ids"] = listOfIds
		exclusionsSet = true
	}

	// IBeacons
	if len(resource.Scope.Exclusions.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range resource.Scope.Exclusions.IBeacons {
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

	// Set Scope to state
	err := d.Set("scope", out_scope)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Self Service

	out_self_service := make([]map[string]interface{}, 0)
	out_self_service = append(out_self_service, make(map[string]interface{}, 1))
	var selfServiceSet bool

	// Fix the stupid broken double key issue
	err = FixStupidDoubleKey(resource, &out_self_service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// TODO this is problematic and will be solved another day
	// if len(resource.SelfService.SelfServiceCategories) > 0 {
	// 	var listOfIds []int
	// 	for _, v := range resource.SelfService.SelfServiceCategories {
	// 		listOfIds = append(listOfIds, v.ID)
	// 	}
	// 	out_self_service[0]["self_service_categories"] = listOfIds
	// 	selfServiceSet = true
	// }

	if selfServiceSet {
		err = d.Set("self_service", out_self_service)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("no self service") // TODO logging
	}

	return diags

}
