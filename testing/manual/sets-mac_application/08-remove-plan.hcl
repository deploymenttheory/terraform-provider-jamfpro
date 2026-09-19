resource "jamfpro_category" "alpha" {
 name = "tf-sets-mac-alpha-20260919"
 priority = 9
}
resource "jamfpro_category" "bravo" {
 name = "tf-sets-mac-bravo-20260919"
 priority = 9
}
resource "jamfpro_mac_application" "test" {
 name = "tf-sets-mac-application-20260919"
 version = "1.0"
 bundle_id = "com.tinyspeck.slackmacgap"
 url = "https://apps.apple.com/gb/app/slack-for-desktop/id803453959?mt=12"
 deployment_type = "Make Available in Self Service"
 site_id = -1
 category_id = -1
 scope {
  all_computers = false
  all_jss_users = false
 }
 self_service {
  install_button_text = "Install"
  self_service_description = "Temporary unscoped set migration test"

  self_service_category {
   id = jamfpro_category.bravo.id
   name = jamfpro_category.bravo.name
   display_in = true
   feature_in = false
  }
 }
}
