// ========================================================================== //
// Smart Mobile Device Groups
// ========================================================================== //

resource "jamfpro_smart_mobile_device_group_v1" "name" {
  name        = "tf-testing-${var.testing_id}-script-max-v1${random_id.rng.hex}"
  description = "Terraform testing smart mobile device group."
  criteria {
    name        = "Serial Number"
    search_type = "not like"
    value       = "C0"
  }
  criteria {
    name          = "OS Version"
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
    value         = "iPhone"
    closing_paren = true
  }
}
