
resource "jamfpro_computer_checkin" "jamfpro_computer_checkin" {
  check_in_frequency                        = 15
  create_startup_script                     = true // 
  log_startup_event                         = false //  requires create_startup_script
  ensure_ssh_is_enabled                     = false //  requires create_startup_script
  check_for_policies_at_startup             = false //  requires create_startup_script
  create_login_logout_hooks                 = true // 
  log_username                              = false //  requires create_login_logout_hooks 
  check_for_policies_at_login_logout        = false //  requires create_login_logout_hooks 
  
}