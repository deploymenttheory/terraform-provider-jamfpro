# Go Style Guide

## Base Style Guide

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go) as its foundation.

## Provider-Specific Code Style

The `internal/resources/policy` directory represents our latest iteration of code organization and patterns. Follow these conventions:

1. Schema Organization:
   - Split complex resource schemas into logical components using separate `schema_*.go` files
   - Examples: `schema_account_maintenance.go`, `schema_network_limitations.go`, `schema_reboot.go`
   - Each schema file should focus on one specific aspect of the resource

2. State Management:
   - Separate state handling into focused files with `state_*.go` naming
   - Use `state_payloads.go` for API request/response structures
   - Use `state_migration.go` for version migrations
   - Split complex state operations into logical groups (e.g., `state_general.go`, `state_scope.go`)

3. Core Files:
   - `constructor.go` - Resource schema assembly and initialization
   - `crud.go` - Basic CRUD operations
   - `data_source.go` - Data source implementation
   - `resource.go` - Resource type definitions and main logic
   - `data_validator.go` - Custom validation functions

4. File Organization:
   - Keep files focused and single-purpose
   - Use clear, descriptive file names that indicate their content
   - Break down large schemas into manageable components

## Repository Structure

```
terraform-provider-jamfpro/
├── docs/               # Documentation for data sources and resources
├── examples/          # Example configurations for each resource/data source
├── internal/          # Internal provider code
│   ├── provider/     # Core provider implementation
│   └── resources/    # Individual resource implementations
├── scripts/          # Maintenance and utility scripts
├── testing/          # Test configurations and fixtures
└── tools/            # Development tools and utilities
```

## Resource Implementation Structure

Each resource in `internal/resources/` follows a consistent file structure:

- `constructor.go` - Resource schema and constructor functions
- `crud.go` - Create, Read, Update, Delete operations
- `data_source.go` - Data source implementation
- `resource.go` - Resource type definitions and main implementation
- `state.go` - State management functions
- `data_customdiff.go` - Custom diff functions (if needed)

The `internal/resources/common/` directory contains shared code and utilities used across multiple resources.

## Testing Guidelines

1. All tests are located in the `testing/` directory
2. Follow the test guide in `docs/` when submitting a PR

