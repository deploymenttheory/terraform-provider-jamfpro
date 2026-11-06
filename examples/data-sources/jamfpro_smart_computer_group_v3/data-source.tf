data "jamfpro_smart_computer_group_v3" "by_id" {
  id = "123"
}

# Specify either id or name. Name matching is case-insensitive.
data "jamfpro_smart_computer_group_v3" "by_name" {
  name = "Managed Macs running macOS 26"
}

output "group_criteria" {
  value = data.jamfpro_smart_computer_group_v3.by_name.criteria
}
