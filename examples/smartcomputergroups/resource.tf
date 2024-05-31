resource "jamfpro_smart_computer_group" "smart_example" {
  name = "Example Smart Computer Group"
  # Optional: Specify site details
  site {
    id   = 123  # Replace with the actual site ID
    name = "Site Name"  # Replace with the actual site name
  }

  # Optional: Define criteria for Smart groups
  criteria {
    name          = "Criterion Name #1"
    priority      = 1
    and_or        = "and"  # or "or", defaults to "and" if not provided
    search_type   = "is"   # or any other supported search type
    value         = "Criterion Value"
    opening_paren = false  # true or false, defaults to false if not provided
    closing_paren = false  # true or false, defaults to false if not provided
  }

  criteria {
    name          = "Criterion Name #n+1 etc"
    priority      = 1
    and_or        = "and"  # or "or", defaults to "and" if not provided
    search_type   = "is"   # or any other supported search type
    value         = "Criterion Value"
    opening_paren = false  # true or false, defaults to false if not provided
    closing_paren = false  # true or false, defaults to false if not provided
  }

}
