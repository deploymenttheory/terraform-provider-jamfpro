# Terraform Provider for Jamf Pro

> [!WARNING]
> This code is in preview and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice until it reaches general availability, and no guarantees or support is provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

## Introduction

This repository hosts the Jamf Pro Community Provider, built to integrate Jamf Pro's robust configuration management capabilities with Terraform's Infrastructure as Code (IaC) approach. Utilizing a comprehensive JAMF Pro SDK [go-api-sdk-jamfpro](https://github.com/deploymenttheory/go-api-sdk-jamfpro), which serves as a cohesive abstraction layer over both Jamf Pro and Jamf Pro Classic APIs, this provider ensures seamless API interactions and brings a wide array of resources under Terraform's management umbrella. The jamfpro provider is engineered to enrich your CI/CD workflows with Jamf Pro's extensive device management functionalities, encompassing device enrollment, inventory tracking, security compliance, and streamlined software deployment. Its primary goal is to enhance the efficiency of managing, deploying, and maintaining Apple devices across your infrastructure, fostering a synchronized and effective IT ecosystem.

## Demo Implementation

To help you get started and understand the practical implementation of this provider, we've created a comprehensive demo repository:

- **Demo Repository**: [Terraform Demo Jamf Pro](https://github.com/deploymenttheory/terraform-demo-jamfpro)

This demo repository showcases a real-world implementation of the Jamf Pro Terraform provider. It's designed to:

1. Illustrate best practices for integrating Jamf Pro with Terraform.
2. Demonstrate a GitLab-flow based workflow for multi environment setups, integrating with Terraform Cloud.
3. Provide practical examples of managing Jamf Pro resources as code.
4. Offer a starting point for your own infrastructure-as-code initiatives with Jamf Pro.

We encourage you to explore this repository to:

- Understand how to structure your Terraform configurations for Jamf Pro.
- See examples of defining and managing various Jamf Pro resources.
- Learn how to integrate Terraform workflows with your CI/CD pipeline.
- Get insights into version control strategies for your Jamf Pro configurations.

Whether you're new to Terraform or looking to enhance your existing Jamf Pro management, this demo repository serves as a valuable resource to kickstart your infrastructure-as-code journey with Jamf Pro.

## Quick Start Guide

- Minimum Requirements:

```hcl
provider "jamfpro" {
  jamfpro_instance_fqdn = "https://yourserver.jamfcloud.com"
  auth_method     = "oauth2"
  client_id       = "your client id"
  client_secret   = "your client secret"
  jamfpro_load_balancer_lock = true
}
```

- Full Configuration:

```hcl

provider "jamfpro" {
  jamfpro_instance_fqdn = "https://yourserver.jamfcloud.com"
  auth_method     = "oauth2"
  client_id       = "your client id"
  client_secret   = "your client secret"
  enable_client_sdk_logs = false
  client_sdk_log_export_path = "/path/to/logfile.json"
  hide_sensitive_data = true
  custom_cookies {
    // Cookie URL is set to jamfpro_instance_fqdn
    name = "cookie name"
    value = "cookie value"
  }
  jamfpro_load_balancer_lock = true
  token_refresh_buffer_period_seconds = 300
  mandatory_request_delay_milliseconds = 100
  
}

```

The provider contains:

- Resources and data sources for Jamf Pro entities (`internal/provider/`),
- Examples [examples](https://github.com/deploymenttheory/terraform-provider-jamfpro/tree/main/examples) directory for sample configurations and usage scenarios of the `terraform-provider-jamfpro` provider.
- Documentation [docs](https://github.com/deploymenttheory/terraform-provider-jamfpro/tree/main/docs)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22.4
- [Jamf Pro](https://www.jamf.com/) >= 11.9.1

## Community & Support

For further community support and to engage with other users of the Jamf Pro Terraform Provider, please join us on the Mac Admins Slack channel. You can ask questions, provide feedback, and share best practices with the community. Join us at:

- [Mac Admins Slack Channel](https://macadmins.slack.com/archives/C06R172PUV6) - #terraform-provider-jamfpro

## Getting Started with Examples

# Provider Configuration for Jamf Pro in Terraform

This documentation provides a detailed explanation of the configuration options available in the `provider.tf` file for setting up the Jamf Pro provider in Terraform.


### Jamf Cloud Load Balancing and Cookies

- Jamf Cloud uses a load balancer to distribute traffic across multiple web app members. This abstraction can cause issues with Terraform's http client when multiple instances are running in parallel. Where freqently the client will be communicating with different web app members, causing stating unfound resource issues.
- To mitigate this please use the `jamfpro_load_balancer_lock` (which enforces a single cookie across all parallel instances of Terraform operations) which extracts the correct cookie on first run from Jamf Pro and utilises it thereafter across all other sessions. This is eqivalent to a sticky session.
- For non Jamf Cloud customers, with load balanced configurations please use `custom_cookies` and configure a custom cookie to be used in all requests instead.

### Concurrency

- You can adjust paralellism by setting the Terraform parallelism count using `terraform apply -parallelism=X` (the default is 10). [HashiCorp Docs](https://developer.hashicorp.com/terraform/cli/commands/apply#parallelism-n)
- The provider remains stable using paralellism of up to 50, going beyond is at your own risk!


## Configuration Schema

### `jamfpro_instance_fqdn`
- **Type:** String
- **Required:** Yes
- **Default:** Fetched from environment variable `envKeyJamfProUrlRoot` if not provided
- **Description:** The base URL for the Jamf Pro instance. Example: `https://mycompany.jamfcloud.com`. This URL is used to interact with the Jamf Pro API.

### `auth_method`
- **Type:** String
- **Required:** Yes
- **Description:** The authentication method to use for connecting to Jamf Pro.
- **Valid Values:** 
  - `basic`: Use basic authentication with a username and password.
  - `oauth2`: Use OAuth2 for authentication.
- **Validation:** Ensures the value is one of the specified valid values.

### `client_id`
- **Type:** String
- **Optional:** Yes
- **Default:** Fetched from environment variable `envKeyOAuthClientSecret` if not provided
- **Description:** The OAuth2 Client ID used for authentication with Jamf Pro. Required if `auth_method` is `oauth2`.

### `client_secret`
- **Type:** String
- **Optional:** Yes
- **Sensitive:** Yes
- **Default:** Fetched from environment variable `envKeyOAuthClientSecret` if not provided
- **Description:** The OAuth2 Client Secret used for authentication with Jamf Pro. This field is sensitive and required if `auth_method` is `oauth2`.

### `basic_auth_username`
- **Type:** String
- **Optional:** Yes
- **Default:** Fetched from environment variable `envKeyBasicAuthUsername` if not provided
- **Description:** The username for basic authentication with Jamf Pro. Required if `auth_method` is `basic`.

### `basic_auth_password`
- **Type:** String
- **Optional:** Yes
- **Sensitive:** Yes
- **Default:** Fetched from environment variable `envKeyBasicAuthPassword` if not provided
- **Description:** The password for basic authentication with Jamf Pro. This field is sensitive and required if `auth_method` is `basic`.


### `enable_client_sdk_logs`
- **Type:** bool
- **Optional:** Yes
- **Default:** false
- **Description:** Enables Client and SDK logs to appear in the tf output.

### `client_sdk_log_export_path`
- **Type:** String
- **Optional:** Yes
- **Default:** `""`
- **Description:** The file path to export HTTP client logs to. If set, logs will be saved to this path. If omitted, logs will not be exported.

### `hide_sensitive_data`
- **Type:** Boolean
- **Optional:** Yes
- **Default:** `true`
- **Description:** Determines whether sensitive information (like passwords) should be hidden in logs. Defaults to hiding sensitive data for security reasons.

### `custom_cookies`
- **Type:** List of Objects
- **Optional:** Yes
- **Default:** `nil`
- **Description:** A list of custom cookies to be included in HTTP requests. Each cookie object should have a `name` and a `value`.
  - **name**: 
    - **Type:** String
    - **Required:** Yes
    - **Description:** The name of the cookie.
  - **value**: 
    - **Type:** String
    - **Required:** Yes
    - **Description:** The value of the cookie.

### `jamf_load_balancer_lock`
- **Type:** Boolean
- **Optional:** Yes
- **Default:** `false`
- **Description:** Temporarily locks all HTTP client instances to a specific web app member in the load balancer for faster execution. This is a temporary solution until Jamf provides an official load balancing solution.

### `token_refresh_buffer_period_seconds`
- **Type:** Integer
- **Optional:** Yes
- **Default:** `300`
- **Description:** The buffer period in seconds before the token expires during which the token will be refreshed. Helps ensure continuous authentication.

### `mandatory_request_delay_milliseconds`
- **Type:** Integer
- **Optional:** Yes
- **Default:** `100`
- **Description:** A mandatory delay after each request before returning to reduce high volume of requests in a short time.


For those new to using Terraform with Jamf Pro, we provide a comprehensive demo example that serves as an excellent starting point. This demo implementation utilizes:

- Terraform Cloud as the remote backend
- GitHub Actions pipelines
- A simple PR process for managing changes
- Sample hcl files for creating and managing Jamf Pro resources

This repository is specifically designed to kickstart your Terraform projects by providing practical, easy-to-follow examples of how to configure and deploy resources within Jamf Pro using Terraform.

# Supported Jamf Pro Resources

[Supported Resources](docs/supported-resources/README.md)
