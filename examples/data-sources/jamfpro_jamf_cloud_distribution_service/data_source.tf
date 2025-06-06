data "jamfpro_jamf_cloud_distribution_service" "current" {}

output "jcds_file_stream_endpoint_enabled" {
  value = data.jamfpro_jamf_cloud_distribution_service.current.file_stream_endpoint_enabled
}

output "jcds2_enabled" {
  value = data.jamfpro_jamf_cloud_distribution_service.current.jcds2_enabled
}

output "jcds_files" {
  value = data.jamfpro_jamf_cloud_distribution_service.current.files
}
