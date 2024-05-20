package policies

import (
	"encoding/xml"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Primary

func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, resourceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// General/Root level
	stateGeneral(d, resp, &diags)

	// Scope
	stateScope(d, resp, &diags)

	// Self Service
	stateSelfService(d, resp, &diags)

	// Payloads
	statePayloads(d, resp, &diags)

	log.Println("STATE-FLAG-28")
	log.Printf("%+v", diags)

	return diags
}

// Child funcs

func stateGeneral(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error

	err = d.Set("name", resp.General.Name)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("enabled", resp.General.Enabled)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_checkin", resp.General.TriggerCheckin)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_enrollment_complete", resp.General.TriggerEnrollmentComplete)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_login", resp.General.TriggerLogin)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_network_state_changed", resp.General.TriggerNetworkStateChanged)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_startup", resp.General.TriggerStartup)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_other", resp.General.TriggerOther)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("frequency", resp.General.Frequency)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_event", resp.General.RetryEvent)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_attempts", resp.General.RetryAttempts)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("notify_on_each_failed_retry", resp.General.NotifyOnEachFailedRetry)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("offline", resp.General.Offline)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	// Site
	// TODO Review this logic
	if resp.General.Site.ID != -1 && resp.General.Site.Name != "None" {
		out_site := []map[string]interface{}{
			{
				"id": resp.General.Site.ID,
			},
		}

		if err := d.Set("site", out_site); err != nil {
			if diags == nil {
				diags = &diag.Diagnostics{}
			}
			*diags = append(*diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default site response") // TODO Logging
	}

	log.Println("STATE-FLAG-4")

	// Category
	if resp.General.Category.ID != -1 && resp.General.Category.Name != "No category assigned" {
		out_category := []map[string]interface{}{
			{
				"id": resp.General.Category.ID,
			},
		}
		if err := d.Set("category", out_category); err != nil {
			if diags == nil {
				diags = &diag.Diagnostics{}
			}
			*diags = append(*diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default category response") // TODO logging
	}

}

func stateScope(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error

	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))
	out_scope[0]["all_computers"] = resp.Scope.AllComputers
	out_scope[0]["all_jss_users"] = resp.Scope.AllJSSUsers

	// DEBUG // TODO remove this
	log.Println("XMLIN")
	log.Println(resp.Scope.AllJSSUsers)
	xmlData, _ := xml.MarshalIndent(resp, "", "    ")
	log.Println(string(xmlData))

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
	out_scope_limitations := make([]map[string]interface{}, 0)
	out_scope_limitations = append(out_scope_limitations, make(map[string]interface{}))
	var limitationsSet bool

	// Users
	if resp.Scope.Limitations.Users != nil && len(*resp.Scope.Limitations.Users) > 0 {
		var listOfNames []string
		for _, v := range *resp.Scope.Limitations.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_limitations[0]["user_names"] = listOfNames
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
		var listOfIds []int
		for _, v := range *resp.Scope.Limitations.UserGroups {
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
		if diags == nil {
			diags = &diag.Diagnostics{}
		}
		*diags = append(*diags, diag.FromErr(err)...)
	}

}

func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error
	out_ss := make([]map[string]interface{}, 0)
	out_ss = append(out_ss, make(map[string]interface{}, 1))

	if resp.SelfService != nil {
		log.Println("STATE-FLAG_RESP_SELFERVICE")
		log.Printf("%+v", resp.SelfService)
		out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
		out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
		out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
		out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
		out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
		out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage

		err = d.Set("self_service", out_ss)
		if err != nil {
			if diags == nil {
				diags = &diag.Diagnostics{}
			}
			*diags = append(*diags, diag.FromErr(err)...)
		}
	}
}

func statePayloads(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error
	out := make([]map[string]interface{}, 0)
	out = append(out, make(map[string]interface{}, 1))
	out[0]["packages"] = make([]map[string]interface{}, 0)
	log.Println("LOGHERE")
	log.Printf("%+v", out)

	if resp.PackageConfiguration != nil {
		log.Println("NOT NIL")
		for _, v := range *resp.PackageConfiguration.Packages {
			log.Println("LOOPING")
			outMap := make(map[string]interface{})
			outMap["id"] = v.ID
			outMap["action"] = v.Action
			outMap["fill_user_template"] = v.FillUserTemplate
			outMap["fill_existing_user_template"] = v.FillExistingUsers
			out[0]["packages"] = append(out[0]["packages"].([]map[string]interface{}), outMap)
		}

	}

	log.Println("OUT DONE:")
	log.Printf("%+v", out)

	err = d.Set("payloads", out)
	if err != nil {
		log.Println("ERROR FOUND", err)
		if diags == nil {
			diags = &diag.Diagnostics{}
		}
		*diags = append(*diags, diag.FromErr(err)...)
	}

}
