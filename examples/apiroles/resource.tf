
resource "jamfpro_api_role" "jamfpro_api_role_001" {
  display_name = "tf-localtest-apirole-cicdpipeline1-crud"
  privileges   = []
}

resource "jamfpro_api_role" "jamfpro_api_role_002" {
  display_name = "tf-localtest-apirole-apiroles-crud"
  privileges   = ["Create API Roles", "Update API Roles", "Read API Roles", "Delete API Roles"]
}

resource "jamfpro_api_role" "jamfpro_api_role_003" {
  display_name = "tf-localtest-apirole-apiintegrations-crud"
  privileges   = ["Create API Integrations","Update API Integrations", "Read API Integrations", "Delete API Integrations"]
}