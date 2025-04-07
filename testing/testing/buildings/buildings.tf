// ========================================================================== //
// Buildings
// ========================================================================== //

// ========================================================================== //
// Single buildings

resource "jamfpro_building" "building_min" {
  name = "tf-testing-${var.testing_id}-min-${random_id.rng.hex}"
}

resource "jamfpro_building" "building_max" {
  name            = "tf-testing-${var.testing_id}-max-${random_id.rng.hex}"

  street_address1 = "unit 1"
  street_address2 = "1 example drive"
  city            = "Jamftown"
  state_province  = "Jamfshire"
  zip_postal_code = "JM22 5AM"
  country         = "Jamf Republic"
}

// ========================================================================== //
// Multiple buildings
resource "jamfpro_building" "building_multiple_max" {
  count           = 100
  name            = "tf-testing-${var.testing_id}-max-${count.index}-${random_id.rng.hex}"
  street_address1 = "unit 1"
  street_address2 = "1 example drive"
  city            = "Jamftown"
  state_province  = "Jamfshire"
  zip_postal_code = "JM22 5AM"
  country         = "Jamf Republic"
}

resource "jamfpro_building" "building_multiple_min" {
  count = 100
  name  = "tf-testing-${var.testing_id}-min-${count.index}-${random_id.rng.hex}"
}