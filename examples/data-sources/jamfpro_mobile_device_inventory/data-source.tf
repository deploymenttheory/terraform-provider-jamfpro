# Example 1: Basic Mobile Device Inventory Lookup by ID
data "jamfpro_mobile_device_inventory" "example_basic" {
  id = "123"
}

# Example 1a: Lookup by Device Name
data "jamfpro_mobile_device_inventory" "example_by_name" {
  name = "John's iPhone"
}

# Example 1b: Lookup by Serial Number
data "jamfpro_mobile_device_inventory" "example_by_serial" {
  serial_number = "DMQVGC0DHLA0"
}

# Example 2: Output Common Mobile Device Information
data "jamfpro_mobile_device_inventory" "example_detailed" {
  id = "456"
}

output "device_name" {
  value       = data.jamfpro_mobile_device_inventory.example_detailed.display_name
  description = "The display name of the mobile device"
}

output "device_serial" {
  value       = data.jamfpro_mobile_device_inventory.example_detailed.serial_number
  description = "Serial number of the mobile device"
}

output "device_model" {
  value       = data.jamfpro_mobile_device_inventory.example_detailed.model
  description = "Device model"
}

output "last_ip_address" {
  value       = data.jamfpro_mobile_device_inventory.example_detailed.ip_address
  description = "Last known IP address"
}

# Example 3: Check MDM and Enrollment Status
data "jamfpro_mobile_device_inventory" "mdm_check" {
  id = "789"
}

output "mdm_status" {
  value = {
    supervised                       = data.jamfpro_mobile_device_inventory.mdm_check.supervised
    managed                          = data.jamfpro_mobile_device_inventory.mdm_check.managed
    declarative_mgmt_enabled         = data.jamfpro_mobile_device_inventory.mdm_check.declarative_device_management_enabled
    enrollment_session_token_valid   = data.jamfpro_mobile_device_inventory.mdm_check.enrollment_session_token_valid
    device_ownership_type            = data.jamfpro_mobile_device_inventory.mdm_check.device_ownership_type
  }
  description = "MDM and enrollment status information"
}

# Example 4: Hardware and Device Information
data "jamfpro_mobile_device_inventory" "hardware_info" {
  id = "101"
}

output "hardware_details" {
  value = {
    model              = data.jamfpro_mobile_device_inventory.hardware_info.model
    model_identifier   = data.jamfpro_mobile_device_inventory.hardware_info.model_identifier
    model_number       = data.jamfpro_mobile_device_inventory.hardware_info.model_number
    capacity_mb        = data.jamfpro_mobile_device_inventory.hardware_info.capacity_mb
    available_space_mb = data.jamfpro_mobile_device_inventory.hardware_info.available_space_mb
    battery_level      = data.jamfpro_mobile_device_inventory.hardware_info.battery_level
    battery_health     = data.jamfpro_mobile_device_inventory.hardware_info.battery_health
  }
  description = "Device hardware specifications"
}

# Example 5: Security and Encryption Status
data "jamfpro_mobile_device_inventory" "security_check" {
  id = "202"
}

output "security_status" {
  value = {
    activation_lock_enabled              = data.jamfpro_mobile_device_inventory.security_check.activation_lock_enabled
    passcode_present                     = data.jamfpro_mobile_device_inventory.security_check.passcode_present
    passcode_compliant                   = data.jamfpro_mobile_device_inventory.security_check.passcode_compliant
    passcode_compliant_with_profile      = data.jamfpro_mobile_device_inventory.security_check.passcode_compliant_with_profile
    data_protection                      = data.jamfpro_mobile_device_inventory.security_check.data_protection
    block_encryption_capable             = data.jamfpro_mobile_device_inventory.security_check.block_encryption_capable
    file_encryption_capable              = data.jamfpro_mobile_device_inventory.security_check.file_encryption_capable
    hardware_encryption_supported        = data.jamfpro_mobile_device_inventory.security_check.hardware_encryption_supported
    jailbreak_status                     = data.jamfpro_mobile_device_inventory.security_check.jailbreak_status
    lost_mode_enabled                    = data.jamfpro_mobile_device_inventory.security_check.lost_mode_enabled
  }
  description = "Security and encryption status"
}

