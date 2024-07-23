// Webhook example
resource "jamfpro_webhook" "jamfpro_webhook_001" {
  name                = "ExampleWebhook"
  enabled             = true
  url                 = "https://example.com/webhook"
  content_type        = "application/json"
  event               = "DeviceAddedToDEP"
  connection_timeout  = 5
  read_timeout        = 5
  authentication_type = "BASIC"
  username            = "exampleUser"
  password            = "examplePassword"

}

// Webhook example with smart group
resource "jamfpro_webhook" "jamfpro_webhook_001" {
  name                = "ExampleWebhook"
  enabled             = true
  url                 = "https://example.com/webhook"
  content_type        = "application/json"
  event               = "SmartGroupComputerMembershipChange"
  connection_timeout  = 5
  read_timeout        = 5
  authentication_type = "BASIC"
  username            = "exampleUser"
  password            = "examplePassword"
  smart_group_id      = jamfpro_smart_group.jamfpro_smart_group_001.id

}