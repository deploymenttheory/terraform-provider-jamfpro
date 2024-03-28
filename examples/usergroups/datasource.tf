data "jamfpro_user_group" "usergroup_001_data" {
  id = jamfpro_user_group.jamfpro_user_group_001.id
}

output "jamfpro_user_group_001_id" {
  value = data.jamfpro_user_group.usergroup_001_data.id
}

output "jamfpro_user_group_001_name" {
  value = data.jamfpro_user_group.usergroup_001_data.name
}

data "jamfpro_user_group" "usergroup_002_data" {
  id = jamfpro_user_group.jamfpro_user_group_002.id
}

output "jamfpro_user_group_002_id" {
  value = data.jamfpro_user_group.usergroup_002_data.id
}

output "jamfpro_user_group_002_name" {
  value = data.jamfpro_user_group.usergroup_002_data.name
}
