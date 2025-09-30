resource "jamfpro_self_service_branding_ios" "example" {
  main_header                  = "My Organization (iOS)"
  icon_id                      = 5
  header_background_color_code = "0066CC"
  menu_icon_color_code         = "FFFFFF"
  branding_name_color_code     = "FFFFFF"
  status_bar_text_color        = "light"
}
