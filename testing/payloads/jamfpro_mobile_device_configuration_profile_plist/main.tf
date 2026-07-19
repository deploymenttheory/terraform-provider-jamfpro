resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_002" {
  name = "test-profile"
  // Regression test for issue #1145 - heredoc strings in HCL always include
  // a trailing newline before EOT, but the API strips it server-side.
  description        = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Restrictions-Baseline.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}
