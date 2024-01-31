data "jamfpro_sites" "example_site" {
  name = "tf-example-site-01"  # Replace this with the actual name of the site you want to retrieve
}

output "site_id" {
  value = data.jamfpro_sites.example_site.id
}

output "site_name" {
  value = data.jamfpro_sites.example_site.name
}
