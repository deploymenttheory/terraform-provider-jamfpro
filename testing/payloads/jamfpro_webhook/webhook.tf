// ========================================================================== //
// Webhooks
// ========================================================================== //

// ========================================================================== //
// Smart-group-membership-change webhook with smart_group_id referenced from a
// resource being created in the same apply. This is the case fixed by the
// custom-diff change to validateSmartGroupIDRequirement — without that fix,
// plan errors with "smart_group_id must be provided and must be a valid
// non-zero integer" because diff.GetOk() can't distinguish unknown from unset.

resource "jamfpro_smart_computer_group_v2" "webhook_target_group" {
  name        = "tf-testing-${var.testing_id}-webhook-target-${random_id.rng.hex}"
  description = "Smart group used as a webhook target in the same apply."
  criteria {
    name        = "Serial Number"
    search_type = "not like"
    value       = "C0"
  }
}

resource "jamfpro_webhook" "smart_group_event_referenced_id" {
  name               = "tf-testing-${var.testing_id}-sg-referenced-${random_id.rng.hex}"
  enabled            = true
  url                = "https://example.com/webhook/sg-referenced"
  content_type       = "application/json"
  event              = "SmartGroupComputerMembershipChange"
  smart_group_id     = jamfpro_smart_computer_group_v2.webhook_target_group.id
  connection_timeout = 5
  read_timeout       = 5
}

// ========================================================================== //
// Non-smart-group event — sanity check that the unrelated path still applies
// without requiring a smart_group_id.

resource "jamfpro_webhook" "non_smart_group_event" {
  name               = "tf-testing-${var.testing_id}-inv-completed-${random_id.rng.hex}"
  enabled            = true
  url                = "https://example.com/webhook/inventory"
  content_type       = "application/json"
  event              = "ComputerInventoryCompleted"
  connection_timeout = 5
  read_timeout       = 5
}
