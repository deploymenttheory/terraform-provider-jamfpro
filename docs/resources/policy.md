---
page_title: "jamfpro_policy"
description: |-
  
---

# jamfpro_policy (Resource)


## Example Usage
```terraform
resource "jamfpro_policy" "jamfpro_policy_001" {
  name                          = "tf-localtest-policy-template-001"
  enabled                       = false
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "EVENT" // "USER_INITIATED" for self service trigger , "EVENT" for an event trigger
  frequency                     = "Once per computer"
  retry_event                   = "none"
  retry_attempts                = -1
  notify_on_each_failed_retry   = false
  target_drive                  = "/"
  offline                       = false
  category_id                   = -1
  site_id                       = -1

  date_time_limitations {
    activation_date       = "2026-12-25 01:00:00"
    activation_date_epoch = 1798160400000
    activation_date_utc   = "2026-12-25T01:00:00.000+0000"
    expiration_date       = "2028-04-01 16:02:00"
    expiration_date_epoch = 1838217720000
    expiration_date_utc   = "2028-04-01T16:02:00.000+0000"
    no_execute_start      = "1:00 AM"
    no_execute_end        = "1:03 PM"
  }

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      //directory_service_usergroup_ids = [3, 4]
    }

    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([118, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
      ibeacon_ids                          = [3, 4]
    }
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = ""
    install_button_text             = "Install"
    reinstall_button_text           = "Reinstall"
    self_service_description        = ""
    force_users_to_view_description = false
    feature_on_main_page = false
  }

  payloads {
    packages {
      distribution_point = "default" // Set the appropriate distribution point
      package {
        id                          = 123       // The ID of the package in Jamf Pro
        action                      = "Install" // The action to perform with the package (e.g., Install, Cache, etc.)
        fill_user_template          = false     // Whether to fill the user template
        fill_existing_user_template = false     // Whether to fill existing user templates
      }
    }
    scripts {
      id          = 123
      priority    = "After"
      parameter4  = "param_value_4"
      parameter5  = "param_value_5"
      parameter6  = "param_value_6"
      parameter7  = "param_value_7"
      parameter8  = "param_value_8"
      parameter9  = "param_value_9"
      parameter10 = "param_value_10"
      parameter11 = "param_value_11"
    }

    disk_encryption {
      action                                     = "apply"
      disk_encryption_configuration_id           = 1
      auth_restart                               = false
      remediate_key_type                         = "Individual"
      remediate_disk_encryption_configuration_id = 2
    }

    printers {
      id           = 1
      name         = "Printer1"
      action       = "install"
      make_default = true
    }

    dock_items {
      id     = 1
      name   = "Safari"
      action = "Add To End"
    }

    account_maintenance {
      local_accounts {
        account {
          action                    = "Create"
          username                  = "newuser"
          realname                  = "New User"
          password                  = "password123"
          archive_home_directory    = false
          archive_home_directory_to = ""
          home                      = "/Users/newuser"
          hint                      = "This is a hint"
          picture                   = "/Library/User Pictures/Animals/Butterfly.tif"
          admin                     = true
          filevault_enabled         = true
        }
      }
      directory_bindings {
        binding {
          id = 1
        }
      }

      management_account {
        action                  = "rotate"
        managed_password        = "newmanagedpassword"
        managed_password_length = 15
      }
      open_firmware_efi_password {
        of_mode     = "command"
        of_password = "firmwarepassword"
      }
    }
    reboot {
      message                        = "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu."
      specify_startup                = "Standard Restart"
      startup_disk                   = "Current Startup Disk"
      no_user_logged_in              = "Do not restart"
      user_logged_in                 = "Do not restart"
      minutes_until_reboot           = 5
      start_reboot_timer_immediately = false
      file_vault_2_reboot            = false
    }
    maintenance {
      recon                       = true
      reset_name                  = false
      install_all_cached_packages = false
      heal                        = false
      prebindings                 = false
      permissions                 = false
      byhost                      = false
      system_cache                = false
      user_cache                  = false
      verify                      = false
    }
    files_processes {
      search_by_path         = "/Applications/SomeApp.app"
      delete_file            = true
      locate_file            = "SomeFile.txt"
      update_locate_database = false
      spotlight_search       = "SomeApp"
      search_for_process     = "SomeProcess"
      kill_process           = true
      run_command            = "echo 'Hello, World!'"
    }
    user_interaction {
      message_start            = "Policy is about to run."
      allow_users_to_defer     = true
      allow_deferral_until_utc = "2024-12-31T23:59:59Z"
      allow_deferral_minutes   = 1440
      message_finish           = "Policy has completed."
    }
    disk_encryption {
      action                                     = "apply"
      disk_encryption_configuration_id           = 1
      auth_restart                               = false
      remediate_key_type                         = "Individual"
      remediate_disk_encryption_configuration_id = 2
    }

  }

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Define whether the policy is enabled.
- `name` (String) The name of the policy.
- `payloads` (Block List, Min: 1) All payloads container (see [below for nested schema](#nestedblock--payloads))
- `scope` (Block List, Min: 1, Max: 1) Scope configuration for the profile. (see [below for nested schema](#nestedblock--scope))

### Optional

- `category_id` (Number) Jamf Pro category-related settings of the policy.
- `date_time_limitations` (Block List, Max: 1) Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time. (see [below for nested schema](#nestedblock--date_time_limitations))
- `frequency` (String) Frequency of policy execution.
- `network_limitations` (Block List, Max: 1) Network limitations for the policy. (see [below for nested schema](#nestedblock--network_limitations))
- `network_requirements` (String) Network requirements for the policy.
- `notify_on_each_failed_retry` (Boolean) Send notifications for each failed policy retry attempt.
- `offline` (Boolean) Make policy available offline by caching the policy to the macOS device to ensure it runs when Jamf Pro is unavailable. Only used when execution policy is set to 'ongoing'.
- `package_distribution_point` (String) repository of which packages are collected from
- `retry_attempts` (Number) Number of retry attempts for the jamf pro policy. Valid values are -1 (not configured) and 1 through 10.
- `retry_event` (String) Event on which to retry policy execution.
- `self_service` (Block List, Max: 1) Self-service settings of the policy. (see [below for nested schema](#nestedblock--self_service))
- `site_id` (Number) Jamf Pro Site-related settings of the policy.
- `target_drive` (String) The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `trigger_checkin` (Boolean) Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro
- `trigger_enrollment_complete` (Boolean) Trigger policy when device enrollment is complete.
- `trigger_login` (Boolean) Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work
- `trigger_network_state_changed` (Boolean) Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)
- `trigger_other` (String) Any other trigger for the policy.
- `trigger_startup` (Boolean) Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work

### Read-Only

- `id` (String) The unique identifier of the Jamf Pro policy.

<a id="nestedblock--payloads"></a>
### Nested Schema for `payloads`

Optional:

- `account_maintenance` (Block List) Account maintenance settings of the policy. Use this section to create and delete local accounts, and to reset local account passwords. Also use this section to disable an existing local account for FileVault 2. (see [below for nested schema](#nestedblock--payloads--account_maintenance))
- `disk_encryption` (Block List) Disk encryption settings of the policy. Use this section to enable FileVault 2 or to issue a new recovery key. (see [below for nested schema](#nestedblock--payloads--disk_encryption))
- `dock_items` (Block List) Dock items settings of the policy. (see [below for nested schema](#nestedblock--payloads--dock_items))
- `files_processes` (Block List) Files and processes settings of the policy. Use this section to search for and log specific files and processes. Also use this section to execute a command. (see [below for nested schema](#nestedblock--payloads--files_processes))
- `maintenance` (Block List) Maintenance settings of the policy. Use this section to update inventory, reset computer names, install all cached packages, and run common maintenance tasks. (see [below for nested schema](#nestedblock--payloads--maintenance))
- `override_default_settings` (Block List) Settings to override default configurations. (see [below for nested schema](#nestedblock--payloads--override_default_settings))
- `packages` (Block List) Package configuration settings of the policy. (see [below for nested schema](#nestedblock--payloads--packages))
- `printers` (Block List) Printers settings of the policy. (see [below for nested schema](#nestedblock--payloads--printers))
- `reboot` (Block List) Use this section to restart computers and specify the disk to boot them to (see [below for nested schema](#nestedblock--payloads--reboot))
- `scripts` (Block List) Scripts settings of the policy. (see [below for nested schema](#nestedblock--payloads--scripts))
- `user_interaction` (Block List) User interaction settings of the policy. (see [below for nested schema](#nestedblock--payloads--user_interaction))

Read-Only:

- `network_requirements` (String) Network requirements for the policy.

<a id="nestedblock--payloads--account_maintenance"></a>
### Nested Schema for `payloads.account_maintenance`

Optional:

- `directory_bindings` (Block List) Directory binding settings for the policy. Use this section to bind computers to a directory service (see [below for nested schema](#nestedblock--payloads--account_maintenance--directory_bindings))
- `local_accounts` (Block List, Max: 1) Local user account configurations (see [below for nested schema](#nestedblock--payloads--account_maintenance--local_accounts))
- `management_account` (Block List) Management account settings for the policy. Use this section to change or reset the management account password. (see [below for nested schema](#nestedblock--payloads--account_maintenance--management_account))
- `open_firmware_efi_password` (Block List) Open Firmware/EFI password settings for the policy. Use this section to set or remove an Open Firmware/EFI password on computers with Intel-based processors. (see [below for nested schema](#nestedblock--payloads--account_maintenance--open_firmware_efi_password))

<a id="nestedblock--payloads--account_maintenance--directory_bindings"></a>
### Nested Schema for `payloads.account_maintenance.directory_bindings`

Optional:

- `binding` (Block List) Details of the directory binding. (see [below for nested schema](#nestedblock--payloads--account_maintenance--directory_bindings--binding))

<a id="nestedblock--payloads--account_maintenance--directory_bindings--binding"></a>
### Nested Schema for `payloads.account_maintenance.directory_bindings.binding`

Optional:

- `id` (Number) The unique identifier of the binding.

Read-Only:

- `name` (String) The name of the binding.



<a id="nestedblock--payloads--account_maintenance--local_accounts"></a>
### Nested Schema for `payloads.account_maintenance.local_accounts`

Optional:

- `account` (Block List) Details of each account configuration. (see [below for nested schema](#nestedblock--payloads--account_maintenance--local_accounts--account))

<a id="nestedblock--payloads--account_maintenance--local_accounts--account"></a>
### Nested Schema for `payloads.account_maintenance.local_accounts.account`

Optional:

- `action` (String) Action to be performed on the account (e.g., Create, Reset, Delete, DisableFileVault).
- `admin` (Boolean) Whether the account has admin privileges.Setting this to true will set the user administrator privileges to the computer
- `archive_home_directory` (Boolean) Permanently delete home directory. If set to true will archive the home directory.
- `archive_home_directory_to` (String) Path in which to archive the home directory to.
- `filevault_enabled` (Boolean) Allow the user to unlock the FileVault 2-encrypted drive
- `hint` (String) Hint to help the user remember the password
- `home` (String) Full path in which to create the home directory (e.g. /Users/username/ or /private/var/username/)
- `password` (String) Set a new account password. This does not update the account's login keychain password or FileVault 2 password.
- `picture` (String) Full path to the account picture (e.g. /Library/User Pictures/Animals/Butterfly.tif )
- `realname` (String) Real name associated with the account.
- `username` (String) Username/short name for the account



<a id="nestedblock--payloads--account_maintenance--management_account"></a>
### Nested Schema for `payloads.account_maintenance.management_account`

Optional:

- `action` (String) Action to perform on the management account.Rotates management account password at next policy execution. Valid values are 'rotate' or 'doNotChange'.
- `managed_password` (String) Managed password for the account. Management account passwords will be automatically randomized with 29 characters by jamf pro.
- `managed_password_length` (Number) Length of the managed password. Only necessary when utilizing the random action


<a id="nestedblock--payloads--account_maintenance--open_firmware_efi_password"></a>
### Nested Schema for `payloads.account_maintenance.open_firmware_efi_password`

Optional:

- `of_mode` (String) Mode for the open firmware/EFI password. Valid values are 'command' or 'none'.
- `of_password` (String) Password for the open firmware/EFI.



<a id="nestedblock--payloads--disk_encryption"></a>
### Nested Schema for `payloads.disk_encryption`

Optional:

- `action` (String) The action to perform for disk encryption (e.g., apply, remediate).
- `auth_restart` (Boolean) Whether to allow authentication restart.
- `disk_encryption_configuration_id` (Number) ID of the disk encryption configuration to apply.
- `remediate_disk_encryption_configuration_id` (Number) Disk encryption ID to utilize for remediating institutional recovery key types.
- `remediate_key_type` (String) Type of key to use for remediation (e.g., Individual, Institutional, Individual And Institutional).


<a id="nestedblock--payloads--dock_items"></a>
### Nested Schema for `payloads.dock_items`

Required:

- `action` (String) Action to be performed for the dock item (e.g., Add To Beginning, Add To End, Remove).
- `id` (Number) Unique identifier of the dock item.
- `name` (String) Name of the dock item.


<a id="nestedblock--payloads--files_processes"></a>
### Nested Schema for `payloads.files_processes`

Optional:

- `delete_file` (Boolean) Whether to delete the file found at the specified path.
- `kill_process` (Boolean) Whether to kill the process if found. This works with exact matches only
- `locate_file` (String) Path of the file to locate. Name of the file, including the file extension. This field is case-sensitive and returns partial matches
- `run_command` (String) Command to execute on computers. This command is executed as the 'root' user
- `search_by_path` (String) Path of the file to search for.
- `search_for_process` (String) Name of the process to search for. This field is case-sensitive and returns partial matches
- `spotlight_search` (String) Search For File Using Spotlight. File to search for. This field is not case-sensitive and returns partial matches
- `update_locate_database` (Boolean) Whether to update the locate database. Update the locate database before searching for the file


<a id="nestedblock--payloads--maintenance"></a>
### Nested Schema for `payloads.maintenance`

Optional:

- `byhost` (Boolean) Whether to fix ByHost files andnpreferences.
- `heal` (Boolean) Whether to heal the policy.
- `install_all_cached_packages` (Boolean) Whether to install all cached packages. Installs packages cached by Jamf Pro
- `permissions` (Boolean) Whether to fix Disk Permissions (Not compatible with macOS v10.12 or later)
- `prebindings` (Boolean) Whether to update prebindings.
- `recon` (Boolean) Whether to run recon (inventory update) as part of the maintenance. Forces computers to submit updated inventory information to Jamf Pro
- `reset_name` (Boolean) Whether to reset the computer name to the name stored in Jamf Pro. Changes the computer name on computers to match the computer name in Jamf Pro
- `system_cache` (Boolean) Whether to flush caches from /Library/Caches/ and /System/Library/Caches/, except for any com.apple.LaunchServices caches
- `user_cache` (Boolean) Whether to flush caches from ~/Library/Caches/, ~/.jpi_cache/, and ~/Library/Preferences/Microsoft/Office version #/Office Font Cache. Enabling this may cause problems with system fonts displaying unless a restart option is configured.
- `verify` (Boolean) Whether to verify system files and structure on the Startup Disk


<a id="nestedblock--payloads--override_default_settings"></a>
### Nested Schema for `payloads.override_default_settings`

Optional:

- `any_ip_address` (Boolean) Whether the policy applies to any IP address.
- `minimum_network_connection` (String) Minimum network connection required for the policy.


<a id="nestedblock--payloads--packages"></a>
### Nested Schema for `payloads.packages`

Required:

- `distribution_point` (String) Distribution point for the package.
- `package` (Block List, Min: 1) List of packages. (see [below for nested schema](#nestedblock--payloads--packages--package))

<a id="nestedblock--payloads--packages--package"></a>
### Nested Schema for `payloads.packages.package`

Required:

- `id` (Number) Unique identifier of the package.

Optional:

- `action` (String) Action to be performed for the package.
- `fill_existing_user_template` (Boolean) Fill Existing Users (FEU).
- `fill_user_template` (Boolean) Fill User Template (FUT).



<a id="nestedblock--payloads--printers"></a>
### Nested Schema for `payloads.printers`

Required:

- `action` (String) Action to be performed for the printer (e.g., install, uninstall).
- `id` (Number) Unique identifier of the printer.
- `name` (String) Name of the printer.

Optional:

- `make_default` (Boolean) Whether to set the printer as the default.


<a id="nestedblock--payloads--reboot"></a>
### Nested Schema for `payloads.reboot`

Optional:

- `file_vault_2_reboot` (Boolean) Perform authenticated restart on computers with FileVault 2 enabled. Restart FileVault 2-encrypted computers without requiring an unlock during the next startup
- `message` (String) The reboot message displayed to the user.
- `minutes_until_reboot` (Number) Amount of time to wait before the restart begins.
- `no_user_logged_in` (String) Action to take if no user is logged in to the computer
- `specify_startup` (String) Reboot Method
- `start_reboot_timer_immediately` (Boolean) Defines if the reboot timer should start immediately once the policy applies to a macOS device.
- `startup_disk` (String) Disk to boot computers to
- `user_logged_in` (String) Action to take if a user is logged in to the computer


<a id="nestedblock--payloads--scripts"></a>
### Nested Schema for `payloads.scripts`

Optional:

- `id` (String) Unique identifier of the script.
- `parameter10` (String) Custom parameter 10 for the script.
- `parameter11` (String) Custom parameter 11 for the script.
- `parameter4` (String) Custom parameter 4 for the script.
- `parameter5` (String) Custom parameter 5 for the script.
- `parameter6` (String) Custom parameter 6 for the script.
- `parameter7` (String) Custom parameter 7 for the script.
- `parameter8` (String) Custom parameter 8 for the script.
- `parameter9` (String) Custom parameter 9 for the script.
- `priority` (String) Execution priority of the script.


<a id="nestedblock--payloads--user_interaction"></a>
### Nested Schema for `payloads.user_interaction`

Optional:

- `allow_deferral_minutes` (Number) Number of minutes after the user was first prompted by the policy at which the policy runs and deferrals are prohibited. Must be a multiple of 1440 (minutes in day)
- `allow_deferral_until_utc` (String) Date/time at which deferrals are prohibited and the policy runs. Uses time zone settings of your hosting server. Standard environments hosted in Jamf Cloud use Coordinated Universal Time (UTC)
- `allow_users_to_defer` (Boolean) Allow user deferral and configure deferral type. A deferral limit must be specified for this to work.
- `message_finish` (String) Message to display when the policy is complete.
- `message_start` (String) Message to display before the policy runs



<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Required:

- `all_computers` (Boolean) Whether the configuration profile is scoped to all computers.

Optional:

- `all_jss_users` (Boolean) Whether the configuration profile is scoped to all JSS users.
- `building_ids` (List of Number) The buildings to which the configuration profile is scoped by Jamf ID
- `computer_group_ids` (List of Number) The computer groups to which the configuration profile is scoped by Jamf ID
- `computer_ids` (List of Number) The computers to which the configuration profile is scoped by Jamf ID
- `department_ids` (List of Number) The departments to which the configuration profile is scoped by Jamf ID
- `exclusions` (Block List, Max: 1) The scope exclusions from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--exclusions))
- `jss_user_group_ids` (List of Number) The jss user groups to which the configuration profile is scoped by Jamf ID
- `jss_user_ids` (List of Number) The jss users to which the configuration profile is scoped by Jamf ID
- `limitations` (Block List, Max: 1) The scope limitations from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--limitations))

<a id="nestedblock--scope--exclusions"></a>
### Nested Schema for `scope.exclusions`

Optional:

- `building_ids` (List of Number) Buildings excluded from scope by Jamf ID.
- `computer_group_ids` (List of Number) Computer Groups excluded from scope by Jamf ID.
- `computer_ids` (List of Number) Computers excluded from scope by Jamf ID.
- `department_ids` (List of Number) Departments excluded from scope by Jamf ID.
- `directory_service_or_local_usernames` (List of String) A list of directory service / local usernames for scoping limitations.
- `directory_service_usergroup_ids` (List of Number) A list of directory service / local user group IDs for limitations.
- `ibeacon_ids` (List of Number) Ibeacons excluded from scope by Jamf ID.
- `jss_user_group_ids` (List of Number) JSS User Groups excluded from scope by Jamf ID.
- `jss_user_ids` (List of Number) JSS Users excluded from scope by Jamf ID.
- `network_segment_ids` (List of Number) Network segments excluded from scope by Jamf ID.


<a id="nestedblock--scope--limitations"></a>
### Nested Schema for `scope.limitations`

Optional:

- `directory_service_or_local_usernames` (List of String) A list of directory service / local usernames for scoping limitations.
- `directory_service_usergroup_ids` (List of Number) A list of directory service user group IDs for limitations.
- `ibeacon_ids` (List of Number) A list of iBeacon IDs for limitations.
- `network_segment_ids` (List of Number) A list of network segment IDs for limitations.



<a id="nestedblock--date_time_limitations"></a>
### Nested Schema for `date_time_limitations`

Optional:

- `activation_date` (String) The activation date of the policy in 'YYYY-MM-DD HH:mm:ss' format. This is when the policy becomes active and starts executing. Example: '2026-12-25 01:00:00'
- `activation_date_epoch` (Number) The epoch time (Unix timestamp) in milliseconds of the activation date. This represents the number of milliseconds since January 1, 1970, 00:00:00 UTC. Example: 1798160400000 (represents December 25, 2026, 01:00:00)
- `activation_date_utc` (String) The UTC time of the activation date in ISO 8601 format with timezone offset. Format: 'YYYY-MM-DDThh:mm:ss.sss+0000'. Example: '2026-12-25T01:00:00.000+0000'
- `expiration_date` (String) The expiration date of the policy in 'YYYY-MM-DD HH:mm:ss' format. After this date, the policy will no longer be active or execute. Example: '2028-04-01 16:02:00'
- `expiration_date_epoch` (Number) The epoch time (Unix timestamp) in milliseconds of the expiration date. This represents the number of milliseconds since January 1, 1970, 00:00:00 UTC. Example: 1838217720000 (represents April 1, 2028, 16:02:00)
- `expiration_date_utc` (String) The UTC time of the expiration date in ISO 8601 format with timezone offset. Format: 'YYYY-MM-DDThh:mm:ss.sss+0000'. Example: '2028-04-01T16:02:00.000+0000'
- `no_execute_end` (String) The daily end time when the policy should not execute, in '12-hour clock' format (h:mm AM/PM). This is part of client-side limitations enforced based on computer settings. Example: '1:03 PM'
- `no_execute_start` (String) The daily start time when the policy should not execute, in '12-hour clock' format (h:mm AM/PM). This is part of client-side limitations enforced based on computer settings. Example: '1:00 AM'


<a id="nestedblock--network_limitations"></a>
### Nested Schema for `network_limitations`

Optional:

- `any_ip_address` (Boolean) Whether the policy applies to any IP address.
- `minimum_network_connection` (String) Minimum network connection required for the policy.


<a id="nestedblock--self_service"></a>
### Nested Schema for `self_service`

Optional:

- `feature_on_main_page` (Boolean) Whether to feature the policy on the main page of self-service.
- `force_users_to_view_description` (Boolean) Whether to force users to view the policy description in self-service.
- `install_button_text` (String) Text displayed on the install button in self-service.
- `reinstall_button_text` (String) Text displayed on the re-install button in self-service.
- `self_service_category` (Block List) Category settings for the policy in self-service. (see [below for nested schema](#nestedblock--self_service--self_service_category))
- `self_service_description` (String) Description of the policy displayed in self-service.
- `self_service_display_name` (String) Display name of the policy in self-service.
- `self_service_icon_id` (Number) Icon for policy to use in self-service
- `use_for_self_service` (Boolean) Whether the policy is available for self-service.

<a id="nestedblock--self_service--self_service_category"></a>
### Nested Schema for `self_service.self_service_category`

Required:

- `display_in` (Boolean) Whether to display the category in self-service.
- `feature_in` (Boolean) Whether to feature the category in self-service.
- `id` (Number) Category ID for the policy in self-service.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)