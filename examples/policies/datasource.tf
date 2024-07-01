data "jamfpro_policy" "jamfpro_policys_001_data" {
  id = jamfpro_policy.jamfpro_policys_001.id
}

output "jamfpro_jamfpro_policys_001_id" {
  value = data.jamfpro_policy.jamfpro_policys_001_data.id
}

output "jamfpro_jamfpro_policys_001_name" {
  value = data.jamfpro_policy.jamfpro_policys_001_data.name
}
