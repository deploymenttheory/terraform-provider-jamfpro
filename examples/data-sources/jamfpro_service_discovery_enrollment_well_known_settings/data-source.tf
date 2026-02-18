data "jamfpro_service_discovery_enrollment_well_known_settings" "all" {
  id = "service_discovery_enrollment_well_known_settings"
}

output "service_discovery_enrollment_well_known_settings" {
  value = data.jamfpro_service_discovery_enrollment_well_known_settings.all
}
