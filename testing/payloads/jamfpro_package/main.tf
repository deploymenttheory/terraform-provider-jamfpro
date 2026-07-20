// Definition a Jamf Pro Package Resource
resource "jamfpro_package" "google_chrome_enterprise" {
  package_name          = "Nudge_Essentials-2.0.12.81807.pkg"                // Required
  package_file_source   = "${path.module}/Nudge_Essentials-2.0.12.81807.pkg" // Required
  priority              = 10                                                 // Required
  reboot_required       = false                                              // Required
  fill_user_template    = false                                              // Required
  fill_existing_users   = false                                              // Required
  os_install            = false                                              // Required
  suppress_updates      = false                                              // Required
  suppress_from_dock    = false                                              // Required
  suppress_eula         = false                                              // Required
  suppress_registration = false                                              // Required
  // Regression test for issue #1145 - heredoc strings in HCL always include
  // a trailing newline before EOT, but the API strips it server-side.
  info  = <<-EOT
    Multi-line info field used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  notes = <<-EOT
    Multi-line notes field used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  timeouts {
    create = "5m"
  }
}
