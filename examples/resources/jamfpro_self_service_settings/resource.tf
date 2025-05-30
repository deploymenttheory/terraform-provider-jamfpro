resource "jamfpro_self_service_settings" "example" {
  install_automatically    = true
  install_location         = "/Applications"
  user_login_level         = "Anonymous"
  allow_remember_me        = true
  use_fido2                = false
  auth_type                = "Basic"
  notifications_enabled    = true
  alert_user_approved_mdm  = true
  default_landing_page     = "HOME"
  default_home_category_id = -1
  bookmarks_name           = "Bookmarks"
}
