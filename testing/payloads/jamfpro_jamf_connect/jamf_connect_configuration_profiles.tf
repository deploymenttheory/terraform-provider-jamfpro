// ========================================================================== //
// Jamf Connect Configuration Profiles
// ========================================================================== //

// ========================================================================== //
// Data source by configuration profile ID

data "jamfpro_jamf_connect" "by_id" {
  profile_id = 429
}

// ========================================================================== //
// Data source by configuration profile name

data "jamfpro_jamf_connect" "by_name" {
  profile_name = "cp_jamfconnect_license"
}

