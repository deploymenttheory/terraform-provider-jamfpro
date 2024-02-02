data "jamfpro_disk_encryption_configuration" "example_disk_encryption_configuration" {
  name = "jamfpro-tf-example-InstitutionalRecoveryKey-config"  # Replace this with the actual name of the disk encryption configuration you want to retrieve
}

output "disk_encryption_configuration_id" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.id
}

output "disk_encryption_configuration_name" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.name
}

output "disk_encryption_configuration_key_type" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.key_type
}

output "disk_encryption_configuration_file_vault_enabled_users" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.file_vault_enabled_users
}

output "disk_encryption_configuration_institutional_recovery_key_certificate_type" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.institutional_recovery_key.certificate_type
}

output "disk_encryption_configuration_institutional_recovery_key_password" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.institutional_recovery_key.password
}

output "disk_encryption_configuration_institutional_recovery_key_data" {
  value = data.jamfpro_disk_encryption_configuration.example_disk_encryption_configuration.institutional_recovery_key.data
}
