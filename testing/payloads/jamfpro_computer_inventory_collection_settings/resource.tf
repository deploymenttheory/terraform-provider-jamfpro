resource "jamfpro_computer_inventory_collection_settings" "jamfpro_computer_inventory_collection_settings" {
  computer_inventory_collection_preferences {
    monitor_application_usage                          = false
    include_packages                                   = true
    include_software_updates                           = false
    include_software_id                                = true
    include_accounts                                   = true
    calculate_sizes                                    = false
    include_hidden_accounts                            = false
    include_printers                                   = true
    include_services                                   = true
    collect_synced_mobile_device_info                  = false
    update_ldap_info_on_computer_inventory_submissions = true
    monitor_beacons                                    = false
    allow_changing_user_and_location                   = true
    use_unix_user_paths                                = true
    collect_unmanaged_certificates                     = true
  }
}
