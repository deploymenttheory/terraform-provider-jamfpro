data "jamfpro_webhook" "webhook_001_data" {
  id = jamfpro_webhook.jamfpro_webhook_001.id
}

output "jamfpro_webhook_001_id" {
  value = data.jamfpro_webhook.webhook_001_data.id
}

output "jamfpro_webhook_001_name" {
  value = data.jamfpro_webhook.webhook_001_data.name
}

data "jamfpro_webhook" "webhook_002_data" {
  id = jamfpro_webhook.jamfpro_webhook_002.id
}

output "jamfpro_webhook_002_id" {
  value = data.jamfpro_webhook.webhook_002_data.id
}

output "jamfpro_webhook_002_name" {
  value = data.jamfpro_webhook.webhook_002_data.name
}
