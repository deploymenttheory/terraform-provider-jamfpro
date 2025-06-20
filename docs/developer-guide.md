# Go Style Guide for terraform-provider-jamfpro

## Base Style Guide

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go) as its foundation.

## Development Setup & Workflow

- Use the provided **GNUmakefile** commands for all build and test tasks:
  - `make build` to compile the provider code.
  - `make install` to build and install the provider locally.
  - `make lint` to run linters and ensure code style compliance.
  - `make test` to run all tests.
  - `make testacc` to run acceptance tests (integration tests).
  - `make docs` to regenerate documentation.
  - `make fmt` to format the code.
- Always run the above `make` commands from the repository root.

## File and Folder Structure

### Repository Structure

```
terraform-provider-jamfpro/
├── docs/               # Documentation for data sources and resources
├── examples/           # Example configurations for each resource/data source
│   ├── data-sources/   # Examples for data sources
│   ├── resources/      # Examples for resources
│   └── provider/       # Provider configuration examples
├── internal/           # Internal provider code
│   ├── provider/       # Core provider implementation
│   └── resources/      # Individual resource implementations
├── scripts/            # Maintenance and utility scripts
├── testing/            # Test configurations and fixtures
└── tools/              # Development tools and utilities
```

### Resource Organization

- All resource implementations are within the `internal/resources` directory, with each resource in its own subdirectory.
- Name resource directories using lowercase words with underscores (e.g., `policy`, `building`, `script`).
- Choose resource names that reflect the Jamf Pro resource domain they represent.

### Resource Files

Each resource directory should contain:

- **Core Files:**
  - `resource.go` - Resource type definitions and main logic
  - `crud.go` - Basic CRUD operations
  - `data_source.go` - Data source implementation
  - `constructor.go` - Resource schema assembly and initialization
  - `data_validator.go` - Custom validation functions (if needed)

- **Schema Files:**
  - Split complex resource schemas into logical components using separate `schema_*.go` files
  - Examples: `schema_account_maintenance.go`, `schema_network_limitations.go`, `schema_reboot.go`

- **State Management Files:**
  - Separate state handling into focused files with `state_*.go` naming
  - Examples: `state_payloads.go`, `state_migration.go`, `state_general.go`

## Naming Conventions

### Resource and Data Source Names

- **Resource Names:** Follow the pattern `jamfpro_resource_name` (e.g., `jamfpro_policy`, `jamfpro_building`).
- **Data Source Names:** Follow the pattern `jamfpro_datasource_name` (e.g., `jamfpro_policy`, `jamfpro_building`).

### Function Naming Conventions

#### Exported Functions (Resource Entry Points)

```go
// Pattern: ResourceJamfPro{ResourceName}() *schema.Resource
func ResourceJamfProPolicies() *schema.Resource       // ✓ Correct
func ResourceJamfProBuildings() *schema.Resource      // ✓ Correct
```

**Rules:**

- Must be the **only exported function** in the resource package
- Always returns `*schema.Resource`
- Uses PascalCase with "ResourceJamfPro" prefix

#### Data Source Functions

```go
// Pattern: DataSourceJamfPro{ResourceName}() *schema.Resource
func DataSourceJamfProPolicies() *schema.Resource     // ✓ Correct
func DataSourceJamfProBuildings() *schema.Resource    // ✓ Correct
```

#### CRUD Operations (Internal Functions)

```go
// All CRUD functions use lowercase naming:
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics  
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

// Standard read variants:
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics
```

#### Helper Functions

```go
// Use descriptive names that clearly indicate purpose:
func construct(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error)
func updateState(resource *jamfpro.ResourcePolicy, d *schema.ResourceData) diag.Diagnostics
func constructPayloads(d *schema.ResourceData, resource *jamfpro.ResourcePolicy)
```

## Schema Definition Guidelines

- Define complete schemas with proper attribute types (String, Int, Bool, etc.).
- Mark attributes explicitly as `Required`, `Optional`, or `Computed`.
- Use `Computed: true` for server-generated fields like IDs.
- Use `Optional: true` with `Computed: true` for fields that can be specified or defaulted by the service.
- Include standard timeouts using `schema.ResourceTimeout`.
- Write clear descriptions for each attribute.

### Schema Descriptions

- Schema key descriptions follow a **common sense approach**
- Focus on **clarity for Terraform users** rather than API documentation verbatim
- Prioritize user understanding over technical accuracy

```go
// ✓ Good: User-friendly description
"name": {
    Type:        schema.TypeString,
    Required:    true,
    Description: "The name of the policy.",
},

// ✗ Avoid: Too technical/API-focused  
"name": {
    Type:        schema.TypeString,
    Required:    true,
    Description: "ResourcePolicy.General.Name field as defined in jamfpro.ResourcePolicy struct",
},
```

## Common Utilities and Shared Code

### Common CRUD Operations

