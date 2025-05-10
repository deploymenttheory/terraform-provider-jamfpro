resource "jamfpro_cloud_ldap" "example" {
  cloud_idp_common {
    provider_name = "GOOGLE"
    display_name  = "Google LDAP"
  }

  server {
    enabled         = true
    use_wildcards   = true
    connection_type = "LDAPS"
    server_url      = "ldap.google.com"
    domain_name     = "jamf.com"
    port            = 636

    connection_timeout = 60
    search_timeout     = 60

    membership_calculation_optimization_enabled = true

    keystore {
      password   = "supersecretpassword"
      file_bytes = filebase64("/path/to/keystore.p12")
      file_name  = "keystore.p12"
    }
  }

  mappings {
    user_mappings {
      object_class_limitation = "ANY_OBJECT_CLASSES"
      object_classes          = "inetOrgPerson"
      search_base             = "ou=Users"
      search_scope            = "ALL_SUBTREES"
      additional_search_base  = ""
      user_id                 = "mail"
      username                = "uid"
      real_name               = "displayName"
      email_address           = "mail"
      department              = "departmentNumber"
      building                = ""
      room                    = ""
      phone                   = ""
      position                = "title"
      user_uuid               = "uid"
    }

    group_mappings {
      object_class_limitation = "ANY_OBJECT_CLASSES"
      object_classes          = "groupOfNames"
      search_base             = "ou=Groups"
      search_scope            = "ALL_SUBTREES"
      group_id                = "cn"
      group_name              = "cn"
      group_uuid              = "gidNumber"
    }

    membership_mappings {
      group_membership_mapping = "memberOf"
    }
  }
}
