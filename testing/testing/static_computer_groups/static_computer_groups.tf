// ========================================================================== //
// Static Computer Groups
// ========================================================================== //

# This test resource relies on a manually created test site. Sites are not yet
# terraformed.

# Site ID= 2859
# Computer ID= 23
variable "site_id" {
  description = "site id used."
}
resource "jamfpro_static_computer_group" "name" {
  name="tf-testing-script-max-${random_id.rng.hex}"
  site_id = var.site_id
  assigned_computer_ids = [23]
}
