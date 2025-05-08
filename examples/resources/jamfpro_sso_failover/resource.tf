resource "jamfpro_sso_failover" "current" {
  regenerate = false # Set to true when you want to regenerate the URL
}

output "failover_url" {
  value     = jamfpro_sso_failover.current.failover_url
  sensitive = true
}
