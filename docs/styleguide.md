# Go Styleguide

This styleguide defines project-specific conventions for Go code in this repository. For any areas not covered here, refer to the [Google Go Style Guide](https://google.github.io/styleguide/go/index.html) and [Effective Go](https://go.dev/doc/effective_go).

---

## 1. Project Structure
- **One logical block per file**: Separate schemas, state, constructors, and validators into their own files.
  - Example: `schema_self_service.go` contains only the self-service schema.
- **File naming**: Use underscores to separate logical parts. E.g., `schema_self_service.go`, `state_general.go`.
- **Package names**: All lowercase, no underscores. E.g., `policy`, `common`.

## 2. Naming Conventions
- **Acronyms**: Capitalize acronyms in names (e.g., `ID`, `URL`).
  - Example: `userID`, `getPolicyByID()`
- **Constants**: Use `CamelCase` or `ALL_CAPS` for acronyms.
  - Example: `DefaultTimeout`, `MAX_RETRIES`

## 3. Functions and Methods
- **Short receiver names**: Use a single letter for method receivers, typically the first letter of the type.
  - Example: `func (p *Policy) Validate() error { ... }`
- **Helpers for repeated logic**: Extract repeated logic into helper functions.
  - Example:
    ```go
    func validateDateTime(v interface{}, k string) (warns []string, errs []error) { ... }
    // Used in multiple schema files
    ```

## 4. Error Handling
- **TODOs for improvements**: If error handling is incomplete or needs review, leave a `TODO` comment.
  - Example:
    ```go
    // TODO: improve error handling for edge cases
    ```

## 5. Comments and Documentation
- **TODOs and FIXMEs**: Use these for future work or known issues.
  - Example:
    ```go
    // TODO: remove log.prints, debug use only
    // FIXME: handle nil pointer dereference
    ```

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
    type User struct { Base; Name string }
    ```
- **No repeated schemas**: If a schema is used more than once, move it to `/common` and import it.
  - Example:
    ```go
    // In /common/schema_shared.go
    func SharedUserSchema() *schema.Resource { ... }
    // In other packages
    import ".../common"
    common.SharedUserSchema()
    ```

## 7. Third-Party Dependencies
- **Pin versions in `go.mod`**: Always specify versions for dependencies in `go.mod`.
  - Example:
    ```go
    require github.com/hashicorp/terraform-plugin-sdk/v2 v2.15.0
    ```
  - See: [Go Modules Reference](https://go.dev/ref/mod)

## 8. Security Best Practices
- **Validate input when required by the API**: If the API is specific about accepted values or formats, always validate user input accordingly. If the API is permissive, validation is optional.
  - Example:
    ```go
    "activation_date": {
        Type: schema.TypeString,
        ValidateFunc: validateDateTime,
    }
    ```

## 9. Code Reviews and Pull Requests
- **Use TODOs/comments for review notes**: Leave clear comments for reviewers about areas needing attention or improvement.
  - Example:
    ```go
    // TODO: review this logic for concurrency issues
    ```

## 10. General Guidance
- **Refer to the Google Go styleguide at all opportunities**: [Google Go Style Guide](https://google.github.io/styleguide/go/index.html)
- **Refer to Effective Go**: [Effective Go](https://go.dev/doc/effective_go)
- **Go Proverbs**: [Go Proverbs](https://go-proverbs.github.io/)

---

## 11. Example: Resource Usage and Structure

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

For any questions or ambiguities, prefer the Google Go styleguide and Effective Go as the final authority.
