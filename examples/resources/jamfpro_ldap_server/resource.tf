resource "jamfpro_ldap_server" "corporate_ldap" {
  name                = "Corporate LDAP"
  hostname            = "ldap.example.com"
  server_type         = "Active Directory"
  port                = 636
  use_ssl             = true
  authentication_type = "simple"

  # Account credentials
  account {
    distinguished_username = "CN=ServiceAccount,DC=example,DC=com"
    password               = "your-secure-password"
  }

  # Optional connection settings
  open_close_timeout = 15
  search_timeout     = 60
  referral_response  = ""
  use_wildcards      = true

  # User mappings configuration
  user_mappings {
    map_object_class_to_any_or_all = "any"
    object_classes                 = "organizationalPerson"
    search_base                    = "DC=example,DC=com"
    search_scope                   = "All Subtrees"
    map_user_id                    = "uSNCreated"
    map_username                   = "sAMAccountName"
    map_realname                   = "displayName"
    map_email_address              = "mail"
    map_department                 = "department"
    map_building                   = "physicalDeliveryOfficeName"
    map_room                       = "room"
    map_phone                      = "telephoneNumber"
    map_position                   = "title"
    map_user_uuid                  = "objectGUID"
  }

  # Optional: User group mappings
  user_group_mappings {
    map_object_class_to_any_or_all = "any"
    object_classes                 = "group, top"
    search_base                    = "DC=example,DC=com"
    search_scope                   = "All Subtrees"
    map_group_id                   = "name"
    map_group_name                 = "name"
    map_group_uuid                 = "objectGUID"
  }

  # Optional: User group membership mappings
  user_group_membership_mappings {
    user_group_membership_stored_in                        = "user object"
    map_group_membership_to_user_field                     = "memberOf"
    use_dn                                                 = true
    recursive_lookups                                      = true
    group_membership_enabled_when_user_membership_selected = false
    map_user_membership_to_group_field                     = false
    map_user_membership_use_dn                             = false
    map_object_class_to_any_or_all                         = "all"
    object_classes                                         = "group"
    search_scope                                           = "All Subtrees"
    user_group_membership_use_ldap_compare                 = true
    membership_scoping_optimization                        = true
  }
}
