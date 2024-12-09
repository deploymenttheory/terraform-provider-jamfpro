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






