# Example 1: Look up by ID
data "jamfpro_jamf_connect" "by_id" {
  profile_id = 1
}

# Example 2: Look up by name
data "jamfpro_jamf_connect" "by_name" {
  profile_name = "Jamf Connect License"
}

# Output examples
output "auto_deployment_type" {
  value = data.jamfpro_jamf_connect.by_name.auto_deployment_type
}

output "config_profile_uuid" {
  value = data.jamfpro_jamf_connect.by_name.config_profile_uuid
}

output "id" {
  value = data.jamfpro_jamf_connect.by_name.id
}

output "jamf_connect_version" {
  value = data.jamfpro_jamf_connect.by_name.jamf_connect_version
}

output "scope_description" {
  value = data.jamfpro_jamf_connect.by_name.scope_description
}

output "site_id" {
  value = data.jamfpro_jamf_connect.by_name.site_id
}
