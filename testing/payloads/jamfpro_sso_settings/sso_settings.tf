resource "jamfpro_sso_settings" "entra_id_example" {
  sso_enabled                                          = true
  configuration_type                                   = "OIDC"
  sso_bypass_allowed                                   = true
  sso_for_enrollment_enabled                           = true
  sso_for_macos_self_service_enabled                   = true
  enrollment_sso_for_account_driven_enrollment_enabled = false
  group_enrollment_access_enabled                      = false
  group_enrollment_access_name                         = ""

  oidc_settings {
    user_mapping = "EMAIL"
  }

  enrollment_sso_config {
    hosts = []
  }
}
