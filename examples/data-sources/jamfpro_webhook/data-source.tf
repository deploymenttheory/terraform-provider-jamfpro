// data source by id
data "jamfpro_webhook" "webhook_001_data" {
  id = jamfpro_webhook.jamfpro_webhook_001.id
}

output "jamfpro_webhook_001_id" {
  value = data.jamfpro_webhook.webhook_001_data.id
}

output "jamfpro_webhook_001_name" {
  value = data.jamfpro_webhook.webhook_001_data.name
}

// data source list

data "jamfpro_webhook_list" "example" {}

output "webhook_ids" {
  value = data.jamfpro_webhook_list.example.ids
}

output "webhooks" {
  value = data.jamfpro_webhook_list.example.webhooks
}