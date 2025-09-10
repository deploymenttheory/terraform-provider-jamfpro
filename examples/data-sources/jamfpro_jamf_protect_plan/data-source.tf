# Example usage of the jamfpro_jamf_protect_plan data source

data "jamfpro_jamf_protect_plan" "by_id" {
  id = "example-plan-id" # Replace with the actual plan ID
}

data "jamfpro_jamf_protect_plan" "by_name" {
  name = "example-plan-name" # Replace with the actual plan name
}

output "jamfpro_jamf_protect_plan_by_id" {
  value = {
    id                = data.jamfpro_jamf_protect_plan.by_id.id
    uuid              = data.jamfpro_jamf_protect_plan.by_id.uuid
    name              = data.jamfpro_jamf_protect_plan.by_id.name
    description       = data.jamfpro_jamf_protect_plan.by_id.description
    profile_id        = data.jamfpro_jamf_protect_plan.by_id.profile_id
    profile_name      = data.jamfpro_jamf_protect_plan.by_id.profile_name
    scope_description = data.jamfpro_jamf_protect_plan.by_id.scope_description
  }
}

output "jamfpro_jamf_protect_plan_by_name" {
  value = {
    id                = data.jamfpro_jamf_protect_plan.by_name.id
    uuid              = data.jamfpro_jamf_protect_plan.by_name.uuid
    name              = data.jamfpro_jamf_protect_plan.by_name.name
    description       = data.jamfpro_jamf_protect_plan.by_name.description
    profile_id        = data.jamfpro_jamf_protect_plan.by_name.profile_id
    profile_name      = data.jamfpro_jamf_protect_plan.by_name.profile_name
    scope_description = data.jamfpro_jamf_protect_plan.by_name.scope_description
  }
}
