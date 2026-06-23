resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_064" {
  name                = "test-profile"
  description         = "An example mobile device configuration profile."
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/Screen Recording - Allow Microsoft Teams.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-CIS-Compliance-ManagedPreferences.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-GoogleChrome-Restrictions.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-Notifications-LockScreen.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-Symantec-WSS-VPN.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-MicrosoftTeams-PPPC.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-Accessibility-Zoom.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-FileVault-DiskEncryption.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-Managed-LoginItems.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-Microsoft-PlatformSSO.mobileconfig")
  payload_validate    = false
  user_removable      = false

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
  payloads            = file("${path.module}/DT-System-Extensions.mobileconfig")
  payload_validate    = false
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }
}
