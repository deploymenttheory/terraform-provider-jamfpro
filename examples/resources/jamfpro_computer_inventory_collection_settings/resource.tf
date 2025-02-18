resource "jamfpro_computer_inventory_collection_settings" "example" {
  computer_inventory_collection_preferences {
    monitor_application_usage                          = true
    include_fonts                                      = true
    include_plugins                                    = true
    include_packages                                   = true
    include_software_updates                           = true
    include_software_id                                = true
    include_accounts                                   = true
    calculate_sizes                                    = true
    include_hidden_accounts                            = true
    include_printers                                   = true
    include_services                                   = true
    collect_synced_mobile_device_info                  = true
    update_ldap_info_on_computer_inventory_submissions = true
    monitor_beacons                                    = true
    allow_changing_user_and_location                   = true
    use_unix_user_paths                                = true
    collect_unmanaged_certificates                     = true
  }

  application_paths {
    path = "/Applications/Custom/App1"
  }
  application_paths {
    path = "/Applications/Custom/App2"
  }
  application_paths {
    path = "/Applications/Adobe/Creative Cloud"
  }

  font_paths {
    path = "/Library/Fonts/Custom/Font1"
  }

  font_paths {
    path = "/Library/Fonts/Custom/Font2"
  }

  plugin_paths {
    path = "/Library/Plugins/Custom/plugin1"
  }

  plugin_paths {
    path = "/Library/Plugins/Custom/plugin2"
  }
}
