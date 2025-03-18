# resource "jamfpro_script" "min_script" {
#   name = "tf-testing-local-bw-script-min"
#   script_contents = "script_contents_field"
#   priority = "BEFORE"
# }

# resource "jamfpro_script" "max_script" {
#   name = "tf-testing-local-bw-script-max"
#   category_id = "9"
#   info = "info_field"
#   notes = "notes_field"
#   os_requirements = "os_requirements_field"
#   priority = "BEFORE"
#   script_contents = "script_contents_field"
#   parameter4 = "parameter4_field"
#   parameter5 = "parameter5_field"
#   parameter6 = "parameter6_field"
#   parameter7 = "parameter7_field"
#   parameter8 = "parameter8_field"
#   parameter9 = "parameter9_field"
#   parameter10 = "parameter10_field"
#   parameter11 = "parametee11_field"
# }

# resource "jamfpro_script" "script1" {
#   count = 10
#   name = "tf-testing-local-bw-${count.index}"
#   script_contents = "echo hello world"
#   priority = "BEFORE"
# }

resource "jamfpro_building" "building1" {
  count = 10
  name = "tf-testing-local-bw-${count.index}"
}

# resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_popup_menu_1" {
#   count = 10
#   name                   = "tf-testing-local-bw-${count.index}"
#   enabled                = true
#   description            = "An attribute collected from a pop-up menu."
#   input_type             = "POPUP"
#   popup_menu_choices     = ["Option 1", "Option 2", "Option 3"]
#   inventory_display_type = "USER_AND_LOCATION"
#   data_type              = "STRING"
# }

// ========================================================================== //
// Packages 

# resource "jamfpro_package" "jamfpro_package_002" {
#   package_name          = "btf-testing-local-bw"                                                   // Required
#   package_file_source   = "https://github.com/obsidianmd/obsidian-releases/releases/download/v1.8.9/Obsidian-1.8.9.dmg" // Required
#   priority              = 10                                                                    // Required
#   reboot_required       = true                                                                  // Required
#   fill_user_template    = false                                                                 // Required
#   fill_existing_users   = false                                                                 // Required
#   os_install            = false                                                                 // Required
#   suppress_updates      = false                                                                 // Required
#   suppress_from_dock    = false                                                                 // Required
#   suppress_eula         = false                                                                 // Required
#   suppress_registration = false                                                                // Required
#   timeouts {
#     create = "90m" // Optional / Useful for large packages uploads
#   }
# }
