package policy

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructPolicy builds the policy object from the HCL. It's composed of several sub-objects, each with their own schema.
func construct(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	var err error
	resource := &jamfpro.ResourcePolicy{}

	constructGeneral(d, resource)

	err = constructScope(d, resource)
	if err != nil {
		return nil, err
	}

	constructSelfService(d, resource)

	constructPayloads(d, resource)

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Policy '%s' to XML: %v", resource.General.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Policy XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// constructGeneral builds the general settings of the jamf pro policy from the HCL.
func constructGeneral(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	resource.General = jamfpro.PolicySubsetGeneral{
		Name:                       d.Get("name").(string),
		Enabled:                    d.Get("enabled").(bool),
		TriggerCheckin:             d.Get("trigger_checkin").(bool),
		TriggerEnrollmentComplete:  d.Get("trigger_enrollment_complete").(bool),
		TriggerLogin:               d.Get("trigger_login").(bool),
		TriggerNetworkStateChanged: d.Get("trigger_network_state_changed").(bool),
		TriggerStartup:             d.Get("trigger_startup").(bool),
		TriggerOther:               d.Get("trigger_other").(string),
		Frequency:                  d.Get("frequency").(string),
		RetryEvent:                 d.Get("retry_event").(string),
		RetryAttempts:              d.Get("retry_attempts").(int),
		NotifyOnEachFailedRetry:    d.Get("notify_on_each_failed_retry").(bool),
		TargetDrive:                d.Get("target_drive").(string),
		Offline:                    d.Get("offline").(bool),
		NetworkRequirements:        d.Get("network_requirements").(string),
	}

	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	setDateTimeLimitations(d, resource)

	setNetworkLimitations(d, resource)

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

}

// Helper function to set DateTime Limitations
func setDateTimeLimitations(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	if dateTimeLimitations, ok := d.GetOk("date_time_limitations"); ok {
		dateTimeLimitationsList := dateTimeLimitations.([]any)
		if len(dateTimeLimitationsList) > 0 {
			dateTimeLimitationsMap := dateTimeLimitationsList[0].(map[string]any)

			var noExecuteOn []string
			if v, ok := dateTimeLimitationsMap["no_execute_on"].(*schema.Set); ok {
				for _, day := range v.List() {
					noExecuteOn = append(noExecuteOn, day.(string))
				}
			}

			resource.General.DateTimeLimitations = &jamfpro.PolicySubsetGeneralDateTimeLimitations{
				ActivationDate:      dateTimeLimitationsMap["activation_date"].(string),
				ActivationDateEpoch: dateTimeLimitationsMap["activation_date_epoch"].(int),
				ActivationDateUTC:   dateTimeLimitationsMap["activation_date_utc"].(string),
				ExpirationDate:      dateTimeLimitationsMap["expiration_date"].(string),
				ExpirationDateEpoch: dateTimeLimitationsMap["expiration_date_epoch"].(int),
				ExpirationDateUTC:   dateTimeLimitationsMap["expiration_date_utc"].(string),
				NoExecuteOn:         noExecuteOn,
				NoExecuteStart:      dateTimeLimitationsMap["no_execute_start"].(string),
				NoExecuteEnd:        dateTimeLimitationsMap["no_execute_end"].(string),
			}
		}
	}
}

// Helper function to set Network Limitations
func setNetworkLimitations(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	if networkLimitations, ok := d.GetOk("network_limitations"); ok {
		networkLimitationsList := networkLimitations.([]any)
		if len(networkLimitationsList) > 0 {
			networkLimitationsMap := networkLimitationsList[0].(map[string]any)
			resource.General.NetworkLimitations = &jamfpro.PolicySubsetGeneralNetworkLimitations{
				MinimumNetworkConnection: networkLimitationsMap["minimum_network_connection"].(string),
				AnyIPAddress:             networkLimitationsMap["any_ip_address"].(bool),
			}
		}
	}
}

// Pulls "scope" settings from HCL and packages into object
func constructScope(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	var err error

	if len(d.Get("scope").([]any)) == 0 {
		return nil
	}

	// Targets
	resource.Scope = jamfpro.PolicySubsetScope{
		Computers:      &[]jamfpro.PolicySubsetComputer{},
		ComputerGroups: &[]jamfpro.PolicySubsetComputerGroup{},
		JSSUsers:       &[]jamfpro.PolicySubsetJSSUser{},
		JSSUserGroups:  &[]jamfpro.PolicySubsetJSSUserGroup{},
		Buildings:      &[]jamfpro.PolicySubsetBuilding{},
		Departments:    &[]jamfpro.PolicySubsetDepartment{},
	}

	// Bools
	resource.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	resource.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	// Computers
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetComputer, int]("scope.0.computer_ids", "ID", d, resource.Scope.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, resource.Scope.ComputerGroups)
	if err != nil {
		return err
	}

	// JSS Users
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, resource.Scope.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, resource.Scope.JSSUserGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetBuilding, int]("scope.0.building_ids", "ID", d, resource.Scope.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetDepartment, int]("scope.0.department_ids", "ID", d, resource.Scope.Departments)
	if err != nil {
		return err
	}

	// Limitations
	resource.Scope.Limitations = &jamfpro.PolicySubsetScopeLimitations{
		Users:           &[]jamfpro.PolicySubsetUser{},
		UserGroups:      &[]jamfpro.PolicySubsetUserGroup{},
		NetworkSegments: &[]jamfpro.PolicySubsetNetworkSegment{},
		IBeacons:        &[]jamfpro.PolicySubsetIBeacon{},
	}

	// Network Segments
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, resource.Scope.Limitations.NetworkSegments)
	if err != nil {
		return err
	}

	// IBeacons
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, resource.Scope.Limitations.IBeacons)
	if err != nil {
		return err
	}

	// User Groups
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetUserGroup, int]("scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, resource.Scope.Limitations.UserGroups)
	if err != nil {
		return err
	}

	// TODO User Limitations

	// Exclusions

	// TODO I don't really want this here but it won't work without it. I think it's defeating the purpose of the struct layout slightly.
	resource.Scope.Exclusions = &jamfpro.PolicySubsetScopeExclusions{
		Computers:       &[]jamfpro.PolicySubsetComputer{},
		ComputerGroups:  &[]jamfpro.PolicySubsetComputerGroup{},
		Users:           &[]jamfpro.PolicySubsetUser{},
		UserGroups:      &[]jamfpro.PolicySubsetUserGroup{},
		Buildings:       &[]jamfpro.PolicySubsetBuilding{},
		Departments:     &[]jamfpro.PolicySubsetDepartment{},
		NetworkSegments: &[]jamfpro.PolicySubsetNetworkSegment{},
		JSSUsers:        &[]jamfpro.PolicySubsetJSSUser{},
		JSSUserGroups:   &[]jamfpro.PolicySubsetJSSUserGroup{},
		IBeacons:        &[]jamfpro.PolicySubsetIBeacon{},
	}

	// Computers
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, resource.Scope.Exclusions.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, resource.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, resource.Scope.Exclusions.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, resource.Scope.Exclusions.Departments)
	if err != nil {
		return err
	}

	// Network Segments
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, resource.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return err
	}

	// JSS Users
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, resource.Scope.Exclusions.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, resource.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return err
	}

	// IBeacons
	err = constructors.MapSetToStructs[jamfpro.PolicySubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, resource.Scope.Exclusions.IBeacons)
	if err != nil {
		return err
	}

	return nil
}

