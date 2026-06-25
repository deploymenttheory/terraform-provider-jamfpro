resource "jamfpro_policy" "recon" {
  name                          = "acc-test-policy-onboarding-recon"
  enabled                       = true
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "USER_INITIATED"
  frequency                     = "Ongoing"
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
    all_computers = true
    all_jss_users = false
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = "Do the thing"
    install_button_text             = "Install"
    reinstall_button_text           = "Reinstall"
    self_service_description        = ""
    force_users_to_view_description = false
    feature_on_main_page            = false

    notification         = false
    notification_type    = "Self Service"
    notification_subject = "Recon"
    notification_message = "This is a message for the Recon policy"
  }

  payloads {
    maintenance {
      recon = true
    }
  }
}

resource "jamfpro_macos_onboarding_settings" "basic_example" {
  enabled = true

  onboarding_items {
    entity_id                = jamfpro_policy.recon.id
    self_service_entity_type = "OS_X_POLICY"
    priority                 = 1
  }
}
