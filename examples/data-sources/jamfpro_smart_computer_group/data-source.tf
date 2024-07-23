data "jamfpro_smart_computer_group" "jamfpro_smart_computer_group_001_data" {
  id = jamfpro_smart_computer_group.jamfpro_smart_computer_group_001.id
}

output "jamfpro_jamfpro_smart_computer_group_001_id" {
  value = data.jamfpro_smart_computer_group.jamfpro_smart_computer_group_001_data.id
}

output "jamfpro_jamfpro_smart_computer_groups_001_name" {
  value = data.jamfpro_smart_computer_group.jamfpro_smart_computer_group_001_data.name
}
