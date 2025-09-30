resource "jamfpro_macos_onboarding_settings" "example" {
  enabled = true

  onboarding_items {
    entity_id                = "1"
    self_service_entity_type = "OS_X_POLICY"
    priority                 = 1
  }

  onboarding_items {
    entity_id                = "5"
    self_service_entity_type = "OS_X_MAC_APP"
    priority                 = 2
  }

  onboarding_items {
    entity_id                = "3"
    self_service_entity_type = "OS_X_CONFIGURATION_PROFILE"
    priority                 = 3
  }
}
