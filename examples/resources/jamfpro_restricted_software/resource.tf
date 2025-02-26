resource "jamfpro_restricted_software" "restricted_software_high_sierra" {
  name                     = "tf-localtest-restrict-high-sierra"
  process_name             = "Install macOS High Sierra.app"
  match_exact_process_name = true
  send_notification        = true
  kill_process             = true
  delete_executable        = true
  display_message          = "This software is restricted and will be terminated."

  # optional
  site_id {
    id = -1
  }

  scope {
    all_computers      = false
    computer_ids       = ([23, 22])
    computer_group_ids = ([13, 12])
    building_ids       = ([1536, 1534])
    department_ids     = ([37501, 37503])
    exclusions {
      computer_ids                         = [14, 15]
      computer_group_ids                   = ([13, 12])
      building_ids                         = ([1536, 1534])
      department_ids                       = ([37501, 37503])
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
    }
  }
}