# Adding Tests to Your PRs

## Requirements

All new contributors **MUST** add tests for their resources and have them pass before PRs can be merged.

## Test Structure

### Location
Place your tests in `testing/payloads/[resource_name]/`

### Naming Convention
- The folder name in `testing/payloads/` **MUST** match the resource name
- The resource folder in `internal/services/` **MUST** also match the resource name
- Example: `jamfpro_building` resource needs:
  - `testing/payloads/jamfpro_building/` folder
  - `internal/services/building/` folder

### Test Content
Each test folder should contain Terraform modules that will be executed automatically by the integration testing workflow.

## Automated Testing

Tests are automatically triggered by the `jamfpro-provider-integration-test.yml` GitHub workflow on:
- Pull requests that modify files in `internal/`
- Manual workflow dispatch

The workflow will:
1. Compile the provider binary
2. Generate test targets based on changed files
3. Execute the relevant Terraform modules
4. Clean up resources after testing

## Getting Started

1. Create your test folder: `testing/payloads/[your_resource_name]/`
2. Add your Terraform configuration files
3. Ensure your resource implementation exists in `internal/services/[resource_name]/`
4. Submit your PR - tests will run automatically
