// ========================================================================== //
// Mobile Device Applications
// ========================================================================== //

resource "jamfpro_mobile_device_application" "app_store_app" {
  name             = "tf-testing-${var.testing_id}-mobile-device-application-${random_id.rng.hex}"
  display_name     = "tf-testing-${var.testing_id}-mobile-device-application-${random_id.rng.hex}"
  bundle_id        = "com.jamfsoftware.selfservice"
  version          = "11.2.1"
  itunes_store_url = "https://apps.apple.com/us/app/jamf-self-service/id718509958"

  deploy_automatically  = false
  deploy_as_managed_app = true

  scope {
    all_mobile_devices = false
    all_jss_users      = false
  }

  // Regression test for issue #1145 - heredoc strings in HCL always include a
  // trailing newline before EOT, but the API strips it server-side. Without
  // DiffSuppressFunc on self_service_description/preferences, this produced
  // perpetual drift on every plan after apply.
  self_service {
    self_service_install_button_text = "Install"
    self_service_description         = <<-EOT
      Multi-line self service description used to verify no
      drift is reported after apply due to the heredoc trailing newline.
    EOT
    feature_on_main_page             = false
  }

  app_configuration {
    preferences = <<-EOT
      <dict>
      <key>INVITATION_STRING</key>
      <string>$MOBILEDEVICEAPPINVITE</string>
      </dict>
    EOT
  }
}
