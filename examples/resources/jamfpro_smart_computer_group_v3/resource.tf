# Requires Jamf Pro 11.28.0 or later.
resource "jamfpro_smart_computer_group_v3" "example" {
  name        = "Managed Macs running macOS 26"
  description = "Computers matching both criteria."
  site_id     = "-1" # None

  criteria {
    name        = "Managed"
    priority    = 0
    search_type = "is"
    value       = "Managed"
  }

  criteria {
    name        = "Operating System Version"
    priority    = 1
    and_or      = "and"
    search_type = "like"
    value       = "26."
  }

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Import using the numeric Jamf Pro group ID, not its Platform UUID:
# terraform import jamfpro_smart_computer_group_v3.example 123
