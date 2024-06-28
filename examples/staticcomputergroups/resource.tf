resource "jamfpro_static_computer_group" "jamfpro_static_computer_group_001" {
  name = "Example Static Computer Group"


  // Optional Block
   site_id = 1


  # Optional: Specify computers for static groups
  assignments {
    computer_ids = [16, 20, 21]
  }
}