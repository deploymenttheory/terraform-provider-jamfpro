# Go Styleguide

This styleguide defines project-specific conventions for Go code in this repository. For any areas not covered here, refer to the [Google Go Style Guide](https://google.github.io/styleguide/go/index.html) and [Effective Go](https://go.dev/doc/effective_go).

---

## 1. File Organization and Structure

### 1.1 Resource File Structure

Each resource should follow this standardized file organization pattern:

```
internal/resources/<resource_name>/
├── resource.go                    # Main resource definition, schema, CRUD operations
├── constructor.go                 # Data constructors and transformations
├── state.go                      # State management functions
├── data_validator.go             # Custom validation functions
├── data_customdiff.go            # Custom diff functions (if needed)
├── state_migration.go            # State migration logic (if needed)
└── schema_<component>.go         # Component-specific schemas (for complex resources)
```

### 1.2 File Purposes

**`resource.go`**
- Main resource definition
- Schema definition (or imports from schema files)
- CRUD operations (`Create`, `Read`, `Update`, `Delete`)
- Resource configuration and metadata

**`constructor.go`**
- Functions that transform Terraform data to API structs
- Functions that transform API responses to Terraform state
- Data mapping and conversion logic
- Example naming: `constructResourceFromState()`, `updateStateFromResponse()`

**`state.go` (or `state_<component>.go`)**
- State management helper functions
- State update and retrieval logic
- For complex resources, split into logical components:
  - `state_general.go` - basic resource state
  - `state_scope.go` - scoping-related state
  - `state_self_service.go` - self-service specific state

**`data_validator.go`**
- Custom validation functions
- Schema validation helpers
- Input sanitization functions

**`data_customdiff.go`**
- Custom diff functions for complex validation
- Inter-field dependency validation
- Conditional validation logic

**`schema_<component>.go`**
- Component-specific schema definitions for complex resources
- Used when the main schema would be too large for a single file
- Examples: `schema_self_service.go`, `schema_scope.go`, `schema_payloads.go`

### 1.3 Shared Components

**Location: `internal/resources/common/`**

```
internal/resources/common/
├── sharedschemas/                 # Reusable schema components
│   ├── category.go               # Category schema
│   ├── site.go                   # Site schema
│   ├── computerscope.go          # Computer scoping schemas
│   ├── mobiledevicescope.go      # Mobile device scoping schemas
│   └── utilities.go              # Schema utility functions
├── constructors/                  # Shared constructor functions
├── configurationprofiles/         # Configuration profile utilities
│   └── plist/                    # Plist handling utilities
└── jamfprivileges/               # Privilege validation utilities
```

**When to Use Shared Components:**
- **Schemas**: When the same schema is used in 2+ resources
- **Constructors**: For common data transformation patterns
- **Validators**: For validation logic used across multiple resources

### 1.4 Directory Structures

**Data Sources:**
```
internal/resources/<resource_name>/
├── data_source.go                # Data source definition and read logic
├── data_source_schema.go         # Schema specific to data source (if different from resource)
└── constructor.go                # Shared with resource (if applicable)
```

**Testing:**
```
internal/resources/<resource_name>/
├── resource_test.go              # Acceptance tests
├── constructor_test.go           # Unit tests for constructors
├── data_validator_test.go        # Unit tests for validators
└── testdata/                     # Test fixtures and mock data
    ├── valid_payload.json
    └── invalid_payload.json
```

**Documentation:**
```
docs/
├── data-sources/
│   └── <resource_name>.md        # Data source documentation
├── resources/
│   └── <resource_name>.md        # Resource documentation
└── styleguide.md                 # This file
```

**Examples:**
```
examples/
├── data-sources/
│   └── <resource_name>/
│       └── data-source.tf        # Example data source usage
├── resources/
│   └── <resource_name>/
│       ├── resource.tf           # Basic resource example
│       ├── complex-example.tf    # Advanced usage example
│       └── payloads/             # Example payload files
└── provider/
    └── provider.tf               # Provider configuration examples
```

---

## 2. Naming Conventions

### 2.1 File Naming

**Go Files:**
- `resource.go` - Main resource file
- `data_source.go` - Data source file
- `constructor.go` - Constructor functions
- `state_<component>.go` - State management
- `schema_<component>.go` - Schema definitions
- `data_validator.go` - Validation functions
- `data_customdiff.go` - Custom diff functions
- `state_migration.go` - State migration
- `helpers.go` - General helper functions

**Test Files:**
- `<filename>_test.go` - Unit tests
- `resource_test.go` - Acceptance tests
- `data_source_test.go` - Data source tests

**Documentation Files:**
- `<resource_name>.md` - Resource/data source documentation
- Use hyphens for multi-word resource names in file paths

### 2.2 Code Naming

**Acronyms:** Capitalize acronyms in names (e.g., `ID`, `URL`, `API`, `HTTP`)
- Example: `userID`, `getPolicyByID()`, `validateGUID()`

**Constants:** Use `CamelCase` for exported constants, `ALL_CAPS` for package-level constants
- Example: `DefaultTimeout`, `MAX_RETRIES`

**Booleans:** Prefix with `is`, `has`, `can`, or `should` when appropriate
- Example: `isEnabled`, `hasPermission`, `canDelete`

---

## 3. Package Organization

### 3.1 Rules

1. **One resource per package**: Each resource gets its own package directory
2. **Shared utilities**: Place in `internal/resources/common/`
3. **No circular dependencies**: Ensure clean dependency hierarchy
4. **Logical grouping**: Group related functionality within files

### 3.2 Import Organization

```go
// Standard library imports
import (
    "context"
    "fmt"
    "strings"
)

// Third-party imports
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Local imports
import (
    "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
    "github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
)
```

---

## 4. Functions and Methods

### 4.1 General Guidelines

**Short receiver names:** Use a single letter for method receivers, typically the first letter of the type
- Example: `func (p *Policy) Validate() error { ... }`

**Helpers for repeated logic:** Extract repeated logic into helper functions

**Function organization:** Group related functions together and separate with comments when logical

### 4.2 Validation Functions

Use consistent naming patterns for validators:

**Schema validation functions:**
```go
func validate<Field>() schema.SchemaValidateFunc
```

**Custom diff validators:**
```go
func validate<Requirement>(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error
```

**Example:**
```go
func validateDateTime(v interface{}, k string) (warns []string, errs []error) { ... }
func validateAuthenticationPrompt(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error { ... }
```

---

## 5. Error Handling

**Descriptive error messages:** Include context about what failed and why
- Example: `fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))`

**Error wrapping:** Use `fmt.Errorf("operation failed: %w", err)` for wrapping errors

**TODOs for improvements:** If error handling is incomplete or needs review, leave a `TODO` comment
- Example:
  ```go
  // TODO: improve error handling for edge cases
  ```

---

## 6. Comments and Documentation

**Package documentation:** Every package should have a package comment

**Function documentation:** Public functions should have doc comments starting with the function name

**TODO formatting:** Use consistent format for TODOs and FIXMEs
- `TODO:` followed by a space and description
- `FIXME:` for bugs that need immediate attention
- Example:
  ```go
  // TODO: remove log.prints, debug use only
  // FIXME: handle nil pointer dereference
  ```

**Inline comments:** Use sparingly, only when the code's purpose isn't self-explanatory

---

## 7. Structs and Interfaces

**Struct tags for schemas:** Always use struct tags for schema definitions
- Example:
  ```go
  type User struct {
      Name string `json:"name"`
      ID   int    `json:"id"`
  }
  ```

**Prefer composition over inheritance:** Embed structs to share functionality
- Example:
  ```go
  type Base struct { ID int }
  type User struct { 
      Base
      Name string 
  }
  ```

**No repeated schemas:** If a schema is used more than once, move it to `/common/sharedschemas` and import it
- Example:
  ```go
  // In /common/sharedschemas/category.go
  func GetCategorySchema() *schema.Schema { ... }
  
  // In other packages
  import "internal/resources/common/sharedschemas"
  "category_id": sharedschemas.GetCategorySchema(),
  ```

---

## 8. Dependencies

**Pin versions in `go.mod`:** Always specify versions for dependencies in `go.mod`
- Example:
  ```go
  require github.com/hashicorp/terraform-plugin-sdk/v2 v2.15.0
  ```

**Minimal dependencies:** Only add dependencies that are absolutely necessary

---

## 9. Security and Validation

**Input validation:** Always validate user input when required by the Jamf Pro API
- If the API has specific requirements, validate accordingly
- If the API is permissive, validation is optional but recommended for UX
- Example:
  ```go
  "activation_date": {
      Type:         schema.TypeString,
      ValidateFunc: validateDateTime,
  }
  ```

**Sensitive data:** Never log sensitive information like passwords or tokens

**GUID validation:** Use the existing `validateGUID()` function for UUID/GUID fields

---

## 10. Testing

**Table-driven tests:** Use table-driven tests for multiple test cases

**Test helpers:** Extract common test setup into helper functions

---

## 11. Code Reviews

**Self-review:** Always review your own code before requesting review

**Clear commit messages:** Use conventional commit format (see PR guidelines)

**Review comments:** Leave clear comments for reviewers about areas needing attention

---

## 12. Performance

**Avoid unnecessary allocations:** Reuse slices and maps where possible

**Context usage:** Always respect context cancellation in long-running operations

**Resource cleanup:** Always clean up resources (close files, connections, etc.)

---

## 13. Terraform Provider Specific

**Schema organization:** Separate complex schemas into logical files

**State management:** Keep state functions focused and testable

**Custom diffs:** Use custom diff functions for complex validation logic

**Resource naming:** Use consistent naming patterns for resource files and functions

---

## 14. Example Resource Structure

For examples of resource usage and structure, see [`examples/resources/jamfpro_mobile_device_configuration_profile_plist/resource.tf`](examples/resources/jamfpro_mobile_device_configuration_profile_plist/resource.tf). This file demonstrates best practices for resource definition, block structure, and field usage in Terraform for this provider.

```hcl
// Example of creating a mobile device configuration profile in Jamf Pro for self service using a plist source file
resource "jamfpro_mobile_device_configuration_profile_plist" "mobile_device_configuration_profile_001" {
  name               = "your-mobile_device_configuration_profile-name"
  description        = "An example mobile device configuration profile."
  deployment_method  = "Install Automatically"
  level              = "Device Level"
  redeploy_on_update = "Newly Assigned"
  payloads           = file("${path.module}/path/to/your.mobileconfig")

  // Optional Block
  site_id = 967
  // Optional Block
  category_id = 5
  
  scope {
    all_mobile_devices      = true
    all_jss_users          = false
    mobile_device_ids      = [101, 102, 103]
    mobile_device_group_ids = [201, 202]
    building_ids           = [301]
    department_ids         = [401, 402]
    jss_user_ids           = [501, 502]
    jss_user_group_ids     = [601, 602]
    
    limitations {
      network_segment_ids                  = [701, 702]
      ibeacon_ids                          = [801]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [1001, 1002]
    }
    
    exclusions {
      mobile_device_ids                    = [1101, 1102]
      mobile_device_group_ids              = [1201]
      building_ids                         = [1301, 1302]
      department_ids                       = [1401]
      network_segment_ids                  = [1501, 1502]
      jss_user_ids                         = [1601, 1602]
      jss_user_group_ids                   = [1701]
      directory_service_or_local_usernames = ["Jane Smith", "John Doe"]
      directory_service_usergroup_ids      = [1001, 1002]
      ibeacon_ids                          = [1801]
    }
  }
}
```
