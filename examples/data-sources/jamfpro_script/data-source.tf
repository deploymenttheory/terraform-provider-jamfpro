# Example using id
data "jamfpro_script" "example_by_id" {
  id = "123"
}

# Example using name
data "jamfpro_script" "example_by_name" {
  name = "Sequoia_cis_lvl1_compliance.sh"
}

# Individual field outputs
output "script_id" {
  value = data.jamfpro_script.example_by_name.id
}

output "script_name" {
  value = data.jamfpro_script.example_by_name.name
}

output "script_category" {
  value = data.jamfpro_script.example_by_name.category_id
}

# Comprehensive output showing all available fields
output "script_all_fields" {
  value = {
    id              = data.jamfpro_script.example_by_name.id
    name            = data.jamfpro_script.example_by_name.name
    category_id     = data.jamfpro_script.example_by_name.category_id
    info            = data.jamfpro_script.example_by_name.info
    notes           = data.jamfpro_script.example_by_name.notes
    os_requirements = data.jamfpro_script.example_by_name.os_requirements
    priority        = data.jamfpro_script.example_by_name.priority
    script_contents = data.jamfpro_script.example_by_name.script_contents
    parameter4      = data.jamfpro_script.example_by_name.parameter4
    parameter5      = data.jamfpro_script.example_by_name.parameter5
    parameter6      = data.jamfpro_script.example_by_name.parameter6
    parameter7      = data.jamfpro_script.example_by_name.parameter7
    parameter8      = data.jamfpro_script.example_by_name.parameter8
    parameter9      = data.jamfpro_script.example_by_name.parameter9
    parameter10     = data.jamfpro_script.example_by_name.parameter10
    parameter11     = data.jamfpro_script.example_by_name.parameter11
  }
}