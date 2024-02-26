data "jamfpro_account_group" "jamfpro_account_group_001_data" {
  id = jamfpro_account_group.jamfpro_account_group_001.id
}

output "jamfpro_jamfpro_account_group_001_id" {
  value = data.jamfpro_account_group.jamfpro_account_group_001_data.id
}

output "jamfpro_jamfpro_account_groups_001_name" {
  value = data.jamfpro_account_group.jamfpro_account_group_001_data.name
}
