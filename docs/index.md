---
page_title: "jamfpro Terraform Provider"
subcategory: ""
description: |-
    The Jamf Pro provider is used to manage a Jamf Pro instance.
---

# Jamf Pro Provider

The Jamf Pro Terraform provider is a plugin for Terraform that allows for the
management of [Jamf Pro](https://www.jamf.com/products/jamf-pro/) configuration and resources.
-----
A typical provider configuration would look something like:
```terraform
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
  platform_environment_id              = var.jamfpro_platform_environment_id
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
  description = "The Jamf Platform gateway base URL. Required when auth_provider is 'platform'. Example: https://us.api.jamfcloud.com"
  type        = string
  default     = ""
}

variable "jamfpro_platform_tenant_id" {
  description = "The Jamf Platform tenant identifier (UUID). Use for tenant-scoped authentication; mutually exclusive with platform_environment_id."
  sensitive   = true
  type        = string
  default     = ""
}

variable "jamfpro_platform_environment_id" {
  description = "The Jamf Platform environment identifier (UUID). Preferred for cross-product platform access. Mutually exclusive with platform_tenant_id."
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
```

## Reverse-proxy authentication

For a reverse proxy that requires Basic credentials in `Authorization` and the Jamf bearer token in another header, use `auth_method = "oauth2"`, set `bearer_token_header` to the proxy's token-header name, and supply the proxy credentials through a `custom_headers` block. The provider obtains and refreshes the Jamf token automatically; do not supply a static Jamf token as a custom header.

```hcl
provider "jamfpro" {
  jamfpro_instance_fqdn = "https://proxy.company.com/proxypath/"
  auth_method          = "oauth2"
  client_id            = var.jamfpro_client_id
  client_secret        = var.jamfpro_client_secret
  bearer_token_header  = "X-Jamf-Authorization"

  custom_headers {
    name  = "Authorization"
    value = "Basic ${base64encode("${var.proxy_username}:${var.proxy_password}")}"
  }
}
```

The proxy base URL may include a path prefix. Custom headers apply to token requests and API calls, including token refresh. `bearer_token_header` defaults to `Authorization` and can also be set using `JAMFPRO_BEARER_TOKEN_HEADER`. Custom header values are sensitive. Duplicate header names and custom headers that collide with the managed bearer-token header are rejected.

Jamf username/password authentication needs `Authorization` for its own token request, so use OAuth2 when the proxy also requires that header. When custom headers or a custom token header are configured, redirects to another origin are blocked; the proxy must forward requests rather than redirect them to Jamf.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `auth_method` (String) The authentication mechanism. Options are 'basic' for username/password or 'oauth2' for client credentials.
- `auth_provider` (String) The authentication provider. 'direct' authenticates against the Jamf Pro instance directly. 'platform' authenticates via the Jamf Platform gateway. Defaults to 'direct'.
- `basic_auth_password` (String, Sensitive) The Jamf Pro password used for authentication when auth_method is 'basic'.
- `basic_auth_username` (String) The Jamf Pro username used for authentication when auth_method is 'basic'.
- `bearer_token_header` (String) HTTP header used for the automatically refreshed Jamf bearer token. Defaults to Authorization. Set a custom header name when a reverse proxy requires Basic credentials in Authorization. Can also be set with JAMFPRO_BEARER_TOKEN_HEADER.
- `client_id` (String) The client ID for OAuth2 authentication.
- `client_sdk_log_export_path` (String) Specify the path to export http client logs to.
- `client_secret` (String, Sensitive) The client secret for OAuth2 authentication.
- `custom_cookies` (Block List) Persistent custom cookies used by HTTP Client in all requests. (see [below for nested schema](#nestedblock--custom_cookies))
- `custom_headers` (Block List) Additional HTTP headers sent to the configured Jamf Pro or gateway origin for token requests and API calls. Use this for reverse-proxy authentication or routing headers. Header values are sensitive. Cross-origin redirects are blocked when custom headers are configured. (see [below for nested schema](#nestedblock--custom_headers))
- `enable_client_sdk_logs` (Boolean) Debug option to propagate logs from the SDK and HttpClient
- `hide_sensitive_data` (Boolean) Define whether sensitive fields should be hidden in logs. Default to hiding sensitive data in logs
- `jamfpro_instance_fqdn` (String) The Jamf Pro FQDN (fully qualified domain name). Required when auth_provider is 'direct'. Example: https://mycompany.jamfcloud.com
- `jamfpro_load_balancer_lock` (Boolean) Programatically determines all available web app members in the load balancer and locks all instances of httpclient to the app for faster executions. 
TEMP SOLUTION UNTIL JAMF PROVIDES SOLUTION
- `mandatory_request_delay_milliseconds` (Number) A mandatory delay after each request before returning to reduce high volume of requests in a short time
- `platform_base_url` (String) The Jamf Platform gateway base URL. Required when auth_provider is 'platform'. Example: https://us.api.jamfcloud.com
- `platform_tenant_id` (String, Sensitive) The Jamf Platform tenant identifier (UUID). Use for tenant-scoped authentication; mutually exclusive with `platform_environment_id`.
- `platform_environment_id` (String, Sensitive) The Jamf Platform environment identifier (UUID). Preferred for cross-product platform access; mutually exclusive with `platform_tenant_id`. Exactly one scope ID is required when `auth_provider` is `platform`.
- `token_refresh_buffer_period_seconds` (Number) The buffer period in seconds for token refresh.

<a id="nestedblock--custom_cookies"></a>
### Nested Schema for `custom_cookies`

Required:

- `name` (String) cookie key
- `value` (String) cookie value


<a id="nestedblock--custom_headers"></a>
### Nested Schema for `custom_headers`

Required:

- `name` (String) Name of the HTTP header, for example Authorization or X-Company-Id. Header names must be unique, ignoring case.
- `value` (String, Sensitive) Value of the HTTP header. May contain proxy credentials or API keys and is treated as sensitive.
