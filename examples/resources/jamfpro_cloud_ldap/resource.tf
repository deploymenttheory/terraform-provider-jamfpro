resource "jamfpro_cloud_ldap" "example" {
  # Cloud IDP Common Settings
  provider_name = "GOOGLE"
  display_name  = "Google LDAP"

  # Server Configuration
  server_enabled = true
  server_url     = "ldap.google.com"
  domain_name    = "jamf.com"
  port           = 636

  # Connection Settings
  connection_type    = "LDAPS"
  connection_timeout = 60
  search_timeout     = 60

  # Keystore Configuration
  keystore_password   = "supersecretpassword"
  keystore_file_bytes = filebase64("/path/to/keystore.p12")
  keystore_file_name  = "keystore.p12"

  # Advanced Server Settings
  use_wildcards                               = true
  membership_calculation_optimization_enabled = true

  # User Mappings
  user_mappings_object_class_limitation = "ANY_OBJECT_CLASSES"
  user_mappings_object_classes          = "inetOrgPerson"
  user_mappings_search_base             = "ou=Users"
  user_mappings_search_scope            = "ALL_SUBTREES"
  user_mappings_additional_search_base  = ""

  # User Attribute Mappings
  user_mappings_id            = "mail"
  user_mappings_username      = "uid"
  user_mappings_real_name     = "displayName"
  user_mappings_email_address = "mail"
  user_mappings_uuid          = "uid"

  # Optional User Attributes
  user_mappings_department = "departmentNumber"
  user_mappings_position   = "title"
  user_mappings_phone      = ""
  user_mappings_building   = ""
  user_mappings_room       = ""

  # Group Mappings
  group_mappings_object_class_limitation = "ANY_OBJECT_CLASSES"
  group_mappings_object_classes          = "groupOfNames"
  group_mappings_search_base             = "ou=Groups"
  group_mappings_search_scope            = "ALL_SUBTREES"

  # Group Attribute Mappings
  group_mappings_id   = "cn"
  group_mappings_name = "cn"
  group_mappings_uuid = "gidNumber"

  # Membership Mapping
  group_membership_mapping = "memberOf"
}
