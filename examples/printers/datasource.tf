data "jamfpro_printers" "example_printer" {
  id = "12345"  # Replace this with the actual ID of the printer you want to retrieve
}

output "printer_id" {
  value = data.jamfpro_printers.example_printer.id
}

output "printer_name" {
  value = data.jamfpro_printers.example_printer.name
}

output "printer_category" {
  value = data.jamfpro_printers.example_printer.category
}

output "printer_uri" {
  value = data.jamfpro_printers.example_printer.uri
}

output "printer_cups_name" {
  value = data.jamfpro_printers.example_printer.cups_name
}

output "printer_location" {
  value = data.jamfpro_printers.example_printer.location
}

output "printer_model" {
  value = data.jamfpro_printers.example_printer.model
}

output "printer_info" {
  value = data.jamfpro_printers.example_printer.info
}

output "printer_notes" {
  value = data.jamfpro_printers.example_printer.notes
}

output "printer_make_default" {
  value = data.jamfpro_printers.example_printer.make_default
}

output "printer_use_generic" {
  value = data.jamfpro_printers.example_printer.use_generic
}

output "printer_ppd" {
  value = data.jamfpro_printers.example_printer.ppd
}

output "printer_ppd_path" {
  value = data.jamfpro_printers.example_printer.ppd_path
}

output "printer_ppd_contents" {
  value = data.jamfpro_printers.example_printer.ppd_contents
}
