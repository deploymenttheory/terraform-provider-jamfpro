# Terraform Provider for Jamf Pro

This repository hosts the Jamf Pro Community Provider, built to integrate Jamf Pro's robust configuration management capabilities with Terraform's Infrastructure as Code (IaC) approach. Utilizing a comprehensive JAMF Pro SDK [go-api-sdk-jamfpro](https://github.com/deploymenttheory/go-api-sdk-jamfpro), which serves as a cohesive abstraction layer over both Jamf Pro and Jamf Pro Classic APIs, this provider ensures seamless API interactions and brings a wide array of resources under Terraform's management umbrella. The jamfpro provider is engineered to enrich your CI/CD workflows with Jamf Pro's extensive device management functionalities, encompassing device enrollment, inventory tracking, security compliance, and streamlined software deployment. Its primary goal is to enhance the efficiency of managing, deploying, and maintaining Apple devices across your infrastructure, fostering a synchronized and effective IT ecosystem.

The provider contains:

- Resources and data sources for Jamf Pro entities (`internal/provider/`),
- Examples [examples](https://github.com/deploymenttheory/terraform-provider-jamfpro/tree/main/examples) directory for sample configurations and usage scenarios of the `terraform-provider-jamfpro` provider.
- Documentation [docs](https://github.com/deploymenttheory/terraform-provider-jamfpro/tree/main/docs)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21
- [Jamf Pro](https://www.jamf.com/) >= 11.2.0

## Resource Completion Status

The follow is a summary of the resources and their completion status.

Resources can have the following statuses:

- **Experimental** - The resource is in the early stages of development and may not be fully functional. It is not recommended for production use.

- **Finished** - The resource is fully functional and has been tested in a production environment.

## Supported Jamf Pro Resources

This section outlines the resources and data sources provided by our Terraform provider for managing various aspects of Jamf Pro. Each resource comes with comprehensive support for the respective Jamf Pro entities, facilitating their management through Terraform.

### Accounts

- **Resource & Data Source**: Enables the management of Account within Jamf Pro, allowing for the configuration of accounts, access levels, privileges, assignment to groups and sites and other details.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.44.`

### Account Groups

- **Resource & Data Source**: Enables the management of Account Groups within Jamf Pro, allowing for the configuration of group names, access levels, privileges, and member details.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.31.`

### API Roles

- **Resource & Data Source**: Enables the management of API roles within Jamf Pro, allowing for the configuration of role names, privileges, and other details. these can be assigned to api integrations.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.44.`

### API Integrations

- **Resource & Data Source**: Enables the management of API integrations within Jamf Pro, allowing for the configuration of integration names, privileges, and other details.

- **Status**: Finished
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

### macOS Configuration Profiles

- **Resource & Data Source**: Facilitates the management of macOS configuration profiles in Jamf Pro. This includes the creation, update, and deletion of configuration profiles, along with the ability to specify profile payloads and associated properties.

- **Status**: Experimental
- **Availability**: Introduced in version `v0.0.37.`

### Mobile Device Configuration Profiles

- **Resource & Data Source**: Facilitates the management of mobile device configuration profiles in Jamf Pro. This includes the creation, update, and deletion of configuration profiles, along with the ability to specify profile payloads and associated properties.

- **Status**: Experimental
- **Availability**: Introduced in version `v0.0.48.`

### Packages

- **Resource & Data Source**: Facilitates the management of Packages in Jamf Pro. This includes the creation, update, and deletion of package entities, along with the ability to specify package payloads and associated properties. It uploads the package to the JCDS 2.0 CDN in AWS S3 and then creates the
package metadata in Jamf Pro.

- **Status**: Experimental
- **Availability**: Introduced in version  `v0.0.34.`

### Scripts

- **Resource & Data Source**: Facilitates the management of Scripts in Jamf Pro. This includes the creation, update, and deletion of script entities, along with the ability to specify script contents and associated properties.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.30.`

### Sites

- **Resource & Data Source**: Provides the ability to manage Sites within Jamf Pro. This resource allows for the specification of site names and details, facilitating the organization of devices and resources across different sites.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.42.`

### User Groups

- **Resource & Data Source**: Enables the handling of User Groups in Jamf Pro. This encompasses the capabilities to create, update, and remove user group entities, as well as the functionality to detail user group attributes and memberships.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.38`.

## Terraform Parallelism and JAMF Pro Resource Creation in Load Balanced Environments

Jamf Pro is frequently hosted in clustered configurations with loadbalancing with two or more web applications. Jamf Pro handles resource propagation between the web applications and in Jamfcloud configurations the propagation time is exactly 60s to align all web applications (2) in the cluster.

When creating resources in Jamf Pro, it is important to consider the propagation time of resources across the Jamf Pro web applications. When managing JAMF with Terraform it's possible for terraform to create a resource successfully, but when it comes to stating the resource it has a 50 / 50 chance that it might reach a web app that hasn't been propagated to yet. This can lead to a resource being created but not stated correctly and leads to orphaned resource scenario's.

On the terraform side, Terraform by default creates a http client for a terraform plan operation and a separate http client for terraform apply. The terraform apply http client once initialised is used for all operations for a given run. This is useful as it means we can implement support for sticky sessions within the http client to ensure that all operations are targeted to the same web app.

Mitigation Strategy:

[1] Utilise sticky sessions within the http client used by this terraform provider.

[2] Enforce a 60 second propagation delay for TF resource creation operations when sticky sessions are disabled.

[3] Ensure that terraform is run with a parallelism of 1 to ensure that resources are created and stated in a controlled manner. (suggested)

Sticky sessions can be enabled like this in the provider configuration:

```bash
provider "jamfpro" {
  enable_cookie_jar = true // or false
}
```

Behaviour Description [False] When disabled, the http client doesn't use sticky sessions and will honor the 60s propagation time of jamf pro in jamf cloud contexts to ensure successful TF resource stating. This results in a given resource creation task taking circa 1 minute to deploy across the board. This approach keeps the load on jamf pro light and when deploying during business hours, this may be the preferred  configuration to ensure that jamf pro api resources are available for various device management activities outside of terraform. The down side of this however is that will take longer for a terraform apply to complete which is pertinent during pipeline runs.

Behaviour Description [True] When enabled, the http client uses sticky sessions and will target all operations to a single jamf pro web app. This negates Jamf Pro's load balancing and results in increased load on the targeted web app. However it provides the benefit that resources can be deployed and stated faster. This is due to the assurance that the web app api that was targeted for resource creation, will always be the same as the web app api used for TF resource stating. The propagation time in this scenario is set to 5 rather than 60 seconds.

### Special note: Terraform parallelism

By default terraform runs 10 operations in parallel. During load testing I have observed that when terraform performs Create operations above 1 against jamf pro it frequently results in unreliable resource deployment behavior. E.g resources deployed with partial configuration leading to stating failure. This is due to the fact that the jamf pro API get's overwhelmed due to the concurrency of the Create requests. Consequently I advise when possible to run terraform with the following

`terraform apply -parallelism=1`

Which restricts terraform to a single operation at a time. From load testing with 500 resource creations, across 10 different resource types with the cookie jar is enabled I was able to deploy successfully and state all resources. Effectively a new resource was created and stated every 5 seconds.

If you are unable to control the parallelism of terraform due to your pipeline design then proceed cautiously when creating jamf resources in batches.
