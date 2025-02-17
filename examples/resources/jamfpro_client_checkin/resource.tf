resource "jamfpro_client_checkin" "jamfpro_client_checkin" {
  check_in_frequency                  = 30 // Valid values: 5, 15, 30, 60
  create_startup_script               = true
  startup_log                         = true // requires create_startup_script
  startup_ssh                         = true // requires create_startup_script
  startup_policies                    = true // requires create_startup_script
  create_hooks                        = true
  hook_log                            = true // requires create_hooks
  hook_policies                       = true // requires create_hooks
  enable_local_configuration_profiles = true
  allow_network_state_change_triggers = true
}