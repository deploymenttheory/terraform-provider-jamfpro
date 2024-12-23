# Create smart group first 
resource "jamfpro_smart_computer_group" "test_group" {
  name = "Test Smart Group"
  is_smart = true
  site_id = "1"

  criteria {
    name = "Operating System Version"
    search_type = "is"
    value = "13.0"
    and_or = "and"
    opening_paren = false
    closing_paren = false
  }

  criteria {
    name = "Computer Model" 
    search_type = "like"
    value = "MacBook Pro"
    and_or = "and"
    opening_paren = false
    closing_paren = false
  }
}

# Query by ID
data "jamfpro_smart_computer_group" "by_id" {
  smart_group_id = jamfpro_smart_computer_group.test_group.id
  depends_on = [jamfpro_smart_computer_group.test_group]
}

# Query by name
data "jamfpro_smart_computer_group" "by_name" {
  name = jamfpro_smart_computer_group.test_group.name
  depends_on = [jamfpro_smart_computer_group.test_group]
}

# Verify output
output "group_verification" {
  value = {
    by_id = {
      id = data.jamfpro_smart_computer_group.by_id.smart_group_id
      name = data.jamfpro_smart_computer_group.by_id.name
      is_smart = data.jamfpro_smart_computer_group.by_id.is_smart
      site_id = data.jamfpro_smart_computer_group.by_id.site_id
      criteria = data.jamfpro_smart_computer_group.by_id.criteria
    }
    by_name = {
      id = data.jamfpro_smart_computer_group.by_name.smart_group_id
      name = data.jamfpro_smart_computer_group.by_name.name
      is_smart = data.jamfpro_smart_computer_group.by_name.is_smart
      site_id = data.jamfpro_smart_computer_group.by_name.site_id
      criteria = data.jamfpro_smart_computer_group.by_name.criteria
    }
  }
}