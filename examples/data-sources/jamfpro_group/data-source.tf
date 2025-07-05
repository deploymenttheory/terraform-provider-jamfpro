
# Example usage of the jamfpro_group data source with all supported lookup methods

# Lookup by name and group_type (COMPUTER)
data "jamfpro_group" "by_computer_name" {
  name       = "All Managed Clients" # Replace with a real computer group name
  group_type = "COMPUTER"
}

# Lookup by name and group_type (MOBILE)
data "jamfpro_group" "by_mobile_name" {
  name       = "All Managed iPads" # Replace with a real mobile group name
  group_type = "MOBILE"
}

# Lookup by group_platform_id (UUID string)
data "jamfpro_group" "by_platform_id" {
  group_platform_id = "4a36a1fe-e45a-430d-a966-a4d3ac993577" # Replace with a real group platform UUID
}

# Lookup by group_jamfpro_id and group_type (COMPUTER)
data "jamfpro_group" "by_computer_jamfpro_id" {
  group_jamfpro_id = "1" # Replace with a real Jamf Pro ID (number as string)
  group_type       = "COMPUTER"
}

# Lookup by group_jamfpro_id and group_type (MOBILE)
data "jamfpro_group" "by_mobile_jamfpro_id" {
  group_jamfpro_id = "1" # Replace with a real Jamf Pro ID (number as string)
  group_type       = "MOBILE"
}

output "group_by_computer_name" {
  value = {
    id               = data.jamfpro_group.by_computer_name.group_jamfpro_id
    name             = data.jamfpro_group.by_computer_name.name
    group_type       = data.jamfpro_group.by_computer_name.group_type
    smart            = data.jamfpro_group.by_computer_name.smart
    membership_count = data.jamfpro_group.by_computer_name.membership_count
  }
}

output "group_by_mobile_name" {
  value = {
    id               = data.jamfpro_group.by_mobile_name.group_jamfpro_id
    name             = data.jamfpro_group.by_mobile_name.name
    group_type       = data.jamfpro_group.by_mobile_name.group_type
    smart            = data.jamfpro_group.by_mobile_name.smart
    membership_count = data.jamfpro_group.by_mobile_name.membership_count
  }
}

output "group_by_platform_id" {
  value = {
    id               = data.jamfpro_group.by_platform_id.group_jamfpro_id
    name             = data.jamfpro_group.by_platform_id.name
    group_type       = data.jamfpro_group.by_platform_id.group_type
    smart            = data.jamfpro_group.by_platform_id.smart
    membership_count = data.jamfpro_group.by_platform_id.membership_count
  }
}

output "group_by_computer_jamfpro_id" {
  value = {
    id               = data.jamfpro_group.by_computer_jamfpro_id.group_jamfpro_id
    name             = data.jamfpro_group.by_computer_jamfpro_id.name
    group_type       = data.jamfpro_group.by_computer_jamfpro_id.group_type
    smart            = data.jamfpro_group.by_computer_jamfpro_id.smart
    membership_count = data.jamfpro_group.by_computer_jamfpro_id.membership_count
  }
}

output "group_by_mobile_jamfpro_id" {
  value = {
    id               = data.jamfpro_group.by_mobile_jamfpro_id.group_jamfpro_id
    name             = data.jamfpro_group.by_mobile_jamfpro_id.name
    group_type       = data.jamfpro_group.by_mobile_jamfpro_id.group_type
    smart            = data.jamfpro_group.by_mobile_jamfpro_id.smart
    membership_count = data.jamfpro_group.by_mobile_jamfpro_id.membership_count
  }
}
