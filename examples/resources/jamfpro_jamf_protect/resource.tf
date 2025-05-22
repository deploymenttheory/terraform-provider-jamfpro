resource "jamfpro_jamf_protect" "settings" {
  protect_url  = "https://myinstance.protect.jamfcloud.com/graphql"
  client_id    = "supersecretclientid"
  password     = "supersecretpassword"
  auto_install = true

  timeouts {
    create = "90s"
  }
}
