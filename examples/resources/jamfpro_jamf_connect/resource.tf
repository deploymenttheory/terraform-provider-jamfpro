resource "jamfpro_jamf_connect" "example" {
  config_profile_uuid  = "e9224719-906e-4879-b393-f302fa40e89d" # UUID of existing profile
  version              = "2.43.0"                               # Version of Jamf Connect to deploy
  auto_deployment_type = "MINOR_AND_PATCH_UPDATES"
}