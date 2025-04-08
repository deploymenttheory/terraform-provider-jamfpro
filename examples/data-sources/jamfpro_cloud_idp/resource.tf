# Example 1: Look up cloud identity provider by ID
data "jamfpro_cloud_idp" "by_id" {
  id = "1"
}

# Example 2: Look up cloud identity provider by display name
data "jamfpro_cloud_idp" "by_name" {
  display_name = "Azure AD"
}

# Example 3: Using variables
variable "idp_display_name" {
  type        = string
  description = "Display name of the cloud identity provider to look up"
  default     = "Google Workspace"
}

data "jamfpro_cloud_idp" "dynamic" {
  display_name = var.idp_display_name
}

# Example 4: Output examples
output "idp_details" {
  value = {
    id           = data.jamfpro_cloud_idp.by_name.id
    display_name = data.jamfpro_cloud_idp.by_name.display_name
    enabled      = data.jamfpro_cloud_idp.by_name.enabled
    provider_name = data.jamfpro_cloud_idp.by_name.provider_name
  }
}

# Example 5: Using with conditions
data "jamfpro_cloud_idp" "active_provider" {
  display_name = "Azure AD"

  lifecycle {
    postcondition {
      condition     = self.enabled == true
      error_message = "The cloud identity provider must be enabled"
    }
  }
}

# Example 6: Using in another resource (hypothetical)
resource "jamfpro_sso_configuration" "example" {
  name        = "SSO Configuration"
  enabled     = true
  idp_id      = data.jamfpro_cloud_idp.by_name.id
  provider    = data.jamfpro_cloud_idp.by_name.provider_name
}

# Example 7: Testing multiple providers
data "jamfpro_cloud_idp" "azure" {
  display_name = "Azure AD"
}

data "jamfpro_cloud_idp" "google" {
  display_name = "Google Workspace"
}

output "idp_comparison" {
  value = {
    azure = {
      id      = data.jamfpro_cloud_idp.azure.id
      enabled = data.jamfpro_cloud_idp.azure.enabled
    }
    google = {
      id      = data.jamfpro_cloud_idp.google.id
      enabled = data.jamfpro_cloud_idp.google.enabled
    }
  }
}