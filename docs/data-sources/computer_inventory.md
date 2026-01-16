---
page_title: "jamfpro_computer_inventory"
description: |-
  
---

# jamfpro_computer_inventory (Data Source)


## Example Usage
```terraform
# Example 1: Basic Computer Inventory Lookup by ID
data "jamfpro_computer_inventory" "example_basic" {
  id = "123"
}

# Example 1a: Lookup by Computer Name
data "jamfpro_computer_inventory" "example_by_name" {
  name = "MacBook-Pro-123"
}

# Example 1b: Lookup by Serial Number
data "jamfpro_computer_inventory" "example_by_serial" {
  serial_number = "C02ABC123DEF"
}

# Example 2: Output Common Computer Information
data "jamfpro_computer_inventory" "example_detailed" {
  id = "456"
}

output "computer_name" {
  value       = data.jamfpro_computer_inventory.example_detailed.general[0].name
  description = "The name of the computer"
}

output "computer_serial" {
  value       = data.jamfpro_computer_inventory.example_detailed.hardware[0].serial_number
  description = "Serial number of the computer"
}

output "last_ip_address" {
  value       = data.jamfpro_computer_inventory.example_detailed.general[0].last_ip_address
  description = "Last known IP address"
}

# Example 3: Check MDM and Enrollment Status
data "jamfpro_computer_inventory" "mdm_check" {
  id = "789"
}

output "mdm_status" {
  value = {
    supervised               = data.jamfpro_computer_inventory.mdm_check.general[0].supervised
    user_approved_mdm        = data.jamfpro_computer_inventory.mdm_check.general[0].user_approved_mdm
    mdm_capable              = data.jamfpro_computer_inventory.mdm_check.general[0].mdm_capable[0].capable
    enrolled_via_ade         = data.jamfpro_computer_inventory.mdm_check.general[0].enrolled_via_automated_device_enrollment
    declarative_mgmt_enabled = data.jamfpro_computer_inventory.mdm_check.general[0].declarative_device_management_enabled
    enrollment_method        = data.jamfpro_computer_inventory.mdm_check.general[0].enrollment_method[0].object_name
  }
  description = "MDM and enrollment status information"
}

# Example 4: Hardware Information
data "jamfpro_computer_inventory" "hardware_info" {
  id = "101"
}

output "hardware_details" {
  value = {
    make                   = data.jamfpro_computer_inventory.hardware_info.hardware[0].make
    model                  = data.jamfpro_computer_inventory.hardware_info.hardware[0].model
    model_identifier       = data.jamfpro_computer_inventory.hardware_info.hardware[0].model_identifier
    processor_type         = data.jamfpro_computer_inventory.hardware_info.hardware[0].processor_type
    processor_architecture = data.jamfpro_computer_inventory.hardware_info.hardware[0].processor_architecture
    total_ram_mb           = data.jamfpro_computer_inventory.hardware_info.hardware[0].total_ram_megabytes
    apple_silicon          = data.jamfpro_computer_inventory.hardware_info.hardware[0].apple_silicon
    battery_capacity       = data.jamfpro_computer_inventory.hardware_info.hardware[0].battery_capacity_percent
  }
  description = "Detailed hardware specifications"
}

# Example 5: Security and Encryption Status
data "jamfpro_computer_inventory" "security_check" {
  id = "202"
}

output "security_status" {
  value = {
    # Security settings
    sip_status        = data.jamfpro_computer_inventory.security_check.security[0].sip_status
    gatekeeper_status = data.jamfpro_computer_inventory.security_check.security[0].gatekeeper_status
    firewall_enabled  = data.jamfpro_computer_inventory.security_check.security[0].firewall_enabled
    activation_lock   = data.jamfpro_computer_inventory.security_check.security[0].activation_lock_enabled
    recovery_lock     = data.jamfpro_computer_inventory.security_check.security[0].recovery_lock_enabled
    secure_boot_level = data.jamfpro_computer_inventory.security_check.security[0].secure_boot_level
    # Disk encryption
    filevault_status   = data.jamfpro_computer_inventory.security_check.operating_system[0].filevault2_status
    recovery_key_valid = data.jamfpro_computer_inventory.security_check.disk_encryption[0].individual_recovery_key_validity_status
    institutional_key  = data.jamfpro_computer_inventory.security_check.disk_encryption[0].institutional_recovery_key_present
  }
  description = "Security and encryption status"
}

# Example 6: Operating System Information
data "jamfpro_computer_inventory" "os_info" {
  id = "303"
}

output "os_details" {
  value = {
    name                    = data.jamfpro_computer_inventory.os_info.operating_system[0].name
    version                 = data.jamfpro_computer_inventory.os_info.operating_system[0].version
    build                   = data.jamfpro_computer_inventory.os_info.operating_system[0].build
    rapid_security_response = data.jamfpro_computer_inventory.os_info.operating_system[0].rapid_security_response
    jamf_binary_version     = data.jamfpro_computer_inventory.os_info.general[0].jamf_binary_version
  }
  description = "Operating system information"
}

# Example 7: Storage and Disk Information
data "jamfpro_computer_inventory" "storage_info" {
  id = "404"
}

output "storage_summary" {
  value = {
    boot_drive_available_mb = data.jamfpro_computer_inventory.storage_info.storage[0].boot_drive_available_space_megabytes
    disk_count              = length(data.jamfpro_computer_inventory.storage_info.storage[0].disks)
  }
  description = "Storage information summary"
}

# Example 8: User and Location Information
data "jamfpro_computer_inventory" "user_location" {
  id = "505"
}

output "user_info" {
  value = {
    username   = data.jamfpro_computer_inventory.user_location.user_and_location[0].username
    realname   = data.jamfpro_computer_inventory.user_location.user_and_location[0].realname
    email      = data.jamfpro_computer_inventory.user_location.user_and_location[0].email
    position   = data.jamfpro_computer_inventory.user_location.user_and_location[0].position
    phone      = data.jamfpro_computer_inventory.user_location.user_and_location[0].phone
    department = data.jamfpro_computer_inventory.user_location.user_and_location[0].department_id
    building   = data.jamfpro_computer_inventory.user_location.user_and_location[0].building_id
    room       = data.jamfpro_computer_inventory.user_location.user_and_location[0].room
  }
  description = "User and location details"
}

# Example 9: Purchasing Information
data "jamfpro_computer_inventory" "purchasing_info" {
  id = "606"
}

output "purchasing_details" {
  value = {
    purchased       = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].purchased
    leased          = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].leased
    po_number       = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].po_number
    vendor          = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].vendor
    purchase_price  = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].purchase_price
    warranty_date   = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].warranty_date
    apple_care_id   = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].apple_care_id
    life_expectancy = data.jamfpro_computer_inventory.purchasing_info.purchasing[0].life_expectancy
  }
  description = "Purchasing and warranty information"
}

# Example 10: List Installed Applications
data "jamfpro_computer_inventory" "app_inventory" {
  id = "707"
}

output "installed_applications" {
  value = [
    for app in data.jamfpro_computer_inventory.app_inventory.applications : {
      name             = app.name
      version          = app.version
      bundle_id        = app.bundle_id
      mac_app_store    = app.mac_app_store
      update_available = app.update_available
    }
  ]
  description = "List of all installed applications"
}

# Example 11: Configuration Profiles Status
data "jamfpro_computer_inventory" "profile_check" {
  id = "808"
}

output "configuration_profiles" {
  value = [
    for profile in data.jamfpro_computer_inventory.profile_check.configuration_profiles : {
      display_name       = profile.display_name
      profile_identifier = profile.profile_identifier
      last_installed     = profile.last_installed
      removable          = profile.removable
    }
  ]
  description = "Installed configuration profiles"
}

# Example 12: Local User Accounts
data "jamfpro_computer_inventory" "user_accounts" {
  id = "909"
}

output "local_users" {
  value = [
    for user in data.jamfpro_computer_inventory.user_accounts.local_user_accounts : {
      username            = user.username
      full_name           = user.full_name
      admin               = user.admin
      filevault_enabled   = user.file_vault2_enabled
      user_account_type   = user.user_account_type
      home_directory_size = user.home_directory_size_mb
    }
  ]
  description = "Local user accounts on the computer"
}

# Example 13: Software Updates Available
data "jamfpro_computer_inventory" "update_check" {
  id = "1010"
}

output "available_updates" {
  value = [
    for update in data.jamfpro_computer_inventory.update_check.software_updates : {
      name         = update.name
      version      = update.version
      package_name = update.package_name
    }
  ]
  description = "Available software updates"
}

# Example 14: Group Memberships
data "jamfpro_computer_inventory" "group_check" {
  id = "1111"
}

output "group_memberships" {
  value = [
    for group in data.jamfpro_computer_inventory.group_check.group_memberships : {
      group_name  = group.group_name
      group_id    = group.group_id
      smart_group = group.smart_group
    }
  ]
  description = "Computer group memberships"
}

# Example 15: Extension Attributes (Custom Fields)
data "jamfpro_computer_inventory" "extension_attrs" {
  id = "1212"
}

output "extension_attributes" {
  value = [
    for attr in data.jamfpro_computer_inventory.extension_attrs.extension_attributes : {
      name        = attr.name
      values      = attr.values
      description = attr.description
      data_type   = attr.data_type
    }
  ]
  description = "Extension attributes (custom inventory fields)"
}

# Example 16: Complete Inventory Export for Reporting
data "jamfpro_computer_inventory" "full_inventory" {
  id = "1313"
}

output "complete_inventory_json" {
  value = jsonencode({
    id                = data.jamfpro_computer_inventory.full_inventory.id
    udid              = data.jamfpro_computer_inventory.full_inventory.udid
    general           = data.jamfpro_computer_inventory.full_inventory.general
    hardware          = data.jamfpro_computer_inventory.full_inventory.hardware
    operating_system  = data.jamfpro_computer_inventory.full_inventory.operating_system
    security          = data.jamfpro_computer_inventory.full_inventory.security
    storage           = data.jamfpro_computer_inventory.full_inventory.storage
    user_and_location = data.jamfpro_computer_inventory.full_inventory.user_and_location
    purchasing        = data.jamfpro_computer_inventory.full_inventory.purchasing
  })
  description = "Complete inventory as JSON for external processing"
}

# Example 17: Use Computer ID from Another Resource
resource "jamfpro_static_computer_group" "example_group" {
  name = "Example Group"
}

# Reference computers in the group
data "jamfpro_computer_inventory" "group_member" {
  # In a real scenario, you would iterate over group members
  id = "1414"
}

# Example 18: Asset Management Check
data "jamfpro_computer_inventory" "asset_check" {
  id = "1515"
}

output "asset_info" {
  value = {
    asset_tag     = data.jamfpro_computer_inventory.asset_check.general[0].asset_tag
    barcode1      = data.jamfpro_computer_inventory.asset_check.general[0].barcode1
    barcode2      = data.jamfpro_computer_inventory.asset_check.general[0].barcode2
    serial        = data.jamfpro_computer_inventory.asset_check.hardware[0].serial_number
    management_id = data.jamfpro_computer_inventory.asset_check.general[0].management_id
  }
  description = "Asset tracking information"
}

# Example 19: Network Information
data "jamfpro_computer_inventory" "network_info" {
  id = "1616"
}

output "network_details" {
  value = {
    primary_mac      = data.jamfpro_computer_inventory.network_info.hardware[0].mac_address
    alt_mac          = data.jamfpro_computer_inventory.network_info.hardware[0].alt_mac_address
    last_ip          = data.jamfpro_computer_inventory.network_info.general[0].last_ip_address
    last_reported_ip = data.jamfpro_computer_inventory.network_info.general[0].last_reported_ip
    nic_speed        = data.jamfpro_computer_inventory.network_info.hardware[0].nic_speed
  }
  description = "Network adapter information"
}

# Example 20: Site Assignment
data "jamfpro_computer_inventory" "site_check" {
  id = "1717"
}

output "site_info" {
  value = {
    site_id   = data.jamfpro_computer_inventory.site_check.general[0].site_id[0].id
    site_name = data.jamfpro_computer_inventory.site_check.general[0].site_id[0].name
  }
  description = "Site assignment information"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String)
- `serial_number` (String)

### Read-Only

- `applications` (List of Object) (see [below for nested schema](#nestedatt--applications))
- `attachments` (List of Object) (see [below for nested schema](#nestedatt--attachments))
- `certificates` (List of Object) (see [below for nested schema](#nestedatt--certificates))
- `configuration_profiles` (List of Object) (see [below for nested schema](#nestedatt--configuration_profiles))
- `disk_encryption` (List of Object) (see [below for nested schema](#nestedatt--disk_encryption))
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedatt--extension_attributes))
- `fonts` (List of Object) (see [below for nested schema](#nestedatt--fonts))
- `general` (List of Object) (see [below for nested schema](#nestedatt--general))
- `group_memberships` (List of Object) (see [below for nested schema](#nestedatt--group_memberships))
- `hardware` (List of Object) (see [below for nested schema](#nestedatt--hardware))
- `ibeacons` (List of Object) (see [below for nested schema](#nestedatt--ibeacons))
- `id` (String) The ID of this resource.
- `licensed_software` (List of Object) (see [below for nested schema](#nestedatt--licensed_software))
- `local_user_accounts` (List of Object) (see [below for nested schema](#nestedatt--local_user_accounts))
- `operating_system` (List of Object) (see [below for nested schema](#nestedatt--operating_system))
- `package_receipts` (List of Object) (see [below for nested schema](#nestedatt--package_receipts))
- `plugins` (List of Object) (see [below for nested schema](#nestedatt--plugins))
- `printers` (List of Object) (see [below for nested schema](#nestedatt--printers))
- `purchasing` (List of Object) (see [below for nested schema](#nestedatt--purchasing))
- `security` (List of Object) (see [below for nested schema](#nestedatt--security))
- `services` (List of Object) (see [below for nested schema](#nestedatt--services))
- `software_updates` (List of Object) (see [below for nested schema](#nestedatt--software_updates))
- `storage` (List of Object) (see [below for nested schema](#nestedatt--storage))
- `udid` (String)
- `user_and_location` (List of Object) (see [below for nested schema](#nestedatt--user_and_location))

<a id="nestedatt--applications"></a>
### Nested Schema for `applications`

Read-Only:

- `bundle_id` (String)
- `external_version_id` (String)
- `mac_app_store` (Boolean)
- `name` (String)
- `path` (String)
- `size_megabytes` (Number)
- `update_available` (Boolean)
- `version` (String)


<a id="nestedatt--attachments"></a>
### Nested Schema for `attachments`

Read-Only:

- `file_type` (String)
- `id` (String)
- `name` (String)
- `size_bytes` (Number)


<a id="nestedatt--certificates"></a>
### Nested Schema for `certificates`

Read-Only:

- `certificate_status` (String)
- `common_name` (String)
- `expiration_date` (String)
- `identity` (Boolean)
- `issued_date` (String)
- `lifecycle_status` (String)
- `serial_number` (String)
- `sha1_fingerprint` (String)
- `subject_name` (String)
- `username` (String)


<a id="nestedatt--configuration_profiles"></a>
### Nested Schema for `configuration_profiles`

Read-Only:

- `display_name` (String)
- `id` (String)
- `last_installed` (String)
- `profile_identifier` (String)
- `removable` (Boolean)
- `username` (String)


<a id="nestedatt--disk_encryption"></a>
### Nested Schema for `disk_encryption`

Read-Only:

- `boot_partition_encryption_details` (List of Object) (see [below for nested schema](#nestedobjatt--disk_encryption--boot_partition_encryption_details))
- `disk_encryption_configuration_name` (String)
- `file_vault2_eligibility_message` (String)
- `file_vault2_enabled_user_names` (List of String)
- `individual_recovery_key_validity_status` (String)
- `institutional_recovery_key_present` (Boolean)

<a id="nestedobjatt--disk_encryption--boot_partition_encryption_details"></a>
### Nested Schema for `disk_encryption.boot_partition_encryption_details`

Read-Only:

- `partition_file_vault2_percent` (Number)
- `partition_file_vault2_state` (String)
- `partition_name` (String)



<a id="nestedatt--extension_attributes"></a>
### Nested Schema for `extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)


