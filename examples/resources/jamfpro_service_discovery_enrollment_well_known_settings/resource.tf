resource "jamfpro_service_discovery_enrollment_well_known_settings" "example" {
  well_known_settings {
    server_uuid     = "1234567890ABCDEF1234567890ABCDEF"
    enrollment_type = "mdm-adde"
  }
}
