terraform {
  required_providers {
    jamfpro = {
      source = "deploymenttheory/jamfpro"
    }
  }
}

provider "jamfpro" {
  jamfpro_instance_fqdn = var.proxy_base_url # Include the proxy's path prefix, if needed.
  auth_method           = "oauth2"
  client_id             = var.jamfpro_client_id
  client_secret         = var.jamfpro_client_secret
  bearer_token_header   = "X-Jamf-Authorization"

  custom_headers {
    name  = "Authorization"
    value = "Basic ${base64encode("${var.proxy_username}:${var.proxy_password}")}"
  }
}

variable "proxy_base_url" {
  type        = string
  description = "Reverse-proxy base URL, for example https://proxy.company.com/proxypath."
}

variable "jamfpro_client_id" {
  type        = string
  description = "Jamf Pro OAuth client ID."
}

variable "jamfpro_client_secret" {
  type        = string
  description = "Jamf Pro OAuth client secret."
  sensitive   = true
}

variable "proxy_username" {
  type        = string
  description = "Service-account username for the reverse proxy."
  sensitive   = true
}

variable "proxy_password" {
  type        = string
  description = "Service-account password for the reverse proxy."
  sensitive   = true
}
