// This test validates a policy with maximal self service configuration including icon and files_processes payload
resource "jamfpro_policy" "jamfpro_policy_self_service_maximal" {
  name                          = "acc-test-policy-self-service-maximal"
  enabled                       = true
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "USER_INITIATED" // Self service trigger
  frequency                     = "Ongoing"
  retry_event                   = "none"
  retry_attempts                = -1
  notify_on_each_failed_retry   = false
  target_drive                  = "/"
  offline                       = false
  category_id                   = jamfpro_category.jamfpro_category_self_service.id
  site_id                       = -1

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers      = false
    all_jss_users      = false
    computer_group_ids = [jamfpro_static_computer_group.jamfpro_static_computer_group_self_service_test.id]
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = "Microsoft Excel"
    install_button_text             = "Install"
    reinstall_button_text           = "Install"
    self_service_description        = "Microsoft Excel is a spreadsheet application featuring calculation, graphing tools, PivotTables, and a macro programming language called Visual Basic for Applications."
    force_users_to_view_description = false
    feature_on_main_page            = true
    self_service_icon_id            = jamfpro_icon.self_service_icon.id

    self_service_category {
      id         = jamfpro_category.jamfpro_category_self_service.id
      display_in = true
      feature_in = true
    }

    notification         = true
    notification_type    = "Self Service"
    notification_subject = "Microsoft Excel Opened"
    notification_message = "Microsoft Excel has been successfully opened on your computer."
  }

  payloads {
    files_processes {
      search_by_path         = ""
      delete_file            = false
      locate_file            = ""
      update_locate_database = false
      spotlight_search       = ""
      search_for_process     = ""
      kill_process           = false
      run_command            = "open /Applications/Microsoft\\ Excel.app"
    }

    reboot {
      message                        = "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu."
      specify_startup                = "Immediately"
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
  }
}

// Icon resource for self service with file path
resource "jamfpro_icon" "self_service_icon" {
  icon_file_base64 = filebase64("${path.module}/support_files/icons/Microsoft Excel.png")
}

// Supporting category resource
resource "jamfpro_category" "jamfpro_category_self_service" {
  name     = "acc-test-category-self-service"
  priority = 5
}

// Supporting computer group for scope
resource "jamfpro_static_computer_group" "jamfpro_static_computer_group_self_service_test" {
  name = "acc-test-static-computer-group-self-service"
}
