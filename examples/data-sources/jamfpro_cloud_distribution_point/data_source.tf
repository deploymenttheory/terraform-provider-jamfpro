data "jamfpro_cloud_distribution_point" "current" {}

output "cdn_type" {
  value = data.jamfpro_cloud_distribution_point.current.cdn_type
}

output "connection_status" {
  value = data.jamfpro_cloud_distribution_point.current.has_connection_succeeded
}

output "upload_capable" {
  value = data.jamfpro_cloud_distribution_point.current.direct_upload_capable
}
