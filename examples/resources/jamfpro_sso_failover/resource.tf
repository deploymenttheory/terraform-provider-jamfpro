resource "jamfpro_sso_failover" "current" {}

output "failover_url" {
  value     = jamfpro_sso_failover.current.failover_url
  sensitive = true
}
