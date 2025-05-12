data "jamfpro_device_enrollments_public_key" "current" {}

output "device_enrollments_public_key" {
  value = data.jamfpro_device_enrollments_public_key.current.public_key
}
