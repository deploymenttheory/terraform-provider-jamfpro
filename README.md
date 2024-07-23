# Terraform Provider for Jamf Pro

> [!WARNING]
> This code is in preview and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice until it reaches general availability, and no guarantees or support is provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

## Introduction

This repository hosts the Jamf Pro Community Provider, built to integrate Jamf Pro's robust configuration management capabilities with Terraform's Infrastructure as Code (IaC) approach. Utilizing a comprehensive JAMF Pro SDK [go-api-sdk-jamfpro](https://github.com/deploymenttheory/go-api-sdk-jamfpro), which serves as a cohesive abstraction layer over both Jamf Pro and Jamf Pro Classic APIs, this provider ensures seamless API interactions and brings a wide array of resources under Terraform's management umbrella. The jamfpro provider is engineered to enrich your CI/CD workflows with Jamf Pro's extensive device management functionalities, encompassing device enrollment, inventory tracking, security compliance, and streamlined software deployment. Its primary goal is to enhance the efficiency of managing, deploying, and maintaining Apple devices across your infrastructure, fostering a synchronized and effective IT ecosystem.

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
- [Jamf Pro](https://www.jamf.com/) >= 11.5.1

## Community & Support

For further community support and to engage with other users of the Jamf Pro Terraform Provider, please join us on the Mac Admins Slack channel. You can ask questions, provide feedback, and share best practices with the community. Join us at:

- [Mac Admins Slack Channel](https://macadmins.slack.com/archives/C06R172PUV6) - #terraform-provider-jamfpro

## Getting Started with Examples

# Provider Configuration for Jamf Pro in Terraform

This documentation provides a detailed explanation of the configuration options available in the `provider.tf` file for setting up the Jamf Pro provider in Terraform.

### Concurrency
- You can adjust paralellism by setting the Terraform parallelism count using `terraform apply -parallelism=X` (the default is 10). [HashiCorp Docs](https://developer.hashicorp.com/terraform/cli/commands/apply#parallelism-n)
- The provider remains stable using paralellism of up to 50, going beyond is at your own risk!

### Cookie Jar
- The cookie jar has been removed as it is redundant with Terraform's parallelism.
- Please use the cookie lock (which enforces a single cookie across all parallel instances of Terraform) or set a custom cookie (also remains consistent across all instances of Terraform).


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

- **Demo Repository**: [Terraform Demo Jamf Pro](https://github.com/deploymenttheory/terraform-demo-jamfpro)

Feel free to explore this repository to better understand the implementation and to get your infrastructure (configuration)-as-code initiatives up and running smoothly.

## Resource Completion Status

The follow is a summary of the resources and their completion status.

Resources can have the following statuses:

- **Beta** - The resource is in the early stages of development and may not be fully functional. It is not recommended for use in production environments as it may contain bugs and undergo significant changes.

- **Community Preview** - The resource is available for public use and feedback. While it has reached a level of stability beyond Beta, it may still undergo changes based on community input and additional testing. Users are encouraged to try it out and provide feedback, but should be cautious when using it in production environments.

- **Finished** - The resource is fully functional and has been tested in a production environment. It is considered stable and reliable for use in live systems. Users can confidently integrate this resource into their production workflows.

## Supported Jamf Pro Resources

This section outlines the resources and data sources provided by our Terraform provider for managing various aspects of Jamf Pro. Each resource comes with comprehensive support for the respective Jamf Pro entities, facilitating their management through Terraform.

### Accounts

- **Resource & Data Source**: Enables the management of Account within Jamf Pro, allowing for the configuration of accounts, access levels, privileges, assignment to groups and sites and other details.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.44.`

### Account Groups

- **Resource & Data Source**: Enables the management of Account Groups within Jamf Pro, allowing for the configuration of group names, access levels, privileges, and member details.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.31.`

### Activation Code

- **Resource & Data Source**: Enables the management of the Activation Code within Jamf Pro, allowing for the configuration of activation code and organization details.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.57.`

### API Roles

- **Resource & Data Source**: Enables the management of API roles within Jamf Pro, allowing for the configuration of role names, privileges, and other details. these can be assigned to api integrations.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.44.`

### API Integrations

- **Resource & Data Source**: Enables the management of API integrations within Jamf Pro, allowing for the configuration of integration names, privileges, and other details.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.44.`

### Buildings

- **Resource & Data Source**: Provides the ability to manage Buildings within Jamf Pro. This resource allows for the specification of building names and addresses, facilitating better organization and segmentation of devices within different physical locations.

- **Status**: Finished
- **Availability**: Introduced in version  `v0.0.30.`

### Categories

- **Resource & Data Source**: Enables the management of Categories within Jamf Pro, allowing for the configuration of category names, used across various Jamf Pro entities to categorize and organize devices, policies, and other resources.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.36.`

### Departments

- **Resource & Data Source**: Provides the ability to manage departments within Jamf Pro. This resource allows for the specification of department names.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.36.`

### Dock Items

- **Resource & Data Source**: Facilitates the management of Dock Items in Jamf Pro. This includes the creation, update, and deletion of dock item entities, along with the ability to specify dock item properties and associated payloads.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.43.`

### macOS Configuration Profiles (Plist)

- **Resource & Data Source**: Facilitates the management of macOS configuration profiles in Jamf Pro. This includes the creation, update, and deletion of configuration profiles, along with the ability to specify profile payloads and associated properties.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.37.`

### Mobile Device Configuration Profiles (Plist)

- **Resource & Data Source**: Facilitates the management of mobile device configuration profiles in Jamf Pro. This includes the creation, update, and deletion of configuration profiles, along with the ability to specify profile payloads and associated properties.

- **Status**: Community Preview
- **Availability**: Introduced in version `v0.0.48.`

### Packages

- **Resource & Data Source**: Facilitates the management of Packages in Jamf Pro. This includes the creation, update, and deletion of package entities, along with the ability to specify package payloads and associated properties. It uploads the package to the JCDS 2.0 CDN in AWS S3 and then creates the
package metadata in Jamf Pro.

- **Status**: Community Preview
- **Availability**: Introduced in version  `v0.0.34.`

### Scripts

- **Resource & Data Source**: Facilitates the management of Scripts in Jamf Pro. This includes the creation, update, and deletion of script entities, along with the ability to specify script contents and associated properties.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.30.`

### Sites

- **Resource & Data Source**: Provides the ability to manage Sites within Jamf Pro. This resource allows for the specification of site names and details, facilitating the organization of devices and resources across different sites.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.42.`

### Restricted Software

- **Resource & Data Source**: Facilitates the management of Restricted Software in Jamf Pro. This includes the creation, update, and deletion of restricted software entities, along with the ability to specify software properties and associated payloads.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.53.`

### User Groups

- **Resource & Data Source**: Enables the handling of User Groups in Jamf Pro. This encompasses the capabilities to create, update, and remove user group entities, as well as the functionality to detail user group attributes and memberships.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.38`.
