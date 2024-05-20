resource "jamfpro_restricted_software" "restricted_software_001" {
  name                  = "tf-localtest-restrict-high-sierra"
  process_name          = "Install macOS High Sierra.app"
  match_exact_process_name = true
  send_notification     = true
  kill_process          = true
  delete_executable     = true
  display_message       = "This software is restricted and will be terminated."

  # site {
  #   id = 967
  # }

  scope { // scope entities will always be stated asending order. User sort() to sort the list if needed.
    all_computers      = false
    computer_ids       = sort([21, 16])
    computer_group_ids = ([55,78])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    exclusions {
      computer_ids       = [14, 15]
      computer_group_ids = [118 ]
      building_ids       = ([1348, 1349])
      department_ids     = ([37287, 37288])
      jss_user_names      = ["Jane Smith","John Doe"]
    }
  }
}