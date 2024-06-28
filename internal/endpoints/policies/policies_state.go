package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parent func for invdividual stating functions
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, resourceID string) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("LOGHERE-RESPONSE")
	// xmlData, _ := xml.MarshalIndent(resp, " ", "	")
	// log.Println(string(xmlData))

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

	return diags
}

// Reads response and states general/root level items
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
	d.Set("site_id", resp.General.Site.ID)

	// Category
	d.Set("category_id", resp.General.Category.ID)
}

// Reads response and states scope items
func stateScope(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error

	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))
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
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// Reads response and states self service items
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error
	out_ss := make([]map[string]interface{}, 0)
	out_ss = append(out_ss, make(map[string]interface{}, 1))

	if resp.SelfService != nil {
		out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
		out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
		out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
		out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
		out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
		out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage

		err = d.Set("self_service", out_ss)
		if err != nil {
			*diags = append(*diags, diag.FromErr(err)...)
		}
	}
}

// Parent func for stating payloads. Constructs var with prep funcs and states as one here.
func statePayloads(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	out := make([]map[string]interface{}, 0)
	out = append(out, make(map[string]interface{}, 1))

	// DiskEncryption
	prepStatePayloadDiskEncryption(&out, resp)

	// Packages
	prepStatePayloadPackages(&out, resp)

	// Scripts
	prepStatePayloadScripts(&out, resp)

	// Account Maintenance
	prepStatePayloadAccountMaintenance(&out, resp)

	// Files Processes
	prepStatePayloadFilesProcesses(&out, resp)

	// User Interaction
	prepStatePayloadUserInteraction(&out, resp)

	// Reboot
	prepStatePayloadReboot(&out, resp)

	// Maintenance
	prepStatePayloadMaintenance(&out, resp)

	// State
	err := d.Set("payloads", out)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// Reads response and preps disk encryption payload items
func prepStatePayloadDiskEncryption(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.DiskEncryption == nil {
		return
	}
	(*out)[0]["disk_encryption"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["action"] = resp.DiskEncryption.Action
	outMap["disk_encryption_configuration_id"] = resp.DiskEncryption.DiskEncryptionConfigurationID
	outMap["auth_restart"] = resp.DiskEncryption.AuthRestart
	(*out)[0]["disk_encryption"] = append((*out)[0]["disk_encryption"].([]map[string]interface{}), outMap)
}

// Reads response and preps package payload items
func prepStatePayloadPackages(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.PackageConfiguration == nil {
		return
	}
	// Packages can be nil but deployment state default
	if resp.PackageConfiguration.Packages == nil {
		return
	}

	(*out)[0]["packages"] = make([]map[string]interface{}, 0)
	for _, v := range *resp.PackageConfiguration.Packages {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["action"] = v.Action
		outMap["fill_user_template"] = v.FillUserTemplate
		outMap["fill_existing_user_template"] = v.FillExistingUsers
		(*out)[0]["packages"] = append((*out)[0]["packages"].([]map[string]interface{}), outMap)
	}

	(*out)[0]["distribution_point"] = resp.PackageConfiguration.DistributionPoint
}

// Reads response and preps script payload items
func prepStatePayloadScripts(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Scripts.Script == nil {
		return
	}

	(*out)[0]["scripts"] = make([]map[string]interface{}, 0)
	for _, v := range *resp.Scripts.Script {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["priority"] = v.Priority

		if v.Parameter4 != "" {
			outMap["parameter4"] = v.Parameter4
		}
		if v.Parameter5 != "" {
			outMap["parameter5"] = v.Parameter5
		}
		if v.Parameter6 != "" {
			outMap["parameter6"] = v.Parameter6
		}
		if v.Parameter7 != "" {
			outMap["parameter7"] = v.Parameter7
		}
		if v.Parameter8 != "" {
			outMap["parameter8"] = v.Parameter8
		}
		if v.Parameter9 != "" {
			outMap["parameter9"] = v.Parameter9
		}
		if v.Parameter10 != "" {
			outMap["parameter10"] = v.Parameter10
		}
		if v.Parameter11 != "" {
			outMap["parameter11"] = v.Parameter11
		}
		(*out)[0]["scripts"] = append((*out)[0]["scripts"].([]map[string]interface{}), outMap)
		log.Println("LOGHERE-SCRIPT OUT")
		log.Println(outMap)
	}
}

// Reads response and preps account maintenance payload items
func prepStatePayloadAccountMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.AccountMaintenance == nil || resp.AccountMaintenance.Accounts == nil {
		return
	}

	(*out)[0]["account_maintenance"] = make([]map[string]interface{}, 0)
	for _, v := range *resp.AccountMaintenance.Accounts {
		outMap := make(map[string]interface{})
		outMap["action"] = v.Action
		outMap["username"] = v.Username
		outMap["realname"] = v.Realname
		outMap["password"] = v.Password
		outMap["archive_home_directory"] = v.ArchiveHomeDirectory
		outMap["archive_home_directory_to"] = v.ArchiveHomeDirectoryTo
		outMap["home"] = v.Home
		outMap["hint"] = v.Hint
		outMap["picture"] = v.Picture
		outMap["admin"] = v.Admin
		outMap["filevault_enabled"] = v.FilevaultEnabled
		(*out)[0]["account_maintenance"] = append((*out)[0]["account_maintenance"].([]map[string]interface{}), outMap)
	}
}

// Reads response and preps files and processes payload items
func prepStatePayloadFilesProcesses(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.FilesProcesses == nil {
		return
	}

	(*out)[0]["files_processes"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["search_by_path"] = resp.FilesProcesses.SearchByPath
	outMap["delete_file"] = resp.FilesProcesses.DeleteFile
	outMap["locate_file"] = resp.FilesProcesses.LocateFile
	outMap["update_locate_database"] = resp.FilesProcesses.UpdateLocateDatabase
	outMap["spotlight_search"] = resp.FilesProcesses.SpotlightSearch
	outMap["search_for_process"] = resp.FilesProcesses.SearchForProcess
	outMap["kill_process"] = resp.FilesProcesses.KillProcess
	outMap["run_command"] = resp.FilesProcesses.RunCommand
	(*out)[0]["files_processes"] = append((*out)[0]["files_processes"].([]map[string]interface{}), outMap)
}

// Reads response and preps user interaction payload items
func prepStatePayloadUserInteraction(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.UserInteraction == nil {
		return
	}

	(*out)[0]["user_interaction"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["message_start"] = resp.UserInteraction.MessageStart
	outMap["allow_user_to_defer"] = resp.UserInteraction.AllowUserToDefer
	outMap["allow_deferral_until_utc"] = resp.UserInteraction.AllowDeferralUntilUtc
	outMap["allow_deferral_minutes"] = resp.UserInteraction.AllowDeferralMinutes
	outMap["message_finish"] = resp.UserInteraction.MessageFinish
	(*out)[0]["user_interaction"] = append((*out)[0]["user_interaction"].([]map[string]interface{}), outMap)
}

// Reads response and preps reboot payload items
func prepStatePayloadReboot(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Reboot == nil {
		return
	}

	(*out)[0]["reboot"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["message"] = resp.Reboot.Message
	outMap["specify_startup"] = resp.Reboot.SpecifyStartup
	outMap["startup_disk"] = resp.Reboot.StartupDisk
	outMap["no_user_logged_in"] = resp.Reboot.NoUserLoggedIn
	outMap["user_logged_in"] = resp.Reboot.UserLoggedIn
	outMap["minutes_until_reboot"] = resp.Reboot.MinutesUntilReboot
	outMap["start_reboot_timer_immediately"] = resp.Reboot.StartRebootTimerImmediately
	outMap["file_vault_2_reboot"] = resp.Reboot.FileVault2Reboot
	(*out)[0]["reboot"] = append((*out)[0]["reboot"].([]map[string]interface{}), outMap)
}

// Reads response and preps maintenance payload items
func prepStatePayloadMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Maintenance == nil {
		return
	}

	(*out)[0]["maintenance"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["recon"] = resp.Maintenance.Recon
	outMap["reset_name"] = resp.Maintenance.ResetName
	outMap["install_all_cached_packages"] = resp.Maintenance.InstallAllCachedPackages
	outMap["heal"] = resp.Maintenance.Heal
	outMap["prebindings"] = resp.Maintenance.Prebindings
	outMap["permissions"] = resp.Maintenance.Permissions
	outMap["byhost"] = resp.Maintenance.Byhost
	outMap["system_cache"] = resp.Maintenance.SystemCache
	outMap["user_cache"] = resp.Maintenance.UserCache
	outMap["verify"] = resp.Maintenance.Verify
	(*out)[0]["maintenance"] = append((*out)[0]["maintenance"].([]map[string]interface{}), outMap)
}
