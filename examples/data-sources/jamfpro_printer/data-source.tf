data "jamfpro_printer" "jamfpro_printers_001_data" {
  id = jamfpro_printer.jamfpro_printers_001.id
}

output "jamfpro_jamfpro_printers_001_id" {
  value = data.jamfpro_printer.jamfpro_printers_001_data.id
}

output "jamfpro_jamfpro_printers_001_name" {
  value = data.jamfpro_printer.jamfpro_printers_001_data.name
}
