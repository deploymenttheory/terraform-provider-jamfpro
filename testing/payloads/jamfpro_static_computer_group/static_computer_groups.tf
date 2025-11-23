// ========================================================================== //
// Static Computer Groups
// ========================================================================== //

resource "jamfpro_static_computer_group" "name" {
  name        = "tf-testing-${var.testing_id}-script-max-${random_id.rng.hex}"
  description = "Terraform Static Computer Group for testing purposes."
}
