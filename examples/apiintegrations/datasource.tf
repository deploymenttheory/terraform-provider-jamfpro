data "jamfpro_api_integration" "jamfpro_api_integration_001_data" {
  id = jamfpro_api_integration.jamfpro_api_integration_001.id
}

output "jamfpro_api_integration_001_data_id" {
  value = data.jamfpro_api_integration.jamfpro_api_integration_001_data.id
}

output "jamfpro_api_integration_001_data_name" {
  value = data.jamfpro_api_integration.jamfpro_api_integration_001_data.name
}
