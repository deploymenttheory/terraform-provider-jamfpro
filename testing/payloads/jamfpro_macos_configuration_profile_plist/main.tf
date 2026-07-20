# PayloadUUID/PayloadIdentifier values are derived per run from var.testing_id
# via uuidv5. The fixtures' static UUIDs collided with profiles orphaned on the
# shared sandbox by earlier failed runs, which the Classic API rejects with
# 409 Duplicate payload uuid. uuidv5 is used rather than random_uuid because it
# is deterministic and known at plan time — an unknown payloads value makes the
# provider's CustomizeDiff plist validators see an empty string and fail.

resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_064" {
  name = "tf-testing-${var.testing_id}-profile-screen-recording-pppc-${random_id.rng.hex}"
  // Regression test for issue #1145 - heredoc strings in HCL always include
  // a trailing newline before EOT, but the API strips it server-side.
  description         = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/Screen Recording - Allow Microsoft Teams.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-jamfpro_macos_configuration_profile_064-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-jamfpro_macos_configuration_profile_064-1"))
    uuid_2 = upper(uuidv5("dns", "${var.testing_id}-jamfpro_macos_configuration_profile_064-2"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }

}

# Regression coverage for the phantom empty <array/> bug on macOS: a
# pretty-printed com.apple.systempreferences payload with two sibling
# arrays-of-strings (whitespace between the <string> children is the same class
# the Classic API mis-parses). Compaction must leave the panes lists intact.
resource "jamfpro_macos_configuration_profile_plist" "system_preferences_pretty" {
  name                = "tf-testing-${var.testing_id}-profile-system-preferences-${random_id.rng.hex}"
  description         = "Pretty-printed System Preferences profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/SystemPreferences-Pretty.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-system_preferences_pretty-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-system_preferences_pretty-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# The remaining pretty-printed payloads exercise structural-whitespace
# compaction across a variety of payload types: dock (deeply nested tile-data
# dicts inside an array), login window (multi-byte UTF-8 in leaf
# strings), and screensaver (scalar booleans/integers).
resource "jamfpro_macos_configuration_profile_plist" "dock_pretty" {
  name                = "tf-testing-${var.testing_id}-profile-dock-${random_id.rng.hex}"
  description         = "Pretty-printed Dock profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/Dock-Pretty.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-dock_pretty-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-dock_pretty-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

resource "jamfpro_macos_configuration_profile_plist" "login_window_pretty" {
  name                = "tf-testing-${var.testing_id}-profile-login-window-${random_id.rng.hex}"
  description         = "Pretty-printed Login Window profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/LoginWindow-Pretty.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-login_window_pretty-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-login_window_pretty-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

resource "jamfpro_macos_configuration_profile_plist" "screensaver_pretty" {
  name                = "tf-testing-${var.testing_id}-profile-screensaver-${random_id.rng.hex}"
  description         = "Pretty-printed Screensaver profile (structural-whitespace compaction regression test)"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/Screensaver-Pretty.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-screensaver_pretty-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-screensaver_pretty-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# ========================================================================== //
# Real-world macOS configuration profiles (sanitized)
# ========================================================================== //
#
# A diverse set of production-shaped .mobileconfig payloads, selected for
# breadth of payload domains, structural depth and plist field types
# (<data>, <real>, <integer>, deeply nested <array>/<dict>). They exercise the
# structural-whitespace compaction path across realistic profiles. All
# organisation-identifying data has been sanitized to "Deployment Theory" and
# the FileVault escrow certificate replaced with a throwaway self-signed cert.

# com.apple.ManagedClient.preferences — CIS Level 1 compliance; extreme dict depth.
resource "jamfpro_macos_configuration_profile_plist" "macos_cis_compliance" {
  name                = "tf-testing-${var.testing_id}-macos-cis-compliance-${random_id.rng.hex}"
  description         = "CIS Level 1 compliance managed preferences (deeply nested dictionaries)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-CIS-Compliance-ManagedPreferences.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_cis_compliance-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_cis_compliance-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.ManagedClient.preferences — Google Chrome restrictions; large arrays + integers.
resource "jamfpro_macos_configuration_profile_plist" "macos_chrome_restrictions" {
  name                = "tf-testing-${var.testing_id}-macos-chrome-restrictions-${random_id.rng.hex}"
  description         = "Google Chrome managed restrictions (array- and integer-heavy managed preferences)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-GoogleChrome-Restrictions.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_chrome_restrictions-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_chrome_restrictions-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.notificationsettings — per-app notification settings incl. lock screen.
resource "jamfpro_macos_configuration_profile_plist" "macos_notifications" {
  name                = "tf-testing-${var.testing_id}-macos-notifications-${random_id.rng.hex}"
  description         = "Per-app notification settings including lock screen (repeated nested dictionaries)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-Notifications-LockScreen.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_notifications-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_notifications-1"))
    uuid_2 = upper(uuidv5("dns", "${var.testing_id}-macos_notifications-2"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.vpn.managed — Symantec WSS Agent VPN; multiple VPN payloads, integer-heavy.
resource "jamfpro_macos_configuration_profile_plist" "macos_wss_vpn" {
  name                = "tf-testing-${var.testing_id}-macos-wss-vpn-${random_id.rng.hex}"
  description         = "Symantec WSS Agent VPN profile (two com.apple.vpn.managed payloads)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-Symantec-WSS-VPN.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_wss_vpn-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_wss_vpn-1"))
    uuid_2 = upper(uuidv5("dns", "${var.testing_id}-macos_wss_vpn-2"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.TCC.configuration-profile-policy — Microsoft Teams PPPC services.
resource "jamfpro_macos_configuration_profile_plist" "macos_teams_pppc" {
  name                = "tf-testing-${var.testing_id}-macos-teams-pppc-${random_id.rng.hex}"
  description         = "Microsoft Teams Privacy Preferences Policy Control (nested services array)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-MicrosoftTeams-PPPC.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_teams_pppc-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_teams_pppc-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.universalaccess — accessibility zoom; exercises <real> float fields.
resource "jamfpro_macos_configuration_profile_plist" "macos_accessibility_zoom" {
  name                = "tf-testing-${var.testing_id}-macos-accessibility-zoom-${random_id.rng.hex}"
  description         = "Universal Access / Zoom accessibility settings (uses <real> float fields)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-Accessibility-Zoom.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_accessibility_zoom-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_accessibility_zoom-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.MCX.FileVault2 + FDERecoveryKeyEscrow + pkcs1 — FileVault; exercises <data>.
resource "jamfpro_macos_configuration_profile_plist" "macos_filevault" {
  name                = "tf-testing-${var.testing_id}-macos-filevault-${random_id.rng.hex}"
  description         = "FileVault disk encryption with recovery key escrow certificate (<data> field, multi-payload)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-FileVault-DiskEncryption.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-1"))
    uuid_2 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-2"))
    uuid_3 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-3"))
    uuid_4 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-4"))
    uuid_5 = upper(uuidv5("dns", "${var.testing_id}-macos_filevault-5"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.servicemanagement — managed login items / background service rules.
resource "jamfpro_macos_configuration_profile_plist" "macos_login_items" {
  name                = "tf-testing-${var.testing_id}-macos-login-items-${random_id.rng.hex}"
  description         = "Managed login items / background service management rules."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-Managed-LoginItems.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_login_items-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_login_items-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.extensiblesso — Microsoft Entra Platform SSO extension.
resource "jamfpro_macos_configuration_profile_plist" "macos_platform_sso" {
  name                = "tf-testing-${var.testing_id}-macos-platform-sso-${random_id.rng.hex}"
  description         = "Microsoft Entra Platform SSO extension (com.apple.extensiblesso)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-Microsoft-PlatformSSO.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_platform_sso-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_platform_sso-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}

# com.apple.system-extension-policy — allowed system extensions.
resource "jamfpro_macos_configuration_profile_plist" "macos_system_extensions" {
  name                = "tf-testing-${var.testing_id}-macos-system-extensions-${random_id.rng.hex}"
  description         = "Allowed system extensions policy (com.apple.system-extension-policy)."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads = templatefile("${path.module}/DT-System-Extensions.mobileconfig", {
    uuid_0 = upper(uuidv5("dns", "${var.testing_id}-macos_system_extensions-0"))
    uuid_1 = upper(uuidv5("dns", "${var.testing_id}-macos_system_extensions-1"))
  })
  payload_validate = false
  user_removable   = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}
