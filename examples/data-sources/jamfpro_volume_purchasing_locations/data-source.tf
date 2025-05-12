
data "jamfpro_volume_purchasing_locations" "example" {
  id = "1"
}

output "vpp_location_details" {
  value = {
    id                                        = data.jamfpro_volume_purchasing_locations.example.id
    name                                      = data.jamfpro_volume_purchasing_locations.example.name
    apple_id                                  = data.jamfpro_volume_purchasing_locations.example.apple_id
    organization_name                         = data.jamfpro_volume_purchasing_locations.example.organization_name
    token_expiration                          = data.jamfpro_volume_purchasing_locations.example.token_expiration
    country_code                              = data.jamfpro_volume_purchasing_locations.example.country_code
    location_name                             = data.jamfpro_volume_purchasing_locations.example.location_name
    client_context_mismatch                   = data.jamfpro_volume_purchasing_locations.example.client_context_mismatch
    automatically_populate_purchased_content  = data.jamfpro_volume_purchasing_locations.example.automatically_populate_purchased_content
    send_notification_when_no_longer_assigned = data.jamfpro_volume_purchasing_locations.example.send_notification_when_no_longer_assigned
    auto_register_managed_users               = data.jamfpro_volume_purchasing_locations.example.auto_register_managed_users
    site_id                                   = data.jamfpro_volume_purchasing_locations.example.site_id
    last_sync_time                            = data.jamfpro_volume_purchasing_locations.example.last_sync_time
    total_purchased_licenses                  = data.jamfpro_volume_purchasing_locations.example.total_purchased_licenses
    total_used_licenses                       = data.jamfpro_volume_purchasing_locations.example.total_used_licenses

    content = [
      for item in data.jamfpro_volume_purchasing_locations.example.content : {
        name                   = item.name
        license_count_total    = item.license_count_total
        license_count_in_use   = item.license_count_in_use
        license_count_reported = item.license_count_reported
        icon_url               = item.icon_url
        device_types           = item.device_types
        content_type           = item.content_type
        pricing_param          = item.pricing_param
        adam_id                = item.adam_id
      }
    ]
  }
}