// Pulls "self service" settings from HCL and packages into object
func constructSelfService(d *schema.ResourceData, out *jamfpro.ResourcePolicy) {
	if len(d.Get("self_service").([]any)) > 0 {
		out.SelfService = jamfpro.PolicySubsetSelfService{
			UseForSelfService:           d.Get("self_service.0.use_for_self_service").(bool),
			SelfServiceDisplayName:      d.Get("self_service.0.self_service_display_name").(string),
			InstallButtonText:           d.Get("self_service.0.install_button_text").(string),
			ReinstallButtonText:         d.Get("self_service.0.reinstall_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			SelfServiceIcon: &jamfpro.SharedResourceSelfServiceIcon{
				ID: d.Get("self_service.0.self_service_icon_id").(int),
			},
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
		}

		categories := d.Get("self_service.0.self_service_category")
		if categories != nil {
			for _, v := range categories.([]any) {
				out.SelfService.SelfServiceCategories = append(out.SelfService.SelfServiceCategories, jamfpro.PolicySubsetSelfServiceCategory{
					ID:        v.(map[string]any)["id"].(int),
					FeatureIn: v.(map[string]any)["feature_in"].(bool),
					DisplayIn: v.(map[string]any)["display_in"].(bool),
				})
			}
		}
	}
}

