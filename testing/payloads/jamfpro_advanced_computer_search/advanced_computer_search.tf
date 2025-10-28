resource "jamfpro_advanced_computer_search" "advanced_computer_search_001" {
  name = "tf-testing-${var.testing_id}-max-script-${random_id.rng.hex}"
}
