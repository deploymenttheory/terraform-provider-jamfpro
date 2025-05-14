data "jamfpro_sso_failover" "current" {}

output "current_failover_url" {
  value     = data.jamfpro_sso_failover.current.failover_url
  sensitive = true
}

output "last_generation_time" {
  value = data.jamfpro_sso_failover.current.generation_time
}
