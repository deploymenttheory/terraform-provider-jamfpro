data "jamfpro_account_groups" "example_account_group" {
  name = "tf-example-account-group-01"  # Replace this with the actual name of the account group you want to retrieve
}

output "account_group_id" {
  value = data.jamfpro_account_group.example_account_group.id
}

output "account_group_name" {
  value = data.jamfpro_account_group.example_account_group.name
}

output "account_group_access_level" {
  value = data.jamfpro_account_group.example_account_group.access_level
}

output "account_group_privilege_set" {
  value = data.jamfpro_account_group.example_account_group.privilege_set
}

output "account_group_site" {
  value = data.jamfpro_account_group.example_account_group.site
}

output "account_group_jss_objects_privileges" {
  value = data.jamfpro_account_group.example_account_group.jss_objects_privileges
}

output "account_group_jss_settings_privileges" {
  value = data.jamfpro_account_group.example_account_group.jss_settings_privileges
}

output "account_group_jss_actions_privileges" {
  value = data.jamfpro_account_group.example_account_group.jss_actions_privileges
}

output "account_group_casper_admin_privileges" {
  value = data.jamfpro_account_group.example_account_group.casper_admin_privileges
}

output "account_group_members" {
  value = data.jamfpro_account_group.example_account_group.members
}
