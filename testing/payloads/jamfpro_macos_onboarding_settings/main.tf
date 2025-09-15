resource "jamfpro_policy" "recon" {
  name      = "Do A Recon for Onboarding and Profit"
  enabled   = true
  frequency = "Ongoing"

  scope {
    all_computers = true
    all_jss_users = false

  }
  self_service {
    use_for_self_service            = true
    self_service_display_name       = "Do the thing"
    install_button_text             = "Install"
    reinstall_button_text           = "Reinstall"
    force_users_to_view_description = false
    feature_on_main_page            = false
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
