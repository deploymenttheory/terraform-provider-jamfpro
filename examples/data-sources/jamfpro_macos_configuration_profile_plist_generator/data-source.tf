data "jamfpro_macos_configuration_profile_plist_generator" "jamfpro_macos_configuration_profile_plist_generator_001_data" {
  id = jamfpro_macos_configuration_profile_plist_generator.jamfpro_macos_configuration_profile_plist_generator_001.id
}

output "jamfpro_macos_configuration_profile_plist_generator_001_data_id" {
  value = data.jamfpro_macos_configuration_profile_plist_generator.jamfpro_macos_configuration_profile_plist_generator_001_data.id
}

output "jamfpro_macos_configuration_profile_plist_generator_001_data_name" {
  value = data.jamfpro_macos_configuration_profile_plist_generator.jamfpro_macos_configuration_profile_plist_generator_001_data.name
}