# Example 6: Operating System Information
data "jamfpro_mobile_device_inventory" "os_info" {
  id = "303"
}

output "os_details" {
  value = {
    os_version                   = data.jamfpro_mobile_device_inventory.os_info.os_version
    os_build                     = data.jamfpro_mobile_device_inventory.os_info.os_build
    os_supplemental_build        = data.jamfpro_mobile_device_inventory.os_info.os_supplemental_build_version
    os_rapid_security_response   = data.jamfpro_mobile_device_inventory.os_info.os_rapid_security_response
  }
  description = "Operating system information"
}

# Example 7: Cellular and Network Information
data "jamfpro_mobile_device_inventory" "network_info" {
  id = "404"
}

output "network_details" {
  value = {
    wifi_mac_address          = data.jamfpro_mobile_device_inventory.network_info.wifi_mac_address
    bluetooth_mac_address     = data.jamfpro_mobile_device_inventory.network_info.bluetooth_mac_address
    ip_address                = data.jamfpro_mobile_device_inventory.network_info.ip_address
    carrier_settings_version  = data.jamfpro_mobile_device_inventory.network_info.carrier_settings_version
    current_carrier_network   = data.jamfpro_mobile_device_inventory.network_info.current_carrier_network
    home_carrier_network      = data.jamfpro_mobile_device_inventory.network_info.home_carrier_network
    cellular_technology       = data.jamfpro_mobile_device_inventory.network_info.cellular_technology
    device_phone_number       = data.jamfpro_mobile_device_inventory.network_info.device_phone_number
    imei                      = data.jamfpro_mobile_device_inventory.network_info.imei
    iccid                     = data.jamfpro_mobile_device_inventory.network_info.iccid
    data_roaming_enabled      = data.jamfpro_mobile_device_inventory.network_info.data_roaming_enabled
    personal_hotspot_enabled  = data.jamfpro_mobile_device_inventory.network_info.personal_hotspot_enabled
  }
  description = "Network and cellular information"
}

# Example 8: User and Location Information
data "jamfpro_mobile_device_inventory" "user_location" {
  id = "505"
}

output "user_info" {
  value = {
    username    = data.jamfpro_mobile_device_inventory.user_location.username
    full_name   = data.jamfpro_mobile_device_inventory.user_location.full_name
    email       = data.jamfpro_mobile_device_inventory.user_location.email_address
    position    = data.jamfpro_mobile_device_inventory.user_location.position
    phone       = data.jamfpro_mobile_device_inventory.user_location.phone_number
    department  = data.jamfpro_mobile_device_inventory.user_location.department
    building    = data.jamfpro_mobile_device_inventory.user_location.building
    room        = data.jamfpro_mobile_device_inventory.user_location.room
  }
  description = "User and location details"
}

# Example 9: Purchasing Information
data "jamfpro_mobile_device_inventory" "purchasing_info" {
  id = "606"
}

output "purchasing_details" {
  value = {
    purchased_or_leased       = data.jamfpro_mobile_device_inventory.purchasing_info.purchased_or_leased
    po_number                 = data.jamfpro_mobile_device_inventory.purchasing_info.po_number
    po_date                   = data.jamfpro_mobile_device_inventory.purchasing_info.po_date
    vendor                    = data.jamfpro_mobile_device_inventory.purchasing_info.vendor
    purchase_price            = data.jamfpro_mobile_device_inventory.purchasing_info.purchase_price
    warranty_expiration_date  = data.jamfpro_mobile_device_inventory.purchasing_info.warranty_expiration_date
    apple_care_id             = data.jamfpro_mobile_device_inventory.purchasing_info.apple_care_id
    lease_expiration_date     = data.jamfpro_mobile_device_inventory.purchasing_info.lease_expiration_date
    life_expectancy_years     = data.jamfpro_mobile_device_inventory.purchasing_info.life_expectancy_years
  }
  description = "Purchasing and warranty information"
}

