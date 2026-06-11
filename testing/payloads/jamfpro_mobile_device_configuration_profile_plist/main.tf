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

# The remaining pretty-printed payloads exercise structural-whitespace
# compaction across a variety of payload types: a certificate (multi-line
# indented base64 inside <data> — non-structural whitespace the compactor
# must leave intact), a web content filter (Jamf-exported, requires
# UserDefinedName for the server-side content-filter record), managed
# domains (sibling string arrays), and WiFi (scalar-heavy single dict).
resource "jamfpro_mobile_device_configuration_profile_plist" "web_content_filter_pretty" {
  name               = "test-web-content-filter-pretty"
  description        = "Web content filter payload (structural-whitespace compaction regression test)"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/WebContentFilter-Pretty.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

resource "jamfpro_mobile_device_configuration_profile_plist" "certificate_pretty" {
  name               = "test-certificate-pretty"
  description        = "Pretty-printed certificate payload (structural-whitespace compaction regression test)"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Certificate-Pretty.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

resource "jamfpro_mobile_device_configuration_profile_plist" "wifi_pretty" {
  name               = "test-wifi-pretty"
  description        = "Pretty-printed WiFi payload (structural-whitespace compaction regression test)"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/WiFi-Pretty.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}

resource "jamfpro_mobile_device_configuration_profile_plist" "domains_pretty" {
  name               = "test-domains-pretty"
  description        = "Pretty-printed managed domains payload (structural-whitespace compaction regression test)"
  level              = "Device Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("${path.module}/Domains-Pretty.mobileconfig")

  scope {
    all_mobile_devices = true
    all_jss_users      = false
  }
}
