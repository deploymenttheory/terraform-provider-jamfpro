// ========================================================================== //
// Departments
// ========================================================================== //

# Departments only consist of a name so cannot be min max tested.

// ========================================================================== //
// Single department

resource "jamfpro_department" "department" {
  name = "tf-testing-${random_id.rng.hex}"
}

// ========================================================================== //
// Multiple departments

resource "jamfpro_department" "multiple_departments" {
  count = 100
  name  = "tf-testing-${count.index}-${random_id.rng.hex}"
}