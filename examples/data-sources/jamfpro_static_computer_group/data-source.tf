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
      id          = data.jamfpro_static_computer_group.by_id.id
      name        = data.jamfpro_static_computer_group.by_id.name
      description = data.jamfpro_static_computer_group.by_id.description
      site_id     = data.jamfpro_static_computer_group.by_id.site_id
    }
    by_name = {
      id          = data.jamfpro_static_computer_group.by_name.id
      name        = data.jamfpro_static_computer_group.by_name.name
      description = data.jamfpro_static_computer_group.by_name.description
      site_id     = data.jamfpro_static_computer_group.by_name.site_id
    }
  }
}
