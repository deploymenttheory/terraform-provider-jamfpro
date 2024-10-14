# Pop-up menu Example
resource "jamfpro_mobile_device_extension_attribute" "popup_menu_example" {
  name               = "Device Location"
  description        = "The primary location where this device is used"
  data_type          = "String"
  inventory_display  = "User and Location"
  
  input_type {
    type = "Pop-up Menu"
    popup_choices = [
      "Head Office",
      "Branch Office",
      "Home Office",
      "Client Site"
    ]
  }
}

# Text Field Example
resource "jamfpro_mobile_device_extension_attribute" "text_field_example" {
  name               = "User Department"
  description        = "The department to which the device user belongs"
  data_type          = "String"
  inventory_display  = "General"
  
  input_type {
    type = "Text Field"
  }
}