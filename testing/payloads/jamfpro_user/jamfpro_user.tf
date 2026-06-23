// ========================================================================== //
// Users (data source)
// ========================================================================== //
//
// The jamfpro_user data source is read-only and there is no jamfpro_user
// resource to create a user with, so these tests exercise the lookup paths
// against whatever users already exist in the test instance. Single lookups
// (user_id / name / email) are guarded with count so the configuration still
// applies cleanly on an instance that has no users.

# List all users (returns id + name only).
data "jamfpro_user" "all" {
  list_all = true
}

# Look up the first listed user by its Jamf Pro ID (full detail).
data "jamfpro_user" "by_id" {
  count   = length(data.jamfpro_user.all.items) > 0 ? 1 : 0
  user_id = data.jamfpro_user.all.items[0].id
}

# Look up the first listed user by name (full detail).
data "jamfpro_user" "by_name" {
  count = length(data.jamfpro_user.all.items) > 0 ? 1 : 0
  name  = data.jamfpro_user.all.items[0].name
}

# Look up by email, only when the first user actually has an email address.
data "jamfpro_user" "by_email" {
  count = (length(data.jamfpro_user.by_id) > 0 && try(data.jamfpro_user.by_id[0].items[0].email_address, "") != "") ? 1 : 0
  email = data.jamfpro_user.by_id[0].items[0].email_address
}

output "jamfpro_user_all_count" {
  description = "Total number of users returned by list_all."
  value       = length(data.jamfpro_user.all.items)
}

output "jamfpro_user_by_id_full_name" {
  description = "Full name of the user resolved by user_id (null when no users exist)."
  value       = length(data.jamfpro_user.by_id) > 0 ? data.jamfpro_user.by_id[0].items[0].full_name : null
}

output "jamfpro_user_by_name_id" {
  description = "ID of the user resolved by name (null when no users exist)."
  value       = length(data.jamfpro_user.by_name) > 0 ? data.jamfpro_user.by_name[0].items[0].id : null
}
