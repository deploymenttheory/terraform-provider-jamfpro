terraform {
  required_providers {
    jamfpro = {
      source  = "deploymenttheory/jamfpro"
      version = "~> 0.0.10"
    }
  }
}

provider "jamfpro" {
  instance_name = var.jamfpro_instance_name
  client_id     = var.jamfpro_client_id
  client_secret = var.jamfpro_client_secret
  log_level     = "debug" # or "debug", "info", "none" depending on the desired verbosity of the http client
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
