data "jamfpro_account" "jamfpro_account_001_data" {
  id = jamfpro_account.jamfpro_account_001.id
}

output "jamfpro_account_001_data_id" {
  value = data.jamfpro_account.jamfpro_account_001_data.id
}

output "jamfpro_account_001_data_name" {
  value = data.jamfpro_account.jamfpro_account_001_data.name
}