<a id="nestedatt--fonts"></a>
### Nested Schema for `fonts`

Read-Only:

- `name` (String)
- `path` (String)
- `version` (String)


<a id="nestedatt--general"></a>
### Nested Schema for `general`

Read-Only:

- `asset_tag` (String)
- `barcode1` (String)
- `barcode2` (String)
- `declarative_device_management_enabled` (Boolean)
- `distribution_point` (String)
- `enrolled_via_automated_device_enrollment` (Boolean)
- `enrollment_method` (List of Object) (see [below for nested schema](#nestedobjatt--general--enrollment_method))
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--general--extension_attributes))
- `initial_entry_date` (String)
- `itunes_store_account_active` (Boolean)
- `jamf_binary_version` (String)
- `last_cloud_backup_date` (String)
- `last_contact_time` (String)
- `last_enrolled_date` (String)
- `last_ip_address` (String)
- `last_reported_ip` (String)
- `management_id` (String)
- `mdm_capable` (List of Object) (see [below for nested schema](#nestedobjatt--general--mdm_capable))
- `mdm_profile_expiration` (String)
- `name` (String)
- `platform` (String)
- `remote_management` (List of Object) (see [below for nested schema](#nestedobjatt--general--remote_management))
- `report_date` (String)
- `site_id` (List of Object) (see [below for nested schema](#nestedobjatt--general--site_id))
- `supervised` (Boolean)
- `user_approved_mdm` (Boolean)

<a id="nestedobjatt--general--enrollment_method"></a>
### Nested Schema for `general.enrollment_method`

Read-Only:

- `id` (String)
- `object_name` (String)
- `object_type` (String)


<a id="nestedobjatt--general--extension_attributes"></a>
### Nested Schema for `general.extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)


<a id="nestedobjatt--general--mdm_capable"></a>
### Nested Schema for `general.mdm_capable`

Read-Only:

- `capable` (Boolean)
- `capable_users` (List of String)


<a id="nestedobjatt--general--remote_management"></a>
### Nested Schema for `general.remote_management`

Read-Only:

- `managed` (Boolean)
- `management_username` (String)


<a id="nestedobjatt--general--site_id"></a>
### Nested Schema for `general.site_id`

Read-Only:

- `id` (String)
- `name` (String)



<a id="nestedatt--group_memberships"></a>
### Nested Schema for `group_memberships`

Read-Only:

- `group_id` (String)
- `group_name` (String)
- `smart_group` (Boolean)


<a id="nestedatt--hardware"></a>
### Nested Schema for `hardware`

Read-Only:

- `alt_mac_address` (String)
- `alt_network_adapter_type` (String)
- `apple_silicon` (Boolean)
- `battery_capacity_percent` (Number)
- `ble_capable` (Boolean)
- `boot_rom` (String)
- `bus_speed_mhz` (Number)
- `cache_size_kilobytes` (Number)
- `core_count` (Number)
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--hardware--extension_attributes))
- `mac_address` (String)
- `make` (String)
- `model` (String)
- `model_identifier` (String)
- `network_adapter_type` (String)
- `nic_speed` (String)
- `open_ram_slots` (Number)
- `optical_drive` (String)
- `processor_architecture` (String)
- `processor_count` (Number)
- `processor_speed_mhz` (Number)
- `processor_type` (String)
- `serial_number` (String)
- `smc_version` (String)
- `supports_ios_app_installs` (Boolean)
- `total_ram_megabytes` (Number)

<a id="nestedobjatt--hardware--extension_attributes"></a>
### Nested Schema for `hardware.extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)



<a id="nestedatt--ibeacons"></a>
### Nested Schema for `ibeacons`

Read-Only:

- `name` (String)


<a id="nestedatt--licensed_software"></a>
### Nested Schema for `licensed_software`

Read-Only:

- `id` (String)
- `name` (String)


<a id="nestedatt--local_user_accounts"></a>
### Nested Schema for `local_user_accounts`

Read-Only:

- `admin` (Boolean)
- `azure_active_directory_id` (String)
- `computer_azure_active_directory_id` (String)
- `file_vault2_enabled` (Boolean)
- `full_name` (String)
- `home_directory` (String)
- `home_directory_size_mb` (Number)
- `password_history_depth` (Number)
- `password_max_age` (Number)
- `password_min_complex_characters` (Number)
- `password_min_length` (Number)
- `password_require_alphanumeric` (Boolean)
- `uid` (String)
- `user_account_type` (String)
- `user_azure_active_directory_id` (String)
- `user_guid` (String)
- `username` (String)


<a id="nestedatt--operating_system"></a>
### Nested Schema for `operating_system`

Read-Only:

- `active_directory_status` (String)
- `build` (String)
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--operating_system--extension_attributes))
- `filevault2_status` (String)
- `name` (String)
- `rapid_security_response` (String)
- `software_update_device_id` (String)
- `supplemental_build_version` (String)
- `version` (String)

<a id="nestedobjatt--operating_system--extension_attributes"></a>
### Nested Schema for `operating_system.extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)



<a id="nestedatt--package_receipts"></a>
### Nested Schema for `package_receipts`

Read-Only:

- `cached` (List of String)
- `installed_by_installer_swu` (List of String)
- `installed_by_jamf_pro` (List of String)


<a id="nestedatt--plugins"></a>
### Nested Schema for `plugins`

Read-Only:

- `name` (String)
- `path` (String)
- `version` (String)


<a id="nestedatt--printers"></a>
### Nested Schema for `printers`

Read-Only:

- `location` (String)
- `name` (String)
- `type` (String)
- `uri` (String)


<a id="nestedatt--purchasing"></a>
### Nested Schema for `purchasing`

Read-Only:

- `apple_care_id` (String)
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--purchasing--extension_attributes))
- `lease_date` (String)
- `leased` (Boolean)
- `life_expectancy` (Number)
- `po_date` (String)
- `po_number` (String)
- `purchase_price` (String)
- `purchased` (Boolean)
- `purchasing_account` (String)
- `purchasing_contact` (String)
- `vendor` (String)
- `warranty_date` (String)

