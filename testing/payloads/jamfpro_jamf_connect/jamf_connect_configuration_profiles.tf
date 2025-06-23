// ========================================================================== //
// Jamf Connect Configuration Profiles
// ========================================================================== //

// ========================================================================== //
// Create configuration profile for testing

resource "jamfpro_macos_configuration_profile_plist" "jamf_connect_license_001" {
  name                = "Jamf Connect License"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("./jamf_connect_license.mobileconfig")
  payload_validate    = true
  user_removable      = false
  scope {
    all_computers = true
    all_jss_users = false
  }
}

// ========================================================================== //
// Data source by configuration profile ID

data "jamfpro_jamf_connect" "by_id" {
  depends_on = [
    jamfpro_macos_configuration_profile_plist.jamf_connect_license_001
  ]
  profile_id = jamfpro_macos_configuration_profile_plist.jamf_connect_license_001.id
}

// ========================================================================== //
// Data source by configuration profile name

data "jamfpro_jamf_connect" "by_name" {
  depends_on = [
    jamfpro_macos_configuration_profile_plist.jamf_connect_license_001
  ]
  profile_name = "Jamf Connect License"
}