Use the shared CRUD operations in `internal/resources/common/crud.go` when the resource type is performing
single api calls per resource.

```go
// Create
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    return common.Create(
        ctx,
        d,
        meta,
        construct,
        meta.(*jamfpro.Client).CreatePolicy,
        readNoCleanup,
    )
}

// Read
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
    return common.Read(
        ctx,
        d,
        meta,
        cleanup,
        meta.(*jamfpro.Client).GetPolicyByID,
        updateState,
    )
}

// Update
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    return common.Update(
        ctx,
        d,
        meta,
        construct,
        meta.(*jamfpro.Client).UpdatePolicyByID,
        readNoCleanup,
    )
}

// Delete
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    return common.Delete(
        ctx,
        d,
        meta,
        meta.(*jamfpro.Client).DeletePolicyByID,
    )
}
```


### Complex Resource CRUD Operations

For resources that require multiple API calls or complex operations, implement custom CRUD functions directly in the resource's `crud.go` file. These resources typically involve:

1. Multiple related API endpoints
2. File uploads or downloads
3. Multi-step operations with dependencies
4. Complex state management
5. Verification or validation steps

#### Example: Package Resource

The package resource is a good example of a complex resource that requires multiple API calls and file handling:

```go
// create handles the creation of a Jamf Pro package resource:
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*jamfpro.Client)
    var diags diag.Diagnostics
    
    // 1. Construct the resource and get file path
    resource, localFilePath, err := construct(d)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
    }
    
    // 2. Calculate initial file hash for verification
    initialHash, err := jamfpro.CalculateSHA3_512(localFilePath)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to calculate SHA3-512: %v", err))
    }
    
    // 3. Create package metadata in Jamf Pro
    err = retry.RetryContext(ctx, PackagesMetaTimeout, func() *retry.RetryError {
        creationResponse, err := client.CreatePackage(*resource)
        if err != nil {
            return retry.RetryableError(fmt.Errorf("failed to create package metadata: %v", err))
        }
        packageID = creationResponse.ID
        return nil
    })
    
    // 4. Upload the actual package file
    err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
        _, err = client.UploadPackage(packageID, []string{localFilePath})
        if err != nil {
            return retry.RetryableError(fmt.Errorf("failed to upload package file: %v", err))
        }
        return nil
    })
    
    // 5. Verify the uploaded file hash matches the original
    if err := verifyPackageUpload(ctx, client, packageID, resource.FileName, initialHash,
        d.Timeout(schema.TimeoutCreate)); err != nil {
        return diag.FromErr(fmt.Errorf("failed to verify upload: %v", err))
    }
    
    // 6. Set ID and read back state
    d.SetId(packageID)
    return append(diags, readNoCleanup(ctx, d, meta)...)
}
```

#### Example: User Initiated Enrollment Resource

The User Initiated Enrollment resource is another complex example that manages multiple related configurations through different API endpoints:

```go
// create is responsible for creating jamf pro User-initiated enrollment settings
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*jamfpro.Client)
    var diags diag.Diagnostics

    // 1. Update main enrollment settings
    enrollmentSettings, err := constructEnrollmentSettings(d)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to construct enrollment settings: %v", err))
    }
    
    err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
        _, apiErr := client.UpdateEnrollment(enrollmentSettings)
        if apiErr != nil {
            return retry.RetryableError(apiErr)
        }
        return nil
    })
    
    // 2. Create language messaging configurations
    messagingList, err := constructEnrollmentMessaging(d, client)
    for i := range messagingList {
        message := &messagingList[i]
        err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
            _, apiErr := client.UpdateEnrollmentMessageByLanguageID(message.LanguageCode, message)
            if apiErr != nil {
                return retry.RetryableError(apiErr)
            }
            return nil
        })
    }
    
    // 3. Create directory service group enrollment settings
    accessGroups, err := constructDirectoryServiceGroupSettings(d)
    for i := range accessGroups {
        group := accessGroups[i]
        err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
            _, apiErr := client.CreateAccountDrivenUserEnrollmentAccessGroup(group)
            if apiErr != nil {
                return retry.RetryableError(apiErr)
            }
            return nil
        })
    }
    
    // 4. Set singleton ID and read back state
    d.SetId(ResourceIDSingleton)
    return append(diags, readNoCleanup(ctx, d, meta)...)
}
```

### Guidelines for Complex Resources

When implementing complex resources that require multiple API calls:

1. **Structure your code clearly:**
   - Break down operations into logical steps
   - Use helper functions for distinct operations
   - Consider using separate files for different aspects of functionality

2. **Error handling:**
   - Implement proper cleanup on partial failures
   - Use descriptive error messages that indicate which step failed
   - Consider implementing rollback mechanisms for multi-step operations

3. **Retry logic:**
   - Apply retry.RetryContext to each API call that might fail transiently
   - Use appropriate timeouts for different operations (e.g., longer timeouts for file uploads)
   - Log retry attempts at appropriate levels

