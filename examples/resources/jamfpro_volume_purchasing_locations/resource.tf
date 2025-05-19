resource "jamfpro_volume_purchasing_locations" "example" {
  name                                      = "Company Apple Business Manager"
  service_token                             = "eyJleHBEYXRlIjo..." # Your base64 encoded service token
  automatically_populate_purchased_content  = true
  send_notification_when_no_longer_assigned = false
  auto_register_managed_users               = true
  site_id                                   = "-1" # Default site
}
