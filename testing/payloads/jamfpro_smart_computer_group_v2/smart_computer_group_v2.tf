// ========================================================================== //
// Smart Computer Groups
// ========================================================================== //

resource "jamfpro_smart_computer_group_v2" "name" {
  name        = "tf-testing-${var.testing_id}-script-max-v2-${random_id.rng.hex}"
  description = "Terraform testing smart computer group V2."
  criteria {
    name        = "Serial Number"
    search_type = "not like"
    value       = "C0"
  }
  criteria {
    name          = "Operating System Version"
    priority      = 1
    and_or        = "and"
    search_type   = "is"
    value         = "15.1"
    opening_paren = true
  }
  criteria {
    name          = "Model"
    priority      = 2
    and_or        = "or"
    search_type   = "like"
    value         = "Book"
    closing_paren = true
  }
}
