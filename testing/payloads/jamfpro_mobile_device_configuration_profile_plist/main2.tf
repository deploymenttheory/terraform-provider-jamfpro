resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_003" {
  name               = "test-profile-2"
  description        = "Description"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Restrictions-Shared.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}
