# Terraform Provider for Jamf Pro

This repository contains the [Terraform](https://www.terraform.io) provider for managing resources in Jamf Pro. It is built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) and intended for creating and managing Jamf Pro entities like departments, sites, policies, etc.

The provider contains:

- Resources and data sources for Jamf Pro entities (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

These files contain the actual code and configurations for the Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform. 

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for guidance on creating a new repository.

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
go get github.com/author/dependency
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