resource "jamfpro_enrollment_customization" "test" {
 display_name = "tf-sets-enrollment-20260919"
 description = "Temporary unattached set migration test"
 enrollment_customization_image_source = ""
 site_id = "-1"
 branding_settings {
  text_color = "000000"
  button_color = "000000"
  button_text_color = "FFFFFF"
  background_color = "FFFFFF"
 }
 ldap_pane {
  display_name = "Temporary LDAP pane"
  rank = 1
  title = "Sign in"
  username_label = "Username"
  password_label = "Password"
  back_button_text = "Back"
  continue_button_text = "Continue"
 }
}
