
data "jamfpro_sites" "example_site" {
  id = resource.jamfpro_sites.example_site.id
}

output "jamfpro_file_share_distribution_point_id" {
  value = data.jamfpro_sites.example_site.id
}

output "jamfpro_file_share_distribution_point_name" {
  value = data.jamfpro_sites.example_site.name
}