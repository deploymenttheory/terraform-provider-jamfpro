# General Style Guide

Where reasonable and in the absence of a more specific style, we will follow "clean code" practices.
https://github.com/Gatjuat-Wicteat-Riek/clean-code-book


# Go Style Guide

## Base Style Guide

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go) as its foundation.

## Provider-Specific Code Style

The `internal/resources/policy` directory represents our latest iteration of code organization and patterns. Follow these conventions:

### 1. Schema Organization

Split complex resource schemas into logical components using separate `schema_*.go` files:

**Examples from policy resource:**
- `schema_account_maintenance.go` - Account management configurations
- `schema_network_limitations.go` - Network restriction settings  
- `schema_reboot.go` - Reboot behavior configurations
- `schema_user_interaction.go` - User interaction prompts

**Guidelines:**
- Each schema file should focus on one specific aspect of the resource
- Use descriptive names that clearly indicate the schema's purpose
- Keep related functionality grouped together

### 2. State Management

Separate state handling into focused files with `state_*.go` naming:

**File Structure:**
- `state_payloads.go` - API request/response structures
- `state_migration.go` - Version migrations
- `state_general.go`, `state_scope.go` - Logical groupings of state operations

**Guidelines:**
- Split complex state operations into logical groups
- Use consistent naming patterns for state management functions
- Keep state construction and updates separate from business logic

### 3. Core Files

Every resource directory must contain these core files:

- **`constructor.go`** - Resource schema assembly and initialization
- **`crud.go`** - Basic CRUD operations
- **`data_source.go`** - Data source implementation  
- **`resource.go`** - Resource type definitions and main logic
- **`data_validator.go`** - Custom validation functions (if needed)

### 4. Function Naming Conventions

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

### 5. File Organization

**Keep files focused and single-purpose:**
- Break down large schemas into manageable, logical components
- Use clear, descriptive file names that indicate their content
- Avoid monolithic files - prefer multiple smaller, focused files

**File Naming Patterns:**
- `schema_*.go` - Schema definitions for specific resource aspects
- `state_*.go` - State management for specific resource aspects  
- `crud.go` - Standard CRUD operations
- `constructor.go` - Resource construction logic
- `data_validator.go` - Custom validation logic

### 6. Comments

- **Use comments sparingly** and only when necessary
- Comments should explain **why** something is done, not **what** is being done
- The code itself should be self-documenting for the "what"
- Avoid obvious comments that just restate the code

### 7. Schema Descriptions

- Schema key descriptions follow a **common sense approach**
- Descriptions **don't always match the API documentation** exactly  
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

**Required Files:**
- `constructor.go` - Resource schema and constructor functions
- `crud.go` - Create, Read, Update, Delete operations
- `data_source.go` - Data source implementation
- `resource.go` - Resource type definitions and main implementation
- `state.go` - State management functions
- `data_custom_diff.go` - Custom diff functions (if needed)

The `internal/resources/common/` directory contains shared code and utilities used across multiple resources.

### Data Source Location

Data sources are currently located within the resource directories rather than in a separate `data-sources/` directory. This is a **historical decision**, and while we recognize that data sources could logically live in their own directory structure, we don't see the value in updating this organization at this time. For now, data sources should continue to live in the resource folder.

## Testing Guidelines

1. All tests are located in the `testing/` directory
2. Follow the test guide in `docs/` when submitting a PR

## Code Examples

### Resource Function Pattern
```go
// Every resource package exports exactly one function following this pattern:
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

### CRUD Function Pattern  
```go
// All CRUD functions follow this exact signature pattern:
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    return common.Create(ctx, d, meta, construct, meta.(*jamfpro.Client).CreatePolicy, readNoCleanup)
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
    return common.Read(ctx, d, meta, cleanup, meta.(*jamfpro.Client).GetPolicyByID, updateState)
}
```

### Singleton Configuration Resources

Some resources manage singleton configurations rather than individual entities (e.g., global settings). These follow a modified CRUD pattern:

**Characteristics:**
- Represent system-wide configuration that always exists
- Cannot be "created" or "deleted" in the traditional sense
- Use hardcoded singleton IDs

**Example Pattern:**
```go
// create calls UPDATE API method since configuration always exists
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*jamfpro.Client)
    
    settings, err := construct(d)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to construct settings: %v", err))
    }

    err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
        _, apiErr := client.UpdateComputerInventoryCollectionSettings(settings)
        if apiErr != nil {
            return retry.RetryableError(apiErr)
        }
        return nil
    })

    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to apply settings: %v", err))
    }

    // Use descriptive singleton ID
    d.SetId("jamfpro_computer_inventory_collection_settings_singleton")
    return readNoCleanup(ctx, d, meta)
}

// delete only removes from Terraform state, doesn't delete from API
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // Since this represents configuration and not an entity that can be deleted,
    // simply remove from Terraform state
    d.SetId("")
    return nil
}
```

**Singleton ID Naming:**
- Format: `jamfpro_{resource_name}_singleton`
- Examples: 
  - `jamfpro_computer_inventory_collection_settings_singleton`
  - `jamfpro_device_communication_settings_singleton`
  - `jamfpro_client_checkin_singleton`

