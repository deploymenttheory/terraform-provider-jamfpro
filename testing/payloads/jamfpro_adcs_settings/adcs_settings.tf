// ========================================================================== //
// ADCS Settings - Outbound Example
// ========================================================================== //

resource "jamfpro_adcs_settings" "outbound" {
  display_name  = "tf-testing-adcs-settings-${var.testing_id}"
  ca_name       = "Contoso Issuing CA"
  fqdn          = "connector.contoso.corp"
  api_client_id = jamfpro_api_integration.adcs_settings_integration.client_id

  revocation_enabled = true
  outbound           = true
}

// ========================================================================== //
// API Role and Integration
// ========================================================================== //

resource "jamfpro_api_role" "adcs_settings_role" {
  display_name = "tf-testing-adcs-settings-role-${var.testing_id}"

  privileges = [
    "Update AD CS Certificate Jobs",
    "Read AD CS Certificate Jobs",
  ]
}

resource "jamfpro_api_integration" "adcs_settings_integration" {
  display_name                  = "tf-testing-adcs-settings-integration-${var.testing_id}"
  enabled                       = true
  authorization_scopes          = [jamfpro_api_role.adcs_settings_role.display_name]
  access_token_lifetime_seconds = 60
}