<a id="nestedobjatt--purchasing--extension_attributes"></a>
### Nested Schema for `purchasing.extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)



<a id="nestedatt--security"></a>
### Nested Schema for `security`

Read-Only:

- `activation_lock_enabled` (Boolean)
- `auto_login_disabled` (Boolean)
- `bootstrap_token_allowed` (Boolean)
- `external_boot_level` (String)
- `firewall_enabled` (Boolean)
- `gatekeeper_status` (String)
- `recovery_lock_enabled` (Boolean)
- `remote_desktop_enabled` (Boolean)
- `secure_boot_level` (String)
- `sip_status` (String)
- `xprotect_version` (String)


<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `name` (String)


<a id="nestedatt--software_updates"></a>
### Nested Schema for `software_updates`

Read-Only:

- `name` (String)
- `package_name` (String)
- `version` (String)


<a id="nestedatt--storage"></a>
### Nested Schema for `storage`

Read-Only:

- `boot_drive_available_space_megabytes` (Number)
- `disks` (List of Object) (see [below for nested schema](#nestedobjatt--storage--disks))

<a id="nestedobjatt--storage--disks"></a>
### Nested Schema for `storage.disks`

Read-Only:

- `device` (String)
- `id` (String)
- `model` (String)
- `partitions` (List of Object) (see [below for nested schema](#nestedobjatt--storage--disks--partitions))
- `revision` (String)
- `serial_number` (String)
- `size_megabytes` (Number)
- `smart_status` (String)
- `type` (String)

<a id="nestedobjatt--storage--disks--partitions"></a>
### Nested Schema for `storage.disks.partitions`

Read-Only:

- `available_megabytes` (Number)
- `file_vault2_progress_percent` (Number)
- `file_vault2_state` (String)
- `lvm_managed` (Boolean)
- `name` (String)
- `partition_type` (String)
- `percent_used` (Number)
- `size_megabytes` (Number)




<a id="nestedatt--user_and_location"></a>
### Nested Schema for `user_and_location`

Read-Only:

- `building_id` (String)
- `department_id` (String)
- `email` (String)
- `extension_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--user_and_location--extension_attributes))
- `phone` (String)
- `position` (String)
- `realname` (String)
- `room` (String)
- `username` (String)

<a id="nestedobjatt--user_and_location--extension_attributes"></a>
### Nested Schema for `user_and_location.extension_attributes`

Read-Only:

- `data_type` (String)
- `definition_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `input_type` (String)
- `multi_value` (Boolean)
- `name` (String)
- `options` (List of String)
- `values` (List of String)