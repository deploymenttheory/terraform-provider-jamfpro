
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_popup_menu_1" {
  name        = "tf-ghatest-cexa-popup-menu-example"
  enabled     = true
  description = "An attribute collected from a pop-up menu."
  data_type   = "String"

  input_type  = "Pop-up Menu"
  input_popup = ["Option 1", "Option 2", "Option 3"]


  inventory_display = "User and Location"
}

# //-------------------------------------------------------------------//

resource "jamfpro_computer_extension_attribute" "computer_extension_attribute_text_field_1" {
  name        = "tf-example-cexa-text-field-example"
  enabled     = true
  description = "An attribute collected from a text field."
  data_type   = "String"

  input_type  = "Text Field"
  inventory_display = "Hardware"
}

# //-------------------------------------------------------------------//

resource "jamfpro_computer_extension_attribute" "computer_extension_attribute_script_1" {
  name        = "tf-example-cexa-hello-world"
  enabled     = true
  description = "An attribute collected via a script."
  data_type   = "String"

  input_type = "script"
  input_script   = "#!/bin/bash\necho 'Hello, World!!!!! :)'"

  inventory_display = "General"
}