---
page_title: "jamfpro_macos_configuration_profile_plist"
description: |-
  
---

# jamfpro_macos_configuration_profile_plist (Resource)


## Example Usage
```terraform
variable "version_number" {
  description = "The version number to include in the name and install button text."
  type        = string
  default     = "v1.0"
}

// Minimum viable example of creating a macOS configuration profile in Jamf Pro for automatic installation using a plist source file

resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_064" {
  name                = "your-name-${var.version_number}"
  description         = "An example mobile device configuration profile."
  level               = "System"
  distribution_method = "Install Automatically" // "Make Available in Self Service", "Install Automatically"
  payloads            = file("${path.module}/path/to/your/file.mobileconfig")
  payload_validate    = true
  user_removable      = false

  scope {
    all_computers = true
    all_jss_users = false
  }

}

// Example of creating a macOS configuration profile in Jamf Pro for self service using a plist source file
resource "jamfpro_macos_configuration_profile_plist" "jamfpro_macos_configuration_profile_001" {
  name                = "your-name-${var.version_number}"
  description         = "An example mobile device configuration profile."
  level               = "User"                           // "User", "Device"
  distribution_method = "Make Available in Self Service" // "Make Available in Self Service", "Install Automatically"
  payloads            = file("${path.module}/path/to/your/file.mobileconfig")
  payload_validate     = true
  user_removable      = false

  // Optional Block

  site_id = 967

  // Optional Block
  category_id = 5
  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    // Optional Block
    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
    }

    // Optional Block
    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([78, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
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

    self_service_categories {
      id         = 10
      display_in = true
      feature_in = true
    }

    self_service_categories {
      id         = 5
      display_in = false
      feature_in = true
    }
  }
}

// Example of creating a macOS configuration profile in Jamf Pro for automatic installation
resource "jamfpro_macos_configuration_profile" "jamfpro_macos_configuration_profile_002" {
  name                = "your-name-${var.version_number}"
  description         = "An example mobile device configuration profile."
  level               = "User"                  // "User", "Device"
  distribution_method = "Install Automatically" // "Make Available in Self Service", "Install Automatically"
  payloads            = file("${path.module}/path/to/your/file.mobileconfig")
  user_removable      = false
  payload_validate     = true

  // Optional Block
  site_id = 1

  // Optional Block
  category_id = 1
  scope {
    all_computers = false
    all_jss_users = false

    computer_ids       = [16, 20, 21]
    computer_group_ids = sort([78, 1])
    building_ids       = ([1348, 1349])
    department_ids     = ([37287, 37288])
    jss_user_ids       = sort([2, 1])
    jss_user_group_ids = [4, 505]

    // Optional Block
    limitations {
      network_segment_ids                  = [4, 5]
      ibeacon_ids                          = [3, 4]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
    }
    // Optional Block
    exclusions {
      computer_ids                         = [16, 20, 21]
      computer_group_ids                   = sort([78, 1])
      building_ids                         = ([1348, 1349])
      department_ids                       = ([37287, 37288])
      network_segment_ids                  = [4, 5]
      jss_user_ids                         = sort([2, 1])
      jss_user_group_ids                   = [4, 505]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [3, 4]
      ibeacon_ids                          = [3, 4]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Jamf UI name for configuration profile.
- `payloads` (String) A MacOS configuration profile as a plist-formatted XML string.
- `scope` (Block List, Min: 1, Max: 1) The scope of the configuration profile. (see [below for nested schema](#nestedblock--scope))

### Optional

- `category_id` (Number) Jamf Pro category-related settings of the policy.
- `description` (String) Description of the configuration profile.
- `distribution_method` (String) The distribution method for the configuration profile. ['Make Available in Self Service','Install Automatically']
- `level` (String) The deployment level of the configuration profile. Available options are: 'User' or 'System'. Note: 'System' is mapped to 'Computer Level' in the Jamf Pro GUI.
- `payload_validate` (Boolean) Validates plist payload XML. Turn off to force malformed XML confguration. Required when the configuration profile is a non Jamf Pro source, e.g iMazing. Removing this may cause unexpected stating behaviour.
- `self_service` (Block List, Max: 1) Self Service Configuration (see [below for nested schema](#nestedblock--self_service))
- `site_id` (Number) Jamf Pro Site-related settings of the policy.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_removable` (Boolean) Whether the configuration profile is user removeable or not.

