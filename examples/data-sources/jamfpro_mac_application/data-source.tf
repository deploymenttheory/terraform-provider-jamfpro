data "jamfpro_mac_application" "by_name" {
  name = "Slack for Desktop"
}

data "jamfpro_mac_application" "by_id" {
  id = 17
}

output "app_id" {
  value = data.jamfpro_mac_application.by_name.id
}

output "existing_app_version" {
  value = data.jamfpro_mac_application.by_id.version
}
