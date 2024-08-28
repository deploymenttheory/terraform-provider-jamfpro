// Example of creating a mobie device configuration profile in Jamf Pro for self service using a plist source file
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_001" {
  name               = "your-mobile_device_configuration_profile-name"
  description        = "An example mobile device configuration profile."
  deployment_method  = "Install Automatically"
  level              = "Device Level"
  redeploy_on_update = "Newly Assigned"
  payloads           = file("${path.module}/path/to/your.mobileconfig")

  // Optional Block

  site_id = 967

  // Optional Block
  category_id = 5
  scope {
    all_mobile_devices = true
    all_jss_users      = false

    mobile_device_ids       = [101, 102, 103]
    mobile_device_group_ids = [201, 202]
    building_ids            = [301]
    department_ids          = [401, 402]
    jss_user_ids            = [501, 502]
    jss_user_group_ids      = [601, 602]

    limitations {
      network_segment_ids                  = [701, 702]
      ibeacon_ids                          = [801]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [1001, 1002]
    }

    exclusions {
      mobile_device_ids                    = [1101, 1102]
      mobile_device_group_ids              = [1201]
      building_ids                         = [1301, 1302]
      department_ids                       = [1401]
      network_segment_ids                  = [1501, 1502]
      jss_user_ids                         = [1601, 1602]
      jss_user_group_ids                   = [1701]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [1001, 1002]
      ibeacon_ids                          = [1801]
    }
  }
}

resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_002" {
  name               = "Profile 002"
  description        = "Description for Profile 002"
  level              = "User Level"
  deployment_method  = "Make Available in Self Service"
  redeploy_on_update = "All"
  payloads           = file("path/to/profile_002.mobileconfig")

  scope {
    all_mobile_devices      = false
    all_jss_users           = true
    mobile_device_ids       = [4, 5, 6]
    mobile_device_group_ids = [3, 4]
    building_ids            = [2]
    department_ids          = [2]
    jss_user_ids            = [3, 4]
    jss_user_group_ids      = [2]
  }
}