// ========================================================================== //
// Categories
// ========================================================================== //

// ========================================================================== //
// Single Categories

resource "jamfpro_category" "category_min" {
  name = "tf-testing-${var.testing_id}-min-${random_id.rng.hex}"
}

resource "jamfpro_category" "category_max" {
  name     = "tf-testing-${var.testing_id}-max-${random_id.rng.hex}"
  priority = 2
}

// ========================================================================== //
// Multiple Categories

resource "jamfpro_category" "category_min_multiple" {
  count = 100
  name  = "tf-testing-${var.testing_id}-min-${count.index}-${random_id.rng.hex}"
}

resource "jamfpro_category" "category_max_multiple" {
  count    = 100
  name     = "tf-testing-${var.testing_id}-max-${count.index}-${random_id.rng.hex}"
  priority = 2
}