resource "jamfpro_macos_configuration_profile_plist_generator" "test" {
 name = "tf-sets-plist-20260919"
 description = "Temporary unscoped set migration test"
 distribution_method = "Make Available in Self Service"
 redeploy_on_update = "Newly Assigned"
 site_id = -1
 category_id = -1
 scope { all_computers = false }
 self_service {
  install_button_text = "Install"
  notification_subject = ""
  self_service_display_name = "Temporary profile"
  self_service_description = "Temporary unscoped set migration test"
 }
 payloads {
  payload_description_header = "Temporary unscoped set migration test"
  payload_display_name_header = "tf-sets-plist-20260919"
  payload_enabled_header = true
  payload_organization_header = "Terraform test"
  payload_type_header = "Configuration"
  payload_version_header = 1
  payload_scope_header = "System"
  payload_content {
   payload_enabled = true
   payload_organization = "Terraform test"
   payload_type = "com.example.terraform.set-test"
   payload_version = 1
   setting {
    key = "Zulu"
    value = "last"
   }
   setting {
    key = "Alpha"
    value = "first"
   }
  }
 }
}
