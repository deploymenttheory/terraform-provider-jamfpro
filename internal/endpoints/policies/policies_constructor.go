package policies

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Returns ResourcePolicy required for client to marshal into api req
func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	var err error
	resource := &jamfpro.ResourcePolicy{}

	constructGeneral(d, resource)

	err = constructScope(d, resource)
	if err != nil {
		return nil, err
	}

	constructSelfService(d, resource)

	err = constructPayloads(d, resource)
	if err != nil {
		return nil, err
	}

	// DEBUG
	log.Println("LOGHERE-CONSTRUCTED")
	policyXML, _ := xml.MarshalIndent(resource, "", "  ")
	log.Println(string(policyXML))

	return resource, nil
}

// Pulls "general" settings from HCL and packages into object
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
	}

	// Category
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	// Site
	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
}

// Pulls "scope" settings from HCL and packages into object
func constructScope(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	var err error

	if len(d.Get("scope").([]interface{})) == 0 {
		return nil
	}

	// Targets
	resource.Scope = &jamfpro.PolicySubsetScope{
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
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputer, int]("scope.0.computer_ids", "ID", d, resource.Scope.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, resource.Scope.ComputerGroups)
	if err != nil {
		return err
	}

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, resource.Scope.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, resource.Scope.JSSUserGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.building_ids", "ID", d, resource.Scope.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.department_ids", "ID", d, resource.Scope.Departments)
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
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, resource.Scope.Limitations.NetworkSegments)
	if err != nil {
		return err
	}

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, resource.Scope.Limitations.IBeacons)
	if err != nil {
		return err
	}

	// User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetUserGroup, int]("scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, resource.Scope.Limitations.UserGroups)
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
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, resource.Scope.Exclusions.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, resource.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, resource.Scope.Exclusions.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, resource.Scope.Exclusions.Departments)
	if err != nil {
		return err
	}

	// Network Segments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, resource.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return err
	}

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, resource.Scope.Exclusions.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, resource.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return err
	}

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, resource.Scope.Exclusions.IBeacons)
	if err != nil {
		return err
	}

	// TODO make this better, it works for now
	if resource.Scope.AllComputers && (resource.Scope.Computers != nil ||
		resource.Scope.ComputerGroups != nil ||
		resource.Scope.Departments != nil ||
		resource.Scope.Buildings != nil) {
		return fmt.Errorf("invalid combination - all computers with scoped endpoints")
	}

	return nil
}

// Pulls "self service" settings from HCL and packages into object
func constructSelfService(d *schema.ResourceData, out *jamfpro.ResourcePolicy) {

	if len(d.Get("self_service").([]interface{})) > 0 {
		out.SelfService = &jamfpro.PolicySubsetSelfService{
			UseForSelfService:      d.Get("self_service.0.use_for_self_service").(bool),
			SelfServiceDisplayName: d.Get("self_service.0.self_service_display_name").(string),
			InstallButtonText:      d.Get("self_service.0.install_button_text").(string),
			// ReinstallButtonText:         d.Get("self_service.0.reinstall_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			// TODO self service icon
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
			// TODO Self service categories
		}
	}
}

// constructPayloads pulls "payload" settings from HCL and packages them into the resource.
func constructPayloads(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	if err := constructPayloadPackages(d, resource); err != nil {
		return err
	}
	if err := constructPayloadScripts(d, resource); err != nil {
		return err
	}
	if err := constructPayloadDiskEncryption(d, resource); err != nil {
		return err
	}
	if err := constructPayloadPrinters(d, resource); err != nil {
		return err
	}
	if err := constructPayloadDockItems(d, resource); err != nil {
		return err
	}
	if err := constructPayloadAccountMaintenance(d, resource); err != nil {
		return err
	}
	if err := constructPayloadFilesProcesses(d, resource); err != nil {
		return err
	}
	if err := constructPayloadUserInteraction(d, resource); err != nil {
		return err
	}
	if err := constructPayloadReboot(d, resource); err != nil {
		return err
	}
	if err := constructPayloadMaintenance(d, resource); err != nil {
		return err
	}

	return nil
}

// Pulls "package" settings from HCL and packages them into the resource.
func constructPayloadPackages(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.packages")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	outBlock := new(jamfpro.PolicySubsetPackageConfiguration)
	outBlock.DistributionPoint = d.Get("payloads.0.packages.0.distribution_point").(string)
	outBlock.Packages = &[]jamfpro.PolicySubsetPackageConfigurationPackage{}
	payload := *outBlock.Packages

	for _, v := range hcl.([]interface{}) {
		newObj := jamfpro.PolicySubsetPackageConfigurationPackage{
			ID:                v.(map[string]interface{})["id"].(int),
			Action:            v.(map[string]interface{})["action"].(string),
			FillUserTemplate:  v.(map[string]interface{})["fill_user_template"].(bool),
			FillExistingUsers: v.(map[string]interface{})["fill_existing_user_template"].(bool),
		}
		payload = append(payload, newObj)
	}

	outBlock.Packages = &payload
	resource.PackageConfiguration = outBlock

	return nil
}

