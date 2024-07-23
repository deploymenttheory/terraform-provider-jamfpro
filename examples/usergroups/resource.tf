// Static user group example
resource "jamfpro_user_group" "jamfpro_user_group_001" {
  name                = "tf-localtest-user-group-static-01"
  is_smart            = false
  is_notify_on_change = false

  site_id = 5

  assigned_user_ids = [1]
}


// Dynamic user group example
resource "jamfpro_user_group" "jamfpro_user_group_002" {
  name                = "tf-ghatest-usergroup-dynamic"
  is_smart            = true
  is_notify_on_change = true

  site_id = 5

  criteria {
    name          = "Email Address"
    priority      = 0
    and_or        = "and"
    search_type   = "like"
    value         = "company.com"
    opening_paren = false
    closing_paren = false
  }
}

// Dynamic user group example with multiple criteria and no site
resource "jamfpro_user_group" "jamfpro_user_group_003" {
  name                = "tf-ghatest-usergroup-dynamic-testing"
  is_smart            = true
  is_notify_on_change = true
  criteria {
    name          = "Email Address"
    priority      = 0
    and_or        = "and"
    search_type   = "like"
    value         = "company.com"
    opening_paren = false
    closing_paren = false
  }
  criteria {
    name          = "Managed Apple ID"
    priority      = 1
    and_or        = "and"
    search_type   = "like"
    value         = "company.com"
    opening_paren = false
    closing_paren = false
  }
}