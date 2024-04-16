resource "jamfpro_mobile_device_configuration_profile" "jamfpro_mobile_device_configuration_profile_001" {
  name              = "your-mobile_device_configuration_profile-name"
  description       = "An example mobile device configuration profile."
  deployment_method = "Install Automatically"
  level             = "Device Level"
  payloads          = file("${path.module}/path/to/your.mobileconfig")

  scope {
    all_mobile_devices = false
    all_jss_users      = false
    mobile_device_groups {
      id = 1
    }
    jss_users {
      id = 3
    }
    jss_user_groups {
      id = 4
    }
    buildings {
      id = 1348
    }
    departments {
      id = 37287
    }
    limitations {
      ibeacons {
        id = 5
      }
      network_segments {
        id = 4
      }
      users {
        id = 3
      }
      user_groups {
        id = 4
      }
    }
    exclusions {
      mobile_device_groups {
        id = 1
      }
      jss_users {
        id = 3
      }
      jss_user_groups {
        id = 4
      }
      buildings {
        id = 1348
      }
      departments {
        id = 37287
      }
      ibeacons {
        id = 5
      }
      network_segments {
        id = 4
      }
      users {
        id = 3
      }
      user_groups {
        id = 4
      }
    }
  }
}
