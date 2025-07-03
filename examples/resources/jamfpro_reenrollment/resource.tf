resource "jamfpro_reenrollment" "settings" {
  flush_location_information         = true
  flush_location_information_history = true
  flush_policy_history               = true
  flush_extension_attributes         = true
  flush_software_update_plans        = true
  flush_mdm_queue                    = "DELETE_NOTHING" // required, valide values: DELETE_NOTHING, DELETE_ERRORS, DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED, DELETE_EVERYTHING
}
