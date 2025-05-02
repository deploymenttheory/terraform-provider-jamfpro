data "jamfpro_mobile_device_prestage_enrollment" "example" {
  id = "1"
}

output "prestage_name" {
  value = data.jamfpro_mobile_device_prestage_enrollment.example.display_name
}
