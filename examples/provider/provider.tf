terraform {
  required_providers {
    jamfpro = {
      source  = "deploymenttheory/jamfpro"
      version = "~> 0.0.43"
    }
  }
}

provider "jamfpro" {
  instance_name               = var.jamfpro_instance_name
  client_id                   = var.jamfpro_client_id
  client_secret               = var.jamfpro_client_secret
  log_level                   = "debug" # or "debug", "info", "none" depending on the desired verbosity of the http client
  log_output_format           = "console" # or "JSON" for JSON format
  log_console_separator       = " " # Separator character for console log output
  hide_sensitive_data         = true # Hides sensitive data in logs
  max_retry_attempts          = 5
  enable_dynamic_rate_limiting = false
  max_concurrent_requests     = 5
  token_refresh_buffer_period = 5 # minutes
  total_retry_duration        = 30 # seconds
  custom_timeout              = 30 # seconds
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
