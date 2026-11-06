# Code review follow-up (2026-09-20)

Three issues were reproduced in local regression tests and fixed:

1. Real-valued array numbers such as `[1.0,1e0,1.00]` were read back as `[1,1,1]`, changing the set identity and losing the integer/real distinction on later construction. Parsing, validation, hashing and serialization now use one numeric representation. Real markers survive readback; equivalent real spellings normalize consistently.
2. `[null,"a"]` serialized to `["a"]` without an error. JSON null and unsupported numeric ranges are now rejected before apply, including nested values. Array order and duplicates remain unchanged.
3. A v0 payload with null `payload_content` panicked during migration. Null payload containers are now preserved without assertions on nil interfaces.

Regression tests include complete HCL-to-plist-to-HCL round trips, exact large integers, integer/real hash identity, order/duplicate sensitivity, invalid values and null migration containers.

## Scoped tests

Command: `go test ./internal/common/... ./internal/services/macos_configuration_profile_plist_generator ./internal/provider`

```text
ok      github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors (cached)
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/crypto   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/enrollment_lock  [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/errors   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/files    [no test files]
ok      github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud   (cached)
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/images   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/jamf_privileges  [no test files]
ok      github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist    0.477s
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist/test/encode    [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist/test/plistsanitizefortfstate   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist/test/removekeys    [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist/test/sortkeys  [no test files]
ok      github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/redact   (cached)
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/validation    [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2_crud   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas   [no test files]
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/utils    [no test files]
ok      github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/macos_configuration_profile_plist_generator    0.802s
?       github.com/deploymenttheory/terraform-provider-jamfpro/internal/provider    [no test files]
```

Provider build also passed. The tenant was not mutated during this code review; the earlier live records describe the earlier build. The new edge cases above were verified locally, not claimed as fresh live-tenant validation.

## Full repository tests

`go test ./...` still fails in the unchanged raw-plist resource test below. The same focused test fails on the pristine `f8efc66ae7a7242ac41ff25727bc3e555dc26c85` baseline. This failure is outside the generator changes and has not been fixed here.

Excerpt from the actual full-suite output (workspace path and whitespace normalized):

```text
--- FAIL: TestDiffSuppressEquivalentPayloads (0.00s)
    --- FAIL: TestDiffSuppressEquivalentPayloads/Validation_disabled_should_not_suppress (0.00s)
        resource_payload_diff_suppress_test.go:255:
                Error Trace:    <test-workspace>/macos_configuration_profile_plist_generator/internal/services/macos_configuration_profile_plist/resource_payload_diff_suppress_test.go:255
                Error:          Not equal:
                                expected: false
                                actual  : true
                Test:           TestDiffSuppressEquivalentPayloads/Validation_disabled_should_not_suppress
                Messages:       DiffSuppressPayloads() with payloadValidate=false returned true, want false
FAIL
FAIL    github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/macos_configuration_profile_plist  0.780s
```

The live nested-dictionary migration and array-plan display limitations already documented in the PR remain open. This follow-up does not remove those verification limits.
