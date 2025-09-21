# Example usage of the jamfpro_jamf_protect_plans data source

data "jamfpro_jamf_protect_plans" "all" {}

output "jamfpro_jamf_protect_plans" {
  value = data.jamfpro_jamf_protect_plans.all.plans
}
