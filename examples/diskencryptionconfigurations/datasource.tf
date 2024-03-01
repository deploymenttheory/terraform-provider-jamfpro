data "jamfpro_disk_encryption_configuration" "disk_encryption_configuration_001_data" {
  id = jamfpro_disk_encryption_configuration.disk_encryption_configuration_001.id
}

output "disk_encryption_configuration_001_id" {
  value = data.jamfpro_disk_encryption_configuration.disk_encryption_configuration_001_data.id
}

output "disk_encryption_configuration_001_name" {
  value = data.jamfpro_disk_encryption_configuration.disk_encryption_configuration_001_data.name
}