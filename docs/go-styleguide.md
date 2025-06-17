# Go Styleguide

This styleguide defines project-specific conventions for Go code in this repository. For any areas not covered here, refer to the [Google Go Style Guide](https://google.github.io/styleguide/go/index.html) and [Effective Go](https://go.dev/doc/effective_go).

---

## 1. Project Structure
- **One logical block per file**: Separate schemas, state, constructors, and validators into their own files.
  - Example: `schema_self_service.go` contains only the self-service schema.
- **File naming**: Use underscores to separate logical parts. E.g., `schema_self_service.go`, `state_general.go`.
- **Package names**: All lowercase, no underscores. E.g., `policy`, `common`.

## 2. Naming Conventions
- **Acronyms**: Capitalize acronyms in names (e.g., `ID`, `URL`, `API`, `HTTP`).
  - Example: `userID`, `getPolicyByID()`, `validateGUID()`
- **Constants**: Use `CamelCase` for exported constants, `ALL_CAPS` for package-level constants.
  - Example: `DefaultTimeout`, `MAX_RETRIES`
- **Booleans**: Prefix with `is`, `has`, `can`, or `should` when appropriate.
  - Example: `isEnabled`, `hasPermission`, `canDelete`

## 3. Functions and Methods
- **Short receiver names**: Use a single letter for method receivers, typically the first letter of the type.
  - Example: `func (p *Policy) Validate() error { ... }`
- **Validation functions**: Use consistent naming patterns for validators.
  - Schema validation functions: `validate<Field>() schema.SchemaValidateFunc`
  - Custom diff validators: `validate<Requirement>(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error`
  - Example:
    ```go
    func validateDateTime(v interface{}, k string) (warns []string, errs []error) { ... }
    func validateAuthenticationPrompt(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error { ... }
    ```
- **Helpers for repeated logic**: Extract repeated logic into helper functions.
- **Function organization**: Group related functions together and separate with comments when logical.

## 4. Error Handling
- **Descriptive error messages**: Include context about what failed and why.
  - Example: `fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))`
- **Error wrapping**: Use `fmt.Errorf("operation failed: %w", err)` for wrapping errors.
- **TODOs for improvements**: If error handling is incomplete or needs review, leave a `TODO` comment.
  - Example:
    ```go
    // TODO: improve error handling for edge cases
    ```

## 5. Comments and Documentation
- **Package documentation**: Every package should have a package comment.
- **Function documentation**: Public functions should have doc comments starting with the function name.
- **TODO formatting**: Use consistent format for TODOs and FIXMEs.
  - `TODO:` followed by a space and description
  - `FIXME:` for bugs that need immediate attention
  - Example:
    ```go
    // TODO: remove log.prints, debug use only
    // FIXME: handle nil pointer dereference
    ```
- **Inline comments**: Use sparingly, only when the code's purpose isn't self-explanatory.

## 6. Structs and Interfaces
- **Struct tags for schemas**: Always use struct tags for schema definitions.
  - Example:
    ```go
    type User struct {
        Name string `json:"name"`
        ID   int    `json:"id"`
    }
    ```
- **Prefer composition over inheritance**: Embed structs to share functionality.
  - Example:
    ```go
    type Base struct { ID int }
    type User struct { 
        Base
        Name string 
    }
    ```
- **No repeated schemas**: If a schema is used more than once, move it to `/common/sharedschemas` and import it.
  - Example:
    ```go
    // In /common/sharedschemas/category.go
    func GetCategorySchema() *schema.Schema { ... }
    
    // In other packages
    import "internal/resources/common/sharedschemas"
    "category_id": sharedschemas.GetCategorySchema(),
    ```

## 7. Third-Party Dependencies
- **Pin versions in `go.mod`**: Always specify versions for dependencies in `go.mod`.
  - Example:
    ```go
    require github.com/hashicorp/terraform-plugin-sdk/v2 v2.15.0
    ```
  - See: [Go Modules Reference](https://go.dev/ref/mod)
- **Minimal dependencies**: Only add dependencies that are absolutely necessary.

## 8. Security and Validation
- **Input validation**: Always validate user input when required by the Jamf Pro API.
  - If the API has specific requirements, validate accordingly
  - If the API is permissive, validation is optional but recommended for UX
  - Example:
    ```go
    "activation_date": {
        Type:         schema.TypeString,
        ValidateFunc: validateDateTime,
    }
    ```
- **Sensitive data**: Never log sensitive information like passwords or tokens.
- **GUID validation**: Use the existing `validateGUID()` function for UUID/GUID fields.

## 9. Testing
- **Test file naming**: Use `_test.go` suffix for test files.
- **Test function naming**: Prefix with `Test` for unit tests, `TestAcc` for acceptance tests.
- **Table-driven tests**: Use table-driven tests for multiple test cases.
- **Test helpers**: Extract common test setup into helper functions.

## 10. Code Reviews and Pull Requests
- **Self-review**: Always review your own code before requesting review.
- **Clear commit messages**: Use conventional commit format (see PR guidelines).
- **Review comments**: Leave clear comments for reviewers about areas needing attention.
  - Example:
    ```go
    // TODO: review this logic for concurrency issues
    ```

## 11. Performance Considerations
- **Avoid unnecessary allocations**: Reuse slices and maps where possible.
- **Context usage**: Always respect context cancellation in long-running operations.
- **Resource cleanup**: Always clean up resources (close files, connections, etc.).

## 12. Terraform Provider Specific
- **Schema organization**: Separate complex schemas into logical files.
- **State management**: Keep state functions focused and testable.
- **Custom diffs**: Use custom diff functions for complex validation logic.
- **Resource naming**: Use consistent naming patterns for resource files and functions.

## 13. General Guidance
- **Refer to the Google Go styleguide**: [Google Go Style Guide](https://google.github.io/styleguide/go/index.html)
- **Refer to Effective Go**: [Effective Go](https://go.dev/doc/effective_go)
- **Go Proverbs**: [Go Proverbs](https://go-proverbs.github.io/)
- **Terraform Plugin Development**: [Terraform Plugin SDK](https://developer.hashicorp.com/terraform/plugin)

---

## 14. Example: Resource Usage and Structure

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
    all_mobile_devices = true
    all_jss_users      = false
    mobile_device_ids       = [101, 102, 103]
    mobile_device_group_ids = [201, 202]
    building_ids            = [301]
    department_ids          = [401, 402]
    jss_user_ids            = [501, 502]
    jss_user_group_ids      = [601, 602]
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

## 15. Quick Reference Checklist

Before submitting code, verify:
- [ ] Functions follow naming conventions
- [ ] Error messages are descriptive
- [ ] No repeated schema definitions (use shared schemas)
- [ ] Validation functions follow established patterns
- [ ] Comments are necessary and well-formatted
- [ ] Dependencies are pinned in `go.mod`
- [ ] Code follows project structure guidelines
- [ ] Tests are included for new functionality

For any questions or ambiguities, prefer the Google Go styleguide and Effective Go as the final authority.
