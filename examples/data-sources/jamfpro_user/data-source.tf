# Lookup user by username
data "jamfpro_user" "user_by_name" {
  name = "john.doe"
}

# Lookup user by email
data "jamfpro_user" "user_by_email" {
  email = "john.doe@example.com"
}

# Lookup user by ID
data "jamfpro_user" "user_by_id" {
  id = "123"
}

# Use in a static user group
resource "jamfpro_user_group" "engineering_team" {
  name     = "Engineering Team"
  is_smart = false

  assigned_user_ids = [
    tonumber(data.jamfpro_user.user_by_name.id),
  ]
}

# Output user details
output "user_info" {
  value = {
    id            = data.jamfpro_user.user_by_name.id
    name          = data.jamfpro_user.user_by_name.name
    full_name     = data.jamfpro_user.user_by_name.full_name
    email_address = data.jamfpro_user.user_by_name.email_address
    position      = data.jamfpro_user.user_by_name.position
  }
}

# Verify that different lookup methods return the same user
output "matching_ids" {
  value = (
    data.jamfpro_user.user_by_name.id == data.jamfpro_user.user_by_email.id &&
    data.jamfpro_user.user_by_email.id == data.jamfpro_user.user_by_id.id
  )
}
