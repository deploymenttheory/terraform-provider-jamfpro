resource "jamfpro_computer_group" "smart_example" {
  name = "Example Smart Computer Group"
  is_smart = true  # Set to true if the group is a Smart group, otherwise set to false

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

resource "jamfpro_computer_group" "static_example" {
  name = "Example Static Computer Group"
  is_smart = false  # Set to true if the group is a Smart group, otherwise set to false

  # Optional: Specify site details
  site {
    id   = 123  # Replace with the actual site ID
    name = "Site Name"  # Replace with the actual site name
  }

  # Optional: Specify computers for static groups
  computers {
    id             = 456  # Replace with the actual computer ID
    name           = "Computer Name"  # Replace with the actual computer name
    serial_number  = "ABC123"         # Replace with the actual serial number
    mac_address    = "00:11:22:33:44:55"  # Replace with the actual MAC address
    alt_mac_address = "AA:BB:CC:DD:EE:FF"  # Replace with the actual alternative MAC address
  }
}