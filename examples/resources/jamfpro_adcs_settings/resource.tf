resource "jamfpro_adcs_settings" "inbound_example" {
  display_name  = "tf-example-inbound-adcs"
  ca_name       = "Contoso Issuing CA"
  fqdn          = "connector.contoso.corp"
  adcs_url      = "https://connector.contoso.corp/certsrv"
  api_client_id = "c1bcec08-5f34-40fa-af52-9d3413ac916d"

  revocation_enabled = true
  outbound           = false

  server_certificate_filename = "server-certificate.pfx"
  server_certificate_data     = filebase64("./secrets/server-certificate.pfx")
  server_certificate_password = "serverCertPassw0rd!"

  client_certificate_filename = "client-certificate.pfx"
  client_certificate_data     = filebase64("./secrets/client-certificate.pfx")
  client_certificate_password = "clientCertPassw0rd!"
}

resource "jamfpro_adcs_settings" "outbound_example" {
  display_name  = "tf-example-outbound-adcs"
  ca_name       = "Contoso Issuing CA"
  fqdn          = "connector.contoso.corp"
  api_client_id = "c1bcec08-5f34-40fa-af52-9d3413ac916d"

  revocation_enabled = true
  outbound           = true
}
