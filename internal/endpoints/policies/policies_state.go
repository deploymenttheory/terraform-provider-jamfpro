package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"
	"reflect"

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

// stateGeneral Reads response and states general/root level item block
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

	err = d.Set("target_drive", resp.General.OverrideDefaultSettings.TargetDrive)
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

	if resp.General.NetworkRequirements != "" {
		err = d.Set("network_requirements", resp.General.NetworkRequirements)
		if err != nil {
			*diags = append(*diags, diag.FromErr(err)...)
		}
	}

	// Site
	err = d.Set("site_id", resp.General.Site.ID)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	// Category
	err = d.Set("category_id", resp.General.Category.ID)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	// Set DateTime Limitations
	setGeneralDateTimeLimitations(d, resp, diags)

	// Set Network Limitations
	setGeneralNetworkLimitations(d, resp, diags)

}

// setGeneralDateTimeLimitations is a helper function to set the date_time_limitations block under general
func setGeneralDateTimeLimitations(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.General.DateTimeLimitations == nil {
		return
	}

	// Check if all values are at their default (empty string or zero value)
	v := reflect.ValueOf(*resp.General.DateTimeLimitations)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if (field.Kind() == reflect.String && field.String() != "") ||
			(field.Kind() == reflect.Int && field.Int() != 0) {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the date_time_limitations block
	dateTimeLimitations := make(map[string]interface{})
	dateTimeLimitations["activation_date"] = resp.General.DateTimeLimitations.ActivationDate
	dateTimeLimitations["activation_date_epoch"] = resp.General.DateTimeLimitations.ActivationDateEpoch
	dateTimeLimitations["activation_date_utc"] = resp.General.DateTimeLimitations.ActivationDateUTC
	dateTimeLimitations["expiration_date"] = resp.General.DateTimeLimitations.ExpirationDate
	dateTimeLimitations["expiration_date_epoch"] = resp.General.DateTimeLimitations.ExpirationDateEpoch
	dateTimeLimitations["expiration_date_utc"] = resp.General.DateTimeLimitations.ExpirationDateUTC
	dateTimeLimitations["no_execute_start"] = resp.General.DateTimeLimitations.NoExecuteStart
	dateTimeLimitations["no_execute_end"] = resp.General.DateTimeLimitations.NoExecuteEnd

	err := d.Set("date_time_limitations", []interface{}{dateTimeLimitations})
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// setGeneralNetworkLimitations is a helper function to set the network_limitations block under general
func setGeneralNetworkLimitations(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.General.NetworkLimitations == nil {
		return
	}

	// Check if all values are at their default (true, "No Minimum", or empty string)
	v := reflect.ValueOf(*resp.General.NetworkLimitations)
	allDefault := true

	defaults := map[string]interface{}{
		"MinimumNetworkConnection": "No Minimum",
		"AnyIPAddress":             true,
		"NetworkSegments":          "",
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() != defaults[fieldName] {
				allDefault = false
			}
		case reflect.String:
			if field.String() != defaults[fieldName] {
				allDefault = false
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the network_limitations block
	networkLimitations := make(map[string]interface{})
	networkLimitations["minimum_network_connection"] = resp.General.NetworkLimitations.MinimumNetworkConnection
	networkLimitations["any_ip_address"] = resp.General.NetworkLimitations.AnyIPAddress
	//Appears to be removed from gui
	//networkLimitations["network_segments"] = resp.General.NetworkLimitations.NetworkSegments

	err := d.Set("network_limitations", []interface{}{networkLimitations})
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
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

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.SelfService == nil {
		return
	}

	defaults := map[string]interface{}{
		"use_for_self_service":            false,
		"self_service_display_name":       "",
		"install_button_text":             "Install",
		"self_service_description":        "",
		"force_users_to_view_description": false,
		"feature_on_main_page":            false,
	}

	current := map[string]interface{}{
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"self_service_description":        resp.SelfService.SelfServiceDescription,
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
	}

	nonDefault := false
	for key, value := range current {
		if value != defaults[key] {
			nonDefault = true
			break
		}
	}

	if !nonDefault {
		log.Println("[DEBUG] Self-service block has only default values, skipping state")
		return
	}

	log.Println("[DEBUG] Initializing self-service block in state")
	out_ss := make([]map[string]interface{}, 0)
	out_ss = append(out_ss, make(map[string]interface{}, 1))

	out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
	out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
	out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
	out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
	out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
	out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage

	err := d.Set("self_service", out_ss)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
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

// Reads response and preps disk encryption payload items
func prepStatePayloadDiskEncryption(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.DiskEncryption == nil {
		return
	}

	// Check if all values are at their default (false, empty string, or specific defaults)
	v := reflect.ValueOf(*resp.DiskEncryption)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() {
				allDefault = false
			}
		case reflect.String:
			if field.String() != "" && field.String() != "none" {
				allDefault = false
			}
		case reflect.Int:
			if field.Int() != 0 {
				allDefault = false
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the disk_encryption block
	(*out)[0]["disk_encryption"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["action"] = resp.DiskEncryption.Action
	outMap["disk_encryption_configuration_id"] = resp.DiskEncryption.DiskEncryptionConfigurationID
	outMap["auth_restart"] = resp.DiskEncryption.AuthRestart
	if resp.DiskEncryption.RemediateKeyType != "" {
		outMap["remediate_key_type"] = resp.DiskEncryption.RemediateKeyType
	}
	if resp.DiskEncryption.RemediateDiskEncryptionConfigurationID != 0 {
		outMap["remediate_disk_encryption_configuration_id"] = resp.DiskEncryption.RemediateDiskEncryptionConfigurationID
	}
	(*out)[0]["disk_encryption"] = append((*out)[0]["disk_encryption"].([]map[string]interface{}), outMap)
}

// Reads response and preps package payload items
func prepStatePayloadPackages(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.PackageConfiguration == nil {
		log.Println("No package configuration found")
		return
	}
	// Packages can be nil but deployment state default
	if resp.PackageConfiguration.Packages == nil {
		log.Println("No packages found in package configuration")
		return
	}

	// Ensure the map is initialized before setting values
	if len((*out)[0]) == 0 {
		(*out)[0] = make(map[string]interface{})
	}

	log.Println("Initializing packages in state")
	packagesMap := make(map[string]interface{})
	packagesMap["distribution_point"] = resp.PackageConfiguration.DistributionPoint
	packagesMap["package"] = make([]map[string]interface{}, 0)

	for _, v := range *resp.PackageConfiguration.Packages {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["action"] = v.Action
		outMap["fill_user_template"] = v.FillUserTemplate
		outMap["fill_existing_user_template"] = v.FillExistingUsers
		log.Printf("Adding package to state: %+v\n", outMap)
		packagesMap["package"] = append(packagesMap["package"].([]map[string]interface{}), outMap)
	}

	(*out)[0]["packages"] = []map[string]interface{}{packagesMap}
	log.Printf("Final state packages: %+v\n", (*out)[0]["packages"])
}

// Reads response and preps script payload items
func prepStatePayloadScripts(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Scripts.Script == nil {
		log.Println("No scripts found")
		return
	}

	// Ensure the map is initialized before setting values
	if len((*out)[0]) == 0 {
		(*out)[0] = make(map[string]interface{})
	}

	log.Println("Initializing scripts in state")
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
		log.Printf("Adding script to state: %+v\n", outMap)
		(*out)[0]["scripts"] = append((*out)[0]["scripts"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state scripts: %+v\n", (*out)[0]["scripts"])
}

// prepStatePayloadPrinters reads response and preps printer payload items for stating
func prepStatePayloadPrinters(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Printers.Printer == nil {
		log.Println("No printers found")
		return
	}

	// Ensure the map is initialized before setting values
	if len((*out)[0]) == 0 {
		(*out)[0] = make(map[string]interface{})
	}

	log.Println("Initializing printers in state")
	(*out)[0]["printers"] = make([]map[string]interface{}, 0)

	for _, v := range *resp.Printers.Printer {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action
		outMap["make_default"] = v.MakeDefault

		log.Printf("Adding printer to state: %+v\n", outMap)
		(*out)[0]["printers"] = append((*out)[0]["printers"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state printers: %+v\n", (*out)[0]["printers"])
}

// Reads response and preps dock items payload items
func prepStatePayloadDockItems(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.DockItems.DockItem == nil {
		log.Println("No dock items found")
		return
	}

	// Ensure the map is initialized before setting values
	if len((*out)[0]) == 0 {
		(*out)[0] = make(map[string]interface{})
	}

	log.Println("Initializing dock items in state")
	(*out)[0]["dock_items"] = make([]map[string]interface{}, 0)

	for _, v := range *resp.DockItems.DockItem {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action

		log.Printf("Adding dock item to state: %+v\n", outMap)
		(*out)[0]["dock_items"] = append((*out)[0]["dock_items"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state dock items: %+v\n", (*out)[0]["dock_items"])
}

// prepStatePayloadAccountMaintenance reads response and preps account maintenance payload items
func prepStatePayloadAccountMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.AccountMaintenance == nil {
		log.Println("No account maintenance configuration found")
		return
	}

	// Ensure the map is initialized before setting values
	if len((*out)[0]) == 0 {
		(*out)[0] = make(map[string]interface{})
	}

	log.Println("Initializing account maintenance in state")
	(*out)[0]["account_maintenance"] = make([]map[string]interface{}, 0)

	// Handle accounts
	if resp.AccountMaintenance.Accounts != nil {
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

			log.Printf("Adding account to state: %+v\n", outMap)
			(*out)[0]["account_maintenance"] = append((*out)[0]["account_maintenance"].([]map[string]interface{}), outMap)
		}
	}

	// Handle directory bindings
	if resp.AccountMaintenance.DirectoryBindings != nil {
		for _, v := range *resp.AccountMaintenance.DirectoryBindings {
			outMap := make(map[string]interface{})
			outMap["id"] = v.ID
			outMap["name"] = v.Name

			log.Printf("Adding directory binding to state: %+v\n", outMap)
			(*out)[0]["account_maintenance"] = append((*out)[0]["account_maintenance"].([]map[string]interface{}), outMap)
		}
	}

	// Handle management account
	if resp.AccountMaintenance.ManagementAccount != nil {
		outMap := make(map[string]interface{})
		outMap["action"] = resp.AccountMaintenance.ManagementAccount.Action
		outMap["managed_password"] = resp.AccountMaintenance.ManagementAccount.ManagedPassword
		outMap["managed_password_length"] = resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength

		log.Printf("Adding management account to state: %+v\n", outMap)
		(*out)[0]["account_maintenance"] = append((*out)[0]["account_maintenance"].([]map[string]interface{}), outMap)
	}

	// Handle open firmware/EFI password
	if resp.AccountMaintenance.OpenFirmwareEfiPassword != nil {
		outMap := make(map[string]interface{})
		outMap["of_mode"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode
		outMap["of_password"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword

		log.Printf("Adding open firmware/EFI password to state: %+v\n", outMap)
		(*out)[0]["account_maintenance"] = append((*out)[0]["account_maintenance"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state account maintenance: %+v\n", (*out)[0]["account_maintenance"])
}

// Reads response and preps files and processes payload items
func prepStatePayloadFilesProcesses(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.FilesProcesses == nil {
		return
	}

	// Do not set the files_processes block if all values are default (false or empty string)
	v := reflect.ValueOf(*resp.FilesProcesses)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() {
				allDefault = false
				break
			}
		case reflect.String:
			if field.String() != "" {
				allDefault = false
				break
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the files_processes block
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

// Reads response and preps user interaction payload items. If all values are default, do not set the user_interaction block
func prepStatePayloadUserInteraction(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.UserInteraction == nil {
		return
	}

	// Do not set the user_interaction block if all values are default (false, empty string, or 0)
	v := reflect.ValueOf(*resp.UserInteraction)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() {
				allDefault = false
				break
			}
		case reflect.String:
			if field.String() != "" {
				allDefault = false
				break
			}
		case reflect.Int:
			if field.Int() != 0 {
				allDefault = false
				break
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the user_interaction block
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

	// Define default values
	defaults := map[string]interface{}{
		"Message":                     "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu.",
		"SpecifyStartup":              "",
		"StartupDisk":                 "Current Startup Disk",
		"NoUserLoggedIn":              "Do not restart",
		"UserLoggedIn":                "Do not restart",
		"MinutesUntilReboot":          5,
		"StartRebootTimerImmediately": false,
		"FileVault2Reboot":            false,
	}

	// Check if all values are default
	v := reflect.ValueOf(*resp.Reboot)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		switch field.Kind() {
		case reflect.String:
			if field.String() != defaults[fieldName].(string) {
				allDefault = false
			}
		case reflect.Int:
			if field.Int() != int64(defaults[fieldName].(int)) {
				allDefault = false
			}
		case reflect.Bool:
			if field.Bool() != defaults[fieldName].(bool) {
				allDefault = false
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the reboot block
	rebootBlock := make(map[string]interface{})
	if resp.Reboot.Message != defaults["Message"].(string) {
		rebootBlock["message"] = resp.Reboot.Message
	}
	if resp.Reboot.SpecifyStartup != defaults["SpecifyStartup"].(string) {
		rebootBlock["specify_startup"] = resp.Reboot.SpecifyStartup
	}
	if resp.Reboot.StartupDisk != defaults["StartupDisk"].(string) {
		rebootBlock["startup_disk"] = resp.Reboot.StartupDisk
	}
	if resp.Reboot.NoUserLoggedIn != defaults["NoUserLoggedIn"].(string) {
		rebootBlock["no_user_logged_in"] = resp.Reboot.NoUserLoggedIn
	}
	if resp.Reboot.UserLoggedIn != defaults["UserLoggedIn"].(string) {
		rebootBlock["user_logged_in"] = resp.Reboot.UserLoggedIn
	}
	if resp.Reboot.MinutesUntilReboot != defaults["MinutesUntilReboot"].(int) {
		rebootBlock["minutes_until_reboot"] = resp.Reboot.MinutesUntilReboot
	}
	if resp.Reboot.StartRebootTimerImmediately != defaults["StartRebootTimerImmediately"].(bool) {
		rebootBlock["start_reboot_timer_immediately"] = resp.Reboot.StartRebootTimerImmediately
	}
	if resp.Reboot.FileVault2Reboot != defaults["FileVault2Reboot"].(bool) {
		rebootBlock["file_vault_2_reboot"] = resp.Reboot.FileVault2Reboot
	}

	if len(rebootBlock) > 0 {
		(*out)[0]["reboot"] = []map[string]interface{}{rebootBlock}
	}
}

// prepStatePayloadMaintenance Reads response and preps maintenance payload items. If all values are default, do not set the maintenance block
func prepStatePayloadMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Maintenance == nil {
		return
	}

	// Do not set the maintenance block if all values are default (false)
	v := reflect.ValueOf(*resp.Maintenance)
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
	// Else, set the maintenance block
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
