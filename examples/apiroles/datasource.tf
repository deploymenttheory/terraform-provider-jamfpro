data "jamfpro_api_role" "jamfpro_api_role_001_data" {
  id = jamfpro_api_role.jamfpro_api_role_001.id
}

output "jamfpro_api_role_001_data_id" {
  value = data.jamfpro_api_role.jamfpro_api_role_001_data.id
}

output "jamfpro_api_role_001_data_name" {
  value = data.jamfpro_api_role.jamfpro_api_role_001_data.name
}