// Pulls "script" settings from HCL and packages them into the resource.
func constructPayloadScripts(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.scripts")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	outBlock := new(jamfpro.PolicySubsetScripts)
	outBlock.Script = &[]jamfpro.PolicySubsetScript{}
	payload := *outBlock.Script

	for _, v := range hcl.([]interface{}) {
		newObj := jamfpro.PolicySubsetScript{
			ID:          v.(map[string]interface{})["id"].(string),
			Priority:    v.(map[string]interface{})["priority"].(string),
			Parameter4:  v.(map[string]interface{})["parameter4"].(string),
			Parameter5:  v.(map[string]interface{})["parameter5"].(string),
			Parameter6:  v.(map[string]interface{})["parameter6"].(string),
			Parameter7:  v.(map[string]interface{})["parameter7"].(string),
			Parameter8:  v.(map[string]interface{})["parameter8"].(string),
			Parameter9:  v.(map[string]interface{})["parameter9"].(string),
			Parameter10: v.(map[string]interface{})["parameter10"].(string),
			Parameter11: v.(map[string]interface{})["parameter11"].(string),
		}

		payload = append(payload, newObj)
	}

	outBlock.Script = &payload
	resource.Scripts = outBlock

	return nil
}

// Pulls "disk encryption" settings from HCL and packages them into the resource.
func constructPayloadDiskEncryption(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.disk_encryption")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	resource.DiskEncryption = &jamfpro.PolicySubsetDiskEncryption{
		Action:                                 d.Get("payloads.0.disk_encryption.0.action").(string),
		DiskEncryptionConfigurationID:          d.Get("payloads.0.disk_encryption.0.disk_encryption_configuration_id").(int),
		AuthRestart:                            d.Get("payloads.0.disk_encryption.0.auth_restart").(bool),
		RemediateKeyType:                       d.Get("payloads.0.disk_encryption.0.remediate_key_type").(string),
		RemediateDiskEncryptionConfigurationID: d.Get("payloads.0.disk_encryption.0.remediate_disk_encryption_configuration_id").(int),
	}

	return nil
}

// Pulls "printers" settings from HCL and packages them into the resource.
func constructPayloadPrinters(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.printers")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	outBlock := new(jamfpro.PolicySubsetPrinters)
	outBlock.Printer = &[]jamfpro.PolicySubsetPrinter{}
	payload := *outBlock.Printer

	for _, v := range hcl.([]interface{}) {
		newObj := jamfpro.PolicySubsetPrinter{
			ID:          v.(map[string]interface{})["id"].(int),
			Action:      v.(map[string]interface{})["action"].(string),
			MakeDefault: v.(map[string]interface{})["make_default"].(bool),
		}
		payload = append(payload, newObj)
	}

	outBlock.Printer = &payload
	resource.Printers = outBlock

	return nil
}

// constructPayloadDockItems builds the dock items payload settings of the policy.
func constructPayloadDockItems(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.dock_items")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	outBlock := new(jamfpro.PolicySubsetDockItems)
	outBlock.DockItem = &[]jamfpro.PolicySubsetDockItem{}
	payload := *outBlock.DockItem

	for _, v := range hcl.([]interface{}) {
		newObj := jamfpro.PolicySubsetDockItem{
			ID:     v.(map[string]interface{})["id"].(int),
			Action: v.(map[string]interface{})["action"].(string),
		}
		payload = append(payload, newObj)
	}

	outBlock.DockItem = &payload
	resource.DockItems = outBlock

	return nil
}

// constructPayloadAccountMaintenance builds the account maintenance payload settings of the policy.
func constructPayloadAccountMaintenance(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.account_maintenance")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	outBlock := new(jamfpro.PolicySubsetAccountMaintenance)
	payload := []jamfpro.PolicySubsetAccountMaintenanceAccount{}

	for _, v := range hcl.([]interface{}) {
		data := v.(map[string]interface{})
		newObj := jamfpro.PolicySubsetAccountMaintenanceAccount{
			Action:                 data["action"].(string),
			Username:               data["username"].(string),
			Realname:               data["realname"].(string),
			Password:               data["password"].(string),
			ArchiveHomeDirectory:   data["archive_home_directory"].(bool),
			ArchiveHomeDirectoryTo: data["archive_home_directory_to"].(string),
			Home:                   data["home"].(string),
			Hint:                   data["hint"].(string),
			Picture:                data["picture"].(string),
			Admin:                  data["admin"].(bool),
			FilevaultEnabled:       data["filevault_enabled"].(bool),
		}
		payload = append(payload, newObj)
	}

	outBlock.Accounts = &payload
	resource.AccountMaintenance = outBlock

	return nil
}

