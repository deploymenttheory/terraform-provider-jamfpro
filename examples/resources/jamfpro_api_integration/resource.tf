resource "jamfpro_api_integration" "jamfpro_api_integration_001" {
  display_name                  = "tf-localtest-api-integration-001"
  enabled                       = true
  access_token_lifetime_seconds = 7200
  authorization_scopes          = [jamfpro_api_role.jamfpro_api_role_001.display_name]
}

resource "jamfpro_api_integration" "jamfpro_api_integration_002" {
  display_name                  = "tf-localtest-api-integration-002"
  enabled                       = true
  access_token_lifetime_seconds = 6000
  authorization_scopes          = [jamfpro_api_role.jamfpro_api_role_002.display_name]
}