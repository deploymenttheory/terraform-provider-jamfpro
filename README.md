# Terraform Provider for Jamf Pro

This repository contains the [Terraform](https://www.terraform.io) provider for managing resources in Jamf Pro. It is built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) and intended for creating and managing Jamf Pro entities like departments, sites, policies, etc.

The provider contains:

- Resources and data sources for Jamf Pro entities (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.


Once the provider is finalized, you may want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go Get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

To use this provider, you need to configure it with your Jamf Pro instance's URL and authentication credentials. The provider allows you to manage various resources such as departments, sites, and policies within Jamf Pro.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

To run the full suite of Acceptance tests, execute `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Benchmarking

Indicative time taken to implement changes using this terraform provider against a jamf cloud instance of jamf pro.

Setup:
Jamf Cloud instance of Jamf Pro
Empty instance of jamf pro with 0 devices enrolled
Deployment of jamf pro departments (Simplement resource type in jamf pro)

| API Calls      | Resource Type | CRUD Operation        | Time Taken   |
| ---------------| ------------- |-----------------------|--------------|
| 1000           | Departments   | Create                | 8m10s        |
| 1000           | Departments   | Read                  | 8ms          |
| 1000           | Departments   | Update                | 8ms          |
| 1000           | Departments   | delete                | 8ms          |
| 10000          | Departments   | Create (TF Plan+Apply)| 8m10s        |
| 10000          | Departments   | Read (TF Plan)        | 2m49s        |
| 10000          | Departments   | Update (TF Plan+Apply)| 8ms          |
| 10000          | Departments   | delete                | 8ms          |


## Providers

No providers.

## Modules

No modules.

## Resources

The provider currently has working coverage of the following jamf pro resource types

- account groups
- departments
- disk encryption configurations
- dock items
- printers
- scripts
- sites

## Inputs

No inputs.

## Outputs

No outputs.
<!-- END_TF_DOCS -->