package policies

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJamfProPoliciesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the policy object
	resource, err := constructJamfProPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy: %v", err))
	}

	// Extract policy name from schema
	var policyName string
	if generalSettings, ok := d.GetOk("general"); ok && len(generalSettings.([]interface{})) > 0 {
		generalMap := generalSettings.([]interface{})[0].(map[string]interface{})
		policyName = generalMap["name"].(string)
	}

	// Retry the API call to create the policy in Jamf Pro
	var creationResponse *jamfpro.ResourcePolicyCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreatePolicy(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Policy '%s' after retries: %v", policyName, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the policy to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPoliciesRead is responsible for reading the current state of a Jamf Pro policy Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var policy *jamfpro.ResourcePolicy

	// Extract policy name from schema
	var policyName string
	if generalSettings, ok := d.GetOk("general"); ok && len(generalSettings.([]interface{})) > 0 {
		generalMap := generalSettings.([]interface{})[0].(map[string]interface{})
		policyName = generalMap["name"].(string)
	}

	// Use the retry function for the read operation
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		policy, apiErr = conn.GetPolicyByID(resourceIDInt)
		if apiErr != nil {
			// If fetching by ID fails and policyName is available, try fetching by Name
			if policyName != "" {
				policy, apiErr = conn.GetPolicyByName(policyName)
			}
			if apiErr != nil {
				// Consider retrying only if it's a retryable error
				return retry.RetryableError(apiErr)
			}
		}
		// Successfully fetched the policy, exit the retry loop
		return nil
	})

	if err != nil {
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Policy '%s' (ID: %d) after retries: %v", policyName, resourceIDInt, err))
	}

	// Update the Terraform state with the fetched data
	// Set 'general' attributes
	generalAttributes := map[string]interface{}{
		"id":                            policy.General.ID,
		"name":                          policy.General.Name,
		"enabled":                       policy.General.Enabled,
		"trigger":                       policy.General.Trigger,
		"trigger_checkin":               policy.General.TriggerCheckin,
		"trigger_enrollment_complete":   policy.General.TriggerEnrollmentComplete,
		"trigger_login":                 policy.General.TriggerLogin,
		"trigger_logout":                policy.General.TriggerLogout,
		"trigger_network_state_changed": policy.General.TriggerNetworkStateChanged,
		"trigger_startup":               policy.General.TriggerStartup,
		"trigger_other":                 policy.General.TriggerOther,
		"frequency":                     policy.General.Frequency,
		"retry_event":                   policy.General.RetryEvent,
		"retry_attempts":                policy.General.RetryAttempts,
		"notify_on_each_failed_retry":   policy.General.NotifyOnEachFailedRetry,
		"location_user_only":            policy.General.LocationUserOnly,
		"target_drive":                  policy.General.TargetDrive,
		"offline":                       policy.General.Offline,
		"category": []interface{}{map[string]interface{}{
			"id":   policy.General.Category.ID,
			"name": policy.General.Category.Name,
		}},
		"date_time_limitations": []interface{}{
			map[string]interface{}{
				"activation_date":       policy.General.DateTimeLimitations.ActivationDate,
				"activation_date_epoch": policy.General.DateTimeLimitations.ActivationDateEpoch,
				"activation_date_utc":   policy.General.DateTimeLimitations.ActivationDateUTC,
				"expiration_date":       policy.General.DateTimeLimitations.ExpirationDate,
				"expiration_date_epoch": policy.General.DateTimeLimitations.ExpirationDateEpoch,
				"expiration_date_utc":   policy.General.DateTimeLimitations.ExpirationDateUTC,
				"no_execute_on": func() []interface{} {
					noExecOnDays := make([]interface{}, len(policy.General.DateTimeLimitations.NoExecuteOn))
					for i, noExecOn := range policy.General.DateTimeLimitations.NoExecuteOn {
						noExecOnDays[i] = map[string]interface{}{"day": noExecOn.Day}
					}
					return noExecOnDays
				}(),
				"no_execute_start": policy.General.DateTimeLimitations.NoExecuteStart,
				"no_execute_end":   policy.General.DateTimeLimitations.NoExecuteEnd,
			},
		},
		"network_limitations": []interface{}{map[string]interface{}{
			"minimum_network_connection": policy.General.NetworkLimitations.MinimumNetworkConnection,
			"any_ip_address":             policy.General.NetworkLimitations.AnyIPAddress,
			"network_segments":           policy.General.NetworkLimitations.NetworkSegments,
		}},
		"override_default_settings": []interface{}{map[string]interface{}{
			"target_drive":       policy.General.OverrideDefaultSettings.TargetDrive,
			"distribution_point": policy.General.OverrideDefaultSettings.DistributionPoint,
			"force_afp_smb":      policy.General.OverrideDefaultSettings.ForceAfpSmb,
			"sus":                policy.General.OverrideDefaultSettings.SUS,
		}},
		"network_requirements": policy.General.NetworkRequirements,
		"site": []interface{}{map[string]interface{}{
			"id":   policy.General.Site.ID,
			"name": policy.General.Site.Name,
		}},
	}

	if err := d.Set("general", []interface{}{generalAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set 'scope' attributes
	scopeAttributes := map[string]interface{}{
		"all_computers": policy.Scope.AllComputers,
		"computers": func() []interface{} {
			computersInterfaces := make([]interface{}, len(policy.Scope.Computers))
			for i, computer := range policy.Scope.Computers {
				computersInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			return computersInterfaces
		}(),
		"computer_groups": func() []interface{} {
			groupInterfaces := make([]interface{}, len(policy.Scope.ComputerGroups))
			for i, group := range policy.Scope.ComputerGroups {
				groupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return groupInterfaces
		}(),
		"jss_users": func() []interface{} {
			userInterfaces := make([]interface{}, len(policy.Scope.JSSUsers))
			for i, user := range policy.Scope.JSSUsers {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			return userInterfaces
		}(),
		"jss_user_groups": func() []interface{} {
			userGroupInterfaces := make([]interface{}, len(policy.Scope.JSSUserGroups))
			for i, group := range policy.Scope.JSSUserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return userGroupInterfaces
		}(),
		"buildings": func() []interface{} {
			buildingInterfaces := make([]interface{}, len(policy.Scope.Buildings))
			for i, building := range policy.Scope.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			return buildingInterfaces
		}(),
		"departments": func() []interface{} {
			departmentInterfaces := make([]interface{}, len(policy.Scope.Departments))
			for i, department := range policy.Scope.Departments {
				departmentInterfaces[i] = map[string]interface{}{
					"id":   department.ID,
					"name": department.Name,
				}
			}
			return departmentInterfaces
		}(),
		"limitations": func() []interface{} {
			limitationInterfaces := make([]interface{}, 0)
			limitationData := map[string]interface{}{}

			// Network Segments
			networkSegmentInterfaces := make([]interface{}, len(policy.Scope.Limitations.NetworkSegments))
			for i, segment := range policy.Scope.Limitations.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			limitationData["network_segments"] = networkSegmentInterfaces

			// Users
			userInterfaces := make([]interface{}, len(policy.Scope.Limitations.Users))
			for i, user := range policy.Scope.Limitations.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			limitationData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(policy.Scope.Limitations.UserGroups))
			for i, group := range policy.Scope.Limitations.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			limitationData["user_groups"] = userGroupInterfaces

			// iBeacons
			iBeaconInterfaces := make([]interface{}, len(policy.Scope.Limitations.IBeacons))
			for i, beacon := range policy.Scope.Limitations.IBeacons {
				iBeaconInterfaces[i] = map[string]interface{}{
					"id":   beacon.ID,
					"name": beacon.Name,
				}
			}
			limitationData["ibeacons"] = iBeaconInterfaces

			limitationInterfaces = append(limitationInterfaces, limitationData)
			return limitationInterfaces
		}(),
		"exclusions": func() []interface{} {
			exclusionsInterfaces := make([]interface{}, 0)
			exclusionsData := map[string]interface{}{}

			// Computers
			computerInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Computers))
			for i, computer := range policy.Scope.Exclusions.Computers {
				computerInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			exclusionsData["computers"] = computerInterfaces

			// Computer Groups
			computerGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.ComputerGroups))
			for i, group := range policy.Scope.Exclusions.ComputerGroups {
				computerGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["computer_groups"] = computerGroupInterfaces

			// Users
			userInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Users))
			for i, user := range policy.Scope.Exclusions.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.UserGroups))
			for i, group := range policy.Scope.Exclusions.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["user_groups"] = userGroupInterfaces

			// Buildings
			buildingInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Buildings))
			for i, building := range policy.Scope.Exclusions.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			exclusionsData["buildings"] = buildingInterfaces

			// Departments
			departmentInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Departments))
			for i, department := range policy.Scope.Exclusions.Departments {
				departmentInterfaces[i] = map[string]interface{}{
					"id":   department.ID,
					"name": department.Name,
				}
			}
			exclusionsData["departments"] = departmentInterfaces

			// Network Segments
			networkSegmentInterfaces := make([]interface{}, len(policy.Scope.Exclusions.NetworkSegments))
			for i, segment := range policy.Scope.Exclusions.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			exclusionsData["network_segments"] = networkSegmentInterfaces

			// JSS Users
			jssUserInterfaces := make([]interface{}, len(policy.Scope.Exclusions.JSSUsers))
			for i, user := range policy.Scope.Exclusions.JSSUsers {
				jssUserInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["jss_users"] = jssUserInterfaces

			// JSS User Groups
			jssUserGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.JSSUserGroups))
			for i, group := range policy.Scope.Exclusions.JSSUserGroups {
				jssUserGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["jss_user_groups"] = jssUserGroupInterfaces

			// IBeacons
			iBeaconInterfaces := make([]interface{}, len(policy.Scope.Exclusions.IBeacons))
			for i, beacon := range policy.Scope.Exclusions.IBeacons {
				iBeaconInterfaces[i] = map[string]interface{}{
					"id":   beacon.ID,
					"name": beacon.Name,
				}
			}
			exclusionsData["ibeacons"] = iBeaconInterfaces

			exclusionsInterfaces = append(exclusionsInterfaces, exclusionsData)
			return exclusionsInterfaces
		}(),
	}

	if err := d.Set("scope", []interface{}{scopeAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set 'self_service' attributes
	selfServiceAttributes := map[string]interface{}{
		"use_for_self_service":            policy.SelfService.UseForSelfService,
		"self_service_display_name":       policy.SelfService.SelfServiceDisplayName,
		"install_button_text":             policy.SelfService.InstallButtonText,
		"reinstall_button_text":           policy.SelfService.ReinstallButtonText,
		"self_service_description":        policy.SelfService.SelfServiceDescription,
		"force_users_to_view_description": policy.SelfService.ForceUsersToViewDescription,
		"self_service_icon": []interface{}{map[string]interface{}{
			"id":       policy.SelfService.SelfServiceIcon.ID,
			"filename": policy.SelfService.SelfServiceIcon.Filename,
			"uri":      policy.SelfService.SelfServiceIcon.URI,
		}},
		"feature_on_main_page": policy.SelfService.FeatureOnMainPage,
		"self_service_categories": func() []interface{} {
			categories := make([]interface{}, len(policy.SelfService.SelfServiceCategories))
			for i, cat := range policy.SelfService.SelfServiceCategories {
				categories[i] = map[string]interface{}{
					"id":         cat.Category.ID,
					"name":       cat.Category.Name,
					"display_in": cat.Category.DisplayIn,
					"feature_in": cat.Category.FeatureIn,
				}
			}
			return categories
		}(),
	}

	if err := d.Set("self_service", []interface{}{selfServiceAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Fetch the package configuration from the policy and set it in Terraform state
	packageConfigurations := make([]interface{}, 0)
	for _, packageItem := range policy.PackageConfiguration.Packages {
		pkg := make(map[string]interface{})
		pkg["id"] = packageItem.ID
		pkg["name"] = packageItem.Name
		pkg["action"] = packageItem.Action
		pkg["fut"] = packageItem.FillUserTemplate
		pkg["feu"] = packageItem.FillExistingUsers
		pkg["update_autorun"] = packageItem.UpdateAutorun
		packageConfigurations = append(packageConfigurations, pkg)
	}

	// Wrap packageConfigurations in a map under the key 'packages'
	packageConfiguration := map[string]interface{}{
		"packages": packageConfigurations,
	}

	// Wrap this map in a slice to set in the Terraform state
	if err := d.Set("package_configuration", []interface{}{packageConfiguration}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the scripts from the policy and set them in Terraform state
	scriptConfigurations := make([]interface{}, 0)
	for _, scriptItem := range policy.Scripts.Script {
		script := make(map[string]interface{})
		script["id"] = scriptItem.ID
		script["name"] = scriptItem.Name
		script["priority"] = scriptItem.Priority
		script["parameter4"] = scriptItem.Parameter4
		script["parameter5"] = scriptItem.Parameter5
		script["parameter6"] = scriptItem.Parameter6
		script["parameter7"] = scriptItem.Parameter7
		script["parameter8"] = scriptItem.Parameter8
		script["parameter9"] = scriptItem.Parameter9
		script["parameter10"] = scriptItem.Parameter10
		script["parameter11"] = scriptItem.Parameter11
		scriptConfigurations = append(scriptConfigurations, script)
	}

	if err := d.Set("scripts", scriptConfigurations); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the printers from the policy and set them in Terraform state
	printerConfigurations := make([]interface{}, 0)
	for _, printerItem := range policy.Printers.Printer {
		printer := make(map[string]interface{})
		printer["id"] = printerItem.ID
		printer["name"] = printerItem.Name
		printer["action"] = printerItem.Action
		printer["make_default"] = printerItem.MakeDefault
		printerConfigurations = append(printerConfigurations, printer)
	}

	if err := d.Set("printers", []interface{}{
		map[string]interface{}{
			"leave_existing_default": policy.Printers.LeaveExistingDefault,
			"printer":                printerConfigurations,
		},
	}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the dock items from the policy and set them in Terraform state
	dockItemConfigurations := make([]interface{}, 0)
	for _, dockItem := range policy.DockItems.DockItem {
		dock := make(map[string]interface{})
		dock["id"] = dockItem.ID
		dock["name"] = dockItem.Name
		dock["action"] = dockItem.Action
		dockItemConfigurations = append(dockItemConfigurations, dock)
	}

	if err := d.Set("dock_items", []interface{}{
		map[string]interface{}{
			"dock_item": dockItemConfigurations,
		},
	}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the account maintenance data from the policy and set it in Terraform state
	accountMaintenanceState := make(map[string]interface{})
	accountMaintenanceState["accounts"] = []interface{}{}

	// Add account data if present
	if len(policy.AccountMaintenance.Accounts) > 0 {
		accountsState := make([]interface{}, len(policy.AccountMaintenance.Accounts))
		for i, account := range policy.AccountMaintenance.Accounts {
			accountMap := map[string]interface{}{
				"action":                    account.Action,
				"username":                  account.Username,
				"realname":                  account.Realname,
				"password":                  account.Password,
				"archive_home_directory":    account.ArchiveHomeDirectory,
				"archive_home_directory_to": account.ArchiveHomeDirectoryTo,
				"home":                      account.Home,
				"hint":                      account.Hint,
				"picture":                   account.Picture,
				"admin":                     account.Admin,
				"filevault_enabled":         account.FilevaultEnabled,
			}
			accountsState[i] = map[string]interface{}{"account": accountMap}
		}
		accountMaintenanceState["accounts"] = accountsState
	}

	// Add directory bindings data if present
	if len(policy.AccountMaintenance.DirectoryBindings) > 0 {
		bindingsState := make([]interface{}, len(policy.AccountMaintenance.DirectoryBindings))
		for i, binding := range policy.AccountMaintenance.DirectoryBindings {
			bindingMap := map[string]interface{}{
				"id":   binding.ID,
				"name": binding.Name,
			}
			bindingsState[i] = map[string]interface{}{"binding": bindingMap}
		}
		accountMaintenanceState["directory_bindings"] = bindingsState
	}

	// Add management account data
	accountMaintenanceState["management_account"] = []interface{}{
		map[string]interface{}{
			"action":                  policy.AccountMaintenance.ManagementAccount.Action,
			"managed_password":        policy.AccountMaintenance.ManagementAccount.ManagedPassword,
			"managed_password_length": policy.AccountMaintenance.ManagementAccount.ManagedPasswordLength,
		},
	}

	// Add open firmware/EFI password data
	accountMaintenanceState["open_firmware_efi_password"] = []interface{}{
		map[string]interface{}{
			"of_mode":     policy.AccountMaintenance.OpenFirmwareEfiPassword.OfMode,
			"of_password": policy.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword,
		},
	}

	// Set the account_maintenance in state
	if err := d.Set("account_maintenance", []interface{}{accountMaintenanceState}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the reboot data from the policy and set it in Terraform state
	rebootConfig := make(map[string]interface{})
	rebootConfig["message"] = policy.Reboot.Message
	rebootConfig["specify_startup"] = policy.Reboot.SpecifyStartup
	rebootConfig["startup_disk"] = policy.Reboot.StartupDisk
	rebootConfig["no_user_logged_in"] = policy.Reboot.NoUserLoggedIn
	rebootConfig["user_logged_in"] = policy.Reboot.UserLoggedIn
	rebootConfig["minutes_until_reboot"] = policy.Reboot.MinutesUntilReboot
	rebootConfig["start_reboot_timer_immediately"] = policy.Reboot.StartRebootTimerImmediately
	rebootConfig["file_vault_2_reboot"] = policy.Reboot.FileVault2Reboot

	if err := d.Set("reboot", []interface{}{rebootConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the maintenance data from the policy and set it in Terraform state
	maintenanceConfig := make(map[string]interface{})
	maintenanceConfig["recon"] = policy.Maintenance.Recon
	maintenanceConfig["reset_name"] = policy.Maintenance.ResetName
	maintenanceConfig["install_all_cached_packages"] = policy.Maintenance.InstallAllCachedPackages
	maintenanceConfig["heal"] = policy.Maintenance.Heal
	maintenanceConfig["prebindings"] = policy.Maintenance.Prebindings
	maintenanceConfig["permissions"] = policy.Maintenance.Permissions
	maintenanceConfig["byhost"] = policy.Maintenance.Byhost
	maintenanceConfig["system_cache"] = policy.Maintenance.SystemCache
	maintenanceConfig["user_cache"] = policy.Maintenance.UserCache
	maintenanceConfig["verify"] = policy.Maintenance.Verify

	if err := d.Set("maintenance", []interface{}{maintenanceConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the files and processes data from the policy and set it in Terraform state
	filesProcessesConfig := make(map[string]interface{})
	filesProcessesConfig["search_by_path"] = policy.FilesProcesses.SearchByPath
	filesProcessesConfig["delete_file"] = policy.FilesProcesses.DeleteFile
	filesProcessesConfig["locate_file"] = policy.FilesProcesses.LocateFile
	filesProcessesConfig["update_locate_database"] = policy.FilesProcesses.UpdateLocateDatabase
	filesProcessesConfig["spotlight_search"] = policy.FilesProcesses.SpotlightSearch
	filesProcessesConfig["search_for_process"] = policy.FilesProcesses.SearchForProcess
	filesProcessesConfig["kill_process"] = policy.FilesProcesses.KillProcess
	filesProcessesConfig["run_command"] = policy.FilesProcesses.RunCommand

	if err := d.Set("files_processes", []interface{}{filesProcessesConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the user interaction data from the policy and set it in Terraform state
	userInteractionConfig := make(map[string]interface{})
	userInteractionConfig["message_start"] = policy.UserInteraction.MessageStart
	userInteractionConfig["allow_user_to_defer"] = policy.UserInteraction.AllowUserToDefer
	userInteractionConfig["allow_deferral_until_utc"] = policy.UserInteraction.AllowDeferralUntilUtc
	userInteractionConfig["allow_deferral_minutes"] = policy.UserInteraction.AllowDeferralMinutes
	userInteractionConfig["message_finish"] = policy.UserInteraction.MessageFinish

	if err := d.Set("user_interaction", []interface{}{userInteractionConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the disk encryption data from the policy and set it in Terraform state
	diskEncryptionConfig := make(map[string]interface{})
	diskEncryptionConfig["action"] = policy.DiskEncryption.Action
	diskEncryptionConfig["disk_encryption_configuration_id"] = policy.DiskEncryption.DiskEncryptionConfigurationID
	diskEncryptionConfig["auth_restart"] = policy.DiskEncryption.AuthRestart
	diskEncryptionConfig["remediate_key_type"] = policy.DiskEncryption.RemediateKeyType
	diskEncryptionConfig["remediate_disk_encryption_configuration_id"] = policy.DiskEncryption.RemediateDiskEncryptionConfigurationID

	if err := d.Set("disk_encryption", []interface{}{diskEncryptionConfig}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// ResourceJamfProPoliciesUpdate is responsible for updating an existing Jamf Pro policy on the remote system.
func ResourceJamfProPoliciesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Extract policy name from schema
	var policyName string
	if generalSettings, ok := d.GetOk("general"); ok && len(generalSettings.([]interface{})) > 0 {
		generalMap := generalSettings.([]interface{})[0].(map[string]interface{})
		policyName = generalMap["name"].(string)
	}

	// Construct the resource object
	resource, err := constructJamfProPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdatePolicyByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Policy '%s' (ID: %d) after retries: %v", policyName, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPoliciesDelete is responsible for deleting a Jamf Pro policy.
func ResourceJamfProPoliciesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Extract policy name for error reporting
	generalSettings := d.Get("general").([]interface{})
	generalMap := generalSettings[0].(map[string]interface{})
	resourceName := generalMap["name"].(string)

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeletePolicyByID(resourceIDInt)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			apiErrByName := conn.DeletePolicyByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the site, exit the retry loop
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro policy '%s' (ID: %s) after retries: %v", resourceName, d.Id(), err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