// constructPayloads builds the policy payload(s) from the HCL
func constructPayloads(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	constructPayloadPackages(d, resource)
	constructPayloadScripts(d, resource)
	constructPayloadDiskEncryption(d, resource)
	constructPayloadPrinters(d, resource)
	constructPayloadDockItems(d, resource)
	constructPayloadAccountMaintenance(d, resource)
	constructPayloadFilesProcesses(d, resource)
	constructPayloadUserInteraction(d, resource)
	constructPayloadReboot(d, resource)
	constructPayloadMaintenance(d, resource)
}

// constructPayloadPackages builds the packages payload settings of the policy.
func constructPayloadPackages(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.packages.0")
	if len(hcl.(map[string]any)) == 0 {
		return
	}
	var payload jamfpro.PolicySubsetPackageConfiguration
	payload.DistributionPoint = hcl.(map[string]any)["distribution_point"].(string)
	packageList := hcl.(map[string]any)["package"].([]any)

	for _, v := range packageList {
		payload.Packages = append(payload.Packages, jamfpro.PolicySubsetPackageConfigurationPackage{
			ID:                v.(map[string]any)["id"].(int),
			Action:            v.(map[string]any)["action"].(string),
			FillUserTemplate:  v.(map[string]any)["fill_user_template"].(bool),
			FillExistingUsers: v.(map[string]any)["fill_existing_user_template"].(bool),
		})
	}

	resource.PackageConfiguration = payload
}

// Pulls "script" settings from HCL and packages them into the resource.
func constructPayloadScripts(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.scripts")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	var payloads []jamfpro.PolicySubsetScript
	for _, v := range hcl.([]any) {
		payloads = append(payloads, jamfpro.PolicySubsetScript{
			ID:          v.(map[string]any)["id"].(string),
			Priority:    v.(map[string]any)["priority"].(string),
			Parameter4:  v.(map[string]any)["parameter4"].(string),
			Parameter5:  v.(map[string]any)["parameter5"].(string),
			Parameter6:  v.(map[string]any)["parameter6"].(string),
			Parameter7:  v.(map[string]any)["parameter7"].(string),
			Parameter8:  v.(map[string]any)["parameter8"].(string),
			Parameter9:  v.(map[string]any)["parameter9"].(string),
			Parameter10: v.(map[string]any)["parameter10"].(string),
			Parameter11: v.(map[string]any)["parameter11"].(string),
		})
	}

	resource.Scripts = payloads
}

// Pulls "disk encryption" settings from HCL and packages them into the resource.
func constructPayloadDiskEncryption(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.disk_encryption")
	if hcl == nil || len(hcl.([]any)) == 0 {
		outBlock := new(jamfpro.PolicySubsetDiskEncryption)
		outBlock.RemediateKeyType = "Individual"
		resource.DiskEncryption = *outBlock
		return
	}

	outBlock := new(jamfpro.PolicySubsetDiskEncryption)
	data := hcl.([]any)[0].(map[string]any)

	outBlock.Action = data["action"].(string)
	outBlock.DiskEncryptionConfigurationID = data["disk_encryption_configuration_id"].(int)
	outBlock.AuthRestart = data["auth_restart"].(bool)
	outBlock.RemediateKeyType = data["remediate_key_type"].(string)
	outBlock.RemediateDiskEncryptionConfigurationID = data["remediate_disk_encryption_configuration_id"].(int)

	resource.DiskEncryption = *outBlock

}

// Pulls "printers" settings from HCL and packages them into the resource.
func constructPayloadPrinters(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.printers")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	outBlock := new(jamfpro.PolicySubsetPrinters)
	outBlock.Printer = []jamfpro.PolicySubsetPrinter{}
	payload := outBlock.Printer
	for _, v := range hcl.([]any) {
		payload = append(payload, jamfpro.PolicySubsetPrinter{
			ID:          v.(map[string]any)["id"].(int),
			Name:        v.(map[string]any)["name"].(string),
			Action:      v.(map[string]any)["action"].(string),
			MakeDefault: v.(map[string]any)["make_default"].(bool),
		})
	}

	outBlock.Printer = payload
	resource.Printers = *outBlock

}

