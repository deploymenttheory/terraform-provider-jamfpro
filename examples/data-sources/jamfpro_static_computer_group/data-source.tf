# Query group using ID
data "jamfpro_static_computer_group" "by_id" {
  id = "1"
}

# Query using name
data "jamfpro_static_computer_group" "by_name" {
  name = "Development Macs"
}

# Verify both lookups
output "group_verification" {
  value = {
    by_id = {
      id = data.jamfpro_static_computer_group.by_id.id
      name = data.jamfpro_static_computer_group.by_id.name
      is_smart = data.jamfpro_static_computer_group.by_id.is_smart
      site_id = data.jamfpro_static_computer_group.by_id.site_id
      computers = data.jamfpro_static_computer_group.by_id.assigned_computer_ids
    }
    by_name = {
      id = data.jamfpro_static_computer_group.by_name.id
      name = data.jamfpro_static_computer_group.by_name.name
      is_smart = data.jamfpro_static_computer_group.by_name.is_smart
      site_id = data.jamfpro_static_computer_group.by_name.site_id
      computers = data.jamfpro_static_computer_group.by_name.assigned_computer_ids
    }
  }
}