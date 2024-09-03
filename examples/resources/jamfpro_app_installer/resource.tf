resource "jamfpro_app_installer" "example" {
  name                = "Example App Deployment"
  enabled             = true
  deployment_type     = "INSTALL_AUTOMATICALLY"
  update_behavior     = "AUTOMATIC"
  category_id         = "-1"
  site_id             = "-1"
  smart_group_id      = "1"

  install_predefined_config_profiles = true
  title_available_in_ais             = true
  trigger_admin_notifications        = true

  notification_settings {
    notification_message  = "A new update is available"
    notification_interval = 24
    deadline_message      = "Update deadline approaching"
    deadline              = 72
    quit_delay            = 30
    complete_message      = "Update completed successfully"
    relaunch              = true
    suppress              = "NONE"
  }

  self_service_settings {
    include_in_featured_category    = true
    include_in_compliance_category  = false
    force_view_description          = true
    description                     = "This is an example app deployment"

    categories {
      id       = "2"
      featured = true
    }
    categories {
      id       = "3"
      featured = false
    }
  }
}