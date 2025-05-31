# Query by name
data "jamfpro_mobile_device_application" "by_name" {
  name = "Keynote"
}

# Query by ID
data "jamfpro_mobile_device_application" "by_id" {
  id = "2"
}

# Example of using the data source outputs
output "app_details_by_name" {
  value = {
    id           = data.jamfpro_mobile_device_application.by_name.id
    name         = data.jamfpro_mobile_device_application.by_name.name
    display_name = data.jamfpro_mobile_device_application.by_name.display_name
    bundle_id    = data.jamfpro_mobile_device_application.by_name.bundle_id
    version      = data.jamfpro_mobile_device_application.by_name.version
  }
}
