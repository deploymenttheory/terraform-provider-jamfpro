data "jamfpro_smart_mobile_device_group" "jamfpro_smart_mobile_device_group_001_data" {
  id = jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001.id
}

output "jamfpro_jamfpro_smart_mobile_device_group_001_id" {
  value = data.jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001_data.id
}

output "jamfpro_jamfpro_smart_mobile_device_groups_001_name" {
  value = data.jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001_data.name
}

output "jamfpro_jamfpro_smart_mobile_device_groups_001_description" {
  value = data.jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001_data.description
}
