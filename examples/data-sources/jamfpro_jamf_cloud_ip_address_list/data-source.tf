# Get all Jamf Cloud IP addresses (no filters)
data "jamfpro_jamf_cloud_ip_address_list" "all" {}

# Get all AWS IPs for the EU Central region (all services)
data "jamfpro_jamf_cloud_ip_address_list" "aws_eu_central" {
  provider_filter = "aws"
  region_filter   = "eu-central-1"
}

# Get all Azure IPs for the Central US region (all services)
data "jamfpro_jamf_cloud_ip_address_list" "azure_central_us" {
  provider_filter = "azure"
  region_filter   = "centralus"
}

# Get Jamf Cloud Distribution Service FQDNs (inbound traffic)
data "jamfpro_jamf_cloud_ip_address_list" "jcds" {
  service_filter = "jamf_cloud_distribution_service"
  traffic_filter = "inbound"
}

# Output examples
output "publish_date" {
  description = "The publish date of the Jamf Cloud IP address list"
  value       = data.jamfpro_jamf_cloud_ip_address_list.all.publish_date
}

output "all_entries_count" {
  description = "Total number of IP entries"
  value       = length(data.jamfpro_jamf_cloud_ip_address_list.all.public_ips)
}

output "aws_eu_central_ips_by_service" {
  description = "AWS EU Central IP prefixes grouped by service"
  value = {
    for service in distinct([for entry in data.jamfpro_jamf_cloud_ip_address_list.aws_eu_central.public_ips : entry.service]) :
    service => flatten([
      for entry in data.jamfpro_jamf_cloud_ip_address_list.aws_eu_central.public_ips :
      entry.ip_prefixes if entry.service == service
    ])
  }
}

output "azure_central_us_ips_by_service" {
  description = "Azure Central US IP prefixes grouped by service"
  value = {
    for service in distinct([for entry in data.jamfpro_jamf_cloud_ip_address_list.azure_central_us.public_ips : entry.service]) :
    service => flatten([
      for entry in data.jamfpro_jamf_cloud_ip_address_list.azure_central_us.public_ips :
      entry.ip_prefixes if entry.service == service
    ])
  }
}

output "jcds_fqdns_by_region" {
  description = "Jamf Cloud Distribution Service FQDNs grouped by region"
  value = {
    for region in distinct([for entry in data.jamfpro_jamf_cloud_ip_address_list.jcds.public_ips : entry.region]) :
    region => flatten([
      for entry in data.jamfpro_jamf_cloud_ip_address_list.jcds.public_ips :
      entry.fqdns if entry.region == region
    ])
  }
}
