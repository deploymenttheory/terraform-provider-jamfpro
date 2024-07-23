data "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_001_data" {
  id = jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001.id
}

output "jamfpro_computer_extension_attribute_001_data_id" {
  value = data.jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001_data.id
}

output "jamfpro_computer_extension_attribute_001_data_name" {
  value = data.jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001_data.name
}