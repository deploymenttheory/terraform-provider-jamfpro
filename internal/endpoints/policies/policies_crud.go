package policies

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Constructs, creates states
func ResourceJamfProPoliciesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics

	resource, err := constructJamfProPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy: %v", err))
	}

	// Retry the API call to create the policy in Jamf Pro
	var creationResponse *jamfpro.ResourcePolicyCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreatePolicy(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Policy '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// Reads and states
func ResourceJamfProPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resp *jamfpro.ResourcePolicy

	// Extract policy name from schema
	var policyName string
	if generalSettings, ok := d.GetOk("general"); ok && len(generalSettings.([]interface{})) > 0 {
		generalMap := generalSettings.([]interface{})[0].(map[string]interface{})
		policyName = generalMap["name"].(string)
	}

	// Use the retry function for the read operation
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resp, apiErr = conn.GetPolicyByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Policy '%s' (ID: %d) after retries: %v", policyName, resourceIDInt, err))
	}

	// Update the Terraform state with the fetched data
	// Set 'general' attributes
	generalAttributes := map[string]interface{}{
		"id":                            resp.General.ID,
		"name":                          resp.General.Name,
		"enabled":                       resp.General.Enabled,
		"trigger":                       resp.General.Trigger,
		"trigger_checkin":               resp.General.TriggerCheckin,
		"trigger_enrollment_complete":   resp.General.TriggerEnrollmentComplete,
		"trigger_login":                 resp.General.TriggerLogin,
		"trigger_logout":                resp.General.TriggerLogout,
		"trigger_network_state_changed": resp.General.TriggerNetworkStateChanged,
		"trigger_startup":               resp.General.TriggerStartup,
		"trigger_other":                 resp.General.TriggerOther,
		"frequency":                     resp.General.Frequency,
		"retry_event":                   resp.General.RetryEvent,
		"retry_attempts":                resp.General.RetryAttempts,
		"notify_on_each_failed_retry":   resp.General.NotifyOnEachFailedRetry,
		"location_user_only":            resp.General.LocationUserOnly,
		"target_drive":                  resp.General.TargetDrive,
		"offline":                       resp.General.Offline,
		"category": []interface{}{map[string]interface{}{
			"id":   resp.General.Category.ID,
			"name": resp.General.Category.Name,
		}},
		"date_time_limitations": []interface{}{
			map[string]interface{}{
				"activation_date":       resp.General.DateTimeLimitations.ActivationDate,
				"activation_date_epoch": resp.General.DateTimeLimitations.ActivationDateEpoch,
				"activation_date_utc":   resp.General.DateTimeLimitations.ActivationDateUTC,
				"expiration_date":       resp.General.DateTimeLimitations.ExpirationDate,
				"expiration_date_epoch": resp.General.DateTimeLimitations.ExpirationDateEpoch,
				"expiration_date_utc":   resp.General.DateTimeLimitations.ExpirationDateUTC,
				"no_execute_on": func() []interface{} {
					noExecOnDays := make([]interface{}, len(resp.General.DateTimeLimitations.NoExecuteOn))
					for i, noExecOn := range resp.General.DateTimeLimitations.NoExecuteOn {
						noExecOnDays[i] = map[string]interface{}{"day": noExecOn.Day}
					}
					return noExecOnDays
				}(),
				"no_execute_start": resp.General.DateTimeLimitations.NoExecuteStart,
				"no_execute_end":   resp.General.DateTimeLimitations.NoExecuteEnd,
			},
		},
		"network_limitations": []interface{}{map[string]interface{}{
			"minimum_network_connection": resp.General.NetworkLimitations.MinimumNetworkConnection,
			"any_ip_address":             resp.General.NetworkLimitations.AnyIPAddress,
			"network_segments":           resp.General.NetworkLimitations.NetworkSegments,
		}},
		"override_default_settings": []interface{}{map[string]interface{}{
			"target_drive":       resp.General.OverrideDefaultSettings.TargetDrive,
			"distribution_point": resp.General.OverrideDefaultSettings.DistributionPoint,
			"force_afp_smb":      resp.General.OverrideDefaultSettings.ForceAfpSmb,
			"sus":                resp.General.OverrideDefaultSettings.SUS,
		}},
		"network_requirements": resp.General.NetworkRequirements,
		"site": []interface{}{map[string]interface{}{
			"id":   resp.General.Site.ID,
			"name": resp.General.Site.Name,
		}},
	}

	if err := d.Set("general", []interface{}{generalAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set 'scope' attributes
	scopeAttributes := map[string]interface{}{
		"all_computers": resp.Scope.AllComputers,
		"computers": func() []interface{} {
			computersInterfaces := make([]interface{}, len(resp.Scope.Computers))
			for i, computer := range resp.Scope.Computers {
				computersInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			return computersInterfaces
		}(),
		"computer_groups": func() []interface{} {
			groupInterfaces := make([]interface{}, len(resp.Scope.ComputerGroups))
			for i, group := range resp.Scope.ComputerGroups {
				groupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return groupInterfaces
		}(),
		"jss_users": func() []interface{} {
			userInterfaces := make([]interface{}, len(resp.Scope.JSSUsers))
			for i, user := range resp.Scope.JSSUsers {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			return userInterfaces
		}(),
		"jss_user_groups": func() []interface{} {
			userGroupInterfaces := make([]interface{}, len(resp.Scope.JSSUserGroups))
			for i, group := range resp.Scope.JSSUserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return userGroupInterfaces
		}(),
		"buildings": func() []interface{} {
			buildingInterfaces := make([]interface{}, len(resp.Scope.Buildings))
			for i, building := range resp.Scope.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			return buildingInterfaces
		}(),
		"departments": func() []interface{} {
			departmentInterfaces := make([]interface{}, len(resp.Scope.Departments))
			for i, department := range resp.Scope.Departments {
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
			networkSegmentInterfaces := make([]interface{}, len(resp.Scope.Limitations.NetworkSegments))
			for i, segment := range resp.Scope.Limitations.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			limitationData["network_segments"] = networkSegmentInterfaces

			// Users
			userInterfaces := make([]interface{}, len(resp.Scope.Limitations.Users))
			for i, user := range resp.Scope.Limitations.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			limitationData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(resp.Scope.Limitations.UserGroups))
			for i, group := range resp.Scope.Limitations.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			limitationData["user_groups"] = userGroupInterfaces

			// iBeacons
			iBeaconInterfaces := make([]interface{}, len(resp.Scope.Limitations.IBeacons))
			for i, beacon := range resp.Scope.Limitations.IBeacons {
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
			computerInterfaces := make([]interface{}, len(resp.Scope.Exclusions.Computers))
			for i, computer := range resp.Scope.Exclusions.Computers {
				computerInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			exclusionsData["computers"] = computerInterfaces

			// Computer Groups
			computerGroupInterfaces := make([]interface{}, len(resp.Scope.Exclusions.ComputerGroups))
			for i, group := range resp.Scope.Exclusions.ComputerGroups {
				computerGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["computer_groups"] = computerGroupInterfaces

			// Users
			userInterfaces := make([]interface{}, len(resp.Scope.Exclusions.Users))
			for i, user := range resp.Scope.Exclusions.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(resp.Scope.Exclusions.UserGroups))
			for i, group := range resp.Scope.Exclusions.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["user_groups"] = userGroupInterfaces

			// Buildings
			buildingInterfaces := make([]interface{}, len(resp.Scope.Exclusions.Buildings))
			for i, building := range resp.Scope.Exclusions.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			exclusionsData["buildings"] = buildingInterfaces

			// Departments
			departmentInterfaces := make([]interface{}, len(resp.Scope.Exclusions.Departments))
			for i, department := range resp.Scope.Exclusions.Departments {
				departmentInterfaces[i] = map[string]interface{}{
					"id":   department.ID,
					"name": department.Name,
				}
			}
			exclusionsData["departments"] = departmentInterfaces

			// Network Segments
			networkSegmentInterfaces := make([]interface{}, len(resp.Scope.Exclusions.NetworkSegments))
			for i, segment := range resp.Scope.Exclusions.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			exclusionsData["network_segments"] = networkSegmentInterfaces

			// JSS Users
			jssUserInterfaces := make([]interface{}, len(resp.Scope.Exclusions.JSSUsers))
			for i, user := range resp.Scope.Exclusions.JSSUsers {
				jssUserInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["jss_users"] = jssUserInterfaces

			// JSS User Groups
			jssUserGroupInterfaces := make([]interface{}, len(resp.Scope.Exclusions.JSSUserGroups))
			for i, group := range resp.Scope.Exclusions.JSSUserGroups {
				jssUserGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["jss_user_groups"] = jssUserGroupInterfaces

			// IBeacons
			iBeaconInterfaces := make([]interface{}, len(resp.Scope.Exclusions.IBeacons))
			for i, beacon := range resp.Scope.Exclusions.IBeacons {
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
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"reinstall_button_text":           resp.SelfService.ReinstallButtonText,
		"self_service_description":        resp.SelfService.SelfServiceDescription,
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"self_service_icon": []interface{}{map[string]interface{}{
			"id":       resp.SelfService.SelfServiceIcon.ID,
			"filename": resp.SelfService.SelfServiceIcon.Filename,
			"uri":      resp.SelfService.SelfServiceIcon.URI,
		}},
		"feature_on_main_page": resp.SelfService.FeatureOnMainPage,
		"self_service_categories": func() []interface{} {
			categories := make([]interface{}, len(resp.SelfService.SelfServiceCategories))
			for i, cat := range resp.SelfService.SelfServiceCategories {
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

	// Fetch the package configuration from the resp and set it in Terraform state
	packageConfigurations := make([]interface{}, 0)
	for _, packageItem := range resp.PackageConfiguration.Packages {
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

	// Fetch the scripts from the resp and set them in Terraform state
	scriptConfigurations := make([]interface{}, 0)
	for _, scriptItem := range resp.Scripts.Script {
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

	// Fetch the printers from the resp and set them in Terraform state
	printerConfigurations := make([]interface{}, 0)
	for _, printerItem := range resp.Printers.Printer {
		printer := make(map[string]interface{})
		printer["id"] = printerItem.ID
		printer["name"] = printerItem.Name
		printer["action"] = printerItem.Action
		printer["make_default"] = printerItem.MakeDefault
		printerConfigurations = append(printerConfigurations, printer)
	}

	if err := d.Set("printers", []interface{}{
		map[string]interface{}{
			"leave_existing_default": resp.Printers.LeaveExistingDefault,
			"printer":                printerConfigurations,
		},
	}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the dock items from the resp and set them in Terraform state
	dockItemConfigurations := make([]interface{}, 0)
	for _, dockItem := range resp.DockItems.DockItem {
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

	// Fetch the account maintenance data from the resp and set it in Terraform state
	accountMaintenanceState := make(map[string]interface{})
	accountMaintenanceState["accounts"] = []interface{}{}

	// Add account data if present
	if len(resp.AccountMaintenance.Accounts) > 0 {
		accountsState := make([]interface{}, len(resp.AccountMaintenance.Accounts))
		for i, account := range resp.AccountMaintenance.Accounts {
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
	if len(resp.AccountMaintenance.DirectoryBindings) > 0 {
		bindingsState := make([]interface{}, len(resp.AccountMaintenance.DirectoryBindings))
		for i, binding := range resp.AccountMaintenance.DirectoryBindings {
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
			"action":                  resp.AccountMaintenance.ManagementAccount.Action,
			"managed_password":        resp.AccountMaintenance.ManagementAccount.ManagedPassword,
			"managed_password_length": resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength,
		},
	}

	// Add open firmware/EFI password data
	accountMaintenanceState["open_firmware_efi_password"] = []interface{}{
		map[string]interface{}{
			"of_mode":     resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode,
			"of_password": resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword,
		},
	}

	// Set the account_maintenance in state
	if err := d.Set("account_maintenance", []interface{}{accountMaintenanceState}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the reboot data from the resp and set it in Terraform state
	rebootConfig := make(map[string]interface{})
	rebootConfig["message"] = resp.Reboot.Message
	rebootConfig["specify_startup"] = resp.Reboot.SpecifyStartup
	rebootConfig["startup_disk"] = resp.Reboot.StartupDisk
	rebootConfig["no_user_logged_in"] = resp.Reboot.NoUserLoggedIn
	rebootConfig["user_logged_in"] = resp.Reboot.UserLoggedIn
	rebootConfig["minutes_until_reboot"] = resp.Reboot.MinutesUntilReboot
	rebootConfig["start_reboot_timer_immediately"] = resp.Reboot.StartRebootTimerImmediately
	rebootConfig["file_vault_2_reboot"] = resp.Reboot.FileVault2Reboot

	if err := d.Set("reboot", []interface{}{rebootConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the maintenance data from the resp and set it in Terraform state
	maintenanceConfig := make(map[string]interface{})
	maintenanceConfig["recon"] = resp.Maintenance.Recon
	maintenanceConfig["reset_name"] = resp.Maintenance.ResetName
	maintenanceConfig["install_all_cached_packages"] = resp.Maintenance.InstallAllCachedPackages
	maintenanceConfig["heal"] = resp.Maintenance.Heal
	maintenanceConfig["prebindings"] = resp.Maintenance.Prebindings
	maintenanceConfig["permissions"] = resp.Maintenance.Permissions
	maintenanceConfig["byhost"] = resp.Maintenance.Byhost
	maintenanceConfig["system_cache"] = resp.Maintenance.SystemCache
	maintenanceConfig["user_cache"] = resp.Maintenance.UserCache
	maintenanceConfig["verify"] = resp.Maintenance.Verify

	if err := d.Set("maintenance", []interface{}{maintenanceConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the files and processes data from the resp and set it in Terraform state
	filesProcessesConfig := make(map[string]interface{})
	filesProcessesConfig["search_by_path"] = resp.FilesProcesses.SearchByPath
	filesProcessesConfig["delete_file"] = resp.FilesProcesses.DeleteFile
	filesProcessesConfig["locate_file"] = resp.FilesProcesses.LocateFile
	filesProcessesConfig["update_locate_database"] = resp.FilesProcesses.UpdateLocateDatabase
	filesProcessesConfig["spotlight_search"] = resp.FilesProcesses.SpotlightSearch
	filesProcessesConfig["search_for_process"] = resp.FilesProcesses.SearchForProcess
	filesProcessesConfig["kill_process"] = resp.FilesProcesses.KillProcess
	filesProcessesConfig["run_command"] = resp.FilesProcesses.RunCommand

	if err := d.Set("files_processes", []interface{}{filesProcessesConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the user interaction data from the resp and set it in Terraform state
	userInteractionConfig := make(map[string]interface{})
	userInteractionConfig["message_start"] = resp.UserInteraction.MessageStart
	userInteractionConfig["allow_user_to_defer"] = resp.UserInteraction.AllowUserToDefer
	userInteractionConfig["allow_deferral_until_utc"] = resp.UserInteraction.AllowDeferralUntilUtc
	userInteractionConfig["allow_deferral_minutes"] = resp.UserInteraction.AllowDeferralMinutes
	userInteractionConfig["message_finish"] = resp.UserInteraction.MessageFinish

	if err := d.Set("user_interaction", []interface{}{userInteractionConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the disk encryption data from the policy and set it in Terraform state
	diskEncryptionConfig := make(map[string]interface{})
	diskEncryptionConfig["action"] = resp.DiskEncryption.Action
	diskEncryptionConfig["disk_encryption_configuration_id"] = resp.DiskEncryption.DiskEncryptionConfigurationID
	diskEncryptionConfig["auth_restart"] = resp.DiskEncryption.AuthRestart
	diskEncryptionConfig["remediate_key_type"] = resp.DiskEncryption.RemediateKeyType
	diskEncryptionConfig["remediate_disk_encryption_configuration_id"] = resp.DiskEncryption.RemediateDiskEncryptionConfigurationID

	if err := d.Set("disk_encryption", []interface{}{diskEncryptionConfig}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// Constructs, updates and reads
func ResourceJamfProPoliciesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceIDInt, err := strconv.Atoi(resourceID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdatePolicyByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Policy '%s' (ID: %d) after retries: %v", resource.General.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// Deletes and removes from state
func ResourceJamfProPoliciesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

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

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro policy '%s' (ID: %s) after retries: %v", resourceName, d.Id(), err))
	}

	d.SetId("")

	return diags
}
