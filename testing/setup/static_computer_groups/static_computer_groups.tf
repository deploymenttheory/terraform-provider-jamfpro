// ========================================================================== //
// Setup for static computer groups
// ========================================================================== //

resource "jamfpro_site" "static_computer_group_site" {
  name = "tf-testing-${random_id.rng.hex}"
}
output "site_id" {
  value = jamfpro_site.static_computer_group_site.id
}