### Read-Only

- `id` (String) The unique identifier of the macOS configuration profile.
- `redeploy_on_update` (String) Defines the redeployment behaviour when a mobile device config profile update occurs.This is always 'Newly Assigned' on new profile objects, but may be set 'All' on profile update requests and in TF state
- `uuid` (String) The universally unique identifier for the profile.

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Required:

- `all_computers` (Boolean) Whether the configuration profile is scoped to all computers.

Optional:

- `all_jss_users` (Boolean) Whether the configuration profile is scoped to all JSS users.
- `building_ids` (List of Number) The buildings to which the configuration profile is scoped by Jamf ID
- `computer_group_ids` (List of Number) The computer groups to which the configuration profile is scoped by Jamf ID
- `computer_ids` (List of Number) The computers to which the configuration profile is scoped by Jamf ID
- `department_ids` (List of Number) The departments to which the configuration profile is scoped by Jamf ID
- `exclusions` (Block List, Max: 1) The scope exclusions from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--exclusions))
- `jss_user_group_ids` (List of Number) The jss user groups to which the configuration profile is scoped by Jamf ID
- `jss_user_ids` (List of Number) The jss users to which the configuration profile is scoped by Jamf ID
- `limitations` (Block List, Max: 1) The scope limitations from the macOS configuration profile. (see [below for nested schema](#nestedblock--scope--limitations))

<a id="nestedblock--scope--exclusions"></a>
### Nested Schema for `scope.exclusions`

Optional:

- `building_ids` (List of Number) Buildings excluded from scope by Jamf ID.
- `computer_group_ids` (List of Number) Computer Groups excluded from scope by Jamf ID.
- `computer_ids` (List of Number) Computers excluded from scope by Jamf ID.
- `department_ids` (List of Number) Departments excluded from scope by Jamf ID.
- `directory_service_or_local_usernames` (List of String) A list of directory service / local usernames for scoping limitations.
- `directory_service_usergroup_ids` (List of Number) A list of directory service / local user group IDs for limitations.
- `ibeacon_ids` (List of Number) Ibeacons excluded from scope by Jamf ID.
- `jss_user_group_ids` (List of Number) JSS User Groups excluded from scope by Jamf ID.
- `jss_user_ids` (List of Number) JSS Users excluded from scope by Jamf ID.
- `network_segment_ids` (List of Number) Network segments excluded from scope by Jamf ID.


<a id="nestedblock--scope--limitations"></a>
### Nested Schema for `scope.limitations`

Optional:

- `directory_service_or_local_usernames` (List of String) A list of directory service / local usernames for scoping limitations.
- `directory_service_usergroup_ids` (List of Number) A list of directory service user group IDs for limitations.
- `ibeacon_ids` (List of Number) A list of iBeacon IDs for limitations.
- `network_segment_ids` (List of Number) A list of network segment IDs for limitations.



<a id="nestedblock--self_service"></a>
### Nested Schema for `self_service`

Optional:

- `feature_on_main_page` (Boolean) Shows Configuration Profile on Self Service main page
- `force_users_to_view_description` (Boolean) Forces users to view the description
- `install_button_text` (String) Text shown on Self Service install button
- `notification` (Boolean) Enables Notification for this profile in self service
- `notification_message` (String) Message body
- `notification_subject` (String) Message Subject
- `self_service_category` (Block List) Self Service category options (see [below for nested schema](#nestedblock--self_service--self_service_category))
- `self_service_description` (String) Description shown in Self Service
- `self_service_icon_id` (Number) Icon for policy to use in self-service

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