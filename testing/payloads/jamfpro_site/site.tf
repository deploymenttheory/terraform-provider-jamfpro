// ========================================================================== //
// Site (resource + data source)
// ========================================================================== //
//
// Deploys a site, then reads it back through the jamfpro_site data source both
// by id and by name. The by-name lookup exercises the fix in this PR (the data
// source previously only supported lookup by id).

resource "jamfpro_site" "test" {
  name = "tf-testing-${var.testing_id}-site-${random_id.rng.hex}"
}

# Read the site back by its id.
data "jamfpro_site" "by_id" {
  id = jamfpro_site.test.id

  depends_on = [jamfpro_site.test]
}

# Read the site back by its name (the behaviour fixed in this PR).
data "jamfpro_site" "by_name" {
  name = jamfpro_site.test.name

  depends_on = [jamfpro_site.test]
}

output "jamfpro_site_by_id_name" {
  description = "Name resolved via the id lookup."
  value       = data.jamfpro_site.by_id.name
}

output "jamfpro_site_by_name_id" {
  description = "Id resolved via the name lookup."
  value       = data.jamfpro_site.by_name.id
}
