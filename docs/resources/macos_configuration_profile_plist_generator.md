---
page_title: "jamfpro_macos_configuration_profile_plist_generator"
description: |-
  
---

# jamfpro_macos_configuration_profile_plist_generator (Resource)


## Example Usage
```terraform
variable "version_number" {
  description = "The version number to include in the name and install button text."
  type        = string
  default     = "1.0"
}

// example hcl generated plist with 1 level of nesting
resource "jamfpro_macos_configuration_profile_plist_generator" "jamfpro_macos_configuration_profile_plist_generator_003" {
  name                = "tf-localtest-generator-accessibility-seeing-${var.plist_version_number}"
  description         = "Base Level Accessibility settings for vision"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  user_removable      = false
  level               = "System"

  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_names    = ["Marketing", "Finance"]
    }

    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([78, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_names    = ["IT", "Engineering"]
      ibeacon_ids                          = [3, 4]
    }
  }

  self_service {
    install_button_text             = "Install - ${var.version_number}"
    self_service_description        = "This is the self service description"
    force_users_to_view_description = true
    feature_on_main_page            = true
    notification                    = true
    notification_subject            = "New Profile Available"
    notification_message            = "A new profile is available for installation."

    self_service_category {
      id         = 10
      display_in = true
      feature_in = true
    }

    self_service_category {
      id         = 5
      display_in = false
      feature_in = true
    }
  }

  payloads {
    payload_root {
      payload_description_root        = "Base Level Accessibility settings for vision"
      payload_enabled_root            = true
      payload_organization_root       = "Deployment Theory"
      payload_removal_disallowed_root = false
      payload_scope_root              = "System"
      payload_type_root               = "Configuration"
      payload_version_root            = var.plist_version_number
    }

    payload_content {
      configuration {
        key   = "closeViewFarPoint"
        value = 2
      }
      configuration {
        key   = "closeViewHotkeysEnabled"
        value = true
      }
      configuration {
        key   = "closeViewNearPoint"
        value = 10
      }
      configuration {
        key   = "closeViewScrollWheelToggle"
        value = true
      }
      configuration {
        key   = "closeViewShowPreview"
        value = true
      }
      configuration {
        key   = "closeViewSmoothImages"
        value = true
      }
      configuration {
        key   = "contrast"
        value = 0
      }
      configuration {
        key   = "flashScreen"
        value = false
      }
      configuration {
        key   = "grayscale"
        value = false
      }
      configuration {
        key   = "mouseDriver"
        value = false
      }
      configuration {
        key   = "mouseDriverCursorSize"
        value = 3
      }
      configuration {
        key   = "mouseDriverIgnoreTrackpad"
        value = false
      }
      configuration {
        key   = "mouseDriverInitialDelay"
        value = 1.0
      }
      configuration {
        key   = "mouseDriverMaxSpeed"
        value = 3
      }
      configuration {
        key   = "slowKey"
        value = false
      }
      configuration {
        key   = "slowKeyBeepOn"
        value = false
      }
      configuration {
        key   = "slowKeyDelay"
        value = 0
      }
      configuration {
        key   = "stereoAsMono"
        value = false
      }
      configuration {
        key   = "stickyKey"
        value = false
      }
      configuration {
        key   = "stickyKeyBeepOnModifier"
        value = false
      }
      configuration {
        key   = "stickyKeyShowWindow"
        value = false
      }
      configuration {
        key   = "voiceOverOnOffKey"
        value = true
      }
      configuration {
        key   = "whiteOnBlack"
        value = false
      }

      payload_description  = ""
      payload_display_name = "Accessibility"
      payload_enabled      = true
      payload_organization = "Deployment Theory"
      payload_type         = "com.apple.universalaccess"
      payload_version      = var.plist_version_number
      payload_scope        = "System"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Jamf UI name for configuration profile.
- `payloads` (Block List, Min: 1, Max: 1) A list of payloads for the macOS configuration profile. (see [below for nested schema](#nestedblock--payloads))
- `redeploy_on_update` (String) Defines the redeployment behaviour when an update to a macOS config profileoccurs. This is always 'Newly Assigned' on new profile objects, but may be set to 'All'on profile update requests once the configuration profile has been deployed to at least one device.
- `scope` (Block List, Min: 1, Max: 1) The scope of the configuration profile. (see [below for nested schema](#nestedblock--scope))

### Optional

- `category_id` (Number) Jamf Pro category-related settings of the policy.
- `description` (String) Description of the configuration profile.
- `distribution_method` (String) The distribution method for the configuration profile. ['Make Available in Self Service','Install Automatically']
- `level` (String) The deployment level of the configuration profile. Available options are: 'User' or 'System'. Note: 'System' is mapped to 'Computer Level' in the Jamf Pro GUI.
- `self_service` (Block List, Max: 1) Self Service Configuration (see [below for nested schema](#nestedblock--self_service))
- `site_id` (Number) Jamf Pro Site-related settings of the policy.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_removable` (Boolean) Whether the configuration profile is user removeable or not.

