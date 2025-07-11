// ========================================================================== //
// Re-enrollment settings
// ========================================================================== //

// ========================================================================== //
// Configure re-enrollment settings

resource "jamfpro_reenrollment" "settings" {
  flush_location_information         = true
  flush_location_information_history = true
  flush_policy_history               = true
  flush_extension_attributes         = true
  flush_software_update_plans        = true
  flush_mdm_queue                    = "DELETE_EVERYTHING"
}
