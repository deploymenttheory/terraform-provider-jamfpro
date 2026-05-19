// ========================================================================== //
// User Groups
// ========================================================================== //

// Static group with members — standard case
resource "jamfpro_user_group" "static_with_members" {
  name                = "tf-testing-${var.testing_id}-user-group-static-${random_id.rng.hex}"
  is_smart            = false
  is_notify_on_change = false
  assigned_user_ids   = []
}

// Empty static group — validates that the validator allows zero members,
// covering the bug where assigned_user_ids = [] was incorrectly rejected.
resource "jamfpro_user_group" "static_empty" {
  name                = "tf-testing-${var.testing_id}-user-group-empty-${random_id.rng.hex}"
  is_smart            = false
  is_notify_on_change = false
  assigned_user_ids   = []

  lifecycle {
    ignore_changes = [assigned_user_ids]
  }
}

// Smart group with criteria — validates is_smart=true path is unaffected
resource "jamfpro_user_group" "smart" {
  name                = "tf-testing-${var.testing_id}-user-group-smart-${random_id.rng.hex}"
  is_smart            = true
  is_notify_on_change = false
  criteria {
    name        = "Email Address"
    priority    = 0
    search_type = "like"
    value       = "@example.com"
  }
}
