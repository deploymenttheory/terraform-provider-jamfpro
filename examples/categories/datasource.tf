data "jamfpro_category" "example_category" {
  name = "tf-example-category-01"  # Replace this with the actual name of the site you want to retrieve
}

output "category_id" {
  value = data.jamfpro_category.example_category.id
}

output "category_name" {
  value = data.jamfpro_category.example_category.name
}
