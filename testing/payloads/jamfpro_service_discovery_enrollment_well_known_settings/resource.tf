resource "jamfpro_service_discovery_enrollment_well_known_settings" "test" {
  well_known_settings {
    server_uuid     = "ABCDEF1234567890ABCDEF1234567890"
    enrollment_type = "none"
  }
}
