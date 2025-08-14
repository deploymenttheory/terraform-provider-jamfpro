// Definition a Jamf Pro Package Resource
resource "jamfpro_package" "google_chrome_enterprise" {
  package_name          = "Nudge_Essentials-2.0.12.81807.pkg"   // Required
  package_file_source   = "./Nudge_Essentials-2.0.12.81807.pkg" // Required
  priority              = 10                                    // Required
  reboot_required       = false                                 // Required
  fill_user_template    = false                                 // Required
  fill_existing_users   = false                                 // Required
  os_install            = false                                 // Required
  suppress_updates      = false                                 // Required
  suppress_from_dock    = false                                 // Required
  suppress_eula         = false                                 // Required
  suppress_registration = false                                 // Required
  timeouts {
    create = "5m"
  }
}
