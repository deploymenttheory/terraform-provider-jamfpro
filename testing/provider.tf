terraform {
  required_providers {
    jamfpro = {
      source = "terraform.local/local/jamfpro"
      # Specifically 0.1.0
      version = "0.1.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.4.3"
    }
  }
}
provider "jamfpro" {
  jamfpro_instance_fqdn                = var.jamfpro_instance_fqdn
  auth_method                          = var.jamfpro_auth_method
  client_id                            = var.jamfpro_client_id
  client_secret                        = var.jamfpro_client_secret
  enable_client_sdk_logs               = var.enable_client_sdk_logs
  client_sdk_log_export_path           = var.client_sdk_log_export_path
  hide_sensitive_data                  = var.jamfpro_hide_sensitive_data
  jamfpro_load_balancer_lock           = var.jamfpro_load_balancer_lock
  token_refresh_buffer_period_seconds  = var.jamfpro_token_refresh_buffer_period_seconds
  mandatory_request_delay_milliseconds = var.jamfpro_mandatory_request_delay_milliseconds
}


variable "jamfpro_instance_fqdn" {
  description = "The Jamf Pro FQDN (fully qualified domain name). Example: https://mycompany.jamfcloud.com"
  sensitive   = true
}

variable "jamfpro_auth_method" {
  description = "Auth method chosen for Jamf. Options are 'basic' or 'oauth2'."
  sensitive   = true
  type        = string
}

variable "jamfpro_client_id" {
  description = "The Jamf Pro Client ID for authentication."
  sensitive   = true
  type        = string
}

variable "jamfpro_client_secret" {
  description = "The Jamf Pro Client Secret for authentication."
  sensitive   = true
  type        = string
}

variable "enable_client_sdk_logs" {
  description = "Enable client SDK logs."
  type        = bool
  default     = false
}

variable "client_sdk_log_export_path" {
  description = "Specify the path to export http client logs to."
  type        = string
  default     = ""
}

variable "jamfpro_hide_sensitive_data" {
  description = "Define whether sensitive fields should be hidden in logs."
  type        = bool
  default     = true
}

variable "jamfpro_custom_cookies" {
  description = "Custom cookies for the HTTP client."
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "jamfpro_load_balancer_lock" {
  description = "Programmatically determines all available web app members in the load balancer and locks all instances of httpclient to the app for faster executions."
  type        = bool
  default     = true
}

variable "jamfpro_token_refresh_buffer_period_seconds" {
  description = "The buffer period in seconds for token refresh."
  type        = number
  default     = 300
}

variable "jamfpro_mandatory_request_delay_milliseconds" {
  description = "A mandatory delay after each request before returning to reduce high volume of requests in a short time."
  type        = number
  default     = 100
}

variable "testing_id" {
  description = "Unique runtime id to differentiate testing objects between runs to avoid conflicts during cleanup phase"
  type        = string
}


resource "random_id" "rng" {
  keepers = {
    first = "${timestamp()}"
  }
  byte_length = 8
}