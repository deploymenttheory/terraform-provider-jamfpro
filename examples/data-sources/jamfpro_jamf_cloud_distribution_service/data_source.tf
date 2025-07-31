data "jamfpro_jamf_cloud_distribution_service" "current" {}

output "jcds_files" {
  value = data.jamfpro_jamf_cloud_distribution_service.current.files
}
