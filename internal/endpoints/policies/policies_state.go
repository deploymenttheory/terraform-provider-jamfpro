package policies

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, resourceID string) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	// General / root level

	log.Println("LOGHERE-STATE")
	log.Println("STATE-FLAG-1")

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	log.Println("STATE-FLAG-2")

	err = d.Set("name", resp.General.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("enabled", resp.General.Enabled)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_checkin", resp.General.TriggerCheckin)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_enrollment_complete", resp.General.TriggerEnrollmentComplete)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_login", resp.General.TriggerLogin)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_network_state_changed", resp.General.TriggerNetworkStateChanged)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_startup", resp.General.TriggerStartup)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_other", resp.General.TriggerOther)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("frequency", resp.General.Frequency)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_event", resp.General.RetryEvent)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_attempts", resp.General.RetryAttempts)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("notify_on_each_failed_retry", resp.General.NotifyOnEachFailedRetry)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("offline", resp.General.Offline)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	log.Println("STATE-FLAG-3")

	// Site
	// TODO Review this logic
	if resp.General.Site.ID != -1 && resp.General.Site.Name != "None" {
		out_site := []map[string]interface{}{
			{
				"id": resp.General.Site.ID,
				// "name": resp.General.Site.Name,
			},
		}

		if err := d.Set("site", out_site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
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
				// "name": resp.General.Category.Name,
			},
		}
		if err := d.Set("category", out_category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default category response") // TODO logging
	}

	log.Println("STATE-FLAG-5")

	// log.Printf("%+v\n", resp.Scope)
	// jsonData, err := json.MarshalIndent(resp.Scope, "", "	")
	// log.Println(string(jsonData))

	// Scope
	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))
	out_scope[0]["all_computers"] = resp.Scope.AllComputers
	out_scope[0]["all_jss_users"] = resp.Scope.AllJSSUsers

	log.Println("STATE-FLAG-6")

	// Computers
	log.Printf("%+v", resp.Scope)
	if resp.Scope.Computers != nil {
		log.Println("LINE 144")
		if len(*resp.Scope.Computers) > 0 {
			log.Println("TEST HERE")
			var listOfIds []int
			for _, v := range *resp.Scope.Computers {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["computer_ids"] = listOfIds
		}
	}

	log.Println("STATE-FLAG-7")

	// Computer Groups
	if resp.Scope.ComputerGroups != nil {
		if len(*resp.Scope.ComputerGroups) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.ComputerGroups {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["computer_group_ids"] = listOfIds
		}
	}

	log.Println("STATE-FLAG-8")

	a := 1

	if a == 2 {

		// JSS Users
		if len(*resp.Scope.JSSUsers) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.JSSUsers {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["jss_user_ids"] = listOfIds
		}

		log.Println("STATE-FLAG-9")

		// JSS User Groups
		if len(*resp.Scope.JSSUserGroups) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.JSSUserGroups {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["jss_user_group_ids"] = listOfIds
		}

		log.Println("STATE-FLAG-10")

		// Buildings
		if len(*resp.Scope.Buildings) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Buildings {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["building_ids"] = listOfIds
		}

		log.Println("STATE-FLAG-11")

		// Departments
		if len(*resp.Scope.Departments) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Departments {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope[0]["department_ids"] = listOfIds
		}

		log.Println("STATE-FLAG-12")

		// Scope Limitations
		out_scope_limitations := make([]map[string]interface{}, 0)
		out_scope_limitations = append(out_scope_limitations, make(map[string]interface{}))
		var limitationsSet bool

		log.Println("STATE-FLAG-13")

		// Users
		if len(*resp.Scope.Limitations.Users) > 0 {
			var listOfNames []string
			for _, v := range *resp.Scope.Limitations.Users {
				listOfNames = append(listOfNames, v.Name)
			}
			out_scope_limitations[0]["user_names"] = listOfNames
			limitationsSet = true
		}

		log.Println("STATE-FLAG-14")

		// Network Segments
		if len(*resp.Scope.Limitations.NetworkSegments) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Limitations.NetworkSegments {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_limitations[0]["network_segment_ids"] = listOfIds
			limitationsSet = true
		}

		log.Println("STATE-FLAG-15")

		// IBeacons
		if len(*resp.Scope.Limitations.IBeacons) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Limitations.IBeacons {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_limitations[0]["ibeacon_ids"] = listOfIds
			limitationsSet = true
		}

		log.Println("STATE-FLAG-16")

		// User Groups
		if len(*resp.Scope.Limitations.UserGroups) > 0 {
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

		log.Println("STATE-FLAG-17")

		// Scope Exclusions
		out_scope_exclusions := make([]map[string]interface{}, 0)
		out_scope_exclusions = append(out_scope_exclusions, make(map[string]interface{}))
		var exclusionsSet bool

		log.Println("STATE-FLAG-18")

		// Computers
		if len(*resp.Scope.Exclusions.Computers) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.Computers {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["computer_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-19")

		// Computer Groups
		if len(*resp.Scope.Exclusions.ComputerGroups) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.ComputerGroups {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["computer_group_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-20")

		// Buildings
		if len(*resp.Scope.Exclusions.Buildings) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.Buildings {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["building_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-21")

		// Departments
		if len(*resp.Scope.Exclusions.Departments) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.Departments {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["department_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-22")

		// Network Segments
		if len(*resp.Scope.Exclusions.NetworkSegments) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.NetworkSegments {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["network_segment_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-23")

		// JSS Users
		if len(*resp.Scope.Exclusions.JSSUsers) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.JSSUsers {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["jss_user_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-24")

		// JSS User Groups
		if len(*resp.Scope.Exclusions.JSSUserGroups) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.JSSUserGroups {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["jss_user_group_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-25")

		// IBeacons
		if len(*resp.Scope.Exclusions.IBeacons) > 0 {
			var listOfIds []int
			for _, v := range *resp.Scope.Exclusions.IBeacons {
				listOfIds = append(listOfIds, v.ID)
			}
			out_scope_exclusions[0]["ibeacon_ids"] = listOfIds
			exclusionsSet = true
		}

		log.Println("STATE-FLAG-26")

		// Append Exclusions if they're set
		if exclusionsSet {
			out_scope[0]["exclusions"] = out_scope_exclusions
		} else {
			log.Println("No exclusions set") // TODO logging
		}

		log.Println("STATE-FLAG-27")

	}

	// Set Scope to state
	err = d.Set("scope", out_scope)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	log.Println("STATE-FLAG-28")

	return diags
}
