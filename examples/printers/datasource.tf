data "jamfpro_printers" "example_printer" {
  id = resource.jamfpro_printers.example_printer.id
}

output "jamfpro_printer_id" {
  value = data.jamfpro_printers.example_printer.id
}

output "jamfpro_script_name" {
  value = data.jamfpro_printers.example_printer.name
}