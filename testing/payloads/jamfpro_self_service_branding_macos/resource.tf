resource "jamfpro_self_service_branding_macos" "example" {
  application_header   = "Self Service"
  sidebar_heading      = "My Organization"
  sidebar_subheading   = "Division"
  home_page_heading    = "Welcome to Self Service"
  home_page_subheading = "Choose an app to install"
}