4. **State management:**
   - Ensure resource state is properly updated after all operations
   - Handle partial state updates carefully
   - Consider using transaction-like patterns for complex updates

5. **Verification:**
   - Implement verification steps for critical operations
   - Add validation checks between steps
   - Log operation progress for debugging

### Shared Schemas

Use shared schemas from `internal/resources/common/sharedschemas` for common attributes:

```go
"category_id": sharedschemas.GetSharedSchemaCategory(),
"site_id": sharedschemas.GetSharedSchemaSite(),
```

## Error Handling and Retry Logic

- Use the retry mechanisms provided by the Terraform SDK for API calls:
  ```go
  err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
      var apiErr error
      response, apiErr = serverOutcomeFunc(resourceID)
      if apiErr != nil {
          return retry.RetryableError(apiErr)
      }
      return nil
  })
  ```

- Return detailed diagnostics with `diag.FromErr()` for user-friendly messages.

## Singleton Configuration Resources

Some resources manage singleton configurations rather than individual entities (e.g., global settings). These follow a modified CRUD pattern:

**Characteristics:**

- Represent system-wide configuration that always exists
- Cannot be "created" or "deleted" in the traditional sense
- Use hardcoded singleton IDs

**Singleton ID Naming:**

- Format: `jamfpro_{resource_name}_singleton`
- Examples: 
  - `jamfpro_computer_inventory_collection_settings_singleton`
  - `jamfpro_device_communication_settings_singleton`

## Testing Guidelines

The project uses Terraform's built-in testing framework for integration testing. All new resources must include tests before PRs can be merged.

### Test Structure

1. **Location:**
   - Place test configurations in `testing/payloads/[resource_name]/`
   - The folder name must match the resource name (e.g., `jamfpro_building`)
   - The resource folder in `internal/resources/` must also match (e.g., `building`)

2. **Test Naming Convention:**
   - All test resources must have a name prefix of `tf-testing-` to ensure they're properly cleaned up
   - Example: `tf-testing-script-min`
   - Use the `${var.testing_id}` and `${random_id.rng.hex}` variables to ensure unique resource names

3. **Test Content:**
   - Create Terraform configurations that exercise different aspects of the resource
   - Include minimal examples (required fields only) and comprehensive examples (all fields)
   - Test edge cases and validation rules
   - For resources that support bulk operations, include tests with multiple instances

4. **Provider Configuration:**
   - Use symlinks to the root `provider.tf` file in each test directory
   - Create the symlink with: `ln -s ../../provider.tf provider.tf`

### Example Test Configuration

```hcl
// Minimal configuration
resource "jamfpro_script" "min_script" {
  name            = "tf-testing-${var.testing_id}-min-${random_id.rng.hex}"
  script_contents = "script_contents_field"
  priority        = "BEFORE"
}

// Comprehensive configuration
resource "jamfpro_script" "max_script" {
  name            = "tf-testing-${var.testing_id}-max-${random_id.rng.hex}"
  category_id     = "9"
  info            = "info_field"
  notes           = "notes_field"
  os_requirements = "os_requirements_field"
  priority        = "BEFORE"
  script_contents = "script_contents_field"
  parameter4      = "parameter4_field"
  parameter5      = "parameter5_field"
  parameter6      = "parameter6_field"
  parameter7      = "parameter7_field"
  parameter8      = "parameter8_field"
  parameter9      = "parameter9_field"
  parameter10     = "parameter10_field"
  parameter11     = "parameter11_field"
}
```

### Running Tests

Tests are automatically triggered by the GitHub workflow on PRs that modify files in the `internal/` directory. To run tests locally:

1. Set up a Python virtual environment in the `testing` directory:
   ```bash
   cd testing
   python -m venv .venv
   source .venv/bin/activate
   pip install -r requirements.txt
   ```

2. Run the test script:
   ```bash
   ./run_tests.sh
   ```

### Test Cleanup

A cleanup process runs at 23:59 daily to remove test resources. Tests running during this time may fail. The cleanup process:

1. Identifies resources with the `tf-testing` prefix
2. Removes these resources from the Jamf Pro instance
3. Ensures the test environment stays clean

## Comments and Documentation

- Write Go comments only on exported functions, types, and methods to explain their purpose, parameters, and return values when it adds clarity.
- Focus comments on **why** something is done if it's not obvious from the code.
- Avoid redundant comments that just restate the code or don't provide additional insight.
- Use descriptive function and variable names to make the code self-documenting.

## Example Files

- All resources and data sources must have an example file.
- Example files must be named `resource.tf` or `data-source.tf` and placed in the appropriate directory under `examples/resources/{resource_name}` or `examples/data-sources/{data_source_name}`.
- Use a single `.tf` file per example
- Include comments explaining the purpose of the example and any placeholders
- If the resource supports import, include an import command in a comment
