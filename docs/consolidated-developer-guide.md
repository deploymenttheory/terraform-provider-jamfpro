# Terraform Provider JamfPro - Developer Guide

This guide provides comprehensive information for developers working on the terraform-provider-jamfpro project. It consists of two main sections: the Style Guide (how we write Go code in this project) and the Developer Guide (project structure, workflow, and implementation patterns).

## Table of Contents

1. [Developer Guide](#developer-guide)
2. [Style Guide](#style-guide)

## Developer Guide

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
│       └── common/     # Shared code and utilities
├── scripts/            # Maintenance and utility scripts
├── testing/            # Test configurations and fixtures
└── tools/              # Development tools and utilities
```

### Development Setup & Workflow

Use the provided **GNUmakefile** commands for all build and test tasks:

- `make build` - compile the provider code
- `make install` - build and install the provider locally
- `make lint` - run linters and ensure code style compliance
- `make test` - run all tests
- `make testacc` - run acceptance tests (integration tests)
- `make docs` - regenerate documentation
- `make fmt` - format the code

Always run these commands from the repository root.

### Resource Organization

#### Directory Structure

- All resource implementations live in `internal/resources/`
- Each resource has its own subdirectory
- Name directories using lowercase words with underscores (e.g., `policy`, `building`, `script`)
- Choose names that reflect the Jamf Pro resource domain they represent

#### File Organization

**Keep files focused and single-purpose:**
- Break down large schemas into manageable, logical components
- Use clear, descriptive file names that indicate their content
- Avoid monolithic files - prefer multiple smaller, focused files

**Required Core Files:**

Every resource directory must contain these core files:

- **`resource.go`** - Resource type definitions and main logic
- **`crud.go`** - CRUD operations (may use common functions or custom implementation)
- **`data_source.go`** - Data source implementation  
- **`constructor.go`** - Resource schema assembly and initialization
- **`data_validator.go`** - Custom validation functions (if needed)

**Optional Schema Files:**

For complex resources, split schemas into logical components:

- **`schema_*.go`** - Schema definitions for specific resource aspects
  - Examples: `schema_account_maintenance.go`, `schema_network_limitations.go`, `schema_reboot.go`
  - Each file should focus on one specific aspect of the resource
  - Use descriptive names that clearly indicate the schema's purpose

**Optional State Files:**

For complex state management, separate into focused files:

- **`state_*.go`** - State management for specific resource aspects
  - Examples: `state_payloads.go`, `state_migration.go`, `state_general.go`
  - Split complex state operations into logical groups
  - Keep state construction and updates separate from business logic

**Additional Helper Files:**

- **`helpers.go`** - Resource-specific utility functions (for complex resources)
- **`data_custom_diff.go`** - Custom diff functions (if needed)

### Resource Naming Conventions

- **Resource Names:** `jamfpro_resource_name` (e.g., `jamfpro_policy`, `jamfpro_building`)
- **Data Source Names:** `jamfpro_datasource_name` (e.g., `jamfpro_policy`, `jamfpro_building`)

### Common Folder Usage

The `internal/resources/common/` directory contains shared code and utilities used across multiple resources:

#### Shared CRUD Operations

Use `common/crud.go` for resources that perform single API calls per operation:

```go
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    return common.Create(
        ctx,
        d,
        meta,
        construct,
        meta.(*jamfpro.Client).CreatePolicy,
        readNoCleanup,
    )
}
```

#### Shared Schemas

Use shared schemas from `common/sharedschemas` for common attributes:

```go
"category_id": sharedschemas.GetSharedSchemaCategory(),
"site_id": sharedschemas.GetSharedSchemaSite(),
```

#### Other Common Utilities

- `common/constructors/` - Shared construction patterns
- `common/configurationprofiles/` - Configuration profile utilities
- `common/jamfprivileges/` - Jamf privilege management helpers

### Complex vs Simple Resources

#### Simple Resources

Use the shared CRUD operations for resources that:
- Make single API calls per CRUD operation
- Have straightforward state management
- Don't require file uploads or multi-step operations

#### Complex Resources

Implement custom CRUD functions for resources that require:
- Multiple API calls per operation
- File uploads or downloads
- Multi-step operations with dependencies
- Complex state management
- Verification or validation steps

**Example: Package Resource**

The package resource demonstrates a complex implementation with multiple steps:

1. **Create metadata** - First API call to create package record
2. **Upload file** - Second API call to upload actual package file  
3. **Verify upload** - Custom verification logic to ensure file integrity
4. **Cleanup** - Remove temporary files and handle rollback on failures

**Other examples of complex resources:**
- `package` - handles file uploads and verification
- `user_initiated_enrollment_settings` - manages multiple configuration endpoints

### Singleton Configuration Resources

Some resources represent system-wide configuration that always exists in Jamf Pro:

**Characteristics:**
- Cannot be "created" or "deleted" in the traditional sense (the configuration always exists in Jamf Pro)
- Use hardcoded singleton IDs to represent the single instance of the configuration
- Create operations call UPDATE API methods since the configuration already exists
- Delete operations only remove from Terraform state, leaving the configuration unchanged in Jamf Pro

**Singleton ID Naming:**
- Format: `jamfpro_{resource_name}_singleton`
- Examples: 
  - `jamfpro_computer_inventory_collection_settings_singleton`
  - `jamfpro_device_communication_settings_singleton`
  - `jamfpro_client_checkin_singleton`

### Data Source Location

Data sources are located within the resource directories rather than in a separate `data-sources/` directory. This is a **historical decision** - while data sources could logically live in their own directory structure, we don't see value in updating this organization currently. New data sources should continue to live in the resource folder.

### Schema Definition Guidelines

- Define complete schemas with proper attribute types (String, Int, Bool, etc.)
- Mark attributes explicitly as `Required`, `Optional`, or `Computed`
- Use `Computed: true` for server-generated fields like IDs
- Use `Optional: true` with `Computed: true` for fields that can be specified or defaulted by the service
- Include standard timeouts using `schema.ResourceTimeout`
- Write clear, user-friendly descriptions for each attribute

### Error Handling and Retry Logic

Use retry mechanisms provided by the Terraform SDK for API calls:

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

Return detailed diagnostics with `diag.FromErr()` for user-friendly messages.

### Testing Guidelines

All new resources **must include tests** before PRs can be merged. 

For comprehensive testing guidelines, including test structure, naming conventions, automated testing workflows, and local testing instructions, see:

**📖 [Testing Implementation Guide](testing-implementation.md)**

**Key Requirements:**
- Tests must be placed in `testing/payloads/[resource_name]/`
- All test resources must have a `tf-testing-` name prefix
- Tests are automatically triggered by GitHub workflows when PRs modify `internal/` files

### Example Files Requirements

- All resources and data sources must have example files
- Place in `examples/resources/{resource_name}/resource.tf` or `examples/data-sources/{data_source_name}/data-source.tf`
- Use a single `.tf` file per example
- Include comments explaining the purpose and any placeholders
- If the resource supports import, include an import command in comments

### Additional Resources

- Follow the testing guide in `docs/` when submitting PRs
- Refer to `internal/resources/policy/` as the current best practice implementation
- Use the `GNUmakefile` commands for all development tasks 

---

## Style Guide

### Base Style Guide

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go) as its foundation. Where reasonable and in the absence of a more specific style, we follow "clean code" practices as outlined in the [Clean Code book](https://github.com/Gatjuat-Wicteat-Riek/clean-code-book).

### Reference Implementation

The `internal/resources/policy` directory represents our **current best practice** for code organization and patterns. All new resources should follow these conventions, and existing resources should work toward this standard over time.

### Function Naming Conventions

#### Exported Functions (Resource Entry Points)

```go
// Pattern: ResourceJamfPro{ResourceName}() *schema.Resource
func ResourceJamfProPolicies() *schema.Resource       // ✓ Correct
func ResourceJamfProBuildings() *schema.Resource      // ✓ Correct
func ResourceJamfProSites() *schema.Resource          // ✓ Correct
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

All CRUD functions use **lowercase naming** and follow exact signature patterns:

```go
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics  
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics

// Standard read variants:
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics
```

**Read Function Cleanup Parameter:**
- `cleanup bool` parameter determines whether to remove the resource from Terraform state if it's not found in the API
- `readWithCleanup` - removes from state when resource not found (used in normal operations)
- `readNoCleanup` - preserves state when resource not found (used after create/update operations)


#### Standard Function Names

These are standardized function names that each resource should implement with consistent signatures:

```go
func construct(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error)
func updateState(resource *jamfpro.ResourcePolicy, d *schema.ResourceData) diag.Diagnostics
func constructPayloads(d *schema.ResourceData, resource *jamfpro.ResourcePolicy)
```

**Guidelines:**
- `construct` - builds the API request object from Terraform ResourceData
- `updateState` - updates Terraform state from API response object
- `constructPayloads` - used for complex resources with multiple payload structures

### Variable Naming

- Follow Go conventions: camelCase for local variables, PascalCase for exported variables
- Use descriptive names that clearly indicate the variable's purpose
- Avoid abbreviations unless they're well-established (e.g., `ctx` for context, `err` for error)

### Schema Descriptions

Schema key descriptions follow a **common sense approach**:

**Naming Convention:** When possible, descriptions should follow the naming and terminology used in the Jamf Pro GUI to provide familiarity for users transitioning from the web interface to Terraform.

- Focus on **clarity for Terraform users** rather than API documentation verbatim
- Prioritize user understanding over technical accuracy
- Descriptions **don't always match the API documentation** exactly

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

### Comments Guidelines

- **Use comments sparingly** and only when necessary
- Comments should explain **why** something is done, not **what** is being done
- The code itself should be self-documenting for the "what"
- Avoid obvious comments that just restate the code
- Write Go comments on exported functions, types, and methods to explain their purpose when it adds clarity

```go
// ✓ Good: Explains why
// Use hardcoded singleton ID since this represents global configuration
d.SetId("jamfpro_computer_inventory_collection_settings_singleton")

// ✗ Avoid: Just restates the code
// Set the ID to the singleton value
d.SetId("jamfpro_computer_inventory_collection_settings_singleton")
```

### Code Patterns

#### Resource Function Pattern

Every resource package exports exactly one function following this pattern:

```go
func ResourceJamfProPolicies() *schema.Resource {
    return &schema.Resource{
        CreateContext: create,           // ← lowercase CRUD functions
        ReadContext:   readWithCleanup,  
        UpdateContext: update,
        DeleteContext: delete,
        // ... schema definition
    }
}
```

#### Standard CRUD Pattern

For simple resources that make single API calls, use the common CRUD operations:

```go
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    return common.Create(ctx, d, meta, construct, meta.(*jamfpro.Client).CreatePolicy, readNoCleanup)
}

func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
    return common.Read(ctx, d, meta, cleanup, meta.(*jamfpro.Client).GetPolicyByID, updateState)
}
```

#### Complex CRUD Pattern

For resources that require multiple API calls within a single CRUD operation (like package uploads with verification):

```go
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    client := meta.(*jamfpro.Client)
    var packageID string
    
    // Step 1: Construct and validate resource
    resource, localFilePath, err := construct(d)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to construct: %v", err))
    }
    
    // Step 2: Create metadata
    err = retry.RetryContext(ctx, PackagesMetaTimeout, func() *retry.RetryError {
        response, err := client.CreatePackage(*resource)
        if err != nil {
            return retry.RetryableError(err)
        }
        packageID = response.ID
        return nil
    })
    
    // Step 3: Upload file
    err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
        _, err = client.UploadPackage(packageID, []string{localFilePath})
        return retry.RetryableError(err)
    })
    
    // Step 4: Verify upload (custom verification step)
    if err := verifyPackageUpload(ctx, client, packageID, resource.FileName, 
        initialHash, d.Timeout(schema.TimeoutCreate)); err != nil {
        return diag.FromErr(fmt.Errorf("verification failed: %v", err))
    }
    
    d.SetId(packageID)
    return readNoCleanup(ctx, d, meta)
}
```

#### Singleton Configuration Pattern

For resources that manage system-wide configuration:

```go
// create calls UPDATE API method since configuration always exists
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    // ... construct and update configuration
    
    // Use descriptive singleton ID
    d.SetId("jamfpro_computer_inventory_collection_settings_singleton")
    return readNoCleanup(ctx, d, meta)
}

// delete only removes from Terraform state, doesn't delete from API
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    d.SetId("")
    return nil
}
```
