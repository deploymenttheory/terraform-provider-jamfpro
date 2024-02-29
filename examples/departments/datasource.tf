data "jamfpro_department" "example_department" {
  name = "tf-example-department-01"  # Replace this with the actual name of the site you want to retrieve
}

output "department_id" {
  value = data.jamfpro_department.example_department.id
}

output "department_name" {
  value = data.jamfpro_department.example_department.name
}
