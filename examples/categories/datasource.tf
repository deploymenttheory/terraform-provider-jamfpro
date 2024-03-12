data "jamfpro_category" "jamfpro_category_001_data" {
  id = jamfpro_category.jamfpro_category_001.id
}

output "jamfpro_category_001_data_id" {
  value = data.jamfpro_category.jamfpro_category_001_data.id
}

output "jamfpro_category_001_data_name" {
  value = data.jamfpro_category.jamfpro_category_001_data.name
}
