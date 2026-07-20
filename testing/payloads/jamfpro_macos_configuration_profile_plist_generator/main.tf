// Regression test for issue #1145 - heredoc strings in HCL always include a
// trailing newline before EOT, but the API strips it server-side. Without a
// DiffSuppressFunc on description/self_service_description, this produced
// perpetual drift on every plan after apply.
resource "jamfpro_macos_configuration_profile_plist_generator" "jamfpro_macos_configuration_profile_plist_generator_heredoc" {
  name                = "tf-testing-${var.testing_id}-heredoc-${random_id.rng.hex}"
  distribution_method = "Make Available in Self Service"
  level               = "System"
  redeploy_on_update  = "Newly Assigned"
  user_removable      = false

  description = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    install_button_text      = "Install"
    self_service_description = <<-EOT
      Multi-line self service description used to verify no
      drift is reported after apply due to the heredoc trailing newline.
    EOT
  }

  payloads {
    payload_description_header  = "Test payload for issue 1145 regression"
    payload_enabled_header      = true
    payload_organization_header = "tf-testing"
    payload_type_header         = "Configuration"
    payload_version_header      = 1

    payload_content {
      payload_enabled      = true
      payload_organization = "tf-testing"
      payload_type         = "com.jamf.tf-testing.issue-1145"
      payload_version      = 1

      setting {
        key   = "TestKey"
        value = "TestValue"
      }
    }
  }
}
