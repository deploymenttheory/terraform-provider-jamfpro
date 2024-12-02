// data source by id

data "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_001_data" {
  id = jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001.id
}

output "jamfpro_computer_extension_attribute_001_data_id" {
  value = data.jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001_data.id
}

output "jamfpro_computer_extension_attribute_001_data_name" {
  value = data.jamfpro_computer_extension_attribute.jamfpro_computer_extension_attribute_001_data.name
}

// data source list

data "jamfpro_computer_extension_attributes_list" "example" {}

output "attribute_ids" {
  value = data.jamfpro_computer_extension_attributes_list.example.ids
}

output "attributes" {
  value = data.jamfpro_computer_extension_attributes_list.example.attributes
}