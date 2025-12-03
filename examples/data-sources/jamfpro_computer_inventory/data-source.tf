# Example 1: Basic Computer Inventory Lookup by ID
data "jamfpro_computer_inventory" "example_basic" {
  id = "123"
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

