data "jamfpro_macos_configuration_profile" "jamfpro_macos_configuration_profile_001_data" {
  id = jamfpro_macos_configuration_profile.jamfpro_macos_configuration_profile_001.id
}

output "jamfpro_macos_configuration_profile_001_data_id" {
  value = data.jamfpro_macos_configuration_profile.jamfpro_macos_configuration_profile_001_data.id
}

output "jamfpro_macos_configuration_profile_001_data_name" {
  value = data.jamfpro_macos_configuration_profile.jamfpro_macos_configuration_profile_001_data.name
}