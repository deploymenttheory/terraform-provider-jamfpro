package policy

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"
	"reflect"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parent func for stating payloads. Constructs var with prep funcs and states as one here.
func statePayloads(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	out := make([]map[string]any, 0)
	out = append(out, make(map[string]any, 1))

	// DiskEncryption
	prepStatePayloadDiskEncryption(&out, resp)

	// Packages
	prepStatePayloadPackages(&out, resp)

	// Scripts
	prepStatePayloadScripts(&out, resp)

	// Printers
	prepStatePayloadPrinters(&out, resp)

	// Dock Items
	prepStatePayloadDockItems(&out, resp)

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

// prepStatePayloadDiskEncryption reads response and preps disk encryption payload items for stating
func prepStatePayloadDiskEncryption(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	defaults := map[string]any{
		"action":                           "none",
		"disk_encryption_configuration_id": 0,
		"auth_restart":                     false,
		"remediate_key_type":               "",
		"remediate_disk_encryption_configuration_id": 0,
	}

	diskEncryptionStatePayload := map[string]any{
		"action":                           resp.DiskEncryption.Action,
		"disk_encryption_configuration_id": resp.DiskEncryption.DiskEncryptionConfigurationID,
		"auth_restart":                     resp.DiskEncryption.AuthRestart,
		"remediate_key_type":               resp.DiskEncryption.RemediateKeyType,
		"remediate_disk_encryption_configuration_id": resp.DiskEncryption.RemediateDiskEncryptionConfigurationID,
	}

	allDefault := true
	for key, value := range diskEncryptionStatePayload {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	(*out)[0]["disk_encryption"] = []map[string]any{diskEncryptionStatePayload}
}

// Reads response and preps package payload items
func prepStatePayloadPackages(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	if len(resp.PackageConfiguration.Packages) == 0 {
		return
	}

	packagesMap := make(map[string]any)
	packagesMap["distribution_point"] = resp.PackageConfiguration.DistributionPoint
	packagesMap["package"] = make([]map[string]any, 0)

	for _, v := range resp.PackageConfiguration.Packages {
		outMap := make(map[string]any)
		outMap["id"] = v.ID
		outMap["action"] = v.Action
		outMap["fill_user_template"] = v.FillUserTemplate
		outMap["fill_existing_user_template"] = v.FillExistingUsers
		packagesMap["package"] = append(packagesMap["package"].([]map[string]any), outMap)
	}

	(*out)[0]["packages"] = []map[string]any{packagesMap}
}

// Reads response and preps script payload items
func prepStatePayloadScripts(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	if resp.Scripts == nil {
		log.Println("No scripts found")
		return
	}

	log.Println("Initializing scripts in state")
	(*out)[0]["scripts"] = make([]map[string]any, 0)

	for _, v := range resp.Scripts {
		outMap := make(map[string]any)
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

		(*out)[0]["scripts"] = append((*out)[0]["scripts"].([]map[string]any), outMap)
	}

}

// prepStatePayloadPrinters reads response and preps printer payload items for stating
func prepStatePayloadPrinters(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	if resp.Printers.Printer == nil {
		return
	}

	log.Println("Initializing printers in state")
	(*out)[0]["printers"] = make([]map[string]any, 0)

	for _, v := range resp.Printers.Printer {
		outMap := make(map[string]any)
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action
		outMap["make_default"] = v.MakeDefault

		(*out)[0]["printers"] = append((*out)[0]["printers"].([]map[string]any), outMap)
	}

}

// Reads response and preps dock items payload items
func prepStatePayloadDockItems(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	if resp.DockItems == nil {
		return
	}

	(*out)[0]["dock_items"] = make([]map[string]any, 0)

	for _, v := range resp.DockItems {
		outMap := make(map[string]any)
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action

		(*out)[0]["dock_items"] = append((*out)[0]["dock_items"].([]map[string]any), outMap)
	}

}

// prepStatePayloadAccountMaintenance reads response and preps account maintenance payload items.
// If all values are default, do not set the account_maintenance block
func prepStatePayloadAccountMaintenance(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	accountMaintenanceMap := make(map[string]any)

	if resp.AccountMaintenance.Accounts != nil {
		localAccounts := make([]map[string]any, 0)
		for _, v := range *resp.AccountMaintenance.Accounts {
			accountMap := make(map[string]any)
			accountMap["action"] = v.Action
			accountMap["username"] = v.Username
			accountMap["realname"] = v.Realname
			accountMap["password"] = v.Password
			accountMap["archive_home_directory"] = v.ArchiveHomeDirectory
			accountMap["archive_home_directory_to"] = v.ArchiveHomeDirectoryTo
			accountMap["home"] = v.Home
			accountMap["hint"] = v.Hint
			accountMap["picture"] = v.Picture
			accountMap["admin"] = v.Admin
			accountMap["filevault_enabled"] = v.FilevaultEnabled

			localAccounts = append(localAccounts, accountMap)
		}

		if len(localAccounts) > 0 {
			accountMaintenanceMap["local_accounts"] = []map[string]any{
				{"account": localAccounts},
			}
		}
	}

	// Handle directory bindings
	if resp.AccountMaintenance.DirectoryBindings != nil {
		directoryBindings := make([]map[string]any, 0)
		for _, v := range *resp.AccountMaintenance.DirectoryBindings {
			bindingMap := make(map[string]any)
			bindingMap["name"] = v.Name

			directoryBindings = append(directoryBindings, bindingMap)
		}

		if len(directoryBindings) > 0 {
			accountMaintenanceMap["directory_bindings"] = []map[string]any{
				{"binding": directoryBindings},
			}
		}
	}

	// Handle management account
	if resp.AccountMaintenance.ManagementAccount != nil {
		managementAccountMap := make(map[string]any)
		if resp.AccountMaintenance.ManagementAccount.Action != "doNotChange" || resp.AccountMaintenance.ManagementAccount.ManagedPassword != "" || resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength != 0 {
			managementAccountMap["action"] = resp.AccountMaintenance.ManagementAccount.Action
			managementAccountMap["managed_password"] = resp.AccountMaintenance.ManagementAccount.ManagedPassword
			managementAccountMap["managed_password_length"] = resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength

			accountMaintenanceMap["management_account"] = []map[string]any{managementAccountMap}
		}
	}

	// Handle open firmware/EFI password
	if resp.AccountMaintenance.OpenFirmwareEfiPassword != nil {
		openFirmwareMap := make(map[string]any)
		if resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode != "none" || resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword != "" {
			openFirmwareMap["of_mode"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode
			openFirmwareMap["of_password"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword

			accountMaintenanceMap["open_firmware_efi_password"] = []map[string]any{openFirmwareMap}
		}
	}

	if len(accountMaintenanceMap) > 0 {
		(*out)[0]["account_maintenance"] = []map[string]any{accountMaintenanceMap}
	}
}

// prepStatePayloadFilesProcesses reads response and preps files and processes payload items.
func prepStatePayloadFilesProcesses(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	defaults := map[string]any{
		"search_by_path":         "",
		"delete_file":            false,
		"locate_file":            "",
		"update_locate_database": false,
		"spotlight_search":       "",
		"search_for_process":     "",
		"kill_process":           false,
		"run_command":            "",
	}

	filesProcessesBlock := map[string]any{
		"search_by_path":         resp.FilesProcesses.SearchByPath,
		"delete_file":            resp.FilesProcesses.DeleteFile,
		"locate_file":            resp.FilesProcesses.LocateFile,
		"update_locate_database": resp.FilesProcesses.UpdateLocateDatabase,
		"spotlight_search":       resp.FilesProcesses.SpotlightSearch,
		"search_for_process":     resp.FilesProcesses.SearchForProcess,
		"kill_process":           resp.FilesProcesses.KillProcess,
		"run_command":            resp.FilesProcesses.RunCommand,
	}

	allDefault := true
	for key, value := range filesProcessesBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	(*out)[0]["files_processes"] = []map[string]any{filesProcessesBlock}
}

// prepStatePayloadUserInteraction Reads response and preps user interaction payload items. If all values are default, do not set the user_interaction block
func prepStatePayloadUserInteraction(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	defaults := map[string]any{
		"message_start":            "",
		"allow_users_to_defer":     false,
		"allow_deferral_until_utc": "",
		"allow_deferral_minutes":   0,
		"message_finish":           "",
	}

	userInteractionBlock := map[string]any{
		"message_start":            resp.UserInteraction.MessageStart,
		"allow_users_to_defer":     resp.UserInteraction.AllowUsersToDefer,
		"allow_deferral_until_utc": resp.UserInteraction.AllowDeferralUntilUtc,
		"allow_deferral_minutes":   resp.UserInteraction.AllowDeferralMinutes,
		"message_finish":           resp.UserInteraction.MessageFinish,
	}

	allDefault := true
	for key, value := range userInteractionBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	(*out)[0]["user_interaction"] = []map[string]any{userInteractionBlock}
}

// Reads response and preps reboot payload items
func prepStatePayloadReboot(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	defaults := map[string]any{
		"message":                        "",
		"specify_startup":                "",
		"startup_disk":                   "Current Startup Disk",
		"no_user_logged_in":              "Do not restart",
		"user_logged_in":                 "Do not restart",
		"minutes_until_reboot":           0,
		"start_reboot_timer_immediately": false,
		"file_vault_2_reboot":            false,
	}

	rebootBlock := map[string]any{
		"message":                        resp.Reboot.Message,
		"specify_startup":                resp.Reboot.SpecifyStartup,
		"startup_disk":                   resp.Reboot.StartupDisk,
		"no_user_logged_in":              resp.Reboot.NoUserLoggedIn,
		"user_logged_in":                 resp.Reboot.UserLoggedIn,
		"minutes_until_reboot":           resp.Reboot.MinutesUntilReboot,
		"start_reboot_timer_immediately": resp.Reboot.StartRebootTimerImmediately,
		"file_vault_2_reboot":            resp.Reboot.FileVault2Reboot,
	}

	allDefault := true
	for key, value := range rebootBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	(*out)[0]["reboot"] = []map[string]any{rebootBlock}
}

// prepStatePayloadMaintenance Reads response and preps maintenance payload items. If all values are default, do not set the maintenance block
func prepStatePayloadMaintenance(out *[]map[string]any, resp *jamfpro.ResourcePolicy) {
	v := reflect.ValueOf(resp.Maintenance)

	allDefault := true
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Bool() {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	(*out)[0]["maintenance"] = make([]map[string]any, 0)

	outMap := make(map[string]any)
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
	(*out)[0]["maintenance"] = append((*out)[0]["maintenance"].([]map[string]any), outMap)
}
