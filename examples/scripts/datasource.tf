data "jamfpro_scripts" "example_script" {
  id = jamfpro_scripts.example_script.id
}

output "jamfpro_printer_id" {
  value = data.jamfpro_scripts.example_script.id
}

output "jamfpro_script_name" {
  value = data.jamfpro_scripts.example_script.name
}