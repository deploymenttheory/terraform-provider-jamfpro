
resource "jamfpro_advanced_user_search" "advanced_user_search_001" {
  name          = "advanced search name"

  criteria {
    name          = "Email Address"
    priority      = 0
    and_or        = "and"
    search_type   = "like"
    value         = "company.com"
    opening_paren = false
    closing_paren = false
  }

  display_fields =[
    "Content ID",
    "Email Address",
    "Managed Apple ID",
    "Roster Class Display Name"
  ]

}


