# Example 1: Look up package by ID
data "jamfpro_package" "by_id" {
 id = "38" 
}

# Example 2: Look up package by name 
data "jamfpro_package" "by_name" {
 package_name = "Firefox 133.0.3.pkg"
}

# Package Details Output
output "firefox_details" {
 value = {
   id                   = data.jamfpro_package.by_name.id
   package_name         = data.jamfpro_package.by_name.package_name
   filename             = data.jamfpro_package.by_name.filename
   package_file_source  = data.jamfpro_package.by_name.package_file_source
   category_id          = data.jamfpro_package.by_name.category_id
   info                 = data.jamfpro_package.by_name.info
   notes                = data.jamfpro_package.by_name.notes 
   priority             = data.jamfpro_package.by_name.priority
   os_requirements      = data.jamfpro_package.by_name.os_requirements
   fill_user_template   = data.jamfpro_package.by_name.fill_user_template
   indexed              = data.jamfpro_package.by_name.indexed
   fill_existing_users  = data.jamfpro_package.by_name.fill_existing_users
   swu                  = data.jamfpro_package.by_name.swu
   reboot_required      = data.jamfpro_package.by_name.reboot_required
   self_heal_notify     = data.jamfpro_package.by_name.self_heal_notify
   self_healing_action  = data.jamfpro_package.by_name.self_healing_action  
   os_install          = data.jamfpro_package.by_name.os_install
   serial_number       = data.jamfpro_package.by_name.serial_number
   parent_package_id   = data.jamfpro_package.by_name.parent_package_id
   base_path          = data.jamfpro_package.by_name.base_path
   suppress_updates   = data.jamfpro_package.by_name.suppress_updates
   cloud_transfer_status = data.jamfpro_package.by_name.cloud_transfer_status
   ignore_conflicts   = data.jamfpro_package.by_name.ignore_conflicts
   suppress_from_dock = data.jamfpro_package.by_name.suppress_from_dock
   suppress_eula     = data.jamfpro_package.by_name.suppress_eula
   suppress_registration = data.jamfpro_package.by_name.suppress_registration
   install_language  = data.jamfpro_package.by_name.install_language
   md5              = data.jamfpro_package.by_name.md5
   sha256           = data.jamfpro_package.by_name.sha256
   hash_type        = data.jamfpro_package.by_name.hash_type
   hash_value       = data.jamfpro_package.by_name.hash_value
   size             = data.jamfpro_package.by_name.size
   os_installer_version = data.jamfpro_package.by_name.os_installer_version
   manifest         = data.jamfpro_package.by_name.manifest
   manifest_file_name = data.jamfpro_package.by_name.manifest_file_name
   format           = data.jamfpro_package.by_name.format
   package_uri      = data.jamfpro_package.by_name.package_uri
   md5_file_hash    = data.jamfpro_package.by_name.md5_file_hash
 }
}

# Example 3: Using variables
variable "package_name" {
 type = string
 description = "Name of the package to look up"
 default = "Firefox 133.0.3.pkg"
}

data "jamfpro_package" "dynamic" {
 package_name = var.package_name
}

# Example 4: Using in another resource
resource "jamfpro_policy" "firefox_install" {
 name = "Install Firefox"
 enabled = true
 package_id = data.jamfpro_package.by_name.id
 description = "Deploys Firefox version ${data.jamfpro_package.by_name.package_name}"
}