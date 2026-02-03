variable "version_number" {
  description = "The version number to include in the name and install button text."
  type        = string
  default     = "1.0"
}

// example hcl generated plist with 1 level of nesting
resource "jamfpro_macos_configuration_profile_plist_generator" "jamfpro_macos_configuration_profile_plist_generator_003" {
  name                = "tf-localtest-generator-accessibility-seeing-${var.plist_version_number}"
  description         = "Base Level Accessibility settings for vision"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  user_removable      = false
  level               = "System"

  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_names    = ["Marketing", "Finance"]
    }

    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([78, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_names    = ["IT", "Engineering"]
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

  payloads {
    payload_root {
      payload_description_root        = "Base Level Accessibility settings for vision"
      payload_enabled_root            = true
      payload_organization_root       = "Deployment Theory"
      payload_removal_disallowed_root = false
      payload_scope_root              = "System"
      payload_type_root               = "Configuration"
      payload_version_root            = var.plist_version_number
    }

    payload_content {
      configuration {
        key   = "closeViewFarPoint"
        value = 2
      }
      configuration {
        key   = "closeViewHotkeysEnabled"
        value = true
      }
      configuration {
        key   = "closeViewNearPoint"
        value = 10
      }
      configuration {
        key   = "closeViewScrollWheelToggle"
        value = true
      }
      configuration {
        key   = "closeViewShowPreview"
        value = true
      }
      configuration {
        key   = "closeViewSmoothImages"
        value = true
      }
      configuration {
        key   = "contrast"
        value = 0
      }
      configuration {
        key   = "flashScreen"
        value = false
      }
      configuration {
        key   = "grayscale"
        value = false
      }
      configuration {
        key   = "mouseDriver"
        value = false
      }
      configuration {
        key   = "mouseDriverCursorSize"
        value = 3
      }
      configuration {
        key   = "mouseDriverIgnoreTrackpad"
        value = false
      }
      configuration {
        key   = "mouseDriverInitialDelay"
        value = 1.0
      }
      configuration {
        key   = "mouseDriverMaxSpeed"
        value = 3
      }
      configuration {
        key   = "slowKey"
        value = false
      }
      configuration {
        key   = "slowKeyBeepOn"
        value = false
      }
      configuration {
        key   = "slowKeyDelay"
        value = 0
      }
      configuration {
        key   = "stereoAsMono"
        value = false
      }
      configuration {
        key   = "stickyKey"
        value = false
      }
      configuration {
        key   = "stickyKeyBeepOnModifier"
        value = false
      }
      configuration {
        key   = "stickyKeyShowWindow"
        value = false
      }
      configuration {
        key   = "voiceOverOnOffKey"
        value = true
      }
      configuration {
        key   = "whiteOnBlack"
        value = false
      }

      payload_description  = ""
      payload_display_name = "Accessibility"
      payload_enabled      = true
      payload_organization = "Deployment Theory"
      payload_type         = "com.apple.universalaccess"
      payload_version      = var.plist_version_number
      payload_scope        = "System"
    }
  }
}
