# Example using id
data "jamfpro_computer_extension_attribute" "example_by_id" {
  id = "123"
}

# Example using name
data "jamfpro_computer_extension_attribute" "example_by_name" {
  name = "RAM_Usage"
}

# Example outputs to show the data usage
output "extension_attribute_id" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.id
}

output "extension_attribute_name" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.name
}

output "extension_attribute_description" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.description
}

output "extension_attribute_data_type" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.data_type
}

output "extension_attribute_enabled" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.enabled
}

output "extension_attribute_inventory_display" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.inventory_display_type
}

output "extension_attribute_input_type" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.input_type
}

output "extension_attribute_script_contents" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.script_contents
}

output "extension_attribute_popup_choices" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.popup_menu_choices
}

output "extension_attribute_ldap_mapping" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.ldap_attribute_mapping
}

output "extension_attribute_ldap_allowed" {
  value = data.jamfpro_computer_extension_attribute.example_by_name.ldap_extension_attribute_allowed
}