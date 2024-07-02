// Definition a Jamf Pro Package Resource
resource "jamfpro_package" "jamfpro_package_002" {
  package_name                  = "your-package-name" // Required
  package_file_source           = "/path/to/your/package/file/.pkg or .dmg , or http(s)://path/to/file" // Required
  category_id                   = "your-category-id" // Required /  jamfpro_category.jamfpro_category_001.id
  info                          = "tf package deployment for demonstration" // Optional
  notes                         = "Uploaded by: terraform-provider-jamfpro plugin." // Optional
  priority                      = 10 // Required
  reboot_required               = true // Required
  fill_user_template            = false // Required
  fill_existing_users           = false // Required
  os_requirements               = "macOS 10.15.7, macOS 11.1" // Required
  swu                           = false // optional
  self_heal_notify              = false // optional
  os_install                    = false // Required
  serial_number                 = "" // optional
  suppress_updates              = false // Required
  ignore_conflicts              = false
  suppress_from_dock            = false // Required
  suppress_eula                 = false // Required
  suppress_registration         = false // Required
  manifest                      = ""
  manifest_file_name            = ""
  timeouts {
    create                      = "90m" // Optional / Useful for large packages uploads
  }
}