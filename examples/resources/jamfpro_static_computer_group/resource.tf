resource "jamfpro_static_computer_group" "jamfpro_static_computer_group_001" {
  name = "Example Static Computer Group"


  # Optional Site ID
  site_id = "1"

  # Optional Description
  description = "Managed by Terraform."

  # Optional: Specify computers for static groups
  assigned_computer_ids = ["1", "2", "3"]
}
