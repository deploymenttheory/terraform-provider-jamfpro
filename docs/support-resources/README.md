# Resource Completion Status

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