// constructPayloadDockItems builds the dock items payload settings of the policy.
func constructPayloadDockItems(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.dock_items")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	var payload []jamfpro.PolicySubsetDockItem

	for _, v := range hcl.([]any) {
		newObj := jamfpro.PolicySubsetDockItem{
			ID:     v.(map[string]any)["id"].(int),
			Name:   v.(map[string]any)["name"].(string),
			Action: v.(map[string]any)["action"].(string),
		}
		payload = append(payload, newObj)
	}

	resource.DockItems = payload

}

// constructPayloadAccountMaintenance builds the account maintenance payload settings of the policy.
func constructPayloadAccountMaintenance(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.account_maintenance")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	outBlock := new(jamfpro.PolicySubsetAccountMaintenance)

	for _, v := range hcl.([]any) {
		data := v.(map[string]any)

		// Handle local accounts
		if localAccounts, ok := data["local_accounts"]; ok && len(localAccounts.([]any)) > 0 {
			localAccountsList := localAccounts.([]any)
			if len(localAccountsList) > 0 {
				accountsData := localAccountsList[0].(map[string]any)["account"].([]any)
				accounts := []jamfpro.PolicySubsetAccountMaintenanceAccount{}
				for _, account := range accountsData {
					accountData := account.(map[string]any)
					accounts = append(accounts, jamfpro.PolicySubsetAccountMaintenanceAccount{
						Action:                 accountData["action"].(string),
						Username:               accountData["username"].(string),
						Realname:               accountData["realname"].(string),
						Password:               accountData["password"].(string),
						ArchiveHomeDirectory:   accountData["archive_home_directory"].(bool),
						ArchiveHomeDirectoryTo: accountData["archive_home_directory_to"].(string),
						Home:                   accountData["home"].(string),
						Hint:                   accountData["hint"].(string),
						Picture:                accountData["picture"].(string),
						Admin:                  accountData["admin"].(bool),
						FilevaultEnabled:       accountData["filevault_enabled"].(bool),
					})
				}
				outBlock.Accounts = &accounts
			}
		}

		// Handle directory bindings
		if directoryBindings, ok := data["directory_bindings"]; ok && len(directoryBindings.([]any)) > 0 {
			directoryBindingsList := directoryBindings.([]any)
			bindings := []jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings{}
			for _, binding := range directoryBindingsList {
				bindingData := binding.(map[string]any)
				bindings = append(bindings, jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings{
					ID:   bindingData["id"].(int),
					Name: bindingData["name"].(string),
				})
			}
			outBlock.DirectoryBindings = &bindings
		}

		// Handle management account
		if managementAccount, ok := data["management_account"]; ok && len(managementAccount.([]any)) > 0 {
			managementAccountList := managementAccount.([]any)
			if len(managementAccountList) > 0 {
				managementAccountData := managementAccountList[0].(map[string]any)
				outBlock.ManagementAccount = &jamfpro.PolicySubsetAccountMaintenanceManagementAccount{
					Action:                managementAccountData["action"].(string),
					ManagedPassword:       managementAccountData["managed_password"].(string),
					ManagedPasswordLength: managementAccountData["managed_password_length"].(int),
				}
			}
		}

		// Handle open firmware/EFI password
		if openFirmwareEfiPassword, ok := data["open_firmware_efi_password"]; ok && len(openFirmwareEfiPassword.([]any)) > 0 {
			openFirmwareEfiPasswordList := openFirmwareEfiPassword.([]any)
			if len(openFirmwareEfiPasswordList) > 0 {
				openFirmwareEfiPasswordData := openFirmwareEfiPasswordList[0].(map[string]any)
				outBlock.OpenFirmwareEfiPassword = &jamfpro.PolicySubsetAccountMaintenanceOpenFirmwareEfiPassword{
					OfMode:     openFirmwareEfiPasswordData["of_mode"].(string),
					OfPassword: openFirmwareEfiPasswordData["of_password"].(string),
				}
			}
		}
	}

	resource.AccountMaintenance = *outBlock
}

