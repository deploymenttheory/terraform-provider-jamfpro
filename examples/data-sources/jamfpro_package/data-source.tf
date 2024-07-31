data "jamfpro_package" "jamfpro_package_001_data" {
  id = jamfpro_package.jamfpro_package_001.id
}

output "jamfpro_package_001_data_id" {
  value = data.jamfpro_package.jamfpro_package_001_data.id
}

output "jamfpro_package_001_data_name" {
  value = data.jamfpro_package.jamfpro_package_001_data.name
}