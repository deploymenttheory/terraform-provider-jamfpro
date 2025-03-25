// ========================================================================== //
// Buildings
// ========================================================================== //

// ========================================================================== //
// Single buildings

resource "jamfpro_building" "building_min" {
  name = "tf-testing-min"
}

resource "jamfpro_building" "building_max" {
  name = "tf-testing-max"
  street_address1 = "unit 1"
  street_address2 = "1 example drive"
  city = "Jamftown"
  state_province = "Jamfshire"
  zip_postal_code = "JM22 5AM"
  country = "Jamf Republic"
}

// ========================================================================== //
// Multiple buildings

resource "jamfpro_building" "building_multiple_max" {
  count = 100
  name = "tf-testing-max-${count.index}"
  street_address1 = "unit 1"
  street_address2 = "1 example drive"
  city = "Jamftown"
  state_province = "Jamfshire"
  zip_postal_code = "JM22 5AM"
  country = "Jamf Republic"
}

resource "jamfpro_building" "building_multiple_min" {
  count = 100
  name = "tf-testing-min-${count.index}"
}