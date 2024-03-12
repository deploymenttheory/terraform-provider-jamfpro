data "jamfpro_department" "jamfpro_department_001_data" {
  id = jamfpro_department.jamfpro_department_001.id
}

output "jamfpro_department_001_data_id" {
  value = data.jamfpro_department.jamfpro_department_001_data.id
}

output "jamfpro_department_001_data_name" {
  value = data.jamfpro_department.jamfpro_department_001_data.name
}