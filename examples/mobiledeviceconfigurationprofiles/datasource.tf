data "jamfpro_mobile_device_configuration_profile" "mobile_device_configuration_profile_001_data" {
  id = jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_001.id
}

output "jamfpro_mobile_device_configuration_profile_001_id" {
  value = data.jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_001_data.id
}

output "jamfpro_mobile_device_configuration_profile_001_name" {
  value = data.jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_001_data.name
}

data "jamfpro_mobile_device_configuration_profile" "mobile_device_configuration_profile_002_data" {
  id = jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_002.id
}

output "jamfpro_mobile_device_configuration_profile_002_id" {
  value = data.jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_002_data.id
}

output "jamfpro_mobile_device_configuration_profile_002_name" {
  value = data.jamfpro_mobile_device_configuration_profile.mobile_device_configuration_profile_002_data.name
}
