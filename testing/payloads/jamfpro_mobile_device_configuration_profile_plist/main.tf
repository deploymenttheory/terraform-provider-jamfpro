resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_002" {
  name = "tf-testing-${var.testing_id}-profile-restrictions-${random_id.rng.hex}"
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

# Regression coverage for the phantom empty <array/> bug: a deliberately
# pretty-printed com.apple.homescreenlayout payload (whitespace between the
# <array> tags). The provider must compact that structural whitespace before
# sending so the Classic API does not inject a blank leading home-screen page.
resource "jamfpro_mobile_device_configuration_profile_plist" "home_screen_layout_pretty" {
  name               = "tf-testing-${var.testing_id}-profile-home-screen-layout-${random_id.rng.hex}"
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

# ========================================================================== //
# Real-world iOS / iPadOS configuration profiles (sanitized)
# ========================================================================== //
#
# A diverse set of production-shaped .mobileconfig payloads, selected for
# breadth of payload domains and plist field types (nested <dict>, <array> of
# <integer>, <array> of <string>, <array> of <dict>, <data>). Every fixture is
# deliberately pretty-printed (whitespace between sibling elements) so it
# exercises the structural-whitespace compaction path on both create and
# update. All organisation-identifying data is sanitized to "Deployment Theory"
# and the certificate is a throwaway self-signed cert.

# com.apple.wifi.managed — WPA2-Enterprise WiFi; nested EAPClientConfiguration
# dict with an <array> of <integer> (AcceptEAPTypes).
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_wifi_wpa2_enterprise" {
  name               = "tf-testing-${var.testing_id}-mobile-wifi-${random_id.rng.hex}"
  description        = "WPA2-Enterprise WiFi (nested EAP dict, integer array)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/WiFi-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# com.apple.webcontent-filter — array-heavy allowlist (PermittedURLs +
# WhitelistedBookmarks array of dicts); same inter-<array> whitespace class as
# the home-screen-layout regression.
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_web_content_filter" {
  name               = "tf-testing-${var.testing_id}-mobile-web-content-filter-${random_id.rng.hex}"
  description        = "Web content filter allowlist (array of bookmark dicts)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/WebContentFilter-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# com.apple.domains — managed domains; multiple <array> of <string>.
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_domains" {
  name               = "tf-testing-${var.testing_id}-mobile-domains-${random_id.rng.hex}"
  description        = "Managed domains (string arrays)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Domains-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# com.apple.security.root — root CA certificate; exercises a base64 <data> leaf.
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_certificate_root" {
  name               = "tf-testing-${var.testing_id}-mobile-certificate-${random_id.rng.hex}"
  description        = "Root CA certificate (base64 <data> leaf preservation)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Certificate-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# com.apple.notificationsettings — per-app notification settings; repeated
# nested dicts with <integer> AlertType inside an <array>.
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_notifications" {
  name               = "tf-testing-${var.testing_id}-mobile-notifications-${random_id.rng.hex}"
  description        = "Per-app notification settings (array of nested dicts, integers)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Notifications-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

# com.apple.eas.account — Exchange ActiveSync mail account; string/boolean leaves.
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_exchange" {
  name               = "tf-testing-${var.testing_id}-mobile-exchange-${random_id.rng.hex}"
  description        = "Exchange ActiveSync account (string and boolean leaves)."
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Exchange-Pretty.mobileconfig")
  payload_validate   = false

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}