### Read-Only

- `id` (String) The unique identifier of the macOS configuration profile.
- `uuid` (String) The universally unique identifier for the profile.

<a id="nestedblock--payloads"></a>
### Nested Schema for `payloads`

Required:

- `payload_content` (Block List, Min: 1, Max: 1) The payload content of the macOS configuration profile plist. Multiple payloads can be defined as needed.Defined as key value pairs and supports nested dictionaries. (see [below for nested schema](#nestedblock--payloads--payload_content))
- `payload_description_header` (String) Description of the payload at the header level of the plist. This provides a human-readable explanation of what the overall profile is intended to do or configure.
- `payload_enabled_header` (Boolean) Indicates whether the payload is enabled at the header level of the plist. If set to false, the overall profile will be disabled.
- `payload_organization_header` (String) The organization associated with the payload at the header level of the plist. This represents the entity that created or is responsible for the overall profile.
- `payload_type_header` (String) The type of the config profile payload at the header level of the plist. This indicates what kind of settings or configurations the overall profile applies.
- `payload_version_header` (Number) The version of the payload at the header level of the plist. This helps in identifying the version of the overall profile settings or configurations being applied.

Optional:

- `payload_display_name_header` (String) The display name of the payload at the header level of the plist. This is shown in user interfaces to identify the overall profile to users and administrators. Jamf Pro matches this to the name of the configuation profile, 'name' at the top of the schema.
- `payload_removal_disallowed_header` (Boolean) Indicates whether the removal of the payload is disallowed. If set to true, the MDM profile cannot be removed by users.
- `payload_scope_header` (String) The scope of the payload at the header level of the plist. This defines the context in which the overall profile settings are applied, can be either 'System' or 'User'.

Read-Only:

- `payload_identifier_header` (String) A unique identifier for the payload within the MDM profile at the header level of the plist. This identifier is used to track and reference the overall profile uniquely.
- `payload_uuid_header` (String) The UUID for the payload within the MDM profile at the header level of the plist. This ensures the uniqueness of the overall profile.

<a id="nestedblock--payloads--payload_content"></a>
### Nested Schema for `payloads.payload_content`

Required:

- `payload_enabled` (Boolean) Whether the payload is enabled.
- `payload_organization` (String) Organization associated with the payload.
- `payload_type` (String) Type of the config profile payload.
- `payload_version` (Number) Version of the payload.

Optional:

- `payload_description` (String) Description of the payload.
- `payload_display_name` (String) Display name of the payload.
- `payload_removal_disallowed` (Boolean) Whether the payload removal is disallowed.
- `payload_scope` (String) Scope of the payload. Computed by what is set by level. 'System' or 'User'.
- `setting` (Block Set) The key and value setting items of the macOS configuration profile plist (see [below for nested schema](#nestedblock--payloads--payload_content--setting))

Read-Only:

- `payload_identifier` (String) Identifier for the payload.A GUID.
- `payload_uuid` (String) UUID of the payload.

<a id="nestedblock--payloads--payload_content--setting"></a>
### Nested Schema for `payloads.payload_content.setting`

Required:

- `key` (String) The key for the xml plist entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the xml plist entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary.dictionary.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary.dictionary.dictionary.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary.dictionary.dictionary.dictionary.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Block Set) A nested dictionary structure. (see [below for nested schema](#nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary--dictionary--dictionary))
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.

<a id="nestedblock--payloads--payload_content--setting--dictionary--dictionary--dictionary--dictionary--dictionary--dictionary"></a>
### Nested Schema for `payloads.payload_content.setting.dictionary.dictionary.dictionary.dictionary.dictionary.dictionary`

Required:

- `key` (String) The key for the dictionary entry.

Optional:

- `dictionary` (Map of String) A nested dictionary structure for xml plist definition.
- `array_json` (String) An ordered plist array encoded as JSON. Preserves duplicate elements. Use instead of value or dictionary.
- `value` (String) The value for the dictionary entry.










<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Required:

- `all_computers` (Boolean) Whether the configuration profile is scoped to all computers.

Optional:

- `all_jss_users` (Boolean) Whether the configuration profile is scoped to all JSS users.
- `building_ids` (Set of Number) The buildings to which the configuration profile is scoped by Jamf ID.
- `computer_group_ids` (Set of Number) The computer groups to which the configuration profile is scoped by Jamf ID.
- `computer_ids` (Set of Number) The computers to which the configuration profile is scoped by Jamf ID.
- `department_ids` (Set of Number) The departments to which the configuration profile is scoped by Jamf ID.
- `exclusions` (Block List, Max: 1) The scope exclusions from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--exclusions))
- `jss_user_group_ids` (Set of Number) The JSS user groups to which the configuration profile is scoped by Jamf ID.
- `jss_user_ids` (Set of Number) The JSS users to which the configuration profile is scoped by Jamf ID.
- `limitations` (Block List, Max: 1) The scope limitations from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--limitations))

<a id="nestedblock--scope--exclusions"></a>
### Nested Schema for `scope.exclusions`

Optional:

- `building_ids` (Set of Number) Buildings excluded from scope by Jamf ID.
- `computer_group_ids` (Set of Number) Computer Groups excluded from scope by Jamf ID.
- `computer_ids` (Set of Number) Computers excluded from scope by Jamf ID.
- `department_ids` (Set of Number) Departments excluded from scope by Jamf ID.
- `directory_service_or_local_usernames` (Set of String) A set of directory service / local usernames for scoping exclusions.
- `directory_service_usergroup_names` (Set of String) A set of directory service user group names for exclusions.
- `ibeacon_ids` (Set of Number) Ibeacons excluded from scope by Jamf ID.
- `jss_user_group_ids` (Set of Number) JSS User Groups excluded from scope by Jamf ID.
- `jss_user_ids` (Set of Number) JSS Users excluded from scope by Jamf ID.
- `network_segment_ids` (Set of Number) Network segments excluded from scope by Jamf ID.


<a id="nestedblock--scope--limitations"></a>
### Nested Schema for `scope.limitations`

Optional:

- `directory_service_or_local_usernames` (Set of String) A set of directory service / local usernames for scoping limitations.
- `directory_service_usergroup_names` (Set of String) A set of directory service user group names for limitations.
- `ibeacon_ids` (Set of Number) A set of iBeacon IDs for limitations.
- `network_segment_ids` (Set of Number) A set of network segment IDs for limitations.



<a id="nestedblock--self_service"></a>
### Nested Schema for `self_service`

Optional:

- `feature_on_main_page` (Boolean) Shows Configuration Profile on Self Service main page
- `force_users_to_view_description` (Boolean) Force users to view the description before the profile installs
- `install_button_text` (String) Name for the button that users click to install the profile
- `notification` (Boolean) TEMPORARILY DISABLED
Enables Notification for this profile in self service
- `notification_message` (String) Message body
- `notification_subject` (String) Message Subject
- `self_service_category` (Block Set) Self Service category options (see [below for nested schema](#nestedblock--self_service--self_service_category))
- `self_service_description` (String) Description to display for the profile in Self Service
- `self_service_display_name` (String) Display name for the profile in Self Service (Self Service 10.0.0 or later)
- `self_service_icon_id` (Number) Icon for policy to use in self-service. Can be used in conjection with the icons resource

<a id="nestedblock--self_service--self_service_category"></a>
### Nested Schema for `self_service.self_service_category`

Required:

- `display_in` (Boolean) Display this profile in this category?
- `feature_in` (Boolean) Feature this profile in this category?
- `id` (Number) ID of category. Both ID and Name are required

Read-Only:

- `name` (String) Name of category. Both ID and Name are required



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)
## Collection state migration

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`setting`, `dictionary`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.
See the [resource migration guide](../resource-migration-guide.md).

Only `payloads.payload_content.setting` and its nested dictionary entry blocks
are sets. `payloads`, `payload_content`, and all plist arrays keep their existing
order semantics. Dictionary keys must be unique.

Use `array_json` instead of `value` or `dictionary` for an ordered array:

```hcl
setting {
  key        = "OrderedValues"
  array_json = jsonencode(["second", "first", "second"])
}
```

`array_json` is an optional string on settings and dictionary entries (up to the
existing six nested block levels). It supports JSON-compatible plist array
values and preserves element order and duplicates. Nested dictionaries inside
an array remain array elements. Reordering an array is a real configuration
change. The legacy string value conversion continues to recognize booleans and
integers. The terminal dictionary remains a string map.

Arrays containing plist dates or binary data are rejected with an explicit error;
they are not coerced to JSON strings. Use the raw plist profile resource for
those payload types.

Array JSON numbers are normalized without reordering or deduplicating elements.
Integers retain their signed/unsigned 64-bit values; real numbers use the plist
64-bit floating-point representation and retain a decimal/exponent marker when
read back (for example, `1.0` remains a real). JSON `null`, integers outside the
supported 64-bit range, and non-finite/out-of-range real values are rejected
before apply rather than silently removed or coerced.
