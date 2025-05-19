resource "jamfpro_static_mobile_device_group" "jamfpro_static_mobile_device_group_001" {
  name = "Example Mobile Device Group"


  # Optional Block
  site_id = 1

  # Optional: Specify computers for static groups
  assigned_mobile_device_ids = [1, 2, 3]
}
