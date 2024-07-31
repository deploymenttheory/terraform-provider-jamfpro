resource "jamfpro_policy" "jamfpro_reboot_policy_001" {
  name                          = "tf-localtest-reboot_policy-001"
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

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = ""
    install_button_text             = "Install"
    self_service_description        = ""
    force_users_to_view_description = false

    feature_on_main_page = false
  }

  payloads {
    reboot {
      message                        = "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu."
      specify_startup                = "MDM Restart with Kernel Cache Rebuild" // Standard Restart | "MDM Restart with Kernel Cache Rebuild"
      startup_disk                   = "Current Startup Disk"
      no_user_logged_in              = "Do not restart"
      user_logged_in                 = "Do not restart"
      minutes_until_reboot           = 10
      start_reboot_timer_immediately = false
      file_vault_2_reboot            = false
    }
  }
}






