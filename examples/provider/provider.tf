terraform {
  required_providers {
    jamfpro = {
      source  = "deploymenttheory/jamfpro"
      version = "~> 0.0.43"
    }
  }
}

provider "jamfpro" {
  jamf_instance_fqdn          = var.jamfpro_instance_name
  auth_method =               "oauth2" // basic
  client_id                   = var.jamfpro_client_id
  client_secret               = var.jamfpro_client_secret
  # basic_auth_username = ""
  # basic_auth_password = ""
  log_level                   = "debug" # or "debug", "info", "none" depending on the desired verbosity of the http client
  log_output_format           = "console" # or "JSON" for JSON format
  log_console_separator       = " " # Separator character for console log output
  log_export_path             = "/path/to/log"
  export_logs                 = false
  hide_sensitive_data         = true # Hides sensitive data in logs
  token_refresh_buffer_period_seconds = 5 # minutes
  jamf_load_balancer_lock     = true
  custom_cookies = {
    name = "jpro-ingress"
    value = "value"
  }
  mandatory_request_delay_milliseconds = 100
}

variable "jamfpro_instance_name" {
  description = "Jamf Pro Instance name."
  default     = ""
}

variable "jamfpro_client_id" {
  description = "Jamf Pro Client ID for authentication."
  default     = ""
}

variable "jamfpro_client_secret" {
  description = "Jamf Pro Client Secret for authentication."
  sensitive   = true
  default     = ""
}

variable "jamfpro_username" {
  description = "Jamf Pro username used for authentication."
  default     = ""
}

variable "jamfpro_password" {
  description = "Jamf Pro password used for authentication."
  sensitive   = true
  default     = ""
}
