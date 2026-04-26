terraform {
  required_providers {
    jamfpro = {
      source  = "deploymenttheory/jamfpro"
      version = "~> 0.18.0"
    }
  }
}

provider "jamfpro" {
  jamfpro_instance_fqdn                = var.jamfpro_instance_fqdn
  auth_method                          = var.jamfpro_auth_method
  auth_provider                        = var.jamfpro_auth_provider
  platform_base_url                    = var.jamfpro_platform_base_url
  platform_tenant_id                   = var.jamfpro_platform_tenant_id
  client_id                            = var.jamfpro_client_id
  client_secret                        = var.jamfpro_client_secret
  basic_auth_username                  = var.jamfpro_basic_auth_username
  basic_auth_password                  = var.jamfpro_basic_auth_password
  enable_client_sdk_logs               = var.enable_client_sdk_logs
  client_sdk_log_export_path           = var.client_sdk_log_export_path
  hide_sensitive_data                  = var.jamfpro_hide_sensitive_data
  jamfpro_load_balancer_lock           = var.jamfpro_jamf_load_balancer_lock
  token_refresh_buffer_period_seconds  = var.jamfpro_token_refresh_buffer_period_seconds
  mandatory_request_delay_milliseconds = var.jamfpro_mandatory_request_delay_milliseconds
}

variable "jamfpro_instance_fqdn" {
  description = "The Jamf Pro FQDN (fully qualified domain name). Required when auth_provider is 'direct'. Example: https://mycompany.jamfcloud.com"
  sensitive   = true
  default     = ""
}

variable "jamfpro_auth_method" {
  description = "The authentication mechanism. Options are 'basic' or 'oauth2'."
  sensitive   = true
  default     = ""
}

variable "jamfpro_auth_provider" {
  description = "The authentication provider. 'direct' authenticates against Jamf Pro directly, 'platform' authenticates via the Jamf Platform gateway. Defaults to 'direct'."
  type        = string
  default     = "direct"
}

variable "jamfpro_platform_base_url" {
  description = "The Jamf Platform gateway base URL. Required when auth_provider is 'platform'. Example: https://us.api.platform.jamf.com"
  type        = string
  default     = ""
}

variable "jamfpro_platform_tenant_id" {
  description = "The Jamf Platform gateway tenant identifier (UUID). Required when auth_provider is 'platform'."
  sensitive   = true
  type        = string
  default     = ""
}

variable "jamfpro_client_id" {
  description = "The client ID for OAuth2 authentication."
  sensitive   = true
  default     = ""
}

variable "jamfpro_client_secret" {
  description = "The client secret for OAuth2 authentication."
  sensitive   = true
  default     = ""
}

variable "jamfpro_basic_auth_username" {
  description = "The Jamf Pro username used for authentication when auth_method is 'basic'."
  default     = ""
}

variable "jamfpro_basic_auth_password" {
  description = "The Jamf Pro password used for authentication when auth_method is 'basic'."
  sensitive   = true
  default     = ""
}

variable "enable_client_sdk_logs" {
  description = "Debug option to propagate logs from the SDK and HttpClient"
  default     = false
}

variable "client_sdk_log_export_path" {
  description = "Specify the path to export http client logs to."
  default     = ""
}

variable "jamfpro_hide_sensitive_data" {
  description = "Define whether sensitive fields should be hidden in logs. Default to hiding sensitive data in logs"
  default     = true
}

variable "jamfpro_custom_cookies" {
  description = "Persistent custom cookies used by HTTP Client in all requests."
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "jamfpro_jamf_load_balancer_lock" {
  description = "Programatically determines all available web app members in the load balancer and locks all instances of httpclient to the app for faster executions."
  default     = true
}

variable "jamfpro_token_refresh_buffer_period_seconds" {
  description = "The buffer period in seconds for token refresh."
  default     = 300
}

variable "jamfpro_mandatory_request_delay_milliseconds" {
  description = "A mandatory delay after each request before returning to reduce high volume of requests in a short time."
  default     = 1000
}
