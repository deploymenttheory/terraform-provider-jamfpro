---
page_title: "jamfpro_cloud_ldap"
description: |-
  
---

# jamfpro_cloud_ldap (Resource)


## Example Usage
```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_type` (String) The type of LDAP connection (LDAPS or START_TLS)
- `display_name` (String) The display name for the cloud LDAP configuration
- `domain_name` (String, Sensitive) The domain name for the LDAP server
- `group_mappings_id` (String) Group ID attribute mapping (e.g., cn)
- `group_mappings_name` (String) Group name attribute mapping (e.g., cn)
- `group_mappings_object_class_limitation` (String) Object class limitation for group mappings
- `group_mappings_object_classes` (String) Object classes for group mappings (e.g., groupOfNames)
- `group_mappings_search_base` (String) Search base for group mappings (e.g., ou=Groups)
- `group_mappings_search_scope` (String) Search scope for group mappings
- `group_mappings_uuid` (String) Group UUID attribute mapping (e.g., gidNumber)
- `group_membership_mapping` (String) Group membership attribute mapping (e.g., memberOf)
- `keystore_file_bytes` (String, Sensitive) Base64 encoded keystore file
- `keystore_file_name` (String) Name of the keystore file
- `keystore_password` (String, Sensitive)
- `port` (Number) The port number for the LDAP server
- `provider_name` (String) The name of the cloud identity provider. Must be 'GOOGLE' or 'AZURE'.
- `server_enabled` (Boolean) Whether the cloud LDAP server is enabled
- `server_url` (String) The URL of the LDAP server
- `user_mappings_email_address` (String) Email address attribute mapping (e.g., mail)
- `user_mappings_id` (String) User ID attribute mapping (e.g., mail)
- `user_mappings_object_class_limitation` (String) Object class limitation for user mappings
- `user_mappings_object_classes` (String) Object classes for user mappings (e.g., inetOrgPerson)
- `user_mappings_real_name` (String) Real name attribute mapping (e.g., displayName)
- `user_mappings_search_base` (String) Search base for user mappings (e.g., ou=Users)
- `user_mappings_search_scope` (String) Search scope for user mappings
- `user_mappings_username` (String) Username attribute mapping (e.g., uid)
- `user_mappings_uuid` (String) User UUID attribute mapping (e.g., uid)

### Optional

- `connection_timeout` (Number) Connection timeout in seconds
- `membership_calculation_optimization_enabled` (Boolean) Enable optimization for membership calculations
- `search_timeout` (Number) Search timeout in seconds
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `use_wildcards` (Boolean) Whether to use wildcards in LDAP queries
- `user_mappings_additional_search_base` (String) Additional search base for user mappings
- `user_mappings_building` (String) Building attribute mapping
- `user_mappings_department` (String) Department attribute mapping (e.g., departmentNumber)
- `user_mappings_phone` (String) Phone attribute mapping
- `user_mappings_position` (String) Position attribute mapping (e.g., title)
- `user_mappings_room` (String) Room attribute mapping

### Read-Only

- `id` (String) The ID of this resource.
- `keystore_expiration_date` (String)
- `keystore_subject` (String)
- `keystore_type` (String)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)