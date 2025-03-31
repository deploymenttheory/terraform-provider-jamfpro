// ========================================================================== //
// Static Computer Groups
// ========================================================================== //

resource "jamfpro_static_computer_group" "name" {
  name="tf-testing-script-max-${random_id.rng.hex}"
  site_id = 234
  assigned_computer_ids = [1,2,3,4,5,6]
}