resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_002" {
  name               = "test-profile"
  description        = "Description"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Restrictions-Baseline.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# Regression coverage for the phantom empty <array/> bug: a deliberately
# pretty-printed com.apple.homescreenlayout payload (whitespace between the
# <array> tags). The provider must compact that structural whitespace before
# sending so the Classic API does not inject a blank leading home-screen page.
resource "jamfpro_mobile_device_configuration_profile_plist" "home_screen_layout_pretty" {
  name               = "test-home-screen-layout-pretty"
  description        = "Pretty-printed home screen layout (structural-whitespace compaction regression test)"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/HomeScreenLayout-Pretty.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}
