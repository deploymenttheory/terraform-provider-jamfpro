resource "jamfpro_enrollment_customization" "example" {
  site_id      = "-1"  # -1 for None
  display_name = "Corporate Enrollment"
  description  = "Default corporate enrollment experience"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"  # Black text
    button_color      = "0066CC"  # Blue buttons
    button_text_color = "FFFFFF"  # White button text
    background_color  = "F5F5F5"  # Light gray background
  }
}
