# Look up a single user by Jamf Pro ID (returns the full user object).
data "jamfpro_user" "by_id" {
  user_id = "1"
}

# Look up a single user by name (username).
data "jamfpro_user" "by_name" {
  name = "jappleseed"
}

# Look up a single user by email address.
data "jamfpro_user" "by_email" {
  email = "jappleseed@example.com"
}

# List all users (returns only the id and name of each user).
data "jamfpro_user" "all" {
  list_all = true
}

output "jamfpro_user_by_id_full_name" {
  value = data.jamfpro_user.by_id.items[0].full_name
}

output "jamfpro_user_count" {
  value = length(data.jamfpro_user.all.items)
}
