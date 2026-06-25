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
}
