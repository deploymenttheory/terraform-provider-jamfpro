data "jamfpro_scripts" "script_001_data" {
  id = jamfpro_scripts.script_001.id
}

output "jamfpro_script_001_id" {
  value = data.jamfpro_scripts.script_001_data.id
}

output "jamfpro_script_001_name" {
  value = data.jamfpro_scripts.script_001_data.name
}

data "jamfpro_scripts" "script_002_data" {
  id = jamfpro_scripts.script_002.id
}

output "jamfpro_script_002_id" {
  value = data.jamfpro_scripts.script_002_data.id
}

output "jamfpro_script_002_name" {
  value = data.jamfpro_scripts.script_002_data.name
}