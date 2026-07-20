// ========================================================================== //
// Restricted Software
// ========================================================================== //

// Minimal - only required fields
resource "jamfpro_restricted_software" "restricted_software_min" {
  name         = "tf-test-rs-min"
  process_name = "SomeUnwantedApp.app"
}

// With all optional fields but no scope
resource "jamfpro_restricted_software" "restricted_software_max" {
  name                     = "tf-test-rs-max"
  process_name             = "BitTorrent.app"
  match_exact_process_name = true
  send_notification        = true
  kill_process             = true
  delete_executable        = true
  display_message          = "This application is not permitted."
}

// Partial options
resource "jamfpro_restricted_software" "restricted_software_notify_only" {
  name                     = "tf-test-rs-notify"
  process_name             = "Install macOS"
  match_exact_process_name = false
  send_notification        = true
  kill_process             = false
  delete_executable        = false
  display_message          = "Please contact IT before installing macOS upgrades."
}

// ========================================================================== //
// Regression test for issue #1145 - heredoc strings in HCL always include a
// trailing newline before EOT, but the API strips it server-side. Without a
// DiffSuppressFunc on display_message, this produced perpetual drift on every
// plan after apply. Assert this stays a clean plan.

resource "jamfpro_restricted_software" "restricted_software_heredoc_message" {
  name                     = "tf-test-rs-heredoc"
  process_name             = "BitTorrent.app"
  match_exact_process_name = true
  send_notification        = true
  kill_process             = true
  delete_executable        = false
  display_message          = <<-EOT
    This application is not permitted on managed devices.
    Contact IT if you believe this is in error.
  EOT
}
