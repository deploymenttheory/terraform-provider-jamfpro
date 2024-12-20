# Example 1: Look up API role by ID
data "jamfpro_api_role" "by_id" {
  id = "168"
}

# Output example for ID lookup
output "role_privileges_by_id" {
  value = data.jamfpro_api_role.by_id.privileges
}

# Example 2: Look up API role by display name
data "jamfpro_api_role" "by_name" {
  display_name = "tf-all-jss-actions-permissions-11.12"
}

# Output examples for name lookup
output "role_id_by_name" {
  value = data.jamfpro_api_role.by_name.id
}

output "role_privileges_by_name" {
  value = data.jamfpro_api_role.by_name.privileges
}

# Example 3: Using the data source in another resource
resource "jamfpro_some_resource" "example" {
  name       = "Example Resource"
  role_id    = data.jamfpro_api_role.by_name.id
  privileges = data.jamfpro_api_role.by_name.privileges
}

# Example 4: Using with variables
variable "role_name" {
  type        = string
  description = "The display name of the API role to look up"
  default     = "Read Only Admin"
}

data "jamfpro_api_role" "dynamic" {
  display_name = var.role_name
}