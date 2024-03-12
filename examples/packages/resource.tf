// Definition a Jamf Pro Package Resource
resource "jamfpro_package" "jamfpro_package_001" {
  name                          = "tf-example-package-your-package-name"
  package_file_path             = file("path/to/your.pkg")
  category                      = "" // Optional / can be explitly set or a reference to a jamfpro_category resource e.g jamfpro_category.jamfpro_category_001.name
  info                          = "tf package deployment for demonstration"
  notes                         = "This package is used for Terraform provider documentation example."
  priority                      = 10
  reboot_required               = false
  fill_user_template            = false // Optional for .dmg files only
  fill_existing_users           = false // Optional for .dmg files only
  boot_volume_required          = false
  allow_uninstalled             = false
  os_requirements               = "macOS 10.15.7, macOS 11.1"
  required_processor            = ""
  switch_with_package           = ""
  install_if_reported_available = false
  reinstall_option              = ""
  triggering_files              = ""
  send_notification             = true
}