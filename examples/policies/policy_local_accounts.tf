resource "jamfpro_policy" "jamfpro_local_accounts_policy_001" {
  name                          = "tf-localtest-local_accounts_policy-001"
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
    }
  }
}






