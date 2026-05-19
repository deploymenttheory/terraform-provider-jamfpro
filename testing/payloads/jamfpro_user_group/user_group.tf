// ========================================================================== //
// User Groups
// ========================================================================== //

resource "jamfpro_user_group" "office_users" {
  name                = "tf-testing-${var.testing_id}-office-users-${random_id.rng.hex}"
  is_smart            = false
  is_notify_on_change = false

  assigned_user_ids = []
}