// constructPayloadFilesProcesses builds the files and processes payload settings of the policy.
func constructPayloadFilesProcesses(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.files_processes")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	payload := []jamfpro.PolicySubsetFilesProcesses{}

	for _, v := range hcl.([]interface{}) {
		data := v.(map[string]interface{})
		newObj := jamfpro.PolicySubsetFilesProcesses{
			SearchByPath:         data["search_by_path"].(string),
			DeleteFile:           data["delete_file"].(bool),
			LocateFile:           data["locate_file"].(string),
			UpdateLocateDatabase: data["update_locate_database"].(bool),
			SpotlightSearch:      data["spotlight_search"].(string),
			SearchForProcess:     data["search_for_process"].(string),
			KillProcess:          data["kill_process"].(bool),
			RunCommand:           data["run_command"].(string),
		}
		payload = append(payload, newObj)
	}

	if len(payload) > 0 {
		resource.FilesProcesses = &payload[0]
	}

	return nil
}

// constructPayloadUserInteraction builds the user interaction payload settings of the policy.
func constructPayloadUserInteraction(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.user_interaction")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	resource.UserInteraction = &jamfpro.PolicySubsetUserInteraction{
		MessageStart:          d.Get("payloads.0.user_interaction.0.message_start").(string),
		AllowUserToDefer:      d.Get("payloads.0.user_interaction.0.allow_user_to_defer").(bool),
		AllowDeferralUntilUtc: d.Get("payloads.0.user_interaction.0.allow_deferral_until_utc").(string),
		AllowDeferralMinutes:  d.Get("payloads.0.user_interaction.0.allow_deferral_minutes").(int),
		MessageFinish:         d.Get("payloads.0.user_interaction.0.message_finish").(string),
	}

	return nil
}

// constructPayloadReboot builds the reboot payload settings of the policy.
func constructPayloadReboot(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.reboot")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	resource.Reboot = &jamfpro.PolicySubsetReboot{
		Message:                     d.Get("payloads.0.reboot.0.message").(string),
		SpecifyStartup:              d.Get("payloads.0.reboot.0.specify_startup").(string),
		StartupDisk:                 d.Get("payloads.0.reboot.0.startup_disk").(string),
		NoUserLoggedIn:              d.Get("payloads.0.reboot.0.no_user_logged_in").(string),
		UserLoggedIn:                d.Get("payloads.0.reboot.0.user_logged_in").(string),
		MinutesUntilReboot:          d.Get("payloads.0.reboot.0.minutes_until_reboot").(int),
		StartRebootTimerImmediately: d.Get("payloads.0.reboot.0.start_reboot_timer_immediately").(bool),
		FileVault2Reboot:            d.Get("payloads.0.reboot.0.file_vault_2_reboot").(bool),
	}

	return nil
}

// constructPayloadMaintenance builds the maintenance payload settings of the policy.
func constructPayloadMaintenance(d *schema.ResourceData, resource *jamfpro.ResourcePolicy) error {
	hcl := d.Get("payloads.0.maintenance")
	if len(hcl.([]interface{})) == 0 {
		return nil
	}

	resource.Maintenance = &jamfpro.PolicySubsetMaintenance{
		Recon:                    d.Get("payloads.0.maintenance.0.recon").(bool),
		ResetName:                d.Get("payloads.0.maintenance.0.reset_name").(bool),
		InstallAllCachedPackages: d.Get("payloads.0.maintenance.0.install_all_cached_packages").(bool),
		Heal:                     d.Get("payloads.0.maintenance.0.heal").(bool),
		Prebindings:              d.Get("payloads.0.maintenance.0.prebindings").(bool),
		Permissions:              d.Get("payloads.0.maintenance.0.permissions").(bool),
		Byhost:                   d.Get("payloads.0.maintenance.0.byhost").(bool),
		SystemCache:              d.Get("payloads.0.maintenance.0.system_cache").(bool),
		UserCache:                d.Get("payloads.0.maintenance.0.user_cache").(bool),
		Verify:                   d.Get("payloads.0.maintenance.0.verify").(bool),
	}

	return nil
}
