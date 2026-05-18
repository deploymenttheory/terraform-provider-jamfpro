// ========================================================================== //
// Webhooks
// ========================================================================== //

// Basic webhook with no authentication
resource "jamfpro_webhook" "basic" {
  name         = "tf-testing-${var.testing_id}-webhook-basic-${random_id.rng.hex}"
  enabled      = true
  url          = "https://example.com/webhook"
  content_type = "application/json"
  event        = "ComputerAdded"
}

// Webhook for a smart group event where smart_group_id is sourced from a
// resource created in the same apply — validates that plan-time validation
// correctly defers when the ID is (known after apply).
resource "jamfpro_smart_computer_group_v2" "webhook_test" {
  name        = "tf-testing-${var.testing_id}-webhook-smart-group-${random_id.rng.hex}"
  description = "Smart group used to test webhook smart_group_id computed value handling."
  criteria {
    name        = "Serial Number"
    search_type = "not like"
    value       = "C0"
  }
}

resource "jamfpro_webhook" "smart_group_membership_change" {
  name           = "tf-testing-${var.testing_id}-webhook-smart-group-${random_id.rng.hex}"
  enabled        = true
  url            = "https://example.com/webhook"
  content_type   = "application/json"
  event          = "SmartGroupComputerMembershipChange"
  smart_group_id = jamfpro_smart_computer_group_v2.webhook_test.id
}
