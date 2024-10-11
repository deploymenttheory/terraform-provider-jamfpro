# Text Field Example
resource "jamfpro_mobile_device_extension_attribute" "mobile_device_extension_attribute_text_field_1" {
  name          = "tf-example-text-field-ea"
  description   = "An attribute collected from a text field."
  data_type     = "STRING"
  inventory_display = "EXTENSION_ATTRIBUTES"
}