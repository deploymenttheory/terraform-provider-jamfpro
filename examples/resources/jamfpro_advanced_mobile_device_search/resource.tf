resource "jamfpro_advanced_mobile_device_search" "advanced_mobile_device_search_001" {
  name    = "tf - Example Advanced mobile device Search 2"
  site_id = "-1" # Optional, defaults to "-1" for None

  criteria {
    name          = "Building"
    priority      = 0
    and_or        = "and"
    search_type   = "is"
    value         = "test"
    opening_paren = true
    closing_paren = false
  }

  criteria {
    name          = "iTunes Store Account"
    priority      = 0
    and_or        = "and"
    search_type   = "is"
    value         = "Active"
    opening_paren = false
    closing_paren = true
  }

  display_fields = [
  "Wi-Fi MAC Address", "Building", "iTunes Store Account", "Managed", "UDID"]
}
