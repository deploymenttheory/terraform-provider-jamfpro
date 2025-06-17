# Go Style Guide

## Base Style Guide

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go) as its foundation.

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

