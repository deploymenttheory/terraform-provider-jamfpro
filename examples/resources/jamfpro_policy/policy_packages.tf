resource "jamfpro_policy" "jamfpro_package_policy_001" {
  name                          = "tf-localtest-policy-packages-001"
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
    self_service_display_name       = "some display name"
    install_button_text             = "Install"
    reinstall_button_text           = "Reinstall"
    self_service_description        = "a description by tf"
    force_users_to_view_description = false

    feature_on_main_page = false
  }

  payloads {
    packages {
      distribution_point = "default" // Set the appropriate distribution point
      package {
        id                          = jamfpro_package.jamfpro_package_003.id
        action                      = "Install" // The action to perform with the package (e.g., Install, Cache, etc.)
        fill_user_template          = false     // Whether to fill the user template
        fill_existing_user_template = false     // Whether to fill existing user templates
      }
    }
  }
}






