data "jamfpro_file_share_distribution_point" "dp_example" {
  id = resource.jamfpro_file_share_distribution_point.dp_example.id
}

output "jamfpro_file_share_distribution_point_id" {
  value = data.jamfpro_file_share_distribution_point.dp_example.id
}

output "jamfpro_file_share_distribution_point_name" {
  value = data.jamfpro_file_share_distribution_point.dp_example.name
}
