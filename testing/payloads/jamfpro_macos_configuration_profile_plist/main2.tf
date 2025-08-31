resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_001" {
  name                = "test-profile-2"
  description         = "An example mobile device configuration profile."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/PPPC - Allow Microsoft Teams.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }

}
