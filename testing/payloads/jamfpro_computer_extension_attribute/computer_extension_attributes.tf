// ========================================================================== //
// Computer Extension Attributes
// ========================================================================== //

// ========================================================================== //
// UNTESTED RESOURCES
# - input_type  = "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"

// ========================================================================== //
// Single extension attributes

resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_script" {
  name                   = "tf-testing-${var.testing_id}-max-script-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "SCRIPT"
  script_contents        = "echo hello"
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_text" {
  name                   = "tf-testing-${var.testing_id}-max-text-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "TEXT"
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_popup" {
  name                   = "tf-testing-${var.testing_id}-max-popup-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "POPUP"
  popup_menu_choices     = ["Option 1", "Option 2", "Option 3"]
}

// Test that minimally defined object is sufficient
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_minimum" {
  name       = "tf-testing-${var.testing_id}-min-test"
  enabled    = true
  input_type = "TEXT"
}

// Regression test for issue #1145 - heredoc strings in HCL always include a
// trailing newline before EOT, but the API strips it server-side. Without a
// DiffSuppressFunc on description, this produced perpetual drift on every
// plan after apply.
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_heredoc_description" {
  name        = "tf-testing-${var.testing_id}-heredoc-${random_id.rng.hex}"
  enabled     = true
  input_type  = "TEXT"
  description = <<-EOT
    Multi-line description used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
}

// ========================================================================== //
// Multiple extension attributes

resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_script_multiple" {
  count                  = 100
  name                   = "tf-testing-${var.testing_id}-max-script-${count.index}-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "SCRIPT"
  script_contents        = "echo hello"
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_text_multiple" {
  count                  = 100
  name                   = "tf-testing-${var.testing_id}-max-text-${count.index}-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "TEXT"
}
resource "jamfpro_computer_extension_attribute" "jamfpro_computer_extension_attribute_max_popup_multiple" {
  count                  = 100
  name                   = "tf-testing-${var.testing_id}-max-popup-${count.index}-${random_id.rng.hex}"
  description            = "description"
  data_type              = "STRING"
  enabled                = true
  inventory_display_type = "GENERAL"
  input_type             = "POPUP"
  popup_menu_choices     = ["Option 1", "Option 2", "Option 3"]
}
