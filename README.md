# Community Terraform Provider for Jamf Pro

[![Release](https://img.shields.io/github/v/release/deploymenttheory/terraform-provider-jamfpro)](https://github.com/deploymenttheory/terraform-provider-jamfpro/releases)
[![Installs](https://img.shields.io/badge/dynamic/json?logo=terraform&label=installs&query=$.data.attributes.downloads&url=https%3A%2F%2Fregistry.terraform.io%2Fv2%2Fproviders%2F4960)](https://registry.terraform.io/providers/deploymenttheory/jamfpro)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/deploymenttheory/jamfpro/latest/docs)
[![Lint Status](https://github.com/deploymenttheory/terraform-provider-jamfpro/workflows/Linter/badge.svg)](https://github.com/deploymenttheory/terraform-provider-jamfpro/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/terraform-provider-jamfpro)](https://goreportcard.com/report/github.com/deploymenttheory/terraform-provider-jamfpro)
[![Go Version](https://img.shields.io/github/go-mod/go-version/deploymenttheory/terraform-provider-jamfpro)](https://go.dev/)
[![License](https://img.shields.io/github/license/deploymenttheory/terraform-provider-jamfpro)](LICENSE)
![Status: Public Preview](https://img.shields.io/badge/status-public%20preview-0078D4)

> [!WARNING]
> This provider is in public preview. While it has been tested extensively, please thoroughly test in non-production environments before production use. Features may contain bugs or undergo changes based on community feedback. Use at your own risk until general availability is reached. No guarantees or official support is provided. By using this provider, you acknowledge and agree to these conditions. For questions or issues, please consult the documentation or contact the maintainer.


> [!TIP]
> This is a community-driven project and is not officially supported by Jamf.
> If you need help, want to ask questions, or connect with other users and contributors, join our community
> [Mac Admins Slack Channel](https://macadmins.slack.com/archives/C06R172PUV6) - #terraform-provider-jamfpro

## Introduction

This repository hosts the Community Jamf Pro terraform Provider, built to integrate Jamf Pro's robust configuration management capabilities with Terraform's Infrastructure as Code (IaC) approach to service life cycle management. Utilizing a comprehensive JAMF Pro SDK [go-api-sdk-jamfpro](https://github.com/deploymenttheory/go-api-sdk-jamfpro), which serves as a cohesive abstraction layer over both Jamf Pro and Jamf Pro Classic APIs, this provider ensures seamless API interactions and brings a wide array of resources under Terraform's management umbrella.

The jamfpro provider is engineered to enrich your CI/CD workflows with Jamf Pro's extensive device management functionalities, encompassing device enrollment, inventory tracking, security compliance, and streamlined software deployment.

Its primary goal is to enhance the efficiency of managing, deploying, and maintaining Apple devices across your infrastructure, fostering an 'everything-as-code' mindset.

## Use Cases

- **Infrastructure as Code for Jamf Pro**  
  Manage Jamf Pro configuration (apps, groups, policies, device management, and more) as code, enabling version control, peer review, and repeatable deployments—just as you would for cloud infrastructure in Azure or GCP.

- **Automated, Auditable Change Management**  
  Use Terraform's plan and apply in gitOps workflows to preview, approve, and track changes to your Jamf Pro environment, ensuring all modifications are intentional, reviewed, and logged.

- **Environment Replication and Drift Detection**
  Reproduce Jamf Pro tenant configurations across multiple environments (development, staging, production) or tenants, and detect configuration drift over time using Terraform’s state management.

- **Disaster Recovery and Rapid Rebuilds**  
  Store your Jamf Pro configuration in code, allowing for rapid recovery or migration of tenant settings, policies, and assignments in the event of accidental changes or tenant loss.

- **Collaboration and Delegation**
  Empower teams to collaborate on Jamf Pro configuration using pull requests, code reviews, and CI/CD pipelines, reducing bottlenecks and enabling safe delegation of administrative tasks.

- **Bulk and Consistent Policy Enforcement**
  Apply security, compliance, and device management policies at scale, ensuring consistency and reducing manual configuration errors across large organizations or multiple tenants.

- **Self-Service via Terraform Modules**  
  Build reusable Terraform modules for common Jamf Pro workloads, enabling service-owning teams to provide self-service provisioning to other engineering teams while maintaining standards and reducing manual effort.

- **Integration with Policy-as-Code (OPA/Conftest)**  
  Integrate with Open Policy Agent (OPA) or Conftest to enforce organizational standards, compliance, and guardrails on Jamf Pro resources before deployment, ensuring only approved configurations are applied in production.

- **Guardrailed Deployments**  
  Implement automated checks and guardrails in CI/CD pipelines to prevent misconfiguration and enforce best practices, reducing risk and improving governance for Jamf Pro administration.

## Getting Started

Please refer to the [Getting Started](https://registry.terraform.io/providers/deploymenttheory/jamfpro/latest/docs) guide in the terraform registry for more information on how to get started.

## Demo Implementation

To help you get started and understand the practical implementation of this provider, we've created a comprehensive demo repository:

- **Demo Repository**: [Terraform Demo Jamf Pro](https://github.com/deploymenttheory/terraform-demo-jamfpro-v2)

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

## Getting Started

Please refer to the [Getting Started](https://registry.terraform.io/providers/deploymenttheory/jamfpro/latest/docs) guide in the terraform registry for more information on how to get started.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.13.0
- [Go](https://golang.org/doc/install) >= 1.22.4
- [Jamf Pro](https://www.jamf.com/) >= 11.20.0

(Tested with production Jamf Pro instances, with and without SSO integratioin with Microsoft Entra ID. We do not test against beta or preview versions of Jamf Pro due to potential data model changes.)

## Jamf Cloud Load Balancing and Cookies

- Jamf Cloud uses a load balancer to distribute traffic across multiple web app members (typically 2). When resource's are manipulated on a given web app member, there is up to a 60 second time box until this resources changes are propagated and reflected onto the other web app(s). This architecture can cause issues with Terraform's http client default behaviour when multiple instances are running in parallel and also due to the speed terraform operates. This results in scenarios where it's very likely that a create by terraform, followed by a read (for stating) will freqently communicate with different web app members during a terraform run. This causes stating 'unfound' resource issues.
- To mitigate this please use the `jamfpro_load_balancer_lock` (which enforces a single cookie across all parallel instances of Terraform operations). This feature on first run obtains all available web cookies (jpro-ingress) from Jamf Pro and selects and applies a single one to the http client for all subsequent api calls during the terraform run. This is eqivalent to a sticky session.
- For non Jamf Cloud customers, with load balanced configurations please use `custom_cookies` and configure a custom cookie to be used in all requests instead.

### Concurrency

> [!WARNING]
> Jamf Pro produces inconsistent behaviour when using the default parallelism setting of 10 with terraform. You can adjust paralellism by setting the Terraform parallelism count using `terraform apply -parallelism=X` to a setting of your choice. [HashiCorp Docs](https://developer.hashicorp.com/terraform/cli/commands/apply#parallelism-n) . It's recconmended to always set parallelism to 1 to guarantee successful CRUD operations and resource stating. What this produces in a moderate performance hit is offset by reliability. Not using `-parallelism=1` is at your own risk!

## Community & Support

For further community support and to engage with other users of the Jamf Pro Terraform Provider, please join us on the Mac Admins Slack channel. You can ask questions, provide feedback, and share best practices with the community. Join us at:

- [Mac Admins Slack Channel](https://macadmins.slack.com/archives/C06R172PUV6) - #terraform-provider-jamfpro


## Disclaimer

> [!IMPORTANT]  
> While every effort is made to maintain accuracy and reliability, users should thoroughly test configurations in non-production environments before deploying to production. Always refer to official Jamf documentation for the most up-to-date information on Jamf Pro services and features.
