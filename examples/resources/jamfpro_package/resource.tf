// Example with package_file_source from URL or local file (package is uploaded to Cloud Distribution Point) 
resource "jamfpro_package" "jamfpro_package_002" {
  package_name          = "your-package-name"                                                   // Required
  package_file_source   = "/path/to/your/package/file/.pkg or .dmg , or http(s)://path/to/file" // Required
  category_id           = "your-category-id"                                                    // Optional /  jamfpro_category.jamfpro_category_001.id
  info                  = "tf package deployment for demonstration"                             // Optional
  notes                 = "Uploaded by: terraform-provider-jamfpro plugin."                     // Optional
  priority              = 10                                                                    // Required
  reboot_required       = true                                                                  // Required
  fill_user_template    = false                                                                 // Required
  fill_existing_users   = false                                                                 // Required
  os_requirements       = "macOS 10.15.7, macOS 11.1"                                           // Optional
  swu                   = false                                                                 // optional
  self_heal_notify      = false                                                                 // optional
  os_install            = false                                                                 // Required
  serial_number         = ""                                                                    // Optional
  suppress_updates      = false                                                                 // Required
  ignore_conflicts      = false                                                                 // Optional
  suppress_from_dock    = false                                                                 // Required
  suppress_eula         = false                                                                 // Required
  suppress_registration = false                                                                 // Required
  manifest              = ""                                                                    // Optional
  manifest_file_name    = ""                                                                    // Optional
  timeouts {
    create = "90m" // Optional / Useful for large packages uploads
  }
}

// Example without package_file_source (when package exists on File Share Distribution Point only)
// Package metadata only is created
resource "jamfpro_package" "jamfpro_package_003" {
  package_name          = "your-package-name"
  filename              = "your-package-name.pkg"
  priority              = 10
  fill_user_template    = false
  reboot_required       = false
  os_install            = false
  suppress_updates      = false
  suppress_from_dock    = false
  suppress_eula         = false
  suppress_registration = false

  // Optional: supply known hashes for package validation
  md5        = "d41d8cd98f00b204e9800998ecf8427e"
  sha256     = "e3b0c44298fc1c149afbf4c8996fb924..."
  sha3512    = "6b86b273ff34fce19d6b804eff5a3f57..."
  hash_type  = "SHA3_512"
  hash_value = "6b86b273ff34fce19d6b804eff5a3f57..."
}
