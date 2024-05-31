resource "jamfpro_static_computer_group" "jamfpro_static_computer_group_001" {
  name = "Example Static Computer Group"

  site {
    id   = 123  # Replace with the actual site ID
    name = "Site Name"  # Replace with the actual site name
  }

  # Optional: Specify computers for static groups
  assignments {
    computer_ids = [16, 20, 21]
  }
}