resource "jamfpro_advanced_computer_search" "advanced_computer_search_001" {
  name    = "tf-testing-${var.testing_id}-max-script-${random_id.rng.hex}"
  view_as = "Standard Web Page"

  sort1 = "Serial Number"
  sort2 = "Username"
  sort3 = "Department"

  criteria {
    name          = "Building"
    priority      = 0
    and_or        = "and"
    search_type   = "is"
    value         = "square"
    opening_paren = true
    closing_paren = false
  }

  criteria {
    name          = "Model"
    priority      = 1
    and_or        = "and"
    search_type   = "is"
    value         = "macbook air"
    opening_paren = false
    closing_paren = true
  }

  criteria {
    name          = "Computer Name"
    priority      = 2
    and_or        = "or"
    search_type   = "matches regex"
    value         = "thing"
    opening_paren = true
    closing_paren = false
  }

  criteria {
    name          = "Licensed Software"
    priority      = 3
    and_or        = "and"
    search_type   = "has"
    value         = "office"
    opening_paren = false
    closing_paren = true
  }

  site_id = "-1"

  display_fields = [
    "Activation Lock Manageable",
    "Apple Silicon",
    "Architecture Type",
    "Available RAM Slots"
  ]

}