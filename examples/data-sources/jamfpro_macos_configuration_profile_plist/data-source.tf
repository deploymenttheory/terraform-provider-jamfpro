# Test data source by ID
data "jamfpro_macos_configuration_profile_plist" "test_by_id" {
  id = jamfpro_macos_configuration_profile_plist.privacy_preferences_policy_control.id
}

# Test data source by name
data "jamfpro_macos_configuration_profile_plist" "test_by_name" {
  name = jamfpro_macos_configuration_profile_plist.privacy_preferences_policy_control.name
}

# Outputs for ID-based lookup
output "profile_by_id" {
  value = {
    id                 = data.jamfpro_macos_configuration_profile_plist.test_by_id.id
    name               = data.jamfpro_macos_configuration_profile_plist.test_by_id.name
    description        = data.jamfpro_macos_configuration_profile_plist.test_by_id.description
    uuid               = data.jamfpro_macos_configuration_profile_plist.test_by_id.uuid
    site_id            = data.jamfpro_macos_configuration_profile_plist.test_by_id.site_id
    category_id        = data.jamfpro_macos_configuration_profile_plist.test_by_id.category_id
    distribution_method = data.jamfpro_macos_configuration_profile_plist.test_by_id.distribution_method
    user_removable     = data.jamfpro_macos_configuration_profile_plist.test_by_id.user_removable
    level              = data.jamfpro_macos_configuration_profile_plist.test_by_id.level
    redeploy_on_update = data.jamfpro_macos_configuration_profile_plist.test_by_id.redeploy_on_update
    payloads           = data.jamfpro_macos_configuration_profile_plist.test_by_id.payloads
    scope              = data.jamfpro_macos_configuration_profile_plist.test_by_id.scope
    self_service       = length(data.jamfpro_macos_configuration_profile_plist.test_by_id.self_service) > 0 ? data.jamfpro_macos_configuration_profile_plist.test_by_id.self_service : null
  }
}

# Outputs for name-based lookup
output "profile_by_name" {
  value = {
    id                 = data.jamfpro_macos_configuration_profile_plist.test_by_name.id
    name               = data.jamfpro_macos_configuration_profile_plist.test_by_name.name
    description        = data.jamfpro_macos_configuration_profile_plist.test_by_name.description
    uuid               = data.jamfpro_macos_configuration_profile_plist.test_by_name.uuid
    site_id            = data.jamfpro_macos_configuration_profile_plist.test_by_name.site_id
    category_id        = data.jamfpro_macos_configuration_profile_plist.test_by_name.category_id
    distribution_method = data.jamfpro_macos_configuration_profile_plist.test_by_name.distribution_method
    user_removable     = data.jamfpro_macos_configuration_profile_plist.test_by_name.user_removable
    level              = data.jamfpro_macos_configuration_profile_plist.test_by_name.level
    redeploy_on_update = data.jamfpro_macos_configuration_profile_plist.test_by_name.redeploy_on_update
    payloads           = data.jamfpro_macos_configuration_profile_plist.test_by_name.payloads
    scope              = data.jamfpro_macos_configuration_profile_plist.test_by_name.scope
    self_service       = length(data.jamfpro_macos_configuration_profile_plist.test_by_name.self_service) > 0 ? data.jamfpro_macos_configuration_profile_plist.test_by_name.self_service : null
  }
}