# Example 10: Shared iPad Information
data "jamfpro_mobile_device_inventory" "shared_ipad" {
  id = "707"
}

output "shared_ipad_details" {
  value = {
    shared_ipad    = data.jamfpro_mobile_device_inventory.shared_ipad.shared_ipad
    quota_size     = data.jamfpro_mobile_device_inventory.shared_ipad.quota_size
    resident_users = data.jamfpro_mobile_device_inventory.shared_ipad.resident_users
  }
  description = "Shared iPad configuration"
}

# Example 11: Extension Attributes (Custom Fields)
data "jamfpro_mobile_device_inventory" "extension_attrs" {
  id = "808"
}

output "extension_attributes" {
  value = [
    for attr in data.jamfpro_mobile_device_inventory.extension_attrs.extension_attributes : {
      display_name = attr.display_name
      value        = attr.value
    }
  ]
  description = "Extension attributes (custom inventory fields)"
}

# Example 12: Complete Inventory Export for Reporting
data "jamfpro_mobile_device_inventory" "full_inventory" {
  id = "909"
}

output "complete_inventory_json" {
  value = jsonencode({
    id                = data.jamfpro_mobile_device_inventory.full_inventory.id
    udid              = data.jamfpro_mobile_device_inventory.full_inventory.udid
    display_name      = data.jamfpro_mobile_device_inventory.full_inventory.display_name
    serial_number     = data.jamfpro_mobile_device_inventory.full_inventory.serial_number
    model             = data.jamfpro_mobile_device_inventory.full_inventory.model
    os_version        = data.jamfpro_mobile_device_inventory.full_inventory.os_version
    managed           = data.jamfpro_mobile_device_inventory.full_inventory.managed
    supervised        = data.jamfpro_mobile_device_inventory.full_inventory.supervised
  })
  description = "Complete inventory as JSON for external processing"
}

# Example 13: Asset Management Check
data "jamfpro_mobile_device_inventory" "asset_check" {
  id = "1010"
}

output "asset_info" {
  value = {
    asset_tag     = data.jamfpro_mobile_device_inventory.asset_check.asset_tag
    serial        = data.jamfpro_mobile_device_inventory.asset_check.serial_number
    device_id     = data.jamfpro_mobile_device_inventory.asset_check.device_id
    management_id = data.jamfpro_mobile_device_inventory.asset_check.management_id
  }
  description = "Asset tracking information"
}

# Example 14: Backup Status
data "jamfpro_mobile_device_inventory" "backup_check" {
  id = "1111"
}

output "backup_status" {
  value = {
    cloud_backup_enabled  = data.jamfpro_mobile_device_inventory.backup_check.cloud_backup_enabled
    last_cloud_backup     = data.jamfpro_mobile_device_inventory.backup_check.last_cloud_backup_date
    last_backup_date      = data.jamfpro_mobile_device_inventory.backup_check.last_backup_date
  }
  description = "Backup status information"
}

# Example 15: Lost Mode Status
data "jamfpro_mobile_device_inventory" "lost_mode_check" {
  id = "1212"
}

output "lost_mode_status" {
  value = {
    lost_mode_enabled             = data.jamfpro_mobile_device_inventory.lost_mode_check.lost_mode_enabled
    lost_mode_enabled_date        = data.jamfpro_mobile_device_inventory.lost_mode_check.lost_mode_enabled_date
    device_locator_service        = data.jamfpro_mobile_device_inventory.lost_mode_check.device_locator_service_enabled
  }
  description = "Lost mode and device locator status"
}
