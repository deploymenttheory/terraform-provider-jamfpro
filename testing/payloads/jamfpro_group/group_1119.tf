data "jamfpro_group" "by_computer_name" {
  name       = "All Managed Clients"
  group_type = "COMPUTER"
}
