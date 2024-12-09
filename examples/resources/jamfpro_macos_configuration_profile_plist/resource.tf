variable "version_number" {
  description = "The version number to include in the name and install button text."
  type        = string
  default     = "v1.0"
}

// Minimum viable example of creating a macOS configuration profile in Jamf Pro for automatic installation using a plist source file

resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_064" {
  name                = "your-name-${var.version_number}"
  description         = "An example mobile device configuration profile."
  level               = "System"
  distribution_method = "Install Automatically" // "Make Available in Self Service", "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/path/to/your/file.mobileconfig")
  payload_validate    = true
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }

}

// Example of creating a macOS configuration profile in Jamf Pro for self service using a plist source file
resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_001" {
  name                = "your-name-${var.version_number}"
  description         = "An example mobile device configuration profile."
  level               = "User"                           // "User", "Device"
  distribution_method = "Make Available in Self Service" // "Make Available in Self Service", "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = file("${path.module}/path/to/your/file.mobileconfig")
  payload_validate    = true
  user_removable      = false

  // Optional Block

  site_id = 967

  // Optional Block
  category_id = 5
  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    // Optional Block
    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
    }

    // Optional Block
    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([78, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
      ibeacon_ids                          = [3, 4]
    }
  }

  self_service {
    install_button_text             = "Install - ${var.version_number}"
    self_service_description        = "This is the self service description"
    force_users_to_view_description = true
    feature_on_main_page            = true
    notification                    = true
    notification_subject            = "New Profile Available"
    notification_message            = "A new profile is available for installation."

    self_service_category {
      id         = 10
      display_in = true
      feature_in = true
    }

    self_service_category {
      id         = 5
      display_in = false
      feature_in = true
    }
  }
}