// Regression test for issue #1145 - heredoc strings in HCL always include a
// trailing newline before EOT, but the API strips it server-side. Without a
// DiffSuppressFunc on description/body, this produced perpetual drift on
// every plan after apply.
resource "jamfpro_enrollment_customization" "jamfpro_enrollment_customization_heredoc" {
  site_id                               = "-1"
  display_name                          = "tf-testing-${var.testing_id}-heredoc-${random_id.rng.hex}"
  enrollment_customization_image_source = "${path.module}/icon.png"

  description = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  text_pane {
    display_name         = "Welcome Message"
    rank                 = 1
    title                = "Welcome"
    body                 = <<-EOT
      Multi-line body used to verify no drift is
      reported after apply due to the heredoc trailing newline.
    EOT
    back_button_text     = "Back"
    continue_button_text = "Continue"
  }
}
