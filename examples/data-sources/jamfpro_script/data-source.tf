// data source by id
data "jamfpro_script" "script_001_data" {
  id = jamfpro_script.script_001.id
}

output "jamfpro_script_001_id" {
  value = data.jamfpro_script.script_001_data.id
}

output "jamfpro_script_001_name" {
  value = data.jamfpro_script.script_001_data.name
}

// data source list

data "jamfpro_script_list" "example" {}

output "script_ids" {
  value = data.jamfpro_script_list.example.ids
}

output "scripts" {
  value = data.jamfpro_script_list.example.scripts
}
