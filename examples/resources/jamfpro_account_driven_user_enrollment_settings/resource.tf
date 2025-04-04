resource "jamfpro_account_driven_user_enrollment_settings" "example" {
  enabled                     = true
  expiration_interval_days    = 30
  # Optional - you can use days OR seconds, but not both
  # expiration_interval_seconds = 2592000 # 30 days in seconds
}

# Output the configured settings
output "adue_settings" {
  value = {
    enabled              = jamfpro_account_driven_user_enrollment_settings.example.enabled
    expiration_days      = jamfpro_account_driven_user_enrollment_settings.example.expiration_interval_days
  }
  description = "Account Driven User Enrollment Settings"
}