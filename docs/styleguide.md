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

For any questions or ambiguities, prefer the Google Go styleguide and Effective Go as the final authority.
