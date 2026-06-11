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

# The remaining pretty-printed payloads exercise structural-whitespace
# compaction across a variety of payload types: dock (deeply nested tile-data
# dicts inside an array), login window (entity references and multi-byte UTF-8
# in leaf strings), and screensaver (scalar booleans/integers).
resource "jamfpro_macos_configuration_profile_plist" "dock_pretty" {
  name                = "test-dock-pretty"
  description         = "Pretty-printed Dock profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/Dock-Pretty.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

resource "jamfpro_macos_configuration_profile_plist" "login_window_pretty" {
  name                = "test-login-window-pretty"
  description         = "Pretty-printed Login Window profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/LoginWindow-Pretty.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

resource "jamfpro_macos_configuration_profile_plist" "screensaver_pretty" {
  name                = "test-screensaver-pretty"
  description         = "Pretty-printed Screensaver profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/Screensaver-Pretty.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}
