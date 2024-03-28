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

### Account Groups

- **Resource & Data Source**: Enables the management of Account Groups within Jamf Pro, allowing for the configuration of group names, access levels, privileges, and member details.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.31.`

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

### macOS Configuration Profiles

- **Resource & Data Source**: Facilitates the management of macOS configuration profiles in Jamf Pro. This includes the creation, update, and deletion of configuration profiles, along with the ability to specify profile payloads and associated properties.

- **Status**: Experimental
- **Availability**: Introduced in version `v0.0.37.`

### Packages

- **Resource & Data Source**: Facilitates the management of Packages in Jamf Pro. This includes the creation, update, and deletion of package entities, along with the ability to specify package payloads and associated properties. It uploads the package to the JCDS 2.0 CDN in AWS S3 and then creates the
package metadata in Jamf Pro.

- **Status**: Experimental
- **Availability**: Introduced in version  `v0.0.34.`

### Scripts

- **Resource & Data Source**: Facilitates the management of Scripts in Jamf Pro. This includes the creation, update, and deletion of script entities, along with the ability to specify script contents and associated properties.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.30.`

### User Groups

- **Resource & Data Source**: Enables the handling of User Groups in Jamf Pro. This encompasses the capabilities to create, update, and remove user group entities, as well as the functionality to detail user group attributes and memberships.

- **Status**: Finished
- **Availability**: Introduced in version `v0.0.38`.