// constructPayloadFilesProcesses builds the files and processes payload settings of the policy.
func constructPayloadFilesProcesses(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.files_processes")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	outBlock := new(jamfpro.PolicySubsetFilesProcesses)
	payload := []jamfpro.PolicySubsetFilesProcesses{}

	for _, v := range hcl.([]any) {
		data := v.(map[string]any)
		payload = append(payload, jamfpro.PolicySubsetFilesProcesses{
			SearchByPath:         data["search_by_path"].(string),
			DeleteFile:           data["delete_file"].(bool),
			LocateFile:           data["locate_file"].(string),
			UpdateLocateDatabase: data["update_locate_database"].(bool),
			SpotlightSearch:      data["spotlight_search"].(string),
			SearchForProcess:     data["search_for_process"].(string),
			KillProcess:          data["kill_process"].(bool),
			RunCommand:           data["run_command"].(string),
		})
	}

	if len(payload) > 0 {
		outBlock = &payload[0]
		resource.FilesProcesses = *outBlock
	}

}

// constructPayloadUserInteraction builds the user interaction payload settings of the policy.
func constructPayloadUserInteraction(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.user_interaction")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	outBlock := new(jamfpro.PolicySubsetUserInteraction)
	payload := []jamfpro.PolicySubsetUserInteraction{}

	for _, v := range hcl.([]any) {
		data := v.(map[string]any)
		payload = append(payload, jamfpro.PolicySubsetUserInteraction{
			MessageStart:          data["message_start"].(string),
			AllowUsersToDefer:     data["allow_users_to_defer"].(bool),
			AllowDeferralUntilUtc: data["allow_deferral_until_utc"].(string),
			AllowDeferralMinutes:  data["allow_deferral_minutes"].(int),
			MessageFinish:         data["message_finish"].(string),
		})
	}

	outBlock = &payload[0]
	resource.UserInteraction = *outBlock

}

// constructPayloadReboot builds the reboot payload settings of the policy.
func constructPayloadReboot(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.reboot")
	if len(hcl.([]any)) == 0 {
		resource.Reboot = jamfpro.PolicySubsetReboot{StartupDisk: "Current Startup Disk"}
		return
	}

	hcl = d.Get("payloads.0.reboot.0")

	var payload jamfpro.PolicySubsetReboot

	payload.Message = hcl.(map[string]any)["message"].(string)
	payload.SpecifyStartup = hcl.(map[string]any)["specify_startup"].(string)
	payload.StartupDisk = hcl.(map[string]any)["startup_disk"].(string)
	payload.NoUserLoggedIn = hcl.(map[string]any)["no_user_logged_in"].(string)
	payload.UserLoggedIn = hcl.(map[string]any)["user_logged_in"].(string)
	payload.MinutesUntilReboot = hcl.(map[string]any)["minutes_until_reboot"].(int)
	payload.StartRebootTimerImmediately = hcl.(map[string]any)["start_reboot_timer_immediately"].(bool)
	payload.FileVault2Reboot = hcl.(map[string]any)["file_vault_2_reboot"].(bool)

	resource.Reboot = payload

}

// constructPayloadMaintenance builds the maintenance payload settings of the policy.
func constructPayloadMaintenance(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) {
	hcl := d.Get("payloads.0.maintenance")
	if hcl == nil || len(hcl.([]any)) == 0 {
		return
	}

	outBlock := new(jamfpro.PolicySubsetMaintenance)
	payload := []jamfpro.PolicySubsetMaintenance{}

	for _, v := range hcl.([]any) {
		data := v.(map[string]any)
		payload = append(payload, jamfpro.PolicySubsetMaintenance{
			Recon:                    data["recon"].(bool),
			ResetName:                data["reset_name"].(bool),
			InstallAllCachedPackages: data["install_all_cached_packages"].(bool),
			Heal:                     data["heal"].(bool),
			Prebindings:              data["prebindings"].(bool),
			Permissions:              data["permissions"].(bool),
			Byhost:                   data["byhost"].(bool),
			SystemCache:              data["system_cache"].(bool),
			UserCache:                data["user_cache"].(bool),
			Verify:                   data["verify"].(bool),
		})
	}

	outBlock = &payload[0]
	resource.Maintenance = *outBlock

}
