resource "jamfpro_device_communication_settings" "example" {
  auto_renew_mobile_device_mdm_profile_when_ca_renewed                    = true
  auto_renew_mobile_device_mdm_profile_when_device_identity_cert_expiring = true
  auto_renew_computer_mdm_profile_when_ca_renewed                         = true
  auto_renew_computer_mdm_profile_when_device_identity_cert_expiring      = true
  mdm_profile_mobile_device_expiration_limit_in_days                      = 90
  mdm_profile_computer_expiration_limit_in_days                           = 90
}