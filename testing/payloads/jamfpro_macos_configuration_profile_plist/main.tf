resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_064" {
  name = "test-profile"
  // Regression test for issue #1145 - heredoc strings in HCL always include
  // a trailing newline before EOT, but the API strips it server-side.
  description         = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/Screen Recording - Allow Microsoft Teams.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }

}

# Regression coverage for the phantom empty <array/> bug on macOS: a
# pretty-printed com.apple.systempreferences payload with two sibling
# arrays-of-strings (whitespace between the <string> children is the same class
# the Classic API mis-parses). Compaction must leave the panes lists intact.
resource "jamfpro_macos_configuration_profile_plist" "system_preferences_pretty" {
  name                = "test-system-preferences-pretty"
  description         = "Pretty-printed System Preferences profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/SystemPreferences-Pretty.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}
