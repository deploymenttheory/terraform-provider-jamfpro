resource "jamfpro_policy" "jamfpro_policy_script_example" {
  name                          = "acc-test-policy-script-with-self-service-enabled"
  enabled                       = false
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "EVENT" // "USER_INITIATED" for self service trigger , "EVENT" for an event trigger
  frequency                     = "Once per computer"
  retry_event                   = "none"
  retry_attempts                = -1
  notify_on_each_failed_retry   = false
  target_drive                  = "/"
  offline                       = false
  category_id                   = -1
  site_id                       = -1

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = ""
    install_button_text             = "Install"
    reinstall_button_text           = "Reinstall"
    self_service_description        = ""
    force_users_to_view_description = false
    feature_on_main_page            = false

    self_service_category {
      id         = jamfpro_category.jamfpro_category_001.id
      display_in = true
      feature_in = false
    }

    notification         = false
    notification_type    = "Self Service"
    notification_subject = "Install Script"
    notification_message = "This is a message for the Script install"
  }

  payloads {
    scripts {
      id          = jamfpro_script.jamfpro_script_001.id
      priority    = "After"
      parameter4  = "param_value_4"
      parameter5  = "param_value_5"
      parameter6  = "param_value_6"
      parameter7  = "param_value_7"
      parameter8  = "param_value_8"
      parameter9  = "param_value_9"
      parameter10 = "param_value_10"
      parameter11 = "param_value_11"

    }
  }

}

resource "jamfpro_script" "jamfpro_script_001" {
  name            = "acc-test-add-or-remove-group-membership-v4.0"
  script_contents = file("${path.module}/support_files/scripts/Add or Remove Group Membership.zsh")
  category_id     = "5"
  os_requirements = "13"
  priority        = "BEFORE"
  info            = "Adds target user or group to specified group membership, or removes said membership."
  notes           = "Jamf Pro script parameters: 4 -> 7"
  parameter4      = "100"           // targetID
  parameter5      = "group"         // Target Type - Must be either "user" or "group"
  parameter6      = "someGroupName" // targetMembership
  parameter7      = "add"           // Script Action - Must be either "add" or "remove"
}

resource "jamfpro_category" "jamfpro_category_001" {
  name     = "acc-test-category-01"
  priority = 1
}