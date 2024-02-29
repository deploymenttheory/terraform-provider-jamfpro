data "jamfpro_script" "script_001_data" {
  id = jamfpro_script.script_001.id
}

output "jamfpro_script_001_id" {
  value = data.jamfpro_script.script_001_data.id
}

output "jamfpro_script_001_name" {
  value = data.jamfpro_script.script_001_data.name
}

data "jamfpro_script" "script_002_data" {
  id = jamfpro_script.script_002.id
}

output "jamfpro_script_002_id" {
  value = data.jamfpro_script.script_002_data.id
}

output "jamfpro_script_002_name" {
  value = data.jamfpro_script.script_002_data.name
}