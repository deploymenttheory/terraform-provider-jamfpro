// ========================================================================== //
// Account Group (resource + data source)
// ========================================================================== //
//
// Deploys an account group, then reads it back through the jamfpro_account_group
// data source both by id and by name. The by-name lookup exercises the fix in
// this PR (the data source previously only supported lookup by id).

resource "jamfpro_account_group" "test" {
  name          = "tf-testing-${var.testing_id}-account-group-${random_id.rng.hex}"
  access_level  = "Full Access"
  privilege_set = "Administrator"
}

# Read the account group back by its id.
data "jamfpro_account_group" "by_id" {
  id = jamfpro_account_group.test.id

  depends_on = [jamfpro_account_group.test]
}

# Read the account group back by its name (the behaviour fixed in this PR).
data "jamfpro_account_group" "by_name" {
  name = jamfpro_account_group.test.name

  depends_on = [jamfpro_account_group.test]
}

output "jamfpro_account_group_by_id_name" {
  description = "Name resolved via the id lookup."
  value       = data.jamfpro_account_group.by_id.name
}

output "jamfpro_account_group_by_name_id" {
  description = "Id resolved via the name lookup."
  value       = data.jamfpro_account_group.by_name.id
}
