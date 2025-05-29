// ========================================================================== //
// Static Computer Groups
// ========================================================================== //

data "local_file" "site_and_computer_ids" {
  filename = "testing/data_sources/site_and_computer_ids.json"
}

resource "jamfpro_static_computer_group" "name" {
  name                  = "tf-testing-${var.testing_id}-script-max-${random_id.rng.hex}"
  site_id               = jsondecode(data.local_file.site_and_computer_ids.content).site
  assigned_computer_ids = jsondecode(data.local_file.site_and_computer_ids.content).computers
}
