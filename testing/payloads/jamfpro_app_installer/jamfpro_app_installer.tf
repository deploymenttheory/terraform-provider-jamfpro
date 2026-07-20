resource "jamfpro_app_installer" "jamfpro_app_installer_test_001" {
  app_title_name  = "Jamf Connect"
  name            = "Jamf Connect"
  enabled         = true
  deployment_type = "INSTALL_AUTOMATICALLY"
  update_behavior = "AUTOMATIC"
  category_id     = "-1"
  site_id         = "-1"
  smart_group_id  = "1"

  install_predefined_config_profiles = false
  trigger_admin_notifications        = false

  notification_settings {
    notification_message  = "A new update is available"
    notification_interval = 1
    deadline_message      = "Update deadline approaching"
    deadline              = 1
    quit_delay            = 1
    complete_message      = "Update completed successfully"
    relaunch              = true
    suppress              = false
  }

  #   self_service_settings {
  #   include_in_featured_category   = true
  #   include_in_compliance_category = false
  #   force_view_description         = true
  #   description                    = "This is an example  cheese app deployment"
  # }


}

resource "jamfpro_app_installer" "jamfpro_app_installer_test_002" {
  app_title_name  = "Jamf Connect"
  name            = "Jamf Connect"
  enabled         = true
  deployment_type = "SELF_SERVICE"
  update_behavior = "AUTOMATIC"
  category_id     = "-1"
  site_id         = "-1"
  smart_group_id  = "1"

  install_predefined_config_profiles = false
  trigger_admin_notifications        = false

  notification_settings {
    notification_message  = "A new update is available"
    notification_interval = 1
    deadline_message      = "Update deadline approaching"
    deadline              = 1
    quit_delay            = 1
    complete_message      = "Update completed successfully"
    relaunch              = true
    suppress              = false
  }

  self_service_settings {
    include_in_featured_category   = true
    include_in_compliance_category = false
    force_view_description         = true
    // Regression test for issue #1145 - heredoc strings in HCL always
    // include a trailing newline before EOT, but the API strips it
    // server-side.
    description = <<-EOT
      This is an example cheese app deployment used to verify
      no drift is reported after apply due to the heredoc trailing newline.
    EOT
  }


}
