resource "jamfpro_policy" "jamfpro_printer_policy_001" {
  name                          = "tf-localtest-printer-policy-001"
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

  payloads {
    printers {
      id           = jamfpro_printer.jamfpro_printer_001.id
      name         = jamfpro_printer.jamfpro_printer_001.name // requires both id and name for req to work
      action       = "install"
      make_default = true
    }
    files_processes {
      search_by_path         = "/Applications/SomeApp.app"
      delete_file            = true
      locate_file            = "SomeFile.txt"
      update_locate_database = false
      spotlight_search       = "SomeApp"
      search_for_process     = "SomeProcess"
      kill_process           = true
      run_command            = "echo 'Hello, World!'"
    }
  }
}






