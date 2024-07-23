data "jamfpro_static_computer_group" "jamfpro_static_computer_group_001_data" {
  id = jamfpro_static_computer_group.jamfpro_static_computer_group_001.id
}

output "jamfpro_jamfpro_static_computer_group_001_id" {
  value = data.jamfpro_static_computer_group.jamfpro_static_computer_group_001_data.id
}

output "jamfpro_jamfpro_static_computer_groups_001_name" {
  value = data.jamfpro_static_computer_group.jamfpro_static_computer_group_001_data.name
}
