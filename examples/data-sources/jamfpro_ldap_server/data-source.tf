# Query LDAP server using ID
data "jamfpro_ldap_server" "by_id" {
  id = "1"
}

# Query using name
data "jamfpro_ldap_server" "by_name" {
  name = "Corporate LDAP"
}

# Verify both lookups
output "ldap_server_verification" {
  value = {
    by_id = {
      id                  = data.jamfpro_ldap_server.by_id.id
      name                = data.jamfpro_ldap_server.by_id.name
      hostname            = data.jamfpro_ldap_server.by_id.hostname
      server_type         = data.jamfpro_ldap_server.by_id.server_type
      port                = data.jamfpro_ldap_server.by_id.port
      use_ssl             = data.jamfpro_ldap_server.by_id.use_ssl
      authentication_type = data.jamfpro_ldap_server.by_id.authentication_type

      # Connection settings
      open_close_timeout = data.jamfpro_ldap_server.by_id.open_close_timeout
      search_timeout     = data.jamfpro_ldap_server.by_id.search_timeout
      referral_response  = data.jamfpro_ldap_server.by_id.referral_response
      use_wildcards      = data.jamfpro_ldap_server.by_id.use_wildcards

      # User mappings
      user_mappings = {
        map_object_class_to_any_or_all = data.jamfpro_ldap_server.by_id.user_mappings[0].map_object_class_to_any_or_all
        object_classes                 = data.jamfpro_ldap_server.by_id.user_mappings[0].object_classes
        search_base                    = data.jamfpro_ldap_server.by_id.user_mappings[0].search_base
        search_scope                   = data.jamfpro_ldap_server.by_id.user_mappings[0].search_scope
        map_user_id                    = data.jamfpro_ldap_server.by_id.user_mappings[0].map_user_id
        map_username                   = data.jamfpro_ldap_server.by_id.user_mappings[0].map_username
        map_realname                   = data.jamfpro_ldap_server.by_id.user_mappings[0].map_realname
        map_email_address              = data.jamfpro_ldap_server.by_id.user_mappings[0].map_email_address
      }

      # User group mappings
      user_group_mappings = {
        map_object_class_to_any_or_all = data.jamfpro_ldap_server.by_id.user_group_mappings[0].map_object_class_to_any_or_all
        object_classes                 = data.jamfpro_ldap_server.by_id.user_group_mappings[0].object_classes
        search_base                    = data.jamfpro_ldap_server.by_id.user_group_mappings[0].search_base
        search_scope                   = data.jamfpro_ldap_server.by_id.user_group_mappings[0].search_scope
        map_group_id                   = data.jamfpro_ldap_server.by_id.user_group_mappings[0].map_group_id
        map_group_name                 = data.jamfpro_ldap_server.by_id.user_group_mappings[0].map_group_name
      }
    }
    by_name = {
      id                  = data.jamfpro_ldap_server.by_name.id
      name                = data.jamfpro_ldap_server.by_name.name
      hostname            = data.jamfpro_ldap_server.by_name.hostname
      server_type         = data.jamfpro_ldap_server.by_name.server_type
      port                = data.jamfpro_ldap_server.by_name.port
      use_ssl             = data.jamfpro_ldap_server.by_name.use_ssl
      authentication_type = data.jamfpro_ldap_server.by_name.authentication_type

      # Connection settings
      open_close_timeout = data.jamfpro_ldap_server.by_name.open_close_timeout
      search_timeout     = data.jamfpro_ldap_server.by_name.search_timeout
      referral_response  = data.jamfpro_ldap_server.by_name.referral_response
      use_wildcards      = data.jamfpro_ldap_server.by_name.use_wildcards

      # User mappings
      user_mappings = {
        map_object_class_to_any_or_all = data.jamfpro_ldap_server.by_name.user_mappings[0].map_object_class_to_any_or_all
        object_classes                 = data.jamfpro_ldap_server.by_name.user_mappings[0].object_classes
        search_base                    = data.jamfpro_ldap_server.by_name.user_mappings[0].search_base
        search_scope                   = data.jamfpro_ldap_server.by_name.user_mappings[0].search_scope
        map_user_id                    = data.jamfpro_ldap_server.by_name.user_mappings[0].map_user_id
        map_username                   = data.jamfpro_ldap_server.by_name.user_mappings[0].map_username
        map_realname                   = data.jamfpro_ldap_server.by_name.user_mappings[0].map_realname
        map_email_address              = data.jamfpro_ldap_server.by_name.user_mappings[0].map_email_address
      }

      # User group mappings
      user_group_mappings = {
        map_object_class_to_any_or_all = data.jamfpro_ldap_server.by_name.user_group_mappings[0].map_object_class_to_any_or_all
        object_classes                 = data.jamfpro_ldap_server.by_name.user_group_mappings[0].object_classes
        search_base                    = data.jamfpro_ldap_server.by_name.user_group_mappings[0].search_base
        search_scope                   = data.jamfpro_ldap_server.by_name.user_group_mappings[0].search_scope
        map_group_id                   = data.jamfpro_ldap_server.by_name.user_group_mappings[0].map_group_id
        map_group_name                 = data.jamfpro_ldap_server.by_name.user_group_mappings[0].map_group_name
      }
    }
  }
}
