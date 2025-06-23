# ADR-0001: Use of Custom Go API SDK for Jamf Pro

## Status

Accepted

## Date

25-10-2023

## Context

Terraform providers require a Go-based SDK to interact with the target service's API. For our Jamf Pro Terraform provider, we need a reliable and comprehensive Go SDK to communicate with the Jamf Pro API. However, Jamf does not offer an official Go SDK for either their Classic API or their Pro API.

The Jamf Pro API is extensive, covering hundreds of endpoints across various resources such as computers, mobile devices, users, policies, scripts, and more. These endpoints follow different API patterns (Classic API, Pro API) and have varying authentication mechanisms.

Without an official SDK, we needed a solution that would:

1. Provide comprehensive coverage of the Jamf Pro API
2. Handle authentication (both Basic Auth and OAuth2)
3. Support proper error handling and pagination
4. Follow Go best practices
5. Be maintainable and extensible as Jamf Pro evolves

## Decision Drivers

* Need for a Go-based SDK to integrate with Terraform's Go-based plugin architecture
* Absence of an official Jamf Pro SDK for Go
* Requirement for complete API coverage to support all Jamf Pro resources
* Need for consistent error handling and authentication flows
* Long-term maintenance considerations

## Considered Options

* Develop a custom Go SDK for Jamf Pro
* Use generic HTTP clients directly in the Terraform provider
* Wait for an official SDK from Jamf

## Decision

Chosen option: "Develop a custom Go SDK for Jamf Pro", because it provides the most control, ensures complete coverage of required endpoints, and allows us to maintain and extend the SDK as needed.

We created the `go-api-sdk-jamfpro` as a separate project that serves as the foundation for the Terraform provider. This SDK implements all the necessary API calls, authentication mechanisms, and data structures required to interact with Jamf Pro.

## Rationale

Creating our own SDK gives us several advantages:

1. **Complete Control**: We can implement exactly the functionality we need for the Terraform provider
2. **Consistent Interface**: We can design a clean, idiomatic Go API that follows best practices
3. **Separation of Concerns**: The SDK handles API communication and data structures, while the provider focuses on Terraform-specific logic
4. **Reusability**: The SDK can be used by other Go projects beyond just our Terraform provider
5. **Extensibility**: We can easily add support for new Jamf Pro API endpoints as they are released

The SDK is designed with a clean separation between HTTP client functionality and the Jamf Pro API-specific code, making it easier to maintain and extend.

## Consequences

### Positive

* Complete control over the SDK implementation and features
* Ability to quickly add support for new Jamf Pro API endpoints
* Consistent error handling and authentication flows
* Reusable code that can be leveraged by other Go projects
* Better separation of concerns in the Terraform provider codebase

### Negative

* Maintenance burden of keeping the SDK updated with Jamf Pro API changes
* Need to implement and test all API endpoints ourselves
* Risk of divergence from future official SDKs if Jamf releases one

### Neutral

* Need for regular updates as Jamf Pro releases monthly updates with API changes
* Documentation requirements for both the SDK and how it's used in the provider

## Implementation

### Action Items

* [x] Create a separate repository for the Go API SDK
* [x] Implement core HTTP client with authentication support
* [x] Add support for all required Jamf Pro API endpoints
* [x] Create comprehensive tests for the SDK
* [x] Document SDK usage
* [x] Establish a process for keeping the SDK updated with Jamf Pro releases

### Timeline

* The SDK is already implemented and in use by the Terraform provider
* Ongoing maintenance will follow Jamf Pro's monthly release cycle

## Validation

The success of this decision will be measured by:

* Ability to implement all planned resources in the Terraform provider
* Ease of adding new resources as needed
* Stability of the SDK across Jamf Pro updates
* Minimal breaking changes in the SDK API

## References

* [go-api-sdk-jamfpro Repository](https://github.com/deploymenttheory/go-api-sdk-jamfpro)
* [terraform-provider-jamfpro Repository](https://github.com/deploymenttheory/terraform-provider-jamfpro)
* [Jamf Pro API Documentation](https://developer.jamf.com/jamf-pro/reference/welcome)

## Notes

Our development flow is:

1. Update the SDK when Jamf makes API changes in their monthly releases
2. Reflect these changes in the SDK repository
3. Cascade the changes into the Terraform provider

This approach ensures we maintain compatibility with the latest Jamf Pro API while providing a stable interface for the Terraform provider. 