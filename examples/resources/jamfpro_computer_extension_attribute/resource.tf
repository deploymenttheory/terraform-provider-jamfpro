# Pop-up Menu Example
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_popup_menu_1" {
  name                    = "tf-ghatest-cexa-popup-menu-example"
  enabled                 = true
  description             = "An attribute collected from a pop-up menu."
  input_type              = "POPUP"
  popup_menu_choices      = ["Option 1", "Option 2", "Option 3"]
  inventory_display_type  = "USER_AND_LOCATION"
  data_type               = "STRING"
}

# Text Field Example
resource "jamfpro_computer_extension_attribute" "computer_extension_attribute_text_field_1" {
  name                    = "tf-example-cexa-text-field-example"
  enabled                 = true
  description             = "An attribute collected from a text field."
  input_type              = "TEXT"
  inventory_display_type  = "HARDWARE"
  data_type               = "STRING"
}

# Script Example
resource "jamfpro_computer_extension_attribute" "computer_extension_attribute_script_1" {
  name                    = "tf-example-cexa-hello-world"
  enabled                 = true
  description             = "An attribute collected via a script."
  input_type              = "SCRIPT"
  script_contents         = "#!/bin/bash\necho 'Hello, World!!!!! :)'"
  inventory_display_type  = "GENERAL"
  data_type               = "STRING"
}