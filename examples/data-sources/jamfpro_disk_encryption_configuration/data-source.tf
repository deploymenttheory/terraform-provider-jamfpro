data "jamfpro_disk_encryption_configuration" "jamfpro_disk_encryption_configuration_002_data" {
  id = jamfpro_disk_encryption_configuration.jamfpro_disk_encryption_configuration_002.id
}

output "jamfpro_disk_encryption_configuration_002_id" {
  value = data.jamfpro_disk_encryption_configuration.jamfpro_disk_encryption_configuration_002_data.id
}

output "jamfpro_disk_encryption_configuration_002_name" {
  value = data.jamfpro_disk_encryption_configuration.jamfpro_disk_encryption_configuration_002_data.name
}