// ========================================================================== //
// Computer Extension Attributes
// ========================================================================== //

// ========================================================================== //
// UNTESTED RESOURCES
# - input_type  = "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"

// ========================================================================== //
// Single extension attributes

resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_script" {
  name = "tf-testing-max-script"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "SCRIPT"
  script_contents  =  "echo hello"
  ldap_extension_attribute_allowed = false
}
  resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_text" {
  name = "tf-testing-max-text"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "TEXT"
  ldap_extension_attribute_allowed = false
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_popup" {
  name = "tf-testing-max-popup"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "POPUP"
  popup_menu_choices  = ["Option 1", "Option 2", "Option 3"]
  ldap_extension_attribute_allowed = false
}

// ========================================================================== //
// Multiple extension attributes

resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_script_multiple" {
  count = 100
  name = "tf-testing-max-script-${count.index}"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "SCRIPT"
  script_contents  =  "echo hello"
  ldap_extension_attribute_allowed = false
}
  resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_text_multiple" {
  count = 100
  name = "tf-testing-max-text-${count.index}"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "TEXT"
  ldap_extension_attribute_allowed = false
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_popup_multiple" {
  count = 100
  name = "tf-testing-max-popup-${count.index}"
  description = "description"
  data_type = "STRING"
  enabled = true
  inventory_display_type  = "GENERAL"
  input_type  = "POPUP"
  popup_menu_choices  = ["Option 1", "Option 2", "Option 3"]
  ldap_extension_attribute_allowed = false
}