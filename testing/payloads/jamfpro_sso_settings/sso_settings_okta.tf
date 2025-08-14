resource "jamfpro_sso_settings" "okta" {
  sso_enabled                                          = true
  configuration_type                                   = "OIDC_WITH_SAML"
  sso_for_enrollment_enabled                           = true
  sso_for_macos_self_service_enabled                   = false
  enrollment_sso_for_account_driven_enrollment_enabled = false
  group_enrollment_access_enabled                      = false
  group_enrollment_access_name                         = ""
  sso_bypass_allowed                                   = true

  oidc_settings {
    user_mapping = "EMAIL"
  }

  saml_settings {
    token_expiration_disabled = true
    user_attribute_enabled    = false
    user_attribute_name       = ""
    user_mapping              = "EMAIL"
    group_attribute_name      = "http://schemas.xmlsoap.org/claims/Group"
    group_rdn_key             = ""
    idp_provider_type         = "OKTA"
    idp_url                   = "https://trial-9344750.okta.com/app/exks2co0x41zYbk8y697/sso/saml/metadata"
    entity_id                 = "/saml/metadata"
    metadata_file_name        = ""
    other_provider_type_name  = ""
    federation_metadata_file  = ""
    metadata_source           = "URL"
    session_timeout           = 480
  }

  enrollment_sso_config {
    hosts = []
  }
}
