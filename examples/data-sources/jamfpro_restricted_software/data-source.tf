data "jamfpro_restricted_software" "restricted_software_001_data" {
  id = jamfpro_restricted_software.restricted_software_001.id
}

output "jamfpro_restricted_software_001_id" {
  value = data.jamfpro_restricted_software.restricted_software_001_data.id
}

output "jamfpro_restricted_software_001_name" {
  value = data.jamfpro_restricted_software.restricted_software_001_data.name
}

data "jamfpro_restricted_software" "restricted_software_002_data" {
  id = jamfpro_restricted_software.restricted_software_002.id
}

output "jamfpro_restricted_software_002_id" {
  value = data.jamfpro_restricted_software.restricted_software_002_data.id
}

output "jamfpro_restricted_software_002_name" {
  value = data.jamfpro_restricted_software.restricted_software_002_data.name
}