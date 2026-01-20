resource "jamfpro_app_installer" "jamfpro_app_installer_test_001" {
  app_title_name  = "010 Editor"
  name            = "010 Editor"
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
  app_title_name  = "010 Editor"
  name            = "010 Editor"
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
    description                    = "This is an example  cheese app deployment"
  }
  

}
