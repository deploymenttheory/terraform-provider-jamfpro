data "jamfpro_sites" "site_001_data" {
  id = jamfpro_sites.site_001.id
}

output "jamfpro_site_001_id" {
  value = data.jamfpro_sites.site_001_data.id
}

output "jamfpro_site_001_name" {
  value = data.jamfpro_sites.site_001_data.name
}

data "jamfpro_sites" "site_002_data" {
  id = jamfpro_sites.site_002.id
}

output "jamfpro_site_002_id" {
  value = data.jamfpro_sites.site_002_data.id
}

output "jamfpro_site_002_name" {
  value = data.jamfpro_sites.site_002_data.name